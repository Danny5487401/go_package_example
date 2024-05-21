package main

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
)

func main() {
	// 自定义函数
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return length, nil
		},
	}

	exprString := "strlen('teststring')"
	exprFunc, _ := govaluate.NewEvaluableExpressionWithFunctions(exprString, functions)
	result, _ := exprFunc.Evaluate(nil)
	fmt.Println(result)

	// 参数解析
	parameters := make(map[string]interface{})
	expr, _ := govaluate.NewEvaluableExpression("(mem_used / total_mem) * 100")
	parameters = make(map[string]interface{})
	parameters["total_mem"] = 1024
	parameters["mem_used"] = 512
	result, _ = expr.Evaluate(parameters)
	fmt.Println(result)

	// 结构体方法访问
	u := User{FirstName: "Xia", LastName: "Danny", Age: 18}
	parameters = make(map[string]interface{})
	parameters["u"] = u

	exprStructMethod, _ := govaluate.NewEvaluableExpression("u.Fullname()")
	result, _ = exprStructMethod.Evaluate(parameters)
	fmt.Println("user", result)

	exprStructField, _ := govaluate.NewEvaluableExpression("u.Age > 18")
	result, _ = exprStructField.Evaluate(parameters)
	fmt.Println("age > 18?", result)

	exprStructMethod2, _ := govaluate.NewEvaluableExpression("FullName")
	result, _ = exprStructMethod2.Eval(u)
	fmt.Println("user", result)
}

type User struct {
	FirstName string
	LastName  string
	Age       int
}

func (u User) Fullname() string {
	return u.FirstName + " " + u.LastName
}

func (u User) Get(name string) (interface{}, error) {
	if name == "FullName" {
		return u.FirstName + " " + u.LastName, nil
	}

	return nil, errors.New("unsupported field " + name)
}
