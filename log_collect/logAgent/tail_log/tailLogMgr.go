package tailLog

import (
	"fmt"
	"go_test_project/log_collect/logAgent/etcd"
	"time"
)

// 管理所有tails-->tailObj
type TailMgr struct {
	logEntryList []*etcd.LogEntries // 当前的日志收集项配置信息
	taskMap map[string]*TailTask  //增删改  配置TailObj
	newConfChan chan []*etcd.LogEntries // 获取最新的配置变化
}

// 初始化总的taskMgr，管理所有tails-->tailObj
var taskMgr *TailMgr

func NewTailMgr(logConfList []*etcd.LogEntries) {
	// 初始化对象
	taskMgr = &TailMgr{
		logEntryList: logConfList,
		taskMap: make(map[string]*TailTask,32), // 32个配置项
		newConfChan: make(chan []*etcd.LogEntries), // 无缓冲区通道
	}
	for _, confValue := range taskMgr.logEntryList {
		// 一个配置项对应一个配置任务
		NewTailTask(confValue.Path,confValue.Topic)
	}
	go taskMgr.run()  //监听变化
}

// 监听自己的配置通道newConfChan，有新的就处理
// 1。 配置新增
// 2。 配置删除
// 3。 配置变更
func (t *TailMgr )run ()  {
	for{
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("配置发生变化",newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

// 向外暴露通道newConfChan
func ExposeChan() chan <- []*etcd.LogEntries {
	return taskMgr.newConfChan  // 方法一：给一个内部私有对象的私有字段
}
//
//func PushChan(newConf []*etcd.LogEntries)  {
//	taskMgr.newConfChan <- newConf // 方法二： 直接内部处理
//}
