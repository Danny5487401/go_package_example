<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Masterminds/squirrel](#mastermindssquirrel)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [基础builder](#%E5%9F%BA%E7%A1%80builder)
    - [插入](#%E6%8F%92%E5%85%A5)
    - [查询](#%E6%9F%A5%E8%AF%A2)
    - [执行](#%E6%89%A7%E8%A1%8C)
  - [第三方使用](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# Masterminds/squirrel
SQL生成库


## 源码分析



### 基础builder
```go
// github.com/!masterminds/squirrel@v1.5.4/statement.go

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType builder.Builder

// StatementBuilder 是其他builder 的父builder 
var StatementBuilder = StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(Question)

```
```go
type Builder struct {
    builderMap ps.Map 
}
var (
	EmptyBuilder      = Builder{ps.NewMap()}
	emptyBuilderValue = reflect.ValueOf(EmptyBuilder)
)
```



squirrel的四大结构体
```go
type SelectBuilder builder.Builder

func init() {
	builder.Register(SelectBuilder{}, selectData{})
}
```
- SelectBuilder
- UpdateBuilder
```go
type InsertBuilder builder.Builder

func init() {
	builder.Register(InsertBuilder{}, insertData{})
}

```
- InsertBuilder
- DeleteBuilder


注册
```go
func Register(builderProto, structProto interface{}) interface{} {
	empty := RegisterType(
		reflect.TypeOf(builderProto),
		reflect.TypeOf(structProto),
	).Interface()
	return empty
}
```
```go
var (
	// 仓库
	registry = make(map[reflect.Type]reflect.Type)
	registryMux sync.RWMutex
)


func RegisterType(builderType reflect.Type, structType reflect.Type) *reflect.Value {
	registryMux.Lock()
	defer registryMux.Unlock()
	structType.NumField() // Panic if structType is not a struct
	registry[builderType] = structType
	emptyValue := emptyBuilderValue.Convert(builderType)
	return &emptyValue
}
```



### 插入
数据结构
```go
type insertData struct {
	PlaceholderFormat PlaceholderFormat
	RunWith           BaseRunner
	Prefixes          []Sqlizer
	StatementKeyword  string
	Options           []string
	Into              string
	Columns           []string
	Values            [][]interface{}
	Suffixes          []Sqlizer
	Select            *SelectBuilder
}

```


```go
func Insert(into string) InsertBuilder {
	return StatementBuilder.Insert(into)
}

func (b StatementBuilderType) Insert(into string) InsertBuilder {
	// 转换成 InsertBuilder类型 
	return InsertBuilder(b).Into(into)
}

// 插入某个表
func (b InsertBuilder) Into(from string) InsertBuilder {
	return builder.Set(b, "Into", from).(InsertBuilder)
}

// 转成 sql 
func (b InsertBuilder) ToSql() (string, []interface{}, error) {
	data := builder.GetStruct(b).(insertData)
	return data.ToSql()
}
```

```go
// 设置值返回副本
func Set(builder interface{}, name string, v interface{}) interface{} {
	b := Builder{getBuilderMap(builder).Set(name, v)}
	return convert(b, builder)
}
```

设置值
```go
// github.com/lann/ps@v0.0.0-20150810152359-62de8c46ede0/map.go
func (self *tree) Set(key string, value Any) Map {
	hash := hashKey(key)
	return setLowLevel(self, hash, hash, key, value)
}
```


转换 sql 及 参数

```go
// 获取注册的 builder 
func GetStruct(builder interface{}) interface{} {
	// 初始化反射结构体
	structVal := newBuilderStruct(reflect.TypeOf(builder))
	if structVal == nil {
		return nil
	}
	return scanStruct(builder, structVal)
}

func newBuilderStruct(builderType reflect.Type) *reflect.Value {
	structType := getBuilderStructType(builderType)
	if structType == nil {
		return nil
	}
	newStruct := reflect.New(*structType).Elem()
	return &newStruct
}

```

```go
// 拼接 sql 
func (d *insertData) ToSql() (sqlStr string, args []interface{}, err error) {
    // 校验

	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSql(d.Prefixes, sql, " ", args)
		if err != nil {
			return
		}

		sql.WriteString(" ")
	}

	if d.StatementKeyword == "" {
		sql.WriteString("INSERT ")
	} else {
		sql.WriteString(d.StatementKeyword)
		sql.WriteString(" ")
	}

	if len(d.Options) > 0 {
		sql.WriteString(strings.Join(d.Options, " "))
		sql.WriteString(" ")
	}

	sql.WriteString("INTO ")
	sql.WriteString(d.Into)
	sql.WriteString(" ")

	if len(d.Columns) > 0 {
		sql.WriteString("(")
		sql.WriteString(strings.Join(d.Columns, ","))
		sql.WriteString(") ")
	}

	if d.Select != nil {
		args, err = d.appendSelectToSQL(sql, args)
	} else {
		args, err = d.appendValuesToSQL(sql, args)
	}
	if err != nil {
		return
	}

	if len(d.Suffixes) > 0 {
		sql.WriteString(" ")
		args, err = appendToSql(d.Suffixes, sql, " ", args)
		if err != nil {
			return
		}
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sql.String())
	return
}
```


### 查询
```go

func Select(columns ...string) SelectBuilder {
	return StatementBuilder.Select(columns...)
}
```




### 执行

```go
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Queryer is the interface that wraps the Query method.
//
// Query executes the given query as implemented by database/sql.Query.
type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// QueryRower is the interface that wraps the QueryRow method.
//
// QueryRow executes the given query as implemented by database/sql.QueryRow.
type QueryRower interface {
	QueryRow(query string, args ...interface{}) RowScanner
}

// BaseRunner groups the Execer and Queryer interfaces.
type BaseRunner interface {
	Execer
	Queryer
}

// Runner groups the Execer, Queryer, and QueryRower interfaces.
type Runner interface {
	Execer
	Queryer
	QueryRower
}

```

内置库实现

```go
func (db *DB) Exec(query string, args ...any) (Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

func (db *DB) Query(query string, args ...any) (*Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}
```


## 第三方使用 
```go
// 自动化种子下载工具  autobrr
// https://github.com/autobrr/autobrr/blob/develop/internal/database/database.go
const (
	DriverSQLite   = "sqlite"
	DriverPostgres = "postgres"
)

type DB struct {
	log     zerolog.Logger
	Handler *sql.DB
	lock    sync.RWMutex
	ctx     context.Context
	cfg     *domain.Config

	cancel func()

	Driver string
	DSN    string

	squirrel sq.StatementBuilderType 
}


func NewDB(cfg *domain.Config, log logger.Logger) (*DB, error) {
	db := &DB{
		// 制定占位符¥,  support both sqlite and postgres
		squirrel: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		log:      log.With().Str("module", "database").Str("type", cfg.DatabaseType).Logger(),
		cfg:      cfg,
	}
	// ...
}
```


## 参考

