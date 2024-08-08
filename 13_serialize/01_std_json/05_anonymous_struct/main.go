package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 使用匿名结构体添加字段
	anonymousStructDemo()

	// 使用匿名结构体组合多个结构体
	anonymousStructDemo2()
}

type UserInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func anonymousStructDemo() {
	u1 := UserInfo{
		ID:   123456,
		Name: "Danny",
	}
	// 使用匿名结构体内嵌User并添加额外字段Token
	b, err := json.Marshal(struct {
		*UserInfo
		Token string `json:"token"`
	}{
		&u1,
		"91je3a4s72d1da96h",
	})
	if err != nil {
		fmt.Printf("json.Marsha failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"id":123456,"name":"Danny","token":"91je3a4s72d1da96h"}
}

type Comment struct {
	Content string
}

type Image struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func anonymousStructDemo2() {
	c1 := Comment{
		Content: "永远不要高估自己",
	}
	i1 := Image{
		Title: "赞赏码",
		URL:   "https://www.liwenzhou.com/images/zanshang_qr.jpg",
	}
	// struct -> json string
	b, err := json.Marshal(struct {
		*Comment
		*Image
	}{&c1, &i1})
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
	// json string -> struct
	jsonStr := `{"Content":"永远不要高估自己","title":"赞赏码","url":"https://danny5487401.github.io/"}`
	var (
		c2 Comment
		i2 Image
	)
	if err := json.Unmarshal([]byte(jsonStr), &struct {
		*Comment
		*Image
	}{&c2, &i2}); err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("c2:%#v i2:%#v\n", c2, i2)
}
