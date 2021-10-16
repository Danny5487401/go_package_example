package tailLog

// 专门从日志文件中收集日记

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"go_grpc_example/18_log_collect/logAgent/kafka"
	"time"
)

var (
	Tails *tail.Tail
	//LogChan chan string //string太占内存
)

// 一个日志收集的任务
type TailTask struct {
	path     string
	topic    string
	instance *tail.Tail
	// 用于取消运行的tailTask run()
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	_ = tailObj.init() //根据路径去打开对应的日志信息
	return
}
func (t *TailTask) init() (err error) { // 小写init内部调用
	config := tail.Config{
		ReOpen:    true,                                 //是否重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的什么地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      false,
	}
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Printf("tailFile [%v]failed err:%v", t.path, err)
		return
	}

	// 后期需要退出
	go t.ReadSendMsg() // 直接采集数据并发送到kafka
	return
}
func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}

func (t *TailTask) ReadSendMsg() {
	for {
		select {
		case <-t.ctx.Done():
			// 退出
			fmt.Printf("tailTask [%s_%s] exit\n", t.path, t.topic)
			return

		case line := <-t.ReadChan():
			// 发送到kafka
			//kafka.SendToKafka(t.topic,line.Text)  //函数调用函数：需要优化-->同步变异步，太多日志，不适合直接go协程
			// 做法：a.先把日志数据直接发送到一个通道中
			kafka.SendToChan(t.topic, line.Text)
			//b. 从通道中取数据发送到kafka:初始化通道就开始

		default:

			time.Sleep(time.Second)
		}

	}
}

// 以下用不到
func Init(fileName string) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 //是否重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的什么地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      false,
	}
	Tails, err = tail.TailFile(fileName, config)

	if err != nil {
		fmt.Printf("tailFile failed err:%v", err)
		return
	}
	return
}

//func ReadLog()  {
//	for{
//		line,ok := <- Tails.Lines
//		if !ok {
//			fmt.Printf("tails lines failed err:%v",ok)
//			time.Sleep(time.Second)
//			continue
//		}
//
//		fmt.Println("读取line:",line.Text)
//	}
//}

func ReadChan() <-chan *tail.Line {
	return Tails.Lines
}
