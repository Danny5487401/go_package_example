package main

import (
	"fmt"
	"go_test_project/log_collect/logAgent/conf"
	"go_test_project/log_collect/logAgent/kafka"
	tailLog "go_test_project/log_collect/logAgent/tail_log"
	"time"

	"gopkg.in/ini.v1"
)

// logAgent 入口

var (
	//cfg *init.File // 方法一: load
	appCfg = new(conf.AppConf) // 方法二：映射

)

func main()  {
	// 0. 加载配置文件
	//var err error
	//cfg,err = ini.Load("./log_collect/logAgent/conf/config.ini")  //路径要从项目目录working directory算
	//if err!=nil{
	//	fmt.Printf("[Load Init]failed:%v\n",err)
	//	return
	//}
	err := ini.MapTo(appCfg,"./log_collect/logAgent/conf/config.ini")
	if err!=nil{
		fmt.Printf("[Load Init]failed:%v\n",err)
		return
	}

	// 1. 初始化kafka连接
	//err = kafka.Init([]string{cfg.Section("kafka").Key("address").String()})  //这是大写的Init,不是导入的小写init
	err = kafka.Init([]string{appCfg.Address})  //这是大写的Init,不是导入的小写init
	if err!=nil{
		fmt.Printf("[InitKafka]failed:%v\n",err)
		return
	}
	fmt.Println("[InitKafka]success")

	// 2.打开日志文件准备收集日记
	//fileName := "/Users/python/Desktop/go_test_project/log_collect/logAgent/log.txt"
	//if err := tailLog.Init(fileName);err!=nil{
	if err := tailLog.Init(appCfg.FilePath);err!=nil{
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
			//kafka.SendToKafka(cfg.Section("kafka").Key("topic").String(),line.Text) // cfg得全局变量
			kafka.SendToKafka(appCfg.Topic,line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}