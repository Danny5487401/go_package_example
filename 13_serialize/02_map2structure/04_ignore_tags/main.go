package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

/*
未映射的值
	如果源数据中有未映射的值（即结构体中无对应的字段），mapstructure默认会忽略它。
	我们可以在结构体中定义一个字段，为其设置mapstructure:",remain"标签。这样未映射的值就会添加到这个字段中。注意，这个字段的类型只能为map[string]interface{}或map[interface{}]interface{}
*/

type Person struct {
	Name  string
	Age   int
	Job   string
	Other map[string]interface{} `mapstructure:",remain"`
}

func main() {
	data := `
    { 
      "name": "dj",
      "age":18,
      "job":"programmer",
      "height":"1.8m",
      "handsome": true
    }
  `

	var m map[string]interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatal(err)
	}

	var p Person
	mapstructure.Decode(m, &p)
	fmt.Println("other", p.Other)
}

// 注意：跟版本有关,这里是1.4.1
