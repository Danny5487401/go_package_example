<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [converter](#converter)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# converter 


ValueConverter接口定义了如何把一个普通的值转化成driver.Value的接口


```go
// go1.22.2/src/database/sql/driver/types.go
type ValueConverter interface {
    ConvertValue(v interface{}) (Value, error)
}

```

在开发的数据库驱动包里面实现这个接口的函数在很多地方会使用到，这个ValueConverter有很多好处：

- 转化driver.value到数据库表相应的字段，例如int64的数据如何转化成数据库表uint16字段
- 把数据库查询结果转化成driver.Value值
- 在scan函数里面如何把driver.Value值转化成用户定义的值


```go
// Value is a value that drivers must be able to handle.
// It is either nil, a type handled by a database driver's [NamedValueChecker]
// interface, or an instance of one of these types:
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
//
// If the driver supports cursors, a returned Value may also implement the [Rows] interface
// in this package. This is used, for example, when a user selects a cursor
// such as "select cursor(select * from my_table) from dual". If the [Rows]
// from the select is closed, the cursor [Rows] will also be closed.
type Value any
```


举例 bool
```go
// Bool is a [ValueConverter] that converts input values to bool.
//
// The conversion rules are:
//   - booleans are returned unchanged
//   - for integer types,
//     1 is true
//     0 is false,
//     other integers are an error
//   - for strings and []byte, same rules as [strconv.ParseBool]
//   - all other types are an error
var Bool boolType

type boolType struct{}

var _ ValueConverter = boolType{}

func (boolType) String() string { return "Bool" }

func (boolType) ConvertValue(src any) (Value, error) {
	switch s := src.(type) {
	case bool:
		return s, nil
	case string:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("sql/driver: couldn't convert %q into type bool", s)
		}
		return b, nil
	case []byte:
		b, err := strconv.ParseBool(string(s))
		if err != nil {
			return nil, fmt.Errorf("sql/driver: couldn't convert %q into type bool", s)
		}
		return b, nil
	}

	sv := reflect.ValueOf(src)
	switch sv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		iv := sv.Int()
		if iv == 1 || iv == 0 {
			return iv == 1, nil
		}
		return nil, fmt.Errorf("sql/driver: couldn't convert %d into type bool", iv)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uv := sv.Uint()
		if uv == 1 || uv == 0 {
			return uv == 1, nil
		}
		return nil, fmt.Errorf("sql/driver: couldn't convert %d into type bool", uv)
	}

	return nil, fmt.Errorf("sql/driver: couldn't convert %v (%T) into type bool", src, src)
}

```