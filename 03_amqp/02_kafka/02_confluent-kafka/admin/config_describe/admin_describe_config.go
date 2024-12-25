// List current configuration for a cluster resource
package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"time"
)

var (
	broker            = "tencent.danny.games"
	topicResourceType = "TOPIC"
	resourceName      = "user-active" // 这是topic名称
)

func main() {
	// 资源类型any, broker, topic, group\n
	trt, err := kafka.ResourceTypeFromString(topicResourceType)
	if err != nil {
		fmt.Printf("Invalid resource type: %s\n", topicResourceType)
		os.Exit(1)
	}
	// 资源名称 broker id or topic name
	resourceName := resourceName

	// 创建管理员客户端
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		fmt.Printf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}
	defer a.Close()

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dur, _ := time.ParseDuration("20s")
	// 获取当前集群的资源配置
	results, err := a.DescribeConfigs(ctx,
		[]kafka.ConfigResource{{Type: trt, Name: resourceName}},
		kafka.SetAdminRequestTimeout(dur))
	if err != nil {
		fmt.Printf("Failed to DescribeConfigs(%s, %s): %s\n",
			trt, resourceName, err)
		os.Exit(1)
	}

	// 打印结果
	for _, result := range results {
		fmt.Printf("类型%s 类型名称%s: 结果%s:\n", result.Type, result.Name, result.Error)
		for _, entry := range result.Config {
			// Truncate the value to 60 chars, if needed, for nicer formatting.
			fmt.Printf("%60s = %-60.60s   %-20s Read-only:%v Sensitive:%v\n",
				entry.Name, entry.Value, entry.Source,
				entry.IsReadOnly, entry.IsSensitive)
		}
	}

}
