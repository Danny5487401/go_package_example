package main


import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // "_"代码不直接使用包，但是底层链接要用，init函数会在main函数之前调用
	// 以上等价于import _ "github.com/go-sql-driver/mysql"x
)
// 创建结构体映射表
// 设置表的属性
type Employor struct {
	gorm.Model  // 匿名成员，继承
	// 修改表属性，只能在第一次建表的时候有效，或则增加新字段
	Name string `gorm:"size:100;default:'danny'"` //默认255
	Age int	`gorm:"not null"`
	// mysql 三种时间格式 date,datetime,timestamp
	// 如果必须使用mysql数据库特有的类型，用type指定
	JoinTime time.Time `gorm:"type:datetime"`
}


// 创建连接池
var GlobalConn *gorm.DB


func main(){
	//gorm不能创建数据库，但可以创建数据库表
	var err error
	// 连接数据库  // mysql 8小时时区问题，默认美国东八区
	GlobalConn,err = gorm.Open("mysql",
		"root:chuanzhi@tcp(127.0.0.1:3306)/orm_test?parseTime=True&loc=Local")
	if err != nil{
		fmt.Println(err)
		return
	}
	// 设置为true之后控制台会输出对应的SQL语句
	GlobalConn.LogMode(true)

	GlobalConn.DB().SetMaxIdleConns(10)
	GlobalConn.DB().SetMaxOpenConns(1000)

	GlobalConn.SingularTable(true)

	//创建表, 通过db对象的error判断是否成功
	err = GlobalConn.AutoMigrate(new(Employor)).Error

	if err != nil{
		fmt.Println("创建失败")
	}

	InsertData()

	//DeleteData()

	//RetrieveData()

}

func InsertData(){
	employer := Employor{
		Age: 100,
		Name: "Jordan",
	}
	// 插入数据
	fmt.Println(GlobalConn.Create(&employer).Error)
}


func DeleteData()  {
	// 删除：物理删除，软删除：逻辑删除
	// 软删除：----数据是无价的，使用自带的ORM.model继承字段，mysql自动维护
	// 多次删除，删除时间不变

	//fmt.Println(GlobalConn.Where("name=?","Curry").
	//	Delete(new(Employee)).Error)  // 逻辑删除

	fmt.Println(GlobalConn.Where("name=?","Jordan").Unscoped().
		Delete(new(Employor)).Error)  // 物理删除

}
func RetrieveData(){
	var employers []Employor

	// 查询软删除数据
	GlobalConn.Unscoped().Find(&employers)

	fmt.Println(employers)
}
