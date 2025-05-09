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
  - [性能提示](#%E6%80%A7%E8%83%BD%E6%8F%90%E7%A4%BA)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/ClickHouse/clickhouse-go

* clickhouse-go - 高级语言客户端，支持Go标准的数据库/sql接口或本地接口。
* ch-go - 低级客户端，仅支持本地接口。

在需要每秒数百万次插入的插入重负载使用案例中，我们建议使用低级客户端 ch-go。
该客户端避免了将数据从面向行的格式转换为列所需的附加开销，因为ClickHouse的本地格式要求。此外，它避免了任何反射或使用 interface{}（any）类型以简化使用。

对于专注于聚合或较低吞吐插入工作负载的查询工作负载，clickhouse-go提供了一种熟悉的database/sql接口以及更简单的行语义。


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

- native interface(TCP)

- std database/sql interface
```go
// github.com/!click!house/clickhouse-go/v2@v2.33.1/clickhouse_std.go
func init() {
	var debugf = func(format string, v ...any) {}
	sql.Register("clickhouse", &stdDriver{debugf: debugf})
}

```

协议 

```go
// github.com/!click!house/clickhouse-go/v2@v2.32.1/clickhouse_options.go
type Protocol int

const (
	Native Protocol = iota
	HTTP
)
```
- http 协议(实验性) ,只支持 `database/sql`

## v1 版本(不建议使用)
v1 驱动程序已被弃用，并将不再提供功能更新或对新ClickHouse类型的支持。

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

	// 默认实现的三种连接方式ConnOpenStrategy
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

```go

func dial(ctx context.Context, addr string, num int, opt *Options) (*connect, error) {
	var (
		err    error
		conn   net.Conn
		debugf = func(format string, v ...any) {}
	)

	switch {
	case opt.DialContext != nil:
		conn, err = opt.DialContext(ctx, addr)
	default:
		switch {
		case opt.TLS != nil:
			conn, err = tls.DialWithDialer(&net.Dialer{Timeout: opt.DialTimeout}, "tcp", addr, opt.TLS)
		default:
			// 默认使用 tcp 建立连接
			conn, err = net.DialTimeout("tcp", addr, opt.DialTimeout)
		}
	}

    // ...

	var (
		compression CompressionMethod
		compressor  *compress.Writer
	)
	if opt.Compression != nil {
		switch opt.Compression.Method {
		case CompressionLZ4, CompressionLZ4HC, CompressionZSTD, CompressionNone:
			compression = opt.Compression.Method
		default:
			return nil, fmt.Errorf("unsupported compression method for native protocol")
		}

		compressor = compress.NewWriter(compress.Level(opt.Compression.Level), compress.Method(opt.Compression.Method))
	} else {
		compression = CompressionNone
		compressor = compress.NewWriter(compress.LevelZero, compress.None)
	}

	var (
		connect = &connect{
			id:                   num,
			opt:                  opt,
			conn:                 conn,
			debugf:               debugf,
			buffer:               new(chproto.Buffer),
			reader:               chproto.NewReader(conn),
			revision:             ClientTCPProtocolVersion,
			structMap:            &structMap{},
			compression:          compression,
			connectedAt:          time.Now(),
			compressor:           compressor,
			readTimeout:          opt.ReadTimeout,
			blockBufferSize:      opt.BlockBufferSize,
			maxCompressionBuffer: opt.MaxCompressionBuffer,
		}
	)

	if err := connect.handshake(opt.Auth.Database, opt.Auth.Username, opt.Auth.Password); err != nil {
		return nil, err
	}

	if connect.revision >= proto.DBMS_MIN_PROTOCOL_VERSION_WITH_ADDENDUM {
		if err := connect.sendAddendum(); err != nil {
			return nil, err
		}
	}

	// warn only on the first connection in the pool
	if num == 1 && !resources.ClientMeta.IsSupportedClickHouseVersion(connect.server.Version) {
		debugf("[handshake] WARNING: version %v of ClickHouse is not supported by this client - client supports %v", connect.server.Version, resources.ClientMeta.SupportedVersions())
	}

	return connect, nil
}
```



### 批量写入

```go
// github.com/!click!house/clickhouse-go/v2@v2.32.1/clickhouse.go
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
		// 解析 reader 数据
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
实际调用解析数据

```go
func (c *connect) readData(ctx context.Context, packet byte, compressible bool) (*proto.Block, error) {
	if c.isClosed() {
		err := errors.New("attempted reading on closed connection")
		c.debugf("[read data] err: %v", err)
		return nil, err
	}

	if c.reader == nil {
		err := errors.New("attempted reading on nil reader")
		c.debugf("[read data] err: %v", err)
		return nil, err
	}

	if _, err := c.reader.Str(); err != nil {
		c.debugf("[read data] str error: %v", err)
		return nil, err
	}

	if compressible && c.compression != CompressionNone {
		c.reader.EnableCompression()
		defer c.reader.DisableCompression()
	}

	opts := queryOptions(ctx)
	location := c.server.Timezone
	if opts.userLocation != nil {
		location = opts.userLocation
	}

	// 解析数据
	block := proto.Block{Timezone: location}
	if err := block.Decode(c.reader, c.revision); err != nil {
		c.debugf("[read data] decode error: %v", err)
		return nil, err
	}

	block.Packet = packet
	c.debugf("[read data] compression=%q. block: columns=%d, rows=%d", c.compression, len(block.Columns), block.Rows())
	return &block, nil
}
```





放入数据方式一: Append(v ...any)
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


放入数据方式二: AppendStruct(v any) error 有字段验证功能,最终也是调用Append(values...)


```go
// github.com/!click!house/clickhouse-go/v2@v2.32.1/conn_batch.go

func (b *batch) AppendStruct(v any) error {
	if b.err != nil {
		return b.err
	}
	// 将 columnsName 字段去结构体中寻找数据
	values, err := b.conn.structMap.Map("AppendStruct", b.block.ColumnsNames(), v, false)
	if err != nil {
		return err
	}
	// 实际调用Append
	return b.Append(values...)
}
```


```go
func (m *structMap) Map(op string, columns []string, s any, ptr bool) ([]any, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return nil, &OpError{
			Op:  op,
			Err: fmt.Errorf("must pass a pointer, not a value, to %s destination", op),
		}
	}
	if v.IsNil() {
		return nil, &OpError{
			Op:  op,
			Err: fmt.Errorf("nil pointer passed to %s destination", op),
		}
	}
	t := reflect.TypeOf(s)
	if v = reflect.Indirect(v); t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, &OpError{
			Op:  op,
			Err: fmt.Errorf("%s expects a struct dest", op),
		}
	}

	var (
		index  map[string][]int
		values = make([]any, 0, len(columns))
	)

	switch idx, found := m.cache.Load(t); {
	case found:
		index = idx.(map[string][]int)
	default:
		index = structIdx(t)
		m.cache.Store(t, index)
	}
	for _, name := range columns {
		idx, found := index[name]
		if !found { // 字段验证
			return nil, &OpError{
				Op:  op,
				Err: fmt.Errorf("missing destination name %q in %T", name, s),
			}
		}
		switch field := v.FieldByIndex(idx); {
		case ptr:
			values = append(values, field.Addr().Interface())
		default:
			values = append(values, field.Interface())
		}
	}
	return values, nil
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

