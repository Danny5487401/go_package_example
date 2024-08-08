package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func main() {

	// json 包使用 RFC3339 标准中定义的时间格式
	// 内置的json包不识别我们常用的字符串时间格式，如2020-04-05 12:25:42
	timeFieldDemo()

	// 自定义的时间格式解析
	timeFieldCustomDemo()

}

type Post struct {
	CreateTime time.Time `json:"create_time"`
}

func timeFieldDemo() {
	p1 := Post{CreateTime: time.Now()}
	b, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("json.Marshal p1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"create_time":"2024-08-07T14:45:16.636296+08:00"}
	jsonStr := `{"create_time":"2020-04-05 12:25:42"}`
	var p2 Post
	if err := json.Unmarshal([]byte(jsonStr), &p2); err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("p2:%#v\n", p2)
}

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

var nilTime = (time.Time{}).UnixNano()

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ctLayout))), nil
}

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}

type Post2 struct {
	CreateTime CustomTime `json:"create_time"`
}

func timeFieldCustomDemo() {
	p1 := Post2{CreateTime: CustomTime{time.Now()}}
	b, err := json.Marshal(p1)
	if err != nil {
		fmt.Printf("json.Marshal p1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
	jsonStr := `{"create_time":"2020-04-05 12:25:42"}`
	var p2 Post2
	if err := json.Unmarshal([]byte(jsonStr), &p2); err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("p2:%#v\n", p2)
}
