package main

import (
	"fmt"
	"github.com/Danny5487401/go_package_example/11_dependency_injection/01_wire/03_wire_return_err/wire"
	"os"
)

/*
安装工具go get github.com/google/wire/cmd/wire
同一目录 执行 wire
*/
// 正常调用
func main() {
	e, err := wire.InitializeEvent()
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}
	e.Start()
}
