package tailLog

import (
	"fmt"
	"go_grpc_example/log_collect/logAgent/etcd"
	"time"
)

// 管理所有tails-->tailObj
type TailMgr struct {
	logEntryList []*etcd.LogEntries      // 当前的日志收集项配置信息
	taskMap      map[string]*TailTask    //增删改  配置TailObj
	newConfChan  chan []*etcd.LogEntries // 获取最新的配置变化
}

// 初始化总的taskMgr，管理所有tails-->tailObj
var taskMgr *TailMgr

func NewTailMgr(logConfList []*etcd.LogEntries) {
	// 初始化对象
	taskMgr = &TailMgr{
		logEntryList: logConfList,
		taskMap:      make(map[string]*TailTask, 32), // 32个配置项,记载配置项
		newConfChan:  make(chan []*etcd.LogEntries),  // 无缓冲区通道
	}
	for _, confValue := range taskMgr.logEntryList {
		// 一个配置项对应一个配置任务

		// 记录初始化起了多少个tailTask，方便新配置来了对比
		taskObj := NewTailTask(confValue.Path, confValue.Topic)
		key := fmt.Sprintf("%s_%s", confValue.Path, confValue.Topic)
		taskMgr.taskMap[key] = taskObj
	}
	go taskMgr.run() //监听变化
}

// 监听自己的配置通道newConfChan，有新的就处理
// 1。 配置新增
// 2。 配置删除
// 3。 配置变更
func (t *TailMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			fmt.Println("配置发生变化", newConf)
			// 1。 新增判断
			for _, value := range newConf {
				// path路径变化 或者 topic变更，所以要合成一个key
				key := fmt.Sprintf("%s_%s", value.Path, value.Topic)
				_, ok := t.taskMap[key]
				if ok {
					// 原来就有
					continue
				} else {
					// 说明新增或则修改--当作新增处理
					taskObj := NewTailTask(value.Path, value.Topic)
					t.taskMap[key] = taskObj
				}

			}
			// 原来配置文件有，新的配置文件没有，需要删除后台运行的任务
			for _, c1 := range t.logEntryList {
				isDelete := true
				for _, c2 := range newConf {
					if c2.Path == c1.Path && c2.Topic == c1.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete {
					// 把c1对应的tailObj停掉
					key := fmt.Sprintf("%s_%s", c1.Path, c1.Topic)
					t.taskMap[key].cancelFunc()
				}
			}

		default:
			time.Sleep(time.Second)
		}
	}
}

// 向外暴露通道newConfChan
func ExposeChan() chan<- []*etcd.LogEntries {
	return taskMgr.newConfChan // 方法一：给一个内部私有对象的私有字段
}

//
//func PushChan(newConf []*etcd.LogEntries)  {
//	taskMgr.newConfChan <- newConf // 方法二： 直接内部处理
//}
