package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// Json Marshal：将数据编码成json字符串

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	//实例化一个数据结构，用于生成json字符串
	stu := Stu{
		Name: "张三",
		Age:  18,
		HIgh: true,
		sex:  "男",
	}

	//指针变量
	cla := new(Class)
	cla.Name = "1班"
	cla.Grade = 3
	stu.Class = cla
	//Marshal失败时err!=nil
	jsonStu, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return
	}

	//jsonStu是[]byte类型，转化成string类型便于查看
	// {"name":"张三","Age":18,"HIgh":true,"class":{"Name":"1班","Grade":3}}
	fmt.Println(string(jsonStu))
	// 在实际项目中，编码成json串的数据结构，往往是切片类型。如下定义了一个[]StuRead类型的切片
	/*
		//方式1：只声明，不分配内存
		var stus1 []*Stu

		//方式2：分配初始值为0的内存
		stus2 := make([]*Stu, 0)
	*/

}

//interface{}类型其实是个空接口，即没有方法的接口。go的每一种类型都实现了该接口。因此，任何其他类型的数据都可以赋值给interface{}类型
type Stu struct {
	Name  interface{} `json:"name"`
	Age   interface{}
	HIgh  interface{}
	sex   interface{}
	Class interface{} `json:"class"`
}

type Class struct {
	Name  string
	Grade int
}
