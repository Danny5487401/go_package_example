<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [jackc/pgx](#jackcpgx)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [pgx 连接池初始化](#pgx-%E8%BF%9E%E6%8E%A5%E6%B1%A0%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [查询过程](#%E6%9F%A5%E8%AF%A2%E8%BF%87%E7%A8%8B)
  - [监控](#%E7%9B%91%E6%8E%A7)
  - [测试用例编写](#%E6%B5%8B%E8%AF%95%E7%94%A8%E4%BE%8B%E7%BC%96%E5%86%99)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# jackc/pgx


github.com/jackc/pgx 是一个高性能的PostgreSQL数据库驱动，它既可作为原生PostgreSQL接口使用，也可作为database/sql兼容驱动使用。


## 特点

- 性能比 database/sql 强
- 原生支持notifications, COPY protocol等
- 连接池


## pgx 连接池初始化

pgx.Conn代表单个数据库连接，不是并发安全的。对于需要并发访问的场景，应当使用专门的连接池实现。


```go
// github.com/jackc/pgx/v5@v5.7.6/pgxpool/pool.go

func New(ctx context.Context, connString string) (*Pool, error) {
	// 解析配置
	config, err := ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	return NewWithConfig(ctx, config)
}

// 根据配置初始化
func NewWithConfig(ctx context.Context, config *Config) (*Pool, error) {
    // ...

	p := &Pool{
		config:                config,
		beforeConnect:         config.BeforeConnect,
		afterConnect:          config.AfterConnect,
		prepareConn:           prepareConn,
		afterRelease:          config.AfterRelease,
		beforeClose:           config.BeforeClose,
		minConns:              config.MinConns,
		minIdleConns:          config.MinIdleConns,
		maxConns:              config.MaxConns,
		maxConnLifetime:       config.MaxConnLifetime,
		maxConnLifetimeJitter: config.MaxConnLifetimeJitter,
		maxConnIdleTime:       config.MaxConnIdleTime,
		healthCheckPeriod:     config.HealthCheckPeriod,
		healthCheckChan:       make(chan struct{}, 1),
		closeChan:             make(chan struct{}),
	}

    // ..

	// 使用 puddle 库创建资源池
	var err error
	p.p, err = puddle.NewPool(
		&puddle.Config[*connResource]{
			Constructor: func(ctx context.Context) (*connResource, error) {
				atomic.AddInt64(&p.newConnsCount, 1)
				connConfig := p.config.ConnConfig.Copy()

				// Connection will continue in background even if Acquire is canceled. Ensure that a connect won't hang forever.
				if connConfig.ConnectTimeout <= 0 {
					connConfig.ConnectTimeout = 2 * time.Minute
				}

				// 连接前
				if p.beforeConnect != nil {
					if err := p.beforeConnect(ctx, connConfig); err != nil {
						return nil, err
					}
				}

				// 实际建立连接
				conn, err := pgx.ConnectConfig(ctx, connConfig)
				if err != nil {
					return nil, err
				}

				// 连接后
				if p.afterConnect != nil {
					err = p.afterConnect(ctx, conn)
					if err != nil {
						conn.Close(ctx)
						return nil, err
					}
				}

				jitterSecs := rand.Float64() * config.MaxConnLifetimeJitter.Seconds()
				maxAgeTime := time.Now().Add(config.MaxConnLifetime).Add(time.Duration(jitterSecs) * time.Second)

				cr := &connResource{
					conn:       conn,
					conns:      make([]Conn, 64),
					poolRows:   make([]poolRow, 64),
					poolRowss:  make([]poolRows, 64),
					maxAgeTime: maxAgeTime,
				}

				return cr, nil
			},
			Destructor: func(value *connResource) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				conn := value.conn
				if p.beforeClose != nil {
					// hook 函数处理
					p.beforeClose(conn)
				}
				conn.Close(ctx)
				select {
				case <-conn.PgConn().CleanupDone():
				case <-ctx.Done():
				}
				cancel()
			},
			MaxSize: config.MaxConns,
		},
	)
	if err != nil {
		return nil, err
	}

	go func() {
		targetIdleResources := max(int(p.minConns), int(p.minIdleConns))
		p.createIdleResources(ctx, targetIdleResources) // 创建资源
		p.backgroundHealthCheck() // 后台健康检查
	}()

	return p, nil
}
```


```go
// ithub.com/jackc/pgx/v5@v5.7.6/pgxpool/pool.go

func (p *Pool) createIdleResources(parentCtx context.Context, targetResources int) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	errs := make(chan error, targetResources)

	for i := 0; i < targetResources; i++ {
		go func() {
			// 创建资源
			err := p.p.CreateResource(ctx)
			// Ignore ErrNotAvailable since it means that the pool has become full since we started creating resource.
			if err == puddle.ErrNotAvailable {
				err = nil
			}
			errs <- err
		}()
	}

	var firstError error
	for i := 0; i < targetResources; i++ {
		err := <-errs
		if err != nil && firstError == nil {
			cancel()
			firstError = err
		}
	}

	return firstError
}
```

```go
// github.com/jackc/puddle/v2@v2.2.2/pool.go

func (p *Pool[T]) CreateResource(ctx context.Context) error {
    // 校验 ...

	res := p.createNewResource()
	p.mux.Unlock()

	// 调用初始化方法
	value, err := p.constructor(ctx)
	p.mux.Lock()
	defer p.mux.Unlock()
	defer p.acquireSem.Release(1)
	if err != nil {
		p.allResources.remove(res)
		p.destructWG.Done()
		return err
	}

	res.value = value
	res.status = resourceStatusIdle

	// If closed while constructing resource then destroy it and return an error
	if p.closed {
		go p.destructResourceValue(res.value)
		return ErrClosedPool
	}

	p.idleResources.Push(res)

	return nil
}


func (p *Pool[T]) createNewResource() *Resource[T] {
	res := &Resource[T]{
		pool:           p,
		creationTime:   time.Now(),
		lastUsedNano:   nanotime(),
		poolResetCount: p.resetCount,
		status:         resourceStatusConstructing,
	}

	p.allResources.append(res)
	p.destructWG.Add(1)

	return res
}
```


## 查询过程


## 监控

```go
// github.com/jackc/pgx/v5@v5.7.6/pgxpool/stat.go
type Stat struct {
	s                    *puddle.Stat
	newConnsCount        int64 // 连接创建计数器
	lifetimeDestroyCount int64 // 因超过MaxConnLifetime被销毁的连接数
	idleDestroyCount     int64 // 因超过MaxConnIdleTime被销毁的连接数
}
```



```go
// github.com/jackc/puddle/v2@v2.2.2/pool.go
type Stat struct {
	constructingResources int32
	acquiredResources     int32
	idleResources         int32
	maxResources          int32
	acquireCount          int64
	acquireDuration       time.Duration
	emptyAcquireCount     int64
	emptyAcquireWaitTime  time.Duration
	canceledAcquireCount  int64
}

```

## 测试用例编写
- github.com/pashagolub/pgxmock

- github.com/jackc/pgmock

## 参考
- [pgx 使用指南](https://betterstack.com/community/guides/scaling-go/postgresql-pgx-golang/)
- [pgx连接池与性能优化实战指南](https://blog.csdn.net/gitblog_01117/article/details/150753396)
- [pgx连接池监控：Stat指标与性能优化](https://blog.csdn.net/gitblog_00159/article/details/151239729)
