package main

// 流量控制：基于qps

import (
	"fmt"
	"log"

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
			TokenCalculateStrategy: flow.Direct, // 流量控制器的Token计算策略
			ControlBehavior:        flow.Reject, // 控制器的控制策略
			Threshold:              10,          // 流控阈值
			StatIntervalInMs:       1000,        // 流量控制器的独立统计结构的统计周，1000(也就是1秒)
		},
		{
			Resource:               "test1",     //资源名
			TokenCalculateStrategy: flow.Direct, // 流量控制器的Token计算策略
			ControlBehavior:        flow.Reject, // 控制器的控制策略
			Threshold:              10,          // 流控阈值
			StatIntervalInMs:       1000,        // 流量控制器的独立统计结构的统计周，1000(也就是1秒)
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}

	// 构建一秒钟多个流量
	for i := 0; i < 12; i++ {
		// 流量入口代码
		e, b := sentinel.Entry("test", sentinel.WithTrafficType(base.Inbound)) //入口流量控制
		if b != nil {
			// 被blocked了
			fmt.Println("限流了", i)
		} else {
			fmt.Println("检查通过", i)
			// 记住退出
			e.Exit()
		}

	}

}
