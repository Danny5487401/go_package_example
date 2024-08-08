package main

import (
	"fmt"
	"reflect"

	jsoniter "github.com/json-iterator/go"
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

	Null interface{} // 数据无
}

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	//json字符中的"引号，需用\进行转义，否则编译出错
	//json字符串沿用上面的结果，但对key进行了大小的修改，并添加了sex数据
	data := "{\"name\":\"张三\",\"Age\":18,\"high\":true,\"sex\":\"男\",\"CLASS\":{\"naME\":\"1班\",\"GradE\":3}}"

	// Unmarshal的第一个参数是json字符串，第二个参数是接受json解析的数据结构。
	// 第二个参数必须是指针，否则无法接收解析的数据，如stu仍为空对象StuRead{}
	stu := StuRead{}

	err := json.Unmarshal([]byte(data), &stu)
	if err != nil {
		fmt.Println("解析失败会报错，如json字符串格式不对，缺\"号，缺}等。", err)
		return
	}
	fmt.Println("--json 解析后数据-----")
	printType(&stu)

}

// 利用反射，打印变量类型
func printType(stu *StuRead) {
	nameType := reflect.TypeOf(stu.Name)
	ageType := reflect.TypeOf(stu.Age)
	highType := reflect.TypeOf(stu.HIgh)
	sexType := reflect.TypeOf(stu.sex)
	classType := reflect.TypeOf(stu.Class)
	nullType := reflect.TypeOf(stu.Null)

	fmt.Printf("nameType: %v,value: %v \n", nameType, stu.Name)
	fmt.Printf("ageType:: %v,value: %v \n", ageType, stu.Age)
	fmt.Printf("highType: %v,value: %v \n", highType, stu.HIgh)
	fmt.Printf("sexType: %v,value: %v \n", sexType, stu.sex)
	fmt.Printf("classType: %v,value: %v \n", classType, stu.Class)
	fmt.Printf("nullType:: %v,value: %v \n", nullType, stu.Null)
}
