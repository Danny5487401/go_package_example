package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type Class struct {
	Name  string
	Grade int
}
type StuRead struct {
	Name interface{} `json:"name"`
	Age  interface{}
	HIgh interface{}
	sex  interface{} //小写
	//普通struct类型
	//Class Class `json:"class"`
	//指针类型
	Class *Class `json:"class"`

	Test interface{}
}

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	//json字符中的"引号，需用\进行转义，否则编译出错
	//json字符串沿用上面的结果，但对key进行了大小的修改，并添加了sex数据
	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	str := []byte(data)

	//1.Unmarshal的第一个参数是json字符串，第二个参数是接受json解析的数据结构。
	//第二个参数必须是指针，否则无法接收解析的数据，如stu仍为空对象StuRead{}
	//2.可以直接stu:=new(StuRead),此时的stu自身就是指针
	stu := StuRead{}
	printType(&stu)
	err := json.Unmarshal(str, &stu)

	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("--------------json 解析后-----------")
	fmt.Println(&stu)
}

/*
结论
	1。json字符串解析时，需要一个“接收体”接受解析后的数据，且Unmarshal时接收体必须传递指针。否则解析虽不报错，但数据无法赋值到接受体中。如这里用的是StuRead{}接收。

	2。解析时，接收体可自行定义。json串中的key自动在接收体中寻找匹配的项进行赋值。匹配规则是：
		a. 先查找与key一样的json标签，找到则赋值给该标签对应的变量(如Name)。
		b. 没有json标签的，就从上往下依次查找变量名与key一样的变量，如Age。或者变量名忽略大小写后与key一样的变量。如HIgh，Class。第一个匹配的就赋值，后面就算有匹配的也忽略。
		(前提是该变量必需是可导出的，即首字母大写)。
	3.不可导出的变量无法被解析（如sex变量，虽然json串中有key为sex的k-v，解析后其值仍为nil,即空值）

	4.当接收体中存在json串中匹配不了的项时，解析会自动忽略该项，该项仍保留原值。如变量Test，保留空值nil。

	5.你一定会发现，变量Class貌似没有解析为我们期待样子。因为此时的Class是个interface{}类型的变量，而json串中key为CLASS的value是个复合结构，
	不是可以直接解析的简单类型数据（如“张三”，18，true等）。所以解析时，由于没有指定变量Class的具体类型，json自动将value为复合结构的数据解析为map[string]interface{}类型的项。
	也就是说，此时的struct Class对象与StuRead中的Class变量没有半毛钱关系，故与这次的json解析没有半毛钱关系。

interface{}类型变量在json解析前，打印出的类型都为nil，就是没有具体类型，这是空接口（interface{}类型）的特点。

json解析后，json串中value，只要是”简单数据”，都会按照默认的类型赋值，如”张三”被赋值成string类型到Name变量中，数字18对应float64，true对应bool类型。
	“简单数据”：是指不能再进行二次json解析的数据，如”name”:”张三”只能进行一次json解析。
	“复合数据”：类似”CLASS\”:{\”naME\”:\”1班\”,\”GradE\”:3}这样的数据，是可进行二次甚至多次json解析的，因为它的value也是个可被解析的独立json。即第一次解析key为CLASS的value，第二次解析value中的key为naME和GradE的value

对于”复合数据”，如果接收体中配的项被声明为interface{}类型，go都会默认解析成map[string]interface{}类型。如果我们想直接解析到struct Class对象中，可以将接受体对应的项定义为该struct类型
type StuRead struct {
...
	//普通struct类型
	Class Class `json:"class"`
	//指针类型
	Class *Class `json:"class"`
	}
打印
	Class类型：{张三 18 true <nil> {1班 3} <nil>}
	*Class类型：{张三 18 true <nil> 0xc42008a0c0 <nil>}
可以看出，传递Class类型的指针时，stu中的Class变量存的是指针，我们可通过该指针直接访问所属的数据，如stu.Class.Name/stu.Class.Grade

*/

//利用反射，打印变量类型
func printType(stu *StuRead) {
	nameType := reflect.TypeOf(stu.Name)
	ageType := reflect.TypeOf(stu.Age)
	highType := reflect.TypeOf(stu.HIgh)
	sexType := reflect.TypeOf(stu.sex)
	classType := reflect.TypeOf(stu.Class)
	testType := reflect.TypeOf(stu.Test)

	fmt.Println("nameType:", nameType)
	fmt.Println("ageType:", ageType)
	fmt.Println("highType:", highType)
	fmt.Println("sexType:", sexType)
	fmt.Println("classType:", classType)
	fmt.Println("testType:", testType)
}
