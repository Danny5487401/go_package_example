package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// many more fields…
}

func main() {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	user := &User{
		Email:    "danny@qq.com",
		Password: "pass",
	}
	token := "token866"

	// 1. 临时添加额外的字段
	info, err := json.Marshal(struct {
		*User
		Token    string `json:"token"`
		Password bool   `json:"password,omitempty"` // 忽略掉空Password字段,可以用omitempty
	}{
		User:  user,
		Token: token,
	})
	if err != nil {
		fmt.Println("失败", err.Error())
		return
	}
	fmt.Println(string(info))
}
