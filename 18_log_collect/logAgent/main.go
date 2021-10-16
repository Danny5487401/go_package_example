package main

import (
	"fmt"
	"go_grpc_example/18_log_collect/logAgent/conf"
	"go_grpc_example/18_log_collect/logAgent/etcd"
	"go_grpc_example/18_log_collect/logAgent/kafka"
	tailLog "go_grpc_example/18_log_collect/logAgent/tail_log"
	"go_grpc_example/18_log_collect/logAgent/utils"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

// logAgent 入口

var (
	//cfg *init.File // 方法一: load
	appCfg = new(conf.AppConf) // 方法二：映射

)

func main() {
	// 0. 加载配置文件
	//var err error
	//cfg,err = ini.Load("./log_collect/logAgent/conf/config.ini")  //路径要从项目目录working directory算
	//if err!=nil{
	//	fmt.Printf("[Load Init]failed:%v\n",err)
	//	return
	//}
	err := ini.MapTo(appCfg, "./log_collect/logAgent/conf/config.ini")
	if err != nil {
		fmt.Printf("[Load Init]failed:%v\n", err)
		return
	}

	// 1. 初始化kafka连接
	//err = kafka.Init([]string{cfg.Section("kafka").Key("address").String()})  //这是大写的Init,不是导入的小写init
	//err = kafka.Init([]string{appCfg.KafkaConf.Address})  //这是大写的Init,不是导入的小写init
	err = kafka.Init([]string{appCfg.KafkaConf.Address}, appCfg.KafkaConf.ChanSize) //带缓冲区大小
	if err != nil {
		fmt.Printf("[InitKafka]failed:%v\n", err)
		return
	}

	fmt.Println("[InitKafka]success.")

	// 1.1 初始化etcd
	err = etcd.Init(appCfg.EtcdConf.Address, time.Duration(appCfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("[InitEtcd]failed:%v\n", err)
		return
	}
	fmt.Println("[InitEtcd]success.")

	// 1.2 从etcd中中获取配置项信息，并监视key变化
	// 为了实现每个logAgent都获取自己的配置，绑定自己的Ip区别
	ip, _ := utils.GetOutboundIP()
	etcdKey := fmt.Sprintf(appCfg.EtcdConf.Key, ip)
	logEntries, err := etcd.GetConf(etcdKey)
	if err != nil {
		fmt.Printf("get [Etcd]failed:%v\n", err)
		return
	}
	fmt.Printf("get [Etcd]success,%v\n", logEntries)

	for index, value := range logEntries {
		fmt.Printf("第一次Etcd获取结果index:%v,value:%v\n", index, value)
	}

	// 1.3 收集多个配置文件发送kafka
	// 1.3.1 循环每个日记收集项 ，创建多个Tails  ---> tailObj
	// 做法；一个管理者，管理所有的Tails---> tailObj,后期需要增加减少配置项tailObj
	// 未优化前
	//for _,confValue := range logEntries{  // 太多逻辑在main函数：需优化，放在tailLogMgr中
	//	//config := tail.Config{
	//	//	ReOpen:true,//是否重新打开
	//	//	Follow:true,//是否跟随
	//	//	Location:&tail.SeekInfo{Offset:0,Whence:2},//从文件的什么地方开始读
	//	//	MustExist:false,//文件不存在不报错
	//	//	Poll:false,
	//	//}
	//	//tailObj, err := tail.TailFile(confValue.Path,config)
	//	//if err != nil{
	//	//	fmt.Printf("tailFile [%v]failed err:%v",confValue.Path,err)
	//	//	return
	//	//}
	//
	//	/*
	//	// 1.3.2 发往kafka
	//	for  {
	//		select {
	//		case line := <- tailObj.Lines:  // 从tailObj的通道中一行一行获取数据
	//			kafka.SendToKafka(confValue.Topic,line.Text)
	//		}
	//	}
	//	 */  //这样写会死循环在第一个
	//
	//	tailLog.NewTailTask(confValue.Path,confValue.Topic)
	//}
	// 优化后：收集日志发送到kafka
	tailLog.NewTailMgr(logEntries) // 进行初始化newConfChan

	// 派一个哨兵去监视logEntries变化
	var wg sync.WaitGroup
	wg.Add(1)
	newConfChan := tailLog.ExposeChan()           // 从tailLog 获取对外暴露的通道
	go etcd.WatchConfChange(etcdKey, newConfChan) //哨兵发现变化会通知到通道中
	wg.Wait()                                     // 一直监听着

	// 2.打开日志文件准备收集日记
	//fileName := "/Users/python/Desktop/go_grpc_example/log_collect/logAgent/log.txt"
	//if err := tailLog.Init(fileName);err!=nil{
	//if err := tailLog.Init(appCfg.FilePath);err!=nil{
	//	fmt.Printf("[Init tailLog failed]:%v",err)
	//	return
	//}
	//fmt.Println("[Init tailLog]success")
	//run()
}

func run() {
	// 1. 不停的读取日志
	for {
		select {
		case line := <-tailLog.ReadChan():
			// 2. 发送到kafka
			//kafka.SendToKafka(cfg.Section("kafka").Key("topic").String(),line.Text) // cfg得全局变量
			kafka.SendToKafka(appCfg.Topic, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}
}
