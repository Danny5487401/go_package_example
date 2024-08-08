package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func main() {
	extra.SupportPrivateFields()
	type TestObject struct {
		field1 string // 私有字段
	}

	obj := TestObject{}
	jsoniter.UnmarshalFromString(`{"field1":"Hello"}`, &obj)
	fmt.Println(obj)
}
