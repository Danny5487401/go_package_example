package main

/*  GORM

1 实现原理
*/

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

// 1. 通过开发一个 QueryBuilder struct来将对像映射为sql语句

type QueryBuilder struct {
	Type reflect.Type // 通常特定类型的元数据可以通过reflect.Type来获取,所以QueryBuilder只有一个成员
}

// 创建查询语句，例如 SELECT ID, FirstName, LastName, Birthday FROM Employee
func (qb *QueryBuilder) CreateSelectQuery() string {
	buffer := bytes.NewBufferString("")
	//应该先判断qb.Type的类型是否为struct, 否则抛出异常
	for index := 0; index < qb.Type.NumField(); index++ {
		field := qb.Type.Field(index) //StructField 描述struct单个字段的信息

		if index == 0 {
			buffer.WriteString("SELECT ")
		} else {
			buffer.WriteString(", ")
		}
		column := field.Name
		// 处理大小的问题和别名的问题，
		if tag := field.Tag.Get("orm"); tag != "" {
			column = tag
		}
		buffer.WriteString(column)
	}

	if buffer.Len() > 0 {
		_, _ = fmt.Fprintf(buffer, " FROM %s", qb.Type.Name())
	}

	return buffer.String()
}

func Validate(obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for index := 0; index < v.NumField(); index++ {
		vField := v.Field(index)
		tField := t.Field(index)

		tag := tField.Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.Float64:
			value := vField.Float()
			if tag == "positive" && value < 0 {
				value = math.Abs(value)
				vField.SetFloat(value) // 修改值
			}
		case reflect.String:
			value := vField.String()
			if tag == "upper_case" {
				value = strings.ToUpper(value)
				vField.SetString(value)
			}
		default:
			return fmt.Errorf("unsupported kind '%s'", vField.Kind())
		}
	}

	return nil
}

// 字段值的读取与较验

type Employee struct {
	ID        uint32
	FirstName string
	LastName  string
	Birthday  time.Time
}

// 假设我们有一个Validator接口
type Validator interface {
	Validate() error
}

type PaymentTransaction struct {
	Amount      float64 `validate:"positive"`
	Description string  `validate:"max_length:250"`
}

func (p *PaymentTransaction) Validate() error {
	fmt.Println("Validating payment transaction")
	return nil
}

// 确定PaymentTransaction是否实现了接口，我们应该调用reflect.Type的Implements方法, 如何返回true
func CustomValidate(obj interface{}) error {
	v := reflect.ValueOf(obj)
	t := v.Type()

	interfaceT := reflect.TypeOf((*Validator)(nil)).Elem()
	if !t.Implements(interfaceT) {
		return fmt.Errorf("the Validator interface is not implemented")
	}

	validateFunc := v.MethodByName("Validate")
	validateFunc.Call(nil)
	return nil
}

// 为了获取Employee的元数据，我们需要先对它进行实例化。
func main() {
	//看源码得知调用Elem之前，应当先判断它的Kind()为Array, Chan, Map, Ptr, or Slice，否则Elem()抛出异常
	//reflect.TypeOf(&Employee{}).Kind() == Reflect.Ptr
	t := reflect.TypeOf(&Employee{}).Elem() // 获取指针的底层数据的类型
	v := reflect.ValueOf(&Employee{}).Elem()
	fmt.Printf("%+v\n", t)
	fmt.Printf("%+v\n", v)
	builder := &QueryBuilder{Type: t}
	fmt.Println(builder)
	// Reflect包中的Type和Value的区别在于，Type只是获取interface{}的类型，Value可以获取到它的值，
	// 同时可以通过v.Type返回的Type类型，所以TypeOf只是用于在获取元数据时使用，而Value即可以获取它的元数据的值和类型

	b := (*Validator)(nil)
	fmt.Printf("%+v\n", reflect.ValueOf(b).Type()) // *main.Validator
	fmt.Printf("%+v\n", reflect.ValueOf(b).Elem()) // <invalid reflect.Value>
	fmt.Printf("%+v\n", reflect.ValueOf(b))        // <nil>
	fmt.Printf("%+v\n", reflect.TypeOf(b).Elem())  // main.Validator

}

// 以上了解后
/*
1. producer.go 程序 ：type DB struct 主结构，用来管理当前的数据库连接, 它会包含下面将要列出的结构，比如Callback, Dialect
2. scope.go 程序: ype Scope struct,包含每一次sql操作的信息，比如db.Create(&user), 创建一个用户时，scope包含了*DB信息,SQL信息， 字段，主键字段等
3. callback.go 程序： type Callback struct 包含所以CRUD的回调，比如beforeSave, afterSave, beforeUpdate等等，
	在实现我们一些自定义功能时（插件)， 就需要了解这个struct
4. dialect.go 程序： type Dialect interface {} 这是一个接口类型，用来实现不同数据库的相同方法，消除不同数据库要写不一样的代码，
	比如HasColumn(tableName string, columnName string) bool方法，在mysql的sql语句，可能跟postgres的sql语句不同，所以分别需要在dialect_mysql和dialect_postgres中实现


*/
