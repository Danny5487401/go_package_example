package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {

	// 处理字符串格式的数字
	intAndStringDemo()

	// json字符串中的数字经过Go语言中的json包反序列化之后都会成为float64类型
	jsonDemo()

}

func jsonDemo() {
	// map[string]interface{} -> json string
	var m = make(map[string]interface{}, 1)
	m["count"] = 1 // int
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
	}
	fmt.Printf("str:%#v\n", string(b)) // str:"{\"count\":1}"
	// json string -> map[string]interface{}
	var m2 map[string]interface{}
	err = json.Unmarshal(b, &m2)
	if err != nil {
		fmt.Printf("unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("value:%v\n", m2["count"]) // 1
	fmt.Printf("type:%T\n", m2["count"])  // float64

	var m3 map[string]interface{}
	// 使用decoder方式反序列化，指定使用number类型
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	err = decoder.Decode(&m3)
	if err != nil {
		fmt.Printf("unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("value:%v\n", m3["count"]) // 1
	fmt.Printf("type:%T\n", m3["count"])  // json.Number
	// 将m2["count"]转为json.Number之后调用Int64()方法获得int64类型的值
	count, err := m3["count"].(json.Number).Int64()
	if err != nil {
		fmt.Printf("parse to int64 failed, err:%v\n", err)
		return
	}
	fmt.Printf("type:%T\n", int(count)) // int
}

type Card struct {
	ID    int64   `json:"id,string"`    // 添加string tag
	Score float64 `json:"score,string"` // 添加string tag
}

func intAndStringDemo() {
	jsonStr1 := `{"id": "1234567","score": "88.50"}`
	var c1 Card
	if err := json.Unmarshal([]byte(jsonStr1), &c1); err != nil {
		fmt.Printf("json.Unmarsha jsonStr1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("c1:%#v\n", c1) // c1:main.Card{ID:1234567, Score:88.5}
}
