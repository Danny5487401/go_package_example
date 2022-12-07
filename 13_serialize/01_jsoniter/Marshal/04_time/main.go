package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

func main() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	info := TimeObject{
		TimeField:    timeImplementedMarshaler(time.Unix(123, 0)),
		OriginalTime: time.Now(), // 默认会把 time.Time 用字符串方式序列化
	}
	bytesInfo, err := json.Marshal(info)
	if err != nil {
		fmt.Println("失败", err.Error())
		return
	}
	fmt.Println(string(bytesInfo))
}

type timeImplementedMarshaler time.Time

func (obj timeImplementedMarshaler) MarshalJSON() ([]byte, error) {
	seconds := time.Time(obj).Unix()
	return []byte(strconv.FormatInt(seconds, 10)), nil
}

type TimeObject struct {
	TimeField    timeImplementedMarshaler
	OriginalTime time.Time
}
