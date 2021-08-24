package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

/*
type DecoderConfig struct {
	ErrorUnused       bool  // 为true时，如果输入中的键值没有与之对应的字段就返回错误；
	ZeroFields        bool  // 为true时，在Decode前清空目标map。为false时，则执行的是map的合并。用在struct到map的转换中
	WeaklyTypedInput  bool  // 实现WeakDecode/WeakDecodeMetadata的功能
	Metadata          *Metadata   // 不为nil时，收集Metadata数据
	Result            interface{}  // 为结果对象，在map到struct的转换中，Result为struct类型。在struct到map的转换中，Result为map类型
	TagName           string  // 默认使用mapstructure作为结构体的标签名，可以通过该字段设置
}
*/

type Person struct {
	Name string
	Age  int
}

func main() {
	m := map[string]interface{}{
		"name": 123, //类型不匹配
		"age":  "18",
		"job":  "programmer",
	}

	var p Person
	var metadata mapstructure.Metadata
	// 声明解析器
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &p,
		Metadata:         &metadata,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = decoder.Decode(m)
	if err == nil {
		fmt.Println("person:", p)
		fmt.Printf("keys:%#v, unused:%#v\n", metadata.Keys, metadata.Unused)
	} else {
		fmt.Println(err.Error())
	}
}
