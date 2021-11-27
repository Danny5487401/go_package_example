package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // "_"代码不直接使用包，但是底层链接要用，init函数会在main函数之前调用
	// 以上等价于import _ "github.com/go-sql-driver/mysql"x
)

// 创建结构体映射表
type Employee struct {
	gorm.Model // 匿名成员，继承

	Name string
	Age  int
}

// 创建连接池
var GlobalConn *gorm.DB

func main() {
	//gorm不能创建数据库，但可以创建数据库表
	var err error
	// 连接数据库  // mysql 8小时时区问题，默认美国东八区
	GlobalConn, err = gorm.Open("mysql",
		"root:chuanzhi@tcp(127.0.0.1:3306)/orm_test?parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置为true之后控制台会输出对应的SQL语句
	GlobalConn.LogMode(true)

	GlobalConn.DB().SetMaxIdleConns(10)
	GlobalConn.DB().SetMaxOpenConns(1000)

	GlobalConn.SingularTable(true) //创建的表叫Student -> students->student

	////创建表, 通过db对象的error判断是否成功
	//err = GlobalConn.AutoMigrate(new(Employee)).Error // 创建的表叫Student -> students
	//
	//if err != nil {
	//	fmt.Println("创建失败")
	//}

	//InsertData()

	//DeleteData()

	RetrieveData()

}

func InsertData() {
	employee := Employee{
		Age:  219,
		Name: "June",
	}
	// 插入数据
	fmt.Println(GlobalConn.Create(&employee).Error)
}

func DeleteData() {
	// 删除：物理删除，软删除：逻辑删除
	// 软删除：----数据是无价的，使用自带的ORM.model继承字段，mysql自动维护
	// 多次删除，删除时间不变

	//fmt.Println(GlobalConn.Where("name=?","Curry").
	//	Delete(new(Employee)).Error)  // 逻辑删除

	fmt.Println(GlobalConn.Where("name=?", "Jordan").Unscoped().
		Delete(new(Employee)).Error) // 物理删除

}
func RetrieveData() {
	var employees []Employee

	// 获取所有的记录
	// 查询软删除数据
	GlobalConn.Unscoped().Find(&employees)
	fmt.Printf("%+v\n", employees)

	// 获取第一条记录，按主键排序
	var employee = Employee{}
	//GlobalConn.First(&employee)
	//fmt.Printf("%+v\n", employee)
	// 通过主键进行查询 (仅适用于主键是数字类型)
	GlobalConn.First(&employee, 6)
	fmt.Printf("%+v\n", employee)

	//// 获取一条记录，不指定排序
	//GlobalConn.Take(&employee)
	//fmt.Printf("%+v\n", employee)

	// 获取最后一条记录，不指定排序
	//GlobalConn.Unscoped().Last(&employee)
	//fmt.Printf("%+v\n", employee)

}
