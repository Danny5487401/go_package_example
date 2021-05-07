package main

/* Redigo特征
1. 支持所有redis命令
2. 支持管道，以及管道所有事务
3. 支持发布/订阅
4. 支持连接池
5. 支持 EVALSHA命令
6. 辅助函数，用于处理命令答复
*/

import (
	// "github.com/garyburd/redigo/redis" // 旧路径
	"fmt"
	"github.com/gomodule/redigo/redis" // 新路径
)

func main() {
	// 连接方式一
	//c, err := net.Dial("tcp", "106.14.35.115:6379")
	//if err != nil {
	//	panic(err)
	//}
	//read := time.Minute * 60
	//writer := time.Minute * 60
	//conn := redis.NewConn(c, read, writer)
	//// 关闭
	//defer conn.Close()

	// 连接方式二 简单，推荐
	conn, err := redis.Dial("tcp", "106.14.35.115:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	/*
		1.  Conn 是使用redis的主要接口，那么Do是对redis的主要操作方法
		2.  Do方法是发送对redis的操作命令并接收redis的答复
	*/
	// 设置值
	//reply, err := conn.Do("", "name", "danny")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(reply)

	// 设置过期
	reply, err := conn.Do("SET", "age", "18", "EX", "30")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}
	fmt.Println(reply)

	// 获取值
	//reply, err := conn.Do("GET", "name")
	//if err != nil {
	//	panic(err)
	//}
	//// 转换方式一
	////fmt.Printf("%s\n", reply)
	//// 转换方式二
	//replyStr, err := redis.String(reply, err)
	//fmt.Println(replyStr)

	// 判断是否存在

	isKeyExit, err := redis.Bool(conn.Do("EXISTS", "age"))
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("exists or not: %v \n", isKeyExit)

	}
}
