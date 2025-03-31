<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/ClickHouse/clickhouse-go](#githubcomclickhouseclickhouse-go)
  - [v1 对比 v2](#v1-%E5%AF%B9%E6%AF%94-v2)
  - [特性](#%E7%89%B9%E6%80%A7)
  - [dsn 配置](#dsn-%E9%85%8D%E7%BD%AE)
  - [两种接口](#%E4%B8%A4%E7%A7%8D%E6%8E%A5%E5%8F%A3)
  - [v1 版本(不建议使用)](#v1-%E7%89%88%E6%9C%AC%E4%B8%8D%E5%BB%BA%E8%AE%AE%E4%BD%BF%E7%94%A8)
  - [v2](#v2)
    - [批量写入](#%E6%89%B9%E9%87%8F%E5%86%99%E5%85%A5)
  - [第三方应用-->grafana clickhouse plugin](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8--grafana-clickhouse-plugin)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/ClickHouse/clickhouse-go



## v1 对比 v2 
https://github.com/ClickHouse/clickhouse-go/blob/main/v1_v2_CHANGES.md

- v1 有精度损失
- strings 在v2不允许插入Date or DateTime columns
- 数组必须类型说明.[]any containing strings cannot be inserted into a string column
- 默认连接策略不同, v2 使用 ConnOpenInOrder
```go
type ConnOpenStrategy uint8

const (
	ConnOpenInOrder ConnOpenStrategy = iota
	ConnOpenRoundRobin
	ConnOpenRandom
)
```

```go
func DefaultDialStrategy(ctx context.Context, connID int, opt *Options, dial Dial) (r DialResult, err error) {
	random := rand.Int()
	for i := range opt.Addr {
		var num int
		switch opt.ConnOpenStrategy {
		case ConnOpenInOrder: // 顺序
			num = i
		case ConnOpenRoundRobin: // 轮询
			num = (int(connID) + i) % len(opt.Addr)
		case ConnOpenRandom: // 随机
			num = (random + i) % len(opt.Addr)
		}

		if r, err = dial(ctx, opt.Addr[num], opt); err == nil {
			return r, nil
		}
	}

	if err == nil {
		err = ErrAcquireConnNoAddress
	}

	return r, err
}
```


## 特性

 
- 可以 rows 反序列化成结构体 (ScanStruct, Select)
- 可以 结构体 反序列化成 row  (AppendStruct)
- 连接池
- 批量写
- LZ4/ZSTD 压缩支持
- 兼容 database/sql (但是比 native interface 慢!)


## dsn 配置
```
clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60

```
https://github.com/ClickHouse/clickhouse-go/?tab=readme-ov-file#dsn 
- dial_timeout 默认30s
- connection_open_strategy: 默认 in_order
- compress : 默认未压缩


解析 dsn
```go
func (o *Options) fromDSN(in string) error {
	dsn, err := url.Parse(in)
	if err != nil {
		return err
	}

	if dsn.Host == "" {
		return errors.New("parse dsn address failed")
	}

	if o.Settings == nil {
		o.Settings = make(Settings)
	}
	if dsn.User != nil {
		o.Auth.Username = dsn.User.Username()
		o.Auth.Password, _ = dsn.User.Password()
	}
	o.Addr = append(o.Addr, strings.Split(dsn.Host, ",")...)
	var (
		secure     bool
		params     = dsn.Query()
		skipVerify bool
	)
	o.Auth.Database = strings.TrimPrefix(dsn.Path, "/")

	for v := range params {
		switch v {
		case "debug":
			o.Debug, _ = strconv.ParseBool(params.Get(v))
		case "compress":
			if on, _ := strconv.ParseBool(params.Get(v)); on {
				if o.Compression == nil {
					o.Compression = &Compression{}
				}

				o.Compression.Method = CompressionLZ4
				continue
			}
			if compressMethod, ok := compressionMap[params.Get(v)]; ok {
				if o.Compression == nil {
					o.Compression = &Compression{
						// default for now same as Clickhouse - https://clickhouse.com/docs/en/operations/settings/settings#settings-http_zlib_compression_level
						Level: 3,
					}
				}

				o.Compression.Method = compressMethod
			}
        // ... 其他参数
		default:
			switch p := strings.ToLower(params.Get(v)); p {
			case "true":
				o.Settings[v] = int(1)
			case "false":
				o.Settings[v] = int(0)
			default:
				if n, err := strconv.Atoi(p); err == nil {
					o.Settings[v] = n
				} else {
					o.Settings[v] = p
				}
			}
		}
	}
	if secure {
		o.TLS = &tls.Config{
			InsecureSkipVerify: skipVerify,
		}
	}
	o.scheme = dsn.Scheme
	switch dsn.Scheme {
	case "http":
		if secure {
			return fmt.Errorf("clickhouse [dsn parse]: http with TLS specify")
		}
		o.Protocol = HTTP
	case "https":
		if !secure {
			return fmt.Errorf("clickhouse [dsn parse]: https without TLS")
		}
		o.Protocol = HTTP
	default:
		o.Protocol = Native
	}
	return nil
}

```

## 两种接口

- native interface

- std database/sql interface
```go
// github.com/!click!house/clickhouse-go/v2@v2.33.1/clickhouse_std.go
func init() {
	var debugf = func(format string, v ...any) {}
	sql.Register("clickhouse", &stdDriver{debugf: debugf})
}

```

协议 
- http 协议(实验性) ,只支持 `database/sql`

## v1 版本(不建议使用)

初始化注册插件
```go
// clickhouse-go@v1.5.1/bootstrap.go
func init() {
	sql.Register("clickhouse", &bootstrap{})
	go func() {
		for tick := time.Tick(time.Second); ; {
			select {
			case <-tick:
				atomic.AddInt64(&unixtime, int64(time.Second))
			}
		}
	}()
}
```
连接
```go
type bootstrap struct{}

func (d *bootstrap) Open(dsn string) (driver.Conn, error) {
	return Open(dsn)
}

// Open the connection
func Open(dsn string) (driver.Conn, error) {
	clickhouse, err := open(dsn)
	if err != nil {
		return nil, err
	}

	return clickhouse, err
}

func open(dsn string) (*clickhouse, error) {
	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	// 配置默认参数
	var (
		hosts             = []string{url.Host}
		query             = url.Query()
		secure            = false
		skipVerify        = false
		tlsConfigName     = query.Get("tls_config")
		noDelay           = true
		compress          = false
		database          = query.Get("database")  // 数据库
		username          = query.Get("username")  //用户名
		password          = query.Get("password") //密码
		blockSize         = 1000000
		connTimeout       = DefaultConnTimeout
		readTimeout       = DefaultReadTimeout
		writeTimeout      = DefaultWriteTimeout
		connOpenStrategy  = connOpenRandom  //连接选择服务：默认随机
		checkConnLiveness = true
	)
    // ...
	
	var (
		ch = clickhouse{
			logf:              func(string, ...interface{}) {},
			settings:          settings,
			compress:          compress,
			blockSize:         blockSize,
			checkConnLiveness: checkConnLiveness,
			ServerInfo: data.ServerInfo{
				Timezone: time.Local,
			},
		}
		logger = log.New(logOutput, "[clickhouse]", 0)
	)
	if debug, err := strconv.ParseBool(url.Query().Get("debug")); err == nil && debug {
		ch.logf = logger.Printf
	}
	ch.logf("host(s)=%s, database=%s, username=%s",
		strings.Join(hosts, ", "),
		database,
		username,
	)
	options := connOptions{
		secure:       secure,
		tlsConfig:    tlsConfig,
		skipVerify:   skipVerify,
		hosts:        hosts,
		connTimeout:  connTimeout,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		noDelay:      noDelay,
		openStrategy: connOpenStrategy,
		logf:         ch.logf,
	}
	if ch.conn, err = dial(options); err != nil {
		return nil, err
	}
	logger.SetPrefix(fmt.Sprintf("[clickhouse][connect=%d]", ch.conn.ident))
	ch.buffer = bufio.NewWriter(ch.conn)

	ch.decoder = binary.NewDecoderWithCompress(ch.conn)
	ch.encoder = binary.NewEncoderWithCompress(ch.buffer)

	if err := ch.hello(database, username, password); err != nil {
		ch.conn.Close()
		return nil, err
	}
	return &ch, nil
}
```


## v2


初始化连接配置

```go
func Open(opt *Options) (driver.Conn, error) {
	if opt == nil {
		opt = &Options{}
	}
	// 设置默认值
	o := opt.setDefaults()
	conn := &clickhouse{
		opt:  o,
		idle: make(chan *connect, o.MaxIdleConns),
		open: make(chan struct{}, o.MaxOpenConns),
		exit: make(chan struct{}),
	}
	// 定期回收过期
	go conn.startAutoCloseIdleConnections()
	return conn, nil
}
```

实际建立连接

```go
func (ch *clickhouse) acquire(ctx context.Context) (conn *connect, err error) {
	timer := time.NewTimer(ch.opt.DialTimeout)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	select {
	case <-timer.C:
		return nil, ErrAcquireConnTimeout
	case <-ctx.Done():
		return nil, ctx.Err()
	case ch.open <- struct{}{}:
	}
	select {
	case <-timer.C:
		select {
		case <-ch.open:
		default:
		}
		return nil, ErrAcquireConnTimeout
	case conn := <-ch.idle:
		if conn.isBad() {
			conn.close()
			if conn, err = ch.dial(ctx); err != nil {
				select {
				case <-ch.open:
				default:
				}
				return nil, err
			}
		}
		conn.released = false
		return conn, nil
	default:
	}
	// 如果没有闲置连接
	if conn, err = ch.dial(ctx); err != nil {
		select {
		case <-ch.open:
		default:
		}
		return nil, err
	}
	return conn, nil
}


func (ch *clickhouse) dial(ctx context.Context) (conn *connect, err error) {
	connID := int(atomic.AddInt64(&ch.connID, 1))

	dialFunc := func(ctx context.Context, addr string, opt *Options) (DialResult, error) {
		conn, err := dial(ctx, addr, connID, opt)

		return DialResult{conn}, err
	}

	dialStrategy := DefaultDialStrategy
	if ch.opt.DialStrategy != nil {
		dialStrategy = ch.opt.DialStrategy
	}

	result, err := dialStrategy(ctx, connID, ch.opt, dialFunc)
	if err != nil {
		return nil, err
	}
	return result.conn, nil
}

```



### 批量写入

```go
func (ch *clickhouse) PrepareBatch(ctx context.Context, query string, opts ...driver.PrepareBatchOption) (driver.Batch, error) {
	// 获取连接
	conn, err := ch.acquire(ctx)
	if err != nil {
		return nil, err
	}
	batch, err := conn.prepareBatch(ctx, query, getPrepareBatchOptions(opts...), ch.release, ch.acquire)
	if err != nil {
		return nil, err
	}
	return batch, nil
}
```

这里选择原生 Native 协议
```go

func (c *connect) prepareBatch(ctx context.Context, query string, opts driver.PrepareBatchOptions, release func(*connect, error), acquire func(context.Context) (*connect, error)) (driver.Batch, error) {
	query, _, queryColumns, verr := extractNormalizedInsertQueryAndColumns(query)
	if verr != nil {
		return nil, verr
	}

	options := queryOptions(ctx)
	if deadline, ok := ctx.Deadline(); ok {
		c.conn.SetDeadline(deadline)
		defer c.conn.SetDeadline(time.Time{})
	}
	if err := c.sendQuery(query, &options); err != nil {
		release(c, err)
		return nil, err
	}
	var (
		// 进展展示之类
		onProcess  = options.onProcess()
		block, err = c.firstBlock(ctx, onProcess)
	)
	if err != nil {
		release(c, err)
		return nil, err
	}
	// resort batch to specified columns
	if err = block.SortColumns(queryColumns); err != nil {
		return nil, err
	}

	b := &batch{
		ctx:          ctx,
		query:        query,
		conn:         c,
		block:        block,
		released:     false,
		connRelease:  release,
		connAcquire:  acquire,
		onProcess:    onProcess,
		closeOnFlush: opts.CloseOnFlush,
	}

	if opts.ReleaseConnection {
		b.release(b.closeQuery())
	}

	return b, nil
}

```

放入数据
```go

func (b *batch) Append(v ...any) error {
	// 判断死否发送
	if b.sent {
		return ErrBatchAlreadySent
	}
	if b.err != nil {
		return b.err
	}

	if len(v) > 0 {
		// row 类型
		if r, ok := v[0].(*rows); ok {
			return b.appendRowsBlocks(r)
		}
	}

	if err := b.block.Append(v...); err != nil {
		b.err = errors.Wrap(ErrBatchInvalid, err.Error())
		b.release(err)
		return err
	}
	return nil
}

```

```go
func (b *Block) Append(v ...any) (err error) {
	columns := b.Columns
	// 判断values数量与 columns 是否相同
	if len(columns) != len(v) {
		return &BlockError{
			Op:  "Append",
			Err: fmt.Errorf("clickhouse: expected %d arguments, got %d", len(columns), len(v)),
		}
	}
	for i, v := range v {
		if err := b.Columns[i].AppendRow(v); err != nil {
			return &BlockError{
				Op:         "AppendRow",
				Err:        err,
				ColumnName: columns[i].Name(),
			}
		}
	}
	return nil
}
```

发送数据
```go
func (b *batch) Send() (err error) {
	stopCW := contextWatchdog(b.ctx, func() {
		// close TCP connection on context cancel. There is no other way simple way to interrupt underlying operations.
		// as verified in the test, this is safe to do and cleanups resources later on
		if b.conn != nil {
			_ = b.conn.conn.Close()
		}
	})

	defer func() {
		stopCW()
		b.sent = true
		b.release(err)
	}()
	if b.err != nil {
		return b.err
	}
	if b.sent || b.released {
		if err = b.resetConnection(); err != nil {
			return err
		}
	}
	// 数据发送
	if b.block.Rows() != 0 {
		if err = b.conn.sendData(b.block, ""); err != nil {
			// there might be an error caused by context cancellation
			// in this case we should return context error instead of net.OpError
			if ctxErr := b.ctx.Err(); ctxErr != nil {
				return ctxErr
			}

			return err
		}
	}
	if err = b.closeQuery(); err != nil {
		return err
	}
	return nil
}
```





## 第三方应用-->grafana clickhouse plugin
```go
// https://github.com/grafana/clickhouse-datasource/blob/28f86d02d120e38a11fff363fac846224580550b/pkg/plugin/driver.go
func (h *Clickhouse) Connect(ctx context.Context, config backend.DataSourceInstanceSettings, message json.RawMessage) (*sql.DB, error) {
	settings, err := LoadSettings(ctx, config)
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if settings.TlsAuthWithCACert || settings.TlsClientAuth {
		tlsConfig, err = getTLSConfig(settings)
		if err != nil {
			return nil, err
		}
	} else if settings.Secure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: settings.InsecureSkipVerify,
		}
	}

	t, err := strconv.Atoi(settings.DialTimeout)
	if err != nil {
		return nil, backend.DownstreamError(errors.New(fmt.Sprintf("invalid timeout: %s", settings.DialTimeout)))
	}
	qt, err := strconv.Atoi(settings.QueryTimeout)
	if err != nil {
		return nil, backend.DownstreamError(errors.New(fmt.Sprintf("invalid query timeout: %s", settings.QueryTimeout)))
	}

	protocol := clickhouse.Native
	if settings.Protocol == "http" {
		protocol = clickhouse.HTTP
	}

	compression := clickhouse.CompressionLZ4
	if protocol == clickhouse.HTTP {
		compression = clickhouse.CompressionGZIP
	}

	customSettings := make(clickhouse.Settings)
	if settings.CustomSettings != nil {
		for _, setting := range settings.CustomSettings {
			customSettings[setting.Setting] = setting.Value
		}
	}

	httpHeaders, err := extractForwardedHeadersFromMessage(message)
	if err != nil {
		return nil, err
	}

	// merge settings.HttpHeaders with message httpHeaders
	for k, v := range settings.HttpHeaders {
		httpHeaders[k] = v
	}

	opts := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", settings.Host, settings.Port)},
		Auth: clickhouse.Auth{
			Database: settings.DefaultDatabase,
			Password: settings.Password,
			Username: settings.Username,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: getClientInfoProducts(ctx),
		},
		Compression: &clickhouse.Compression{
			Method: compression,
		},
		DialTimeout: time.Duration(t) * time.Second,
		HttpHeaders: httpHeaders,
		HttpUrlPath: settings.Path,
		Protocol:    protocol,
		ReadTimeout: time.Duration(qt) * time.Second,
		Settings:    customSettings,
		TLS:         tlsConfig,
	}

	// dialCtx is used to create a connection to PDC, if it is enabled
	dialCtx, err := getPDCDialContext(settings)
	if err != nil {
		return nil, err
	}
	if dialCtx != nil {
		opts.DialContext = dialCtx
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(t)*time.Second)
	defer cancel()

	db := clickhouse.OpenDB(opts)

	// Set connection pool settings
	if i, err := strconv.Atoi(settings.ConnMaxLifetime); err == nil {
		db.SetConnMaxLifetime(time.Duration(i) * time.Minute)
	}
	if i, err := strconv.Atoi(settings.MaxIdleConns); err == nil {
		db.SetMaxIdleConns(i)
	}
	if i, err := strconv.Atoi(settings.MaxOpenConns); err == nil {
		db.SetMaxOpenConns(i)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("the operation was cancelled before starting: %w", ctx.Err())
	default:
		// proceed
	}

	// `sqlds` normally calls `db.PingContext()` to check if the connection is alive,
	// however, as ClickHouse returns its own non-standard `Exception` type, we need
	// to handle it here so that we can log the error code, message and stack trace
	if err := db.PingContext(ctx); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("the operation was cancelled during execution: %w", ctx.Err())
		}

		if exception, ok := err.(*clickhouse.Exception); ok {
			log.DefaultLogger.Error("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			log.DefaultLogger.Error(err.Error())
		}

		return nil, err
	}

	return db, settings.isValid()
}

// Co
```


## 参考
