package main

import (
	"fmt"
	"net"
)

// 获取可用端口
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func main() {
	port, err := GetFreePort()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(port)
}
