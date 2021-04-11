package tailLog

import "go_test_project/log_collect/logAgent/etcd"

// 管理所有tails-->tailObj
type TailMgr struct {
	logEntryList []*etcd.LogEntries // 当前的日志收集项配置信息
	taskMap map[string]*TailTask  //增删改  配置TailObj
}

var taskMgr *TailMgr

func NewTailMgr(logConfList []*etcd.LogEntries) {
	taskMgr = &TailMgr{
		logEntryList: logConfList,
	}
	for _, confValue := range taskMgr.logEntryList {
		// 一个配置项对应一个配置任务
		NewTailTask(confValue.Path,confValue.Topic)
	}
}
