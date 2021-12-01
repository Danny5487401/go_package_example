package main

import (
	"bufio"
	"fmt"
	"net"
)

var r *bufio.Reader
var w *bufio.Writer

func init() {
	conn, err := net.Dial("tcp", "ali.danny.games:6379")
	if err != nil {
		panic(err)
	}
	r = bufio.NewReader(conn)
	w = bufio.NewWriter(conn)
}

func write(format string, args ...interface{}) {
	if _, err := w.WriteString(fmt.Sprintf(format, args...)); err != nil {
		panic(err)
	}
}

func read() string {
	line, _, err := r.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(line)
}

func set(k, v string) string {
	args := []string{"SET", k, v}
	write("*%d\r\n", len(args))
	for _, arg := range args {
		write("$%d\r\n%s\r\n", len(arg), arg)
	}
	_ = w.Flush()
	return read()
}

// func get(k string) ... // other cmds

func main() {
	fmt.Println(set("keyX", "valueX")) // +OK
}

/*
如上实现明显存在的缺陷：

	连接管理：可用性低（连接断开未重连）、吞吐低（全局仅一条 TCP 连接）等
	逻辑混乱：模块间紧耦合（命令逻辑与协议细节未解耦）、可观测性低（未记录命令执行时长）等

为解决以上问题，go-redis 设计了连接池、将命令分类、把命令逻辑、协议实现与 IO 操作进行了解耦
*/
