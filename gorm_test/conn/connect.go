package main


import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // "_"代码不直接使用包，但是底层链接要用，init函数会在main函数之前调用
	// 以上等价于import _ "github.com/go-sql-driver/mysql"
)
// 创建结构体映射表
type Student struct {
	Id int  //成为主键
	Name string
	Age int
}



func main(){
	//gorm不能创建数据库，但可以创建数据库表

	// 连接数据库 格式 ：用户名：密码@协议（ip:port)/数据库名  详细看源码go-sql-driver/mysql/dsn.go
	conn,err := gorm.Open("mysql","root:chuanzhi@tcp(127.0.0.1:3306)/orm_test")
	if err != nil{
		fmt.Println("gorm open err:",err)
		return
	}
	defer conn.Close()

	conn.SingularTable(true) //创建的表叫Student -> student

	// 创建表, 通过db对象的error判断是否成功
	err = conn.AutoMigrate(new(Student)).Error  // 创建的表叫Student -> students

	if err != nil{
		fmt.Println("创建失败")
	}
}