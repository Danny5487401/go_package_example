package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	// 忽略空值字段
	omitemptyDemo()

	// 忽略嵌套结构体空值字段
	nestedStructDemo()

	// 不修改原结构体忽略空值字段
	omitPasswordDemo()

}

type User struct {
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Hobby []string `json:"hobby"`
}

type User2 struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Hobby []string `json:"hobby,omitempty"`
}

type User3 struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Hobby []string `json:"hobby,omitempty"`
	Profile
}

// 嵌套的json串，需要改为具名嵌套或定义字段tag
type User4 struct {
	Name    string   `json:"name"`
	Email   string   `json:"email,omitempty"`
	Hobby   []string `json:"hobby,omitempty"`
	Profile `json:"profile"`
}

type User5 struct {
	Name    string   `json:"name"`
	Email   string   `json:"email,omitempty"`
	Hobby   []string `json:"hobby,omitempty"`
	Profile `json:"profile,omitempty"`
}

// 使用嵌套的结构体指针

type User6 struct {
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Hobby    []string `json:"hobby,omitempty"`
	*Profile `json:"profile,omitempty"`
}

type Profile struct {
	Website string `json:"site"`
	Slogan  string `json:"slogan"`
}

// 可以使用创建另外一个结构体PublicUser匿名嵌套原User，同时指定Password字段为匿名结构体指针类型，并添加omitempty tag
type User7 struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type PublicUser struct {
	*User7             // 匿名嵌套
	Password *struct{} `json:"password,omitempty"`
}

func omitemptyDemo() {
	u1 := User{
		Name: "Danny",
	}
	// struct -> json string
	b, err := json.Marshal(u1)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"name":"Danny","email":"","hobby":null}

	u2 := User2{
		Name: "Danny",
	}
	// struct -> json string
	b2, err := json.Marshal(u2)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b2) // str:{"name":"Danny"}

}

func nestedStructDemo() {
	u1 := User3{
		Name:  "Danny",
		Hobby: []string{"足球", "双色球"},
	}
	b, err := json.Marshal(u1)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"name":"Danny","hobby":["足球","双色球"],"site":"","slogan":""}

	// 改为具名嵌套或定义字段tag 不够
	u2 := User4{
		Name:  "Danny",
		Hobby: []string{"足球", "双色球"},
	}
	b2, err := json.Marshal(u2)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b2) // str:{"name":"Danny","hobby":["足球","双色球"],"site":"","slogan":""}

	// 仅添加omitempty是不够的
	u3 := User5{
		Name:  "Danny",
		Hobby: []string{"足球", "双色球"},
	}
	b3, err := json.Marshal(u3)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b3) // str:{"name":"Danny","hobby":["足球","双色球"],"site":"","slogan":""}

	// 使用嵌套的结构体指针

	u4 := User6{
		Name:  "Danny",
		Hobby: []string{"足球", "双色球"},
	}
	b4, err := json.Marshal(u4)
	if err != nil {
		fmt.Printf("json.Marshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b4) // str:{"name":"Danny","hobby":["足球","双色球"]}

}

func omitPasswordDemo() {
	u1 := User7{
		Name:     "Danny",
		Password: "123456",
	}
	b, err := json.Marshal(PublicUser{User7: &u1})
	if err != nil {
		fmt.Printf("json.Marshal u1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b) // str:{"name":"七米"}
}
