package main

import "go_grpc_example/11_wire/02_wire/wire"

/*
安装工具go get github.com/google/wire/cmd/wire
同一目录 执行 wire
*/
func main() {
	// 简单初始化
	e := wire.InitializeEvent()

	e.Start()
}