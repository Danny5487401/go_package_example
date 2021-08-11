package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

/*
背景：
	解析来自多种源头的数据流时，我们一般事先并不知道他们对应的具体类型。只有读取到一些字段之后才能做出判断。
做法
	我们可以先使用标准的encoding/json库将数据解码为map[string]interface{}类型，
	然后根据标识字段利用mapstructure库转为相应的 Go 结构体以便使用
*/

type Person struct {
	Name string
	Age  int
	Job  string
}

type Cat struct {
	Name  string
	Age   int
	Breed string
}

func main() {
	datas := []string{`
    { 
      "type": "person",
      "name":"dj",
      "age":18,
      "job": "programmer"
    }
  `,
		`
    {
      "type": "cat",
      "name": "kitty",
      "age": 1,
      "breed": "Ragdoll"
    }
  `,
	}
	for _, data := range datas {
		// 1。使用json反序列化成map[string]interface{}
		var m map[string]interface{}
		err := json.Unmarshal([]byte(data), &m)
		if err != nil {
			log.Fatal(err)
		}
		// 读取type字段
		switch m["type"].(string) {
		case "person":
			// 2。根据标识字段利用mapstructure库转为相应的 Go 结构体
			var p Person
			mapstructure.Decode(m, &p)
			fmt.Println("person", p)

		case "cat":
			var cat Cat
			mapstructure.Decode(m, &cat)
			fmt.Println("cat", cat)
		}
	}
}

/*
流程分析：
	先用json.Unmarshal将字节流解码为map[string]interface{}类型。然后读取里面的type字段。根据type字段的值，再使用mapstructure.Decode将该 JSON 串分别解码为Person和Cat类型的值，并输出。
同理：
	Google Protobuf 通常也使用这种方式。在协议中添加消息 ID 或全限定消息名。接收方收到数据后，先读取协议 ID 或全限定消息名。然后调用 Protobuf 的解码方法将其解码为对应的Message结构。
	从这个角度来看，mapstructure也可以用于网络消息解码，如果你不考虑性能的话😄
*/
