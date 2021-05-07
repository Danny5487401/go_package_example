package main

// 流量控制：基于qps

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

func main() {
	// 初始化
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatalf("初始化unexpected err:%v", err)
	}

	// 定义限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "test",      //资源名
			TokenCalculateStrategy: flow.WarmUp, // 流量控制器的Token计算策略:冷启动
			ControlBehavior:        flow.Reject, // 控制器的控制策略
			Threshold:              1000,        // 流控阈值
			//StatIntervalInMs:       100,         // 流量控制器的独立统计结构的统计周，1000(也就是1秒)
			WarmUpPeriodSec:  30, // 预热的时间长度.30s
			WarmUpColdFactor: 3,  // 预热的因子，默认是3
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	// 做法：每一秒统计一次，通过了多少，总共多少,未通过多少,
	var globalTotal int
	var passTotal int
	var blockedTotal int
	ch := make(chan struct{}, 0)
	for i := 0; i < 100; i++ {
		go func() {
			for {
				globalTotal++
				e, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound)) //入口流量控制
				if b != nil {
					// 被blocked了
					blockedTotal++
					//但是每一秒产生特别多
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					passTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
					// 记住退出
					e.Exit()
				}
			}
		}()
	}
	go func() {
		var oldTotal int // 统计过去一秒总共多少个
		var oldPass int
		var oldBlock int

		for {
			oneSecondTotal := globalTotal - oldTotal
			oldTotal = globalTotal

			oneSecondPass := passTotal - oldPass
			oldPass = passTotal

			oneSecondBlock := blockedTotal - oldBlock
			oldBlock = blockedTotal

			time.Sleep(time.Second)
			//fmt.Printf("total:%d,pass:%d,blocked:%d", globalTotal, passTotal, blockedTotal)
			fmt.Printf("过去一秒total:%d,pass:%d,blocked:%d\n", oneSecondTotal, oneSecondPass, oneSecondBlock) //关心pass
		}
	}()
	<-ch

}
