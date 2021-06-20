package main


import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // "_"代码不直接使用包，但是底层链接要用，init函数会在main函数之前调用
// 以上等价于import _ "github.com/go-sql-driver/mysql"x
)
// 创建结构体映射表
type Student struct {
	Id int  //成为主键
	Name string
	Age int
}

// 创建连接池
var GlobalConn *gorm.DB


func main(){
	//gorm不能创建数据库，但可以创建数据库表
	var err error
	// 连接数据库
	GlobalConn,err = gorm.Open("mysql","root:chuanzhi@tcp(127.0.0.1:3306)/orm_test")
	if err != nil{
		fmt.Println(err)
		return
	}
	// 设置为true之后控制台会输出对应的SQL语句
	GlobalConn.LogMode(true)

	GlobalConn.DB().SetMaxIdleConns(10)
	GlobalConn.DB().SetMaxOpenConns(1000)

	GlobalConn.SingularTable(true) //创建的表叫Student -> student

	//创建表, 通过db对象的error判断是否成功
	//err = GlobalConn.AutoMigrate(new(Student)).Error  // 创建的表叫Student -> students
	//
	//if err != nil{
	//	fmt.Println("创建失败")
	//}

	//InsertData()

	//RetrieveData()

	UpdateData()
}


func InsertData(){
	stud := Student{
		Age: 35,
		Name: "Michael",
	}
	// 插入数据
	fmt.Println(GlobalConn.Create(&stud).Error)
}

func RetrieveData(){
	//var stu Student
	// 方法一 :First(&stu)  按主键排序
	//_ = GlobalConn.First(&stu)

	// 查询部分字段
	//GlobalConn.Select("age").First(&stu)
	//GlobalConn.Select([]string{"age","name"}).First(&stu)
	//fmt.Println(stu)

	// 方法二：Last(&stu)
	//_ = GlobalConn.Last(&stu)
	//fmt.Println(stu)


	// 方法三：Find(&stu) 多条数据
	var stus []Student
	// 注意顺序GlobalConn.Find(&stus).Where("name=?","danny")错误的
	//GlobalConn.Where("name=?","danny").Find(&stus)
	// 两个条件
	GlobalConn.Where("name=? and age =?","danny","12").Find(&stus)


	fmt.Println(stus)
}

func UpdateData()  {
	stud := Student{
		Id: 1,
		Age: 70,
		Name: "Harden",
	}
	// save 根据有没有指定主键，有：更新，无：插入
	fmt.Println(GlobalConn.Save(&stud).Error)

	// update 更新一个字段
	//err := GlobalConn.Model(&stud).Where("name=?","Michael").
	//	Update("age",50).Error
	//if err!= nil {
	//	fmt.Println(err)
	//}

	// updates 更新多个字段
	//err := GlobalConn.Model(&stud).Where("name=?","Michael").
	//	Updates(map[string]interface{}{"name":"Jame","age":9000}).Error
	//if err!= nil {
	//	fmt.Println(err)
	//}
}

func DeleteData()  {
	// 删除：物理删除，软删除：逻辑删除
	// 软删除：----数据是无价的
}