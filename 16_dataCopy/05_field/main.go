package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type User struct {
	Name         string
	Role         string
	Age          int32
	EmployeeCode int64 `copier:"EmployeeNum"` // specify field name

	// Explicitly ignored in the destination struct.
	Salary int
}

// 目标结构体中的标签提供了copy指令。复制忽略
// 或强制复制，如果字段没有被复制则惊慌或返回错误。
type Employee struct {
	//告诉copier。如果没有复制此字段，则复制到panic。
	Name string `copier:"must"`

	//告诉copier。 如果没有复制此字段，则返回错误。
	Age int32 `copier:"must,nopanic"`

	// 告诉copier。 显式忽略复制此字段。
	Salary int `copier:"-"`

	DoubleAge  int32
	EmployeeId int64 `copier:"EmployeeNum"` // 指定字段名
	SuperRole  string
}

func main() {
	var (
		user     = User{Name: "Danny", Age: 18, Role: "Admin", Salary: 200000, EmployeeCode: 10}
		employee = Employee{Salary: 150000}
	)
	copier.Copy(&employee, &user)
	fmt.Printf("%#v\n", employee) // main.Employee{Name:"Danny", Age:18, Salary:150000, DoubleAge:0, EmployeeId:10, SuperRole:""}
}
