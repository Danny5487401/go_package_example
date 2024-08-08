package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// json字符串数组处理

// 接收普通消息结构体(参考)
type articles struct {
	Id    int    //文章id
	Title string //文章标题
}

func main() {
	// 初始化，完全兼容encoding/json
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	//json字符串数组,转换成切片
	articleStrings := `[{"Id":100,"Title":"木华黎"},{"Id":200,"Title":"耶律楚才"},{"Id":300,"Title":"纳呀啊","Test":100}]`
	// 只声明，不分配内存
	var articleSlice []map[string]interface{}
	err := json.Unmarshal([]byte(articleStrings), &articleSlice)
	if err != nil {
		fmt.Println("转换出错：", err)
		return
	}

	for k, v := range articleSlice {
		fmt.Println("第", k, "个数的值是:", v, v["Id"], v["Title"])
	}

	var articleSlice2 []articles
	err = json.Unmarshal([]byte(articleStrings), &articleSlice2)
	if err != nil {
		fmt.Println("转换出错：", err)
		return
	}

	for k, v := range articleSlice2 {
		fmt.Println("第", k, "个数的值是:", v, v.Id, v.Title)
	}

}
