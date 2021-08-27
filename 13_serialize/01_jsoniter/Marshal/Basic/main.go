package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// Json Marshal：将数据编码成json字符串
/*
Go 原生 encoding/json
	使用 json.Unmarshal 和 json.Marshal 函数，可以轻松将 JSON 格式的二进制数据反序列化到指定的 Go 结构体中，以及将 Go 结构体序列化为二进制流。
	而对于未知结构或不确定结构的数据，则支持将二进制反序列化到 map[string]interface{} 类型中，使用 KV 的模式进行数据的存取
json特性
	json 包解析的是一个 JSON 数据，而 JSON 数据既可以是对象（object），也可以是数组（array），同时也可以是字符串（string）、数值（number）、布尔值（boolean）以及空值（null）。
	var s string
	err := json.Unmarshal([]byte(`"Hello, world!"`), &s)
	// 注意字符串中的双引号不能缺，如果仅仅是 `Hello, world`，则这不是一个合法的 JSON 序列，会返回错误。
jsoniter
	从性能上，jsoniter 能够比众多大神联合开发的官方库性能还快的主要原因，一个是尽量减少不必要的内存复制，另一个是减少 reflect 的使用——同一类型的对象，jsoniter 只调用 reflect 解析一次之后即缓存下来。不过随着 go 版本的迭代，原生 json 库的性能也越来越高，jsonter 的性能优势也越来越窄
*/

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
	}

	//jsonStu是[]byte类型，转化成string类型便于查看
	fmt.Println(string(jsonStu))

}

type Stu struct {
	Name  string `json:"name"`
	Age   int
	HIgh  bool // 大小写混乱写
	sex   string
	Class *Class `json:"class"`
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
