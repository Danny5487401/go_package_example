package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type User struct {
	Name string
	Age  int
}

type Employee struct {
	Name string
	Age  int
	Role string
}

func main() {
	users := []User{
		{Name: "Danny1", Age: 18},
		{Name: "Danny2", Age: 28},
	}
	employees := []Employee{}

	copier.Copy(&employees, &users)
	fmt.Printf("%#v\n", employees) // []main.Employee{main.Employee{Name:"Danny1", Age:18, Role:""}, main.Employee{Name:"Danny2", Age:28, Role:""}}
}
