package tailLog
// 专门从日志文件中收集日记

import (
	"fmt"
	"github.com/hpcloud/tail"
)

var(
	Tails *tail.Tail
	//LogChan chan string //string太占内存
)



func Init(fileName string) (err error ){
	config := tail.Config{
		ReOpen:true,//是否重新打开
		Follow:true,//是否跟随
		Location:&tail.SeekInfo{Offset:0,Whence:2},//从文件的什么地方开始读
		MustExist:false,//文件不存在不报错
		Poll:false,
	}
	Tails, err = tail.TailFile(fileName,config)

	if err != nil{
		fmt.Printf("tailFile failed err:%v",err)
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

func ReadChan() <- chan *tail.Line {
	return Tails.Lines
}
