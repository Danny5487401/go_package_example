package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

//接收普通消息结构体
type articles struct {
	Id    int    //文章id
	Title string //文章标题
}

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	//json字符串数组,转换成切片
	articleStrings := `[{"Id2":100,"Title":"木华黎"},{"Id":200,"Title":"耶律楚才"},{"Id":300,"Title":"纳呀啊","Test":100}]`
	// 只声明，不分配内存
	var articleSlide []map[string]interface{}
	multiErr := json.Unmarshal([]byte(articleStrings), &articleSlide)
	if multiErr != nil {
		fmt.Println("转换出错：", multiErr)
	}

	for k, v := range articleSlide {
		fmt.Println("第", k, "个数的值是:", v, v["Id"], v["Title"])
	}

}
