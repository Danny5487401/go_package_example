package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

/*
弱类型输入
	我们并不想对结构体字段类型和map[string]interface{}的对应键值做强类型一致的校验。这时可以使用WeakDecode/WeakDecodeMetadata方法，它们会尝试做类型转换

*/

type Person struct {
	Name   string
	Age    int
	Emails []string
}

func main() {
	m := map[string]interface{}{
		"name":   123,            //类型不一致,string- >int
		"age":    "18",           //类型不一致  int->string
		"emails": []int{1, 2, 3}, // 类型不一致 []string->[]int
	}

	var p Person
	err := mapstructure.WeakDecode(m, &p)
	if err == nil {
		fmt.Printf("person:%+v", p)
	} else {
		fmt.Println("反序列化错误", err.Error())
	}
}
