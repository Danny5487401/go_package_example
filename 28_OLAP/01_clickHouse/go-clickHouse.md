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

## 两种接口

- native interface

- std database/sql interface

协议 
- http 协议(实验性) ,只支持 `database/sql

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
