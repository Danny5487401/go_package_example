package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

var client *api.Client

// Register 注册服务
func Register(address string, port int, name string, tags []string, id string) (err error) {

	// 生成检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf(`http://%s:%d/health`, address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	// 生成注册对象
	registration := api.AgentServiceRegistration{
		Name:  name,
		Tags:  tags,
		ID:    id,
		Check: check,
	}
	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		panic(err)
	}
	return nil
}

// AllService 获取服务 ：服务发现
func AllService() {
	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key := range data {
		fmt.Println(key)
	}
}

// FilterService 过滤服务
func FilterService(name string) {
	filter := fmt.Sprintf(`Service =="%s"`, name)
	data, err := client.Agent().ServicesWithFilter(filter)
	if err != nil {
		panic(err)
	}
	for key, value := range data {
		fmt.Println(key, *value)
	}
}

func initClient() (err error) {
	address := "tencent.danny.games"
	port := 8500
	// 创建默认配置，返回指针可修改
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", address, port)

	client, err = api.NewClient(cfg)
	if err != nil {
		fmt.Println("初始化客户端失败")
		return
	}
	return
}

func main() {
	var err error
	err = initClient()
	if err != nil {
		panic("启动失败")
	}

	// 注册服务
	err = Register("192.168.16.111", 8022, "user_web", []string{
		"danny_shop", "user_web",
	}, "user_web_id")
	if err != nil {
		fmt.Println(err.Error())
	}

	// 列举服务
	AllService()

	// 筛选服务
	FilterService("user-srv")
}
