<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go-clickHouse源码分析](#go-clickhouse%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go-clickHouse 源码分析

## 初始化

注册插件
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
	if len(database) == 0 {
		database = DefaultDatabase
	}
	if len(username) == 0 {
		username = DefaultUsername
	}
	if v, err := strconv.ParseBool(query.Get("no_delay")); err == nil {
		noDelay = v
	}
	tlsConfig := getTLSConfigClone(tlsConfigName)
	if tlsConfigName != "" && tlsConfig == nil {
		return nil, fmt.Errorf("invalid tls_config - no config registered under name %s", tlsConfigName)
	}
	secure = tlsConfig != nil
	if v, err := strconv.ParseBool(query.Get("secure")); err == nil {
		secure = v
	}
	if v, err := strconv.ParseBool(query.Get("skip_verify")); err == nil {
		skipVerify = v
	}
	if duration, err := strconv.ParseFloat(query.Get("timeout"), 64); err == nil {
		connTimeout = time.Duration(duration * float64(time.Second))
	}
	if duration, err := strconv.ParseFloat(query.Get("read_timeout"), 64); err == nil {
		readTimeout = time.Duration(duration * float64(time.Second))
	}
	if duration, err := strconv.ParseFloat(query.Get("write_timeout"), 64); err == nil {
		writeTimeout = time.Duration(duration * float64(time.Second))
	}
	if size, err := strconv.ParseInt(query.Get("block_size"), 10, 64); err == nil {
		blockSize = int(size)
	}
	if altHosts := strings.Split(query.Get("alt_hosts"), ","); len(altHosts) != 0 {
		for _, host := range altHosts {
			if len(host) != 0 {
				hosts = append(hosts, host)
			}
		}
	}
	switch query.Get("connection_open_strategy") {
	case "random":
		connOpenStrategy = connOpenRandom
	case "in_order":
		connOpenStrategy = connOpenInOrder
	case "time_random":
		connOpenStrategy = connOpenTimeRandom
	}

	settings, err := makeQuerySettings(query)
	if err != nil {
		return nil, err
	}

	if v, err := strconv.ParseBool(query.Get("compress")); err == nil {
		compress = v
	}

	if v, err := strconv.ParseBool(query.Get("check_connection_liveness")); err == nil {
		checkConnLiveness = v
	}
	if secure {
		// There is no way to check the liveness of a secure connection, as long as there is no access to raw TCP net.Conn
		checkConnLiveness = false
	}

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