<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [xorm](#xorm)
  - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [设置最大空闲idle连接数和和最大连接数](#%E8%AE%BE%E7%BD%AE%E6%9C%80%E5%A4%A7%E7%A9%BA%E9%97%B2idle%E8%BF%9E%E6%8E%A5%E6%95%B0%E5%92%8C%E5%92%8C%E6%9C%80%E5%A4%A7%E8%BF%9E%E6%8E%A5%E6%95%B0)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# xorm 

## 初始化

集群结构体
```go
// EngineGroup defines an engine group
type EngineGroup struct {
	*Engine  // 主
	slaves []*Engine  //从
	policy GroupPolicy
}

// slave的选择策略
type GroupPolicy interface {
    Slave(*EngineGroup) *Engine
}
```
```go
// Engine is the major struct of xorm, it means a database manager.
// Commonly, an application only need one engine
type Engine struct {
	cacherMgr      *caches.Manager
	defaultContext context.Context
	dialect        dialects.Dialect
	engineGroup    *EngineGroup
	logger         log.ContextLogger
	tagParser      *tags.Parser
	db             *core.DB

	driverName     string
	dataSourceName string

	TZLocation *time.Location // The timezone of the application
	DatabaseTZ *time.Location // The timezone of the database

	logSessionID bool // create session id
}

```

1. 创建单个engine
```go
master, err := xorm.NewEngine("mysql", "root:chuanzhi@tcp(ali.danny.games:3307)/masterSlaveDB")
```
```go
func NewEngine(driverName string, dataSourceName string) (*Engine, error) {
	// 选择一种数据库的dialect，这里拿mysql
	dialect, err := dialects.OpenDialect(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// 打开数据库
	db, err := core.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return newEngine(driverName, dataSourceName, dialect, db)
}
```
```go
// /Users/xiaxin/go/pkg/mod/xorm.io/xorm@v1.1.0/dialects/driver.go

var (
    drivers = map[string]Driver{}
)

// QueryDriver query a driver with name
func QueryDriver(driverName string) Driver {
    return drivers[driverName]
}

// OpenDialect opens a dialect via driver name and connection string
func OpenDialect(driverName, connstr string) (Dialect, error) {
	// 1.选择解析会话的驱动
	driver := QueryDriver(driverName)
	if driver == nil {
		return nil, fmt.Errorf("unsupported driver name: %v", driverName)
	}
    // 2. 这里开始解析地址
	uri, err := driver.Parse(driverName, connstr)
	if err != nil {
		return nil, err
	}

	// 3. 选择DBType,这是是mysql
	dialect := QueryDialect(uri.DBType)
	if dialect == nil {
		return nil, fmt.Errorf("unsupported dialect type: %v", uri.DBType)
	}

	dialect.Init(uri)

	return dialect, nil
}
```

xorm默认注册的dialects
```go
// /Users/xiaxin/go/pkg/mod/xorm.io/xorm@v1.1.0/dialects/dialect.go
func init() {
	regDrvsNDialects()
}

func regDrvsNDialects() bool {
	providedDrvsNDialects := map[string]struct {
		dbType     schemas.DBType
		getDriver  func() Driver
		getDialect func() Dialect
	}{
		"mssql":    {"mssql", func() Driver { return &odbcDriver{} }, func() Dialect { return &mssql{} }},
		"odbc":     {"mssql", func() Driver { return &odbcDriver{} }, func() Dialect { return &mssql{} }}, // !nashtsai! TODO change this when supporting MS Access
		"mysql":    {"mysql", func() Driver { return &mysqlDriver{} }, func() Dialect { return &mysql{} }},
		"mymysql":  {"mysql", func() Driver { return &mymysqlDriver{} }, func() Dialect { return &mysql{} }},
		"postgres": {"postgres", func() Driver { return &pqDriver{} }, func() Dialect { return &postgres{} }},
		"pgx":      {"postgres", func() Driver { return &pqDriverPgx{} }, func() Dialect { return &postgres{} }},
		"sqlite3":  {"sqlite3", func() Driver { return &sqlite3Driver{} }, func() Dialect { return &sqlite3{} }},
		"sqlite":   {"sqlite3", func() Driver { return &sqlite3Driver{} }, func() Dialect { return &sqlite3{} }},
		"oci8":     {"oracle", func() Driver { return &oci8Driver{} }, func() Dialect { return &oracle{} }},
		"goracle":  {"oracle", func() Driver { return &goracleDriver{} }, func() Dialect { return &oracle{} }},
	}

	for driverName, v := range providedDrvsNDialects {
		if driver := QueryDriver(driverName); driver == nil {
			RegisterDriver(driverName, v.getDriver())
			RegisterDialect(v.dbType, v.getDialect)
		}
	}
	return true
}

```
Dialect接口
```go
// Dialect represents a kind of database
type Dialect interface {
	Init(*URI) error  //初始化
	URI() *URI
	SQLType(*schemas.Column) string
	FormatBytes(b []byte) string

	IsReserved(string) bool
	Quoter() schemas.Quoter
	SetQuotePolicy(quotePolicy QuotePolicy)

	AutoIncrStr() string

	GetIndexes(queryer core.Queryer, ctx context.Context, tableName string) (map[string]*schemas.Index, error)
	IndexCheckSQL(tableName, idxName string) (string, []interface{})
	CreateIndexSQL(tableName string, index *schemas.Index) string
	DropIndexSQL(tableName string, index *schemas.Index) string

	GetTables(queryer core.Queryer, ctx context.Context) ([]*schemas.Table, error)
	IsTableExist(queryer core.Queryer, ctx context.Context, tableName string) (bool, error)
	CreateTableSQL(table *schemas.Table, tableName string) ([]string, bool)
	DropTableSQL(tableName string) (string, bool)

	GetColumns(queryer core.Queryer, ctx context.Context, tableName string) ([]string, map[string]*schemas.Column, error)
	IsColumnExist(queryer core.Queryer, ctx context.Context, tableName string, colName string) (bool, error)
	AddColumnSQL(tableName string, col *schemas.Column) string
	ModifyColumnSQL(tableName string, col *schemas.Column) string

	ForUpdateSQL(query string) string

	Filters() []Filter
	SetParams(params map[string]string)
}
```
mysql的dialect
```go
type mysql struct {
	Base
	net               string
	addr              string
	params            map[string]string
	loc               *time.Location
	timeout           time.Duration
	tls               *tls.Config
	allowAllFiles     bool
	allowOldPasswords bool
	clientFoundRows   bool
	rowFormat         string
}

func (db *mysql) Init(uri *URI) error {
	db.quoter = mysqlQuoter
	return db.Base.Init(db, uri)
}
```

```go
// /Users/xiaxin/go/pkg/mod/xorm.io/xorm@v1.1.0/dialects/mysql.go
type mysqlDriver struct {
}

func (p *mysqlDriver) Parse(driverName, dataSourceName string) (*URI, error) {
	dsnPattern := regexp.MustCompile(
		`^(?:(?P<user>.*?)(?::(?P<passwd>.*))?@)?` + // [user[:password]@]
			`(?:(?P<net>[^\(]*)(?:\((?P<addr>[^\)]*)\))?)?` + // [net[(addr)]]
			`\/(?P<dbname>.*?)` + // /dbname
			`(?:\?(?P<params>[^\?]*))?$`) // [?param1=value1&paramN=valueN]
	matches := dsnPattern.FindStringSubmatch(dataSourceName)
	// tlsConfigRegister := make(map[string]*tls.Config)
	names := dsnPattern.SubexpNames()

	uri := &URI{DBType: schemas.MYSQL}

	for i, match := range matches {
		switch names[i] {
		case "dbname":
			uri.DBName = match
		case "params":
			if len(match) > 0 {
				kvs := strings.Split(match, "&")
				for _, kv := range kvs {
					splits := strings.Split(kv, "=")
					if len(splits) == 2 {
						switch splits[0] {
						case "charset":
							uri.Charset = splits[1]
						}
					}
				}
			}

		}
	}
	return uri, nil
}
```

```go
// /Users/xiaxin/go/pkg/mod/xorm.io/xorm@v1.1.0/core/db.go
// Open opens a database
func Open(driverName, dataSourceName string) (*DB, error) {
	// 使用go-sql-driver的插件
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{
		DB:           db,
		Mapper:       names.NewCacheMapper(&names.SnakeMapper{}),
		reflectCache: make(map[reflect.Type]*cacheStruct),
	}, nil
}
```


go-sql-driver插件
```go
// /Users/xiaxin/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/driver.go
// 引入默认注册mysql驱动
func init() {
	sql.Register("mysql", &MySQLDriver{})
}
```

Go标准包
```go
// /Users/xiaxin/go/go1.15.10/src/database/sql/sql.go
var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	// ...
	drivers[name] = driver
}
```

初始化Engine
```go
func newEngine(driverName, dataSourceName string, dialect dialects.Dialect, db *core.DB) (*Engine, error) {
	cacherMgr := caches.NewManager()
	mapper := names.NewCacheMapper(new(names.SnakeMapper))
	tagParser := tags.NewParser("xorm", dialect, mapper, mapper, cacherMgr)  // tag解析器

	engine := &Engine{
		dialect:        dialect,
		TZLocation:     time.Local,
		defaultContext: context.Background(),
		cacherMgr:      cacherMgr,
		tagParser:      tagParser,
		driverName:     driverName,
		dataSourceName: dataSourceName,
		db:             db,
		logSessionID:   false,
	}

	if dialect.URI().DBType == schemas.SQLITE {
		engine.DatabaseTZ = time.UTC
	} else {
		//除了sqlite，DatabaseTZ使用本地
		engine.DatabaseTZ = time.Local
	}

	// 设置日志
	logger := log.NewSimpleLogger(os.Stdout)
	logger.SetLevel(log.LOG_INFO)
	engine.SetLogger(log.NewLoggerAdapter(logger))

	// GC前的关闭连接
	runtime.SetFinalizer(engine, func(engine *Engine) {
		_ = engine.Close()
	})

	return engine, nil
}
```




2. 初始化集群
```go
func NewEngineGroup(args1 interface{}, args2 interface{}, policies ...GroupPolicy) (*EngineGroup, error) {
	var eg EngineGroup
	// 设置策略，默认轮训
	if len(policies) > 0 {
		eg.policy = policies[0]
	} else {
		eg.policy = RoundRobinPolicy()
	}

	// 初始化方式一
	driverName, ok1 := args1.(string)
	conns, ok2 := args2.([]string)
	if ok1 && ok2 {
		engines := make([]*Engine, len(conns))
		for i, conn := range conns {
			engine, err := NewEngine(driverName, conn)
			if err != nil {
				return nil, err
			}
			engine.engineGroup = &eg
			engines[i] = engine
		}

		eg.Engine = engines[0]
		eg.slaves = engines[1:]
		return &eg, nil
	}

    // 初始化方式二
    // 主
	master, ok3 := args1.(*Engine)
	// 从
	slaves, ok4 := args2.([]*Engine)
	if ok3 && ok4 {
		master.engineGroup = &eg
		for i := 0; i < len(slaves); i++ {
			slaves[i].engineGroup = &eg
		}
		eg.Engine = master
		eg.slaves = slaves
		return &eg, nil
	}
	return nil, ErrParamsType
}
```

## 设置最大空闲idle连接数和和最大连接数
```go
// SetMaxIdleConns 设置最大空间连接数目，默认是2
func (eg *EngineGroup) SetMaxIdleConns(conns int) {
	eg.Engine.DB().SetMaxIdleConns(conns)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].DB().SetMaxIdleConns(conns)
	}
}
```
```go
// /Users/xiaxin/go/go1.15.10/src/database/sql/sql.go
func (db *DB) SetMaxIdleConns(n int) {
	db.mu.Lock()
	if n > 0 {
		db.maxIdleCount = n
	} else {
		// No idle connections.
		db.maxIdleCount = -1
	}
	// Make sure maxIdle doesn't exceed maxOpen
	// 确保最大空闲不超过最大的打开数目
	if db.maxOpen > 0 && db.maxIdleConnsLocked() > db.maxOpen {
		db.maxIdleCount = db.maxOpen
	}
	// 需要关闭的连接
	var closing []*driverConn
	idleCount := len(db.freeConn)
	maxIdle := db.maxIdleConnsLocked()
	if idleCount > maxIdle {
		// 如果空闲连接大于最大空闲树木
		closing = db.freeConn[maxIdle:]
		db.freeConn = db.freeConn[:maxIdle]
	}
	db.maxIdleClosed += int64(len(closing))
	db.mu.Unlock()
	for _, c := range closing {
		c.Close()
	}
}


const defaultMaxIdleConns = 2

func (db *DB) maxIdleConnsLocked() int {
	n := db.maxIdleCount
	switch {
	case n == 0:
		// TODO(bradfitz): ask driver, if supported, for its default preference
		return defaultMaxIdleConns
	case n < 0:
		return 0
	default:
		return n
	}
}
```

底层的DB链接
```go
type DB struct {
	// Atomic access only. At top of struct to prevent mis-alignment
	// on 32-bit platforms. Of type time.Duration.
	waitDuration int64 // Total time waited for new connections.

	connector driver.Connector
	// numClosed is an atomic counter which represents a total number of
	// closed connections. Stmt.openStmt checks it before cleaning closed
	// connections in Stmt.css.
	numClosed uint64

	mu           sync.Mutex // protects following fields
	freeConn     []*driverConn
	connRequests map[uint64]chan connRequest
	nextRequest  uint64 // Next key to use in connRequests.
	numOpen      int    // number of opened and pending open connections
	// Used to signal the need for new connections
	// a goroutine running connectionOpener() reads on this chan and
	// maybeOpenNewConnections sends on the chan (one send per needed connection)
	// It is closed during db.Close(). The close tells the connectionOpener
	// goroutine to exit.
	openerCh          chan struct{}
	closed            bool
	dep               map[finalCloser]depSet
	lastPut           map[*driverConn]string // stacktrace of last conn's put; debug only
	maxIdleCount      int                    // zero means defaultMaxIdleConns; negative means 0
	maxOpen           int                    // 小于等于0，代表没有限制
	maxLifetime       time.Duration          // maximum amount of time a connection may be reused
	maxIdleTime       time.Duration          // maximum amount of time a connection may be idle before being closed
	cleanerCh         chan struct{}
	waitCount         int64 // Total number of connections waited for.
	maxIdleClosed     int64 // Total number of connections closed due to idle count.
	maxIdleTimeClosed int64 // Total number of connections closed due to idle time.
	maxLifetimeClosed int64 // Total number of connections closed due to max connection lifetime limit.

	stop func() // stop cancels the connection opener and the session resetter.
}
```

最大连接数目
```go
func (eg *EngineGroup) SetMaxOpenConns(conns int) {
	eg.Engine.DB().SetMaxOpenConns(conns)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].DB().SetMaxOpenConns(conns)
	}
}
```
```go
// /Users/xiaxin/go/go1.15.10/src/database/sql/sql.go
// If n <= 0, 打开连接数没有上限
// 默认没有上限
func (db *DB) SetMaxOpenConns(n int) {
	db.mu.Lock()
	db.maxOpen = n
	if n < 0 {
		db.maxOpen = 0
	}
	syncMaxIdle := db.maxOpen > 0 && db.maxIdleConnsLocked() > db.maxOpen
	db.mu.Unlock()
	if syncMaxIdle {
		db.SetMaxIdleConns(n)
	}
}
```

## 参考

- https://xorm.io/docs/chapter-01/readme/