package main

import(
	"encoding/json"
	"fmt"
	"time"

	"go_test_project/nacos_test/config"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main()  {
	// 服务端信息
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:"81.68.197.3",
			Port:8848,
		},
	}
	// 客户端信息
	clientConfig :=constant.ClientConfig{
		NamespaceId:         "84f5c407-5661-4306-abda-b51a9f02fba1", //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log", //当前项目目录
		CacheDir:            "tmp/nacos/cache",  // 当下次请求不到，可以从缓冲中获取
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	// Create config client for dynamic configuration
	clientCfg, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil{
		panic(err)
	}

	// 获取具体配置信息
	content ,err := clientCfg.GetConfig(vo.ConfigParam{
		//DataId: "user-web-dev.yaml",
		DataId: "user-web-dev.json",
		Group: "dev",

	})
	if err != nil{
		panic(err)
	}
	fmt.Println(content) // yaml格式的字符串
	// go语言本身直接json字符串反射成struct

	// 映射成struct
	srvCfg := config.ServerConfig{
	}
	//想要json字符串转换成struct，需要设置struct的tag
	err = json.Unmarshal([]byte(content),&srvCfg)
	if err != nil{
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
	});err!=nil{
		panic(err)
	}
	time.Sleep(time.Second*3000)
}
