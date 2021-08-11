package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
)

/*
type Person struct {
	Name string
}

// 方式一
type Friend struct {
	Person
}
//方式二 对于mapstructure，与方式一相同
type Friend struct {
	Person Person
}
*/

/*
为了正确解码
方式一：Person结构的数据要在person键下：
map[string]interface{} {
  "person": map[string]interface{}{"name": "dj"},
}
方式二：可以设置mapstructure:",squash"将该结构体的字段提到父结构中
type Friend struct {
  Person `mapstructure:",squash"`
}
这样只需要这样的 JSON 串，无效嵌套person键：
map[string]interface{}{
  "name": "dj",
}
*/
type Person struct {
	Name string
}

type Friend1 struct {
	Person
}

type Friend2 struct {
	Person `mapstructure:",squash"`
}

func main() {
	datas := []string{`
    { 
      "type": "friend1",
      "person": {
        "name":"dj"
      }
    }
  `,
		`
    {
      "type": "friend2",
      "name": "dj2"
    }
  `,
	}

	for _, data := range datas {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(data), &m)
		if err != nil {
			log.Fatal(err)
		}

		switch m["type"].(string) {
		case "friend1":
			var f1 Friend1
			mapstructure.Decode(m, &f1)
			fmt.Println("friend1", f1)

		case "friend2":
			var f2 Friend2
			mapstructure.Decode(m, &f2)
			fmt.Println("friend2", f2)
		}
	}
}

//注意：如果父结构体中有同名的字段，那么mapstructure会将JSON 中对应的值同时设置到这两个字段中，即这两个字段有相同的值
