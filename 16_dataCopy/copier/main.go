package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// User 1。源对象
type User struct {
	Name string
	Age  int
	Role string
}

// DoubleAge 目标对象中的一些字段，源对象中没有，但是源对象有同名的方法。这时Copy会调用这个方法，将返回值赋值给目标对象中的字段
func (u *User) DoubleAge() int {
	return u.Age * 2
}

// Employee 2。目标对象
type Employee struct {
	Name      string
	Age       int
	SuperRole string
}

// 源对象中的某个字段没有出现在目标对象中，但是目标对象有一个同名的方法,这时Copy会以源对象的这个字段作为参数调用目标对象的该方法
func (e *Employee) Role(role string) {
	e.SuperRole = "Super" + role
}

func main() {
	var (
		user  = User{Name: "dj", Age: 18}
		users = []User{
			{Name: "dj", Age: 18, Role: "Admin"},
			{Name: "dj2", Age: 18, Role: "Dev"},
		}
		employee  = Employee{}
		employees = []Employee{}
	)

	copier.Copy(&employee, &user)
	fmt.Printf("结构体复制%#v\n", employee)

	copier.Copy(&employees, &user)
	fmt.Printf("将结构赋值到切片%#v\n", employees)

	// employees = []Employee{}

	copier.Copy(&employees, &users)
	fmt.Printf("切片复制%#v\n", employees)

}
