package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func main() {
	// 原始方式
	output1, err := jsoniter.Marshal(struct {
		UserName      string `json:"user_name"`
		FirstLanguage string `json:"first_language"`
	}{
		UserName:      "taowen",
		FirstLanguage: "Chinese",
	})
	fmt.Println("原生方式", string(output1), err)

	// 优化方式
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
	output2, err := jsoniter.Marshal(struct {
		UserName      string
		FirstLanguage string
	}{
		UserName:      "taowen",
		FirstLanguage: "Chinese",
	})
	fmt.Println("优化方式", string(output2), err)

}
