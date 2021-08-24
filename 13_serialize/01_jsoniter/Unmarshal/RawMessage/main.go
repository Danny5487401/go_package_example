package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

type StuRead struct {
	Name  interface{}
	Age   interface{}
	HIgh  interface{}
	Class json.RawMessage `json:"class"` //注意这里
}

type Class struct {
	Name  string
	Grade int
}

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"
	str := []byte(data)
	stu := StuRead{}
	_ = json.Unmarshal(str, &stu)

	//注意这里：二次解析！
	cla := new(Class)
	json.Unmarshal(stu.Class, cla)

	fmt.Println("stu:", stu)
	fmt.Println("string(stu.Class):", string(stu.Class))
	fmt.Println("class:", cla)
	printType(&stu) //函数实现前面例子有

}

//利用反射，打印变量类型
func printType(stu *StuRead) {
	nameType := reflect.TypeOf(stu.Name)
	ageType := reflect.TypeOf(stu.Age)
	highType := reflect.TypeOf(stu.HIgh)
	classType := reflect.TypeOf(stu.Class)

	fmt.Println("nameType:", nameType)
	fmt.Println("ageType:", ageType)
	fmt.Println("highType:", highType)

	fmt.Println("classType:", classType)

}

/*
结论：
	接收体中，被声明为json.RawMessage类型的变量在json解析时，变量值仍保留json的原值，即未被自动解析为map[string]interface{}类型。如变量Class解析后的值为：{“naME”:”1班”,”GradE”:3}

	从打印的类型也可以看出，在第一次json解析时，变量Class的类型是json.RawMessage。此时，我们可以对该变量进行二次json解析，因为其值仍是个独立且可解析的完整json串。
	我们只需再定义一个新的接受体即可，如json.Unmarshal(stu.Class,cla)

*/
