package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// 调用目标方法
// 源对象中的某个字段没有出现在目标对象中，但是目标对象有一个同名的方法，方法接受一个同类型的参数，这时Copy会以源对象的这个字段作为参数调用目标对象的该方法：

type User struct {
	Name string
	Age  int
	Role string
}

type Employee struct {
	Name      string
	Age       int
	SuperRole string
}

func (e *Employee) Role(role string) {
	e.SuperRole = "Super" + role
}

func main() {
	user := User{Name: "Danny", Age: 18, Role: "Admin"}
	employee := Employee{}

	copier.Copy(&employee, &user)
	fmt.Printf("%#v\n", employee) // main.Employee{Name:"Danny", Age:18, SuperRole:"SuperAdmin"}
}
