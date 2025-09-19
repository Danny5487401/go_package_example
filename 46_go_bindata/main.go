// 安装 go install github.com/go-bindata/go-bindata/v3@v3.1.3
//
//go:generate  go-bindata -o=./public/config_gen.go -pkg=public test.env
package main

import (
	"bytes"
	"fmt"
	"github.com/Danny5487401/go_package_example/46_go_bindata/public"

	"github.com/spf13/viper"
)

func init() {
	fileObj, err := public.Asset("test.env")
	if err != nil {
		fmt.Printf("Asset file err:%v\n", err)
		return
	}
	viper.SetConfigType("env")
	err = viper.ReadConfig(bytes.NewBuffer(fileObj))
	if err != nil {
		fmt.Printf("Read Config err:%v\n", err)
		return
	}
}

func main() {
	fmt.Println("用户为:", viper.GetString("USER"))
	fmt.Println("密码为:", viper.GetString("PASS"))
}
