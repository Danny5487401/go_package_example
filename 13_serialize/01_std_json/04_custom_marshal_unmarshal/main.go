package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	// 自定义的MarshalJSON方法 和 UnmarshalJSON
	customMethodDemo()

}

type Order struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"created_time"`
}

const layout = "2006-01-02 15:04:05"

// MarshalJSON 为Order类型实现自定义的MarshalJSON方法
func (o *Order) MarshalJSON() ([]byte, error) {
	type TempOrder Order // 定义与Order字段一致的新类型
	return json.Marshal(struct {
		CreatedTime string `json:"created_time"`
		*TempOrder         // 避免直接嵌套Order进入死循环
	}{
		CreatedTime: o.CreatedTime.Format(layout),
		TempOrder:   (*TempOrder)(o),
	})
}

// UnmarshalJSON 为Order类型实现自定义的UnmarshalJSON方法
func (o *Order) UnmarshalJSON(data []byte) error {
	type TempOrder Order // 定义与Order字段一致的新类型
	ot := struct {
		CreatedTime string `json:"created_time"` // 注意这是是 string 不是 time.Time 类型
		*TempOrder         // 避免直接嵌套Order进入死循环
	}{
		TempOrder: (*TempOrder)(o),
	}
	if err := json.Unmarshal(data, &ot); err != nil {
		return err
	}
	var err error
	o.CreatedTime, err = time.Parse(layout, ot.CreatedTime)
	if err != nil {
		return err
	}
	return nil
}

// 自定义序列化方法
func customMethodDemo() {
	o1 := Order{
		ID:          123456,
		Title:       "《Danny的Go学习笔记》",
		CreatedTime: time.Now(),
	}
	// 通过自定义的MarshalJSON方法实现struct -> json string
	b, err := json.Marshal(&o1)
	if err != nil {
		fmt.Printf("json.Marshal o1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
	// 通过自定义的UnmarshalJSON方法实现json string -> struct
	jsonStr := `{"created_time":"2020-04-05 10:18:20","id":123456,"title":"《Danny的Go学习笔记》"}`
	var o2 Order
	if err := json.Unmarshal([]byte(jsonStr), &o2); err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("o2:%#v\n", o2)
}
