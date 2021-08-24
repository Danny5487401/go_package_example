package main

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type Person struct {
	Name string
	Age  int
	Job  string `mapstructure:",omitempty"`
}

func main() {
	p := &Person{
		Name: "dj",
		Age:  18,
		//Job:  "gopher",
	}

	var m map[string]interface{}
	mapstructure.Decode(p, &m)

	data, _ := json.Marshal(m)
	fmt.Println(string(data))
}

// 写了job值，{"":"gopher","Age":18,"Name":"dj"}  不建议使用
// 没写job值，{"Age":18,"Name":"dj"}
