package main

import (
	"log"
	"os"
)

func main() {

	// 准备日志文件
	logFile, _ := os.Create("02_log/01_log/demo.log")
	defer func() { _ = logFile.Close() }()

	// 自定义初始化日志对象 `logger`
	logger := log.New(logFile, "[Debug] - ", log.Lshortfile|log.Lmsgprefix)
	logger.Print("Print1")
	logger.Println("Println1")

	// 修改日志配置
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("[Info] - ")
	logger.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	logger.Print("Print2")
	logger.Println("Println2")
}
