package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

//tail的基础使用
func main(){
	fileName := "/Users/python/Desktop/go_test_project/log_collect/logAgent/read_log/log.txt"
	Tails, err := tail.TailFile(fileName,tail.Config{
		ReOpen:true,//是否重新打开
		Follow:true,//是否跟随
		Location:&tail.SeekInfo{Offset:0,Whence:2},//从文件的什么地方开始读
		MustExist:false,//文件不存在不报错
		Poll:false,
	})

	if err != nil{
		fmt.Printf("tailFile failed err:%v",err)
		return
	}

	for{
		line,ok := <- Tails.Lines
		if !ok {
			fmt.Printf("tails lines failed err:%v",ok)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("line:",line.Text)
	}
}
