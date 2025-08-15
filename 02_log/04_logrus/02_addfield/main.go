package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	// 添加字段
	logrus.WithFields(logrus.Fields{
		"name": "Danny",
		"age":  18,
	}).Info("info msg") // time="2025-08-08T11:20:56+08:00" level=info msg="info msg" age=18 name=Danny

	// 在一个函数中的所有日志都需要添加某些字段
	requestLogger := logrus.WithFields(logrus.Fields{
		"user_id": 10010,
		"ip":      "192.168.32.15",
	})

	requestLogger.Info("info msg") // time="2025-08-08T11:24:46+08:00" level=info msg="info msg" ip=192.168.32.15 user_id=10010
	requestLogger.Error("error msg")

}
