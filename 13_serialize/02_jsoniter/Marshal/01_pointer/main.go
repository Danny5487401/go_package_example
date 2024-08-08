package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// 使用标准的配置
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {

	// 实例化一个数据结构，用于生成json字符串
	stu := Stu{
		Name: "张三",
		Age:  18,
		HIgh: true,
		sex:  "男",
	}

	// 指针未初始化,打印有null数值
	stu1 := stu
	jsonStu1, err := json.Marshal(stu1)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return
	}
	fmt.Printf("%s\n", jsonStu1) // {"name":"张三","Age":18,"HIgh":true,"class":null}

	// 指针变量初始化
	cla := new(Class)
	cla.Name = "1班"
	cla.Grade = 3
	stu.Class = cla
	jsonStu, err := json.Marshal(stu)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return
	}

	//jsonStu是[]byte类型，转化成string类型便于查看
	// {"name":"张三","Age":18,"HIgh":true,"class":{"Name":"1班","Grade":3}}
	fmt.Println(string(jsonStu))

}

type Stu struct {
	Name  string `json:"name"`
	Age   int
	HIgh  bool // 大小写混乱写
	sex   string
	Class *Class `json:"class"` //包含指针
}

type Class struct {
	Name  string
	Grade int
}

/*
结论：
	1.只要是可导出成员（变量首字母大写），都可以转成json。因成员变量sex是不可导出的，故无法转成json。

	2.如果变量打上了json标签，如Name旁边的 `json:"name"` ，那么转化成的json key就用该标签“name”，否则取变量名作为key，如“Age”，“HIgh”。

	3.bool类型也是可以直接转换为json的value值。Channel， complex 以及函数不能被编码json字符串。当然，循环的数据结构也不行，它会导致marshal陷入死循环。

	4.指针变量，编码时自动转换为它所指向的值，如cla变量。
	（当然，不传指针，Stu struct的成员Class如果换成Class struct类型，效果也是一模一样的。只不过指针更快，且能节省内存空间。

*/
