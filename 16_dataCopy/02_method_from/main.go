package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

// 方法赋值
// 目标对象中的一些字段，源对象中没有，但是源对象有同名的方法。这时Copy会调用这个方法，将返回值赋值给目标对象中的字段：

type User struct {
	Name string
	Age  int
}

func (u *User) DoubleAge() int {
	return 2 * u.Age
}

// 目标对象
type Employee struct {
	Name      string
	DoubleAge int
	Role      string
}

func main() {
	user := User{Name: "Danny", Age: 18}
	employee := Employee{}

	copier.Copy(&employee, &user)
	fmt.Printf("%#v\n", employee) // main.Employee{Name:"Danny", DoubleAge:36, Role:""}
}