```go
func (c *connect) sendData(block *proto.Block, name string) error {
	if c.isClosed() {
		err := errors.New("attempted sending on closed connection")
		c.debugf("[send data] err: %v", err)
		return err
	}

	c.debugf("[send data] compression=%q", c.compression)
	c.buffer.PutByte(proto.ClientData)
	c.buffer.PutString(name)

	compressionOffset := len(c.buffer.Buf)

	// 头部信息
	if err := block.EncodeHeader(c.buffer, c.revision); err != nil {
		return err
	}

	for i := range block.Columns {
		if err := block.EncodeColumn(c.buffer, c.revision, i); err != nil {
			return err
		}
		if len(c.buffer.Buf) >= c.maxCompressionBuffer {
			if err := c.compressBuffer(compressionOffset); err != nil {
				return err
			}
			c.debugf("[buff compress] buffer size: %d", len(c.buffer.Buf))
			if err := c.flush(); err != nil {
				return err
			}
			compressionOffset = 0
		}
	}

	// 数据压缩
	if err := c.compressBuffer(compressionOffset); err != nil {
		return err
	}

	if err := c.flush(); err != nil {
        // ... 
	}

	defer func() {
		c.buffer.Reset()
	}()

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

```

## 性能提示
- 尽可能利用 ClickHouse API，特别是在处理基本类型时。这可以避免大量的反射和间接调用。
- 如果读取大数据集，考虑修改 BlockBufferSize。这将增加内存占用，但意味着在行迭代期间可以并行解码更多的块。默认值为 2 是保守的，并且最小化内存开销。更高的值将意味着更多的块驻留在内存中。这需要测试，因为不同的查询可能会产生不同的块大小。因此可以在查询级别通过上下文进行设置。
- 插入数据时要明确类型。虽然客户端旨在灵活，例如允许字符串解析为 UUID 或 IP，但这需要数据验证，并在插入时产生开销。
- 在可能的情况下，使用面向列的插入。这些应强类型化，避免客户端转换值的需要。


## 参考
- https://clickhouse.com/docs/zh/integrations/go