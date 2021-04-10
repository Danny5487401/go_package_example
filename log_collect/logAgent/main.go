package main

import (
	"fmt"
	"go_test_project/log_collect/logAgent/kafka"
	tailLog "go_test_project/log_collect/logAgent/tail_log"
	"time"
)

// logAgent 入口

func main()  {
	// 1. 初始化kafka连接
	err := kafka.Init([]string{"81.68.197.3:9092"})  //这是大写的Init,不是导入的小写init
	if err!=nil{
		fmt.Printf("[InitKafka]failed:%v\n",err)
		return
	}
	fmt.Println("[InitKafka]success")
	// 2.打开日志文件准备收集日记
	fileName := "/Users/python/Desktop/go_test_project/log_collect/logAgent/log.txt"
	if err := tailLog.Init(fileName);err!=nil{
		fmt.Printf("[Init tailLog failed]:%v",err)
	}
	fmt.Println("[Init tailLog]success")
	run()

}

func run()  {
	// 1. 不停的读取日志
	for{
		select {
		case line := <- tailLog.ReadChan():
			// 2. 发送到kafka
			kafka.SendToKafka("danny-kafka-test",line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}