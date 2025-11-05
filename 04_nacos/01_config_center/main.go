package main

import (
	"encoding/json"
	"fmt"
	"github.com/Danny5487401/go_package_example/04_nacos/01_config_center/config"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func main() {
	// 1.服务端信息 ，至少一个
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1",
			Port:   8848,
		},
	}
	// 2.客户端信息
	clientConfig := constant.ClientConfig{
		NamespaceId:         "0f83848b-3ce7-41f6-bd45-f21145bbd44a", // we can retrieve multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		//LogDir:              "tmp/nacos/log",   // 当前项目目录
		//CacheDir:            "tmp/nacos/cache", // 当下次请求不到，可以从缓冲中获取
		LogLevel: "debug",
	}

	//  3.创建动态配置客户端
	// Create config client for dynamic configuration
	clientCfg, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}

	// 4.获取具体配置信息
	content, err := clientCfg.GetConfig(vo.ConfigParam{
		//DataId: "user-web-dev.yaml",
		DataId: "user-web-dev.json",
		Group:  "dev",
	})
	if err != nil {
		panic(err)
	}

	// 映射成struct
	srvCfg := config.ServerConfig{}
	//想要json字符串转换成struct，需要设置struct的tag
	err = json.Unmarshal([]byte(content), &srvCfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(srvCfg)

	// 监听配置文件变化
	if err = clientCfg.ListenConfig(vo.ConfigParam{
		//DataId: "user-web-dev.yaml",
		DataId: "user-web-dev.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	}); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 30)
}
