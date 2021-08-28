package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

/*
错误处理
	mapstructure执行转换的过程中不可避免地会产生错误，例如 JSON 中某个键的类型与对应 Go 结构体中的字段类型不一致。Decode/DecodeMetadata会返回这些错误
*/

type Person struct {
	Name   string
	Age    int
	Emails []string
}

func main() {
	m := map[string]interface{}{
		"name":   123,            //与定义的类型不一致
		"age":    "bad value",    //与定义的类型不一致
		"emails": []int{1, 2, 3}, //与定义的类型不一致
	}

	var p Person
	err := mapstructure.Decode(m, &p)
	if err != nil {
		fmt.Println(err.Error())
	}
}
