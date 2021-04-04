package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// 注册服务
func Register(address string,port int, name string,tags []string,id string)error{
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:8500",address)
	fmt.Println(cfg)
	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err.Error())
	}
	// 生成检查对象
	check := &api.AgentServiceCheck{
		HTTP: fmt.Sprintf(`http://%s:%d/health`,address,port),
		Timeout: "5s",
		Interval:"5s",
		DeregisterCriticalServiceAfter: "10s",

	}
	// 生成注册对象
	registration := api.AgentServiceRegistration{
		Name: name,
		Tags: tags,
		ID: id,
		Check:check,
	}
	err = client.Agent().ServiceRegister(&registration)
	if err != nil{
		panic(err)
	}
	return nil
}
// 获取服务 ：服务发现
func AllService(address string)  {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:8500",address)
	fmt.Println(cfg)
	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err.Error())
	}
	data, err := client.Agent().Services()
	if err != nil{
		panic(err)
	}
	for key,_ := range data{
		fmt.Println(key)
	}
}

// 过滤服务
func FilterService(address,name string)  {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:8500",address)
	fmt.Println(*cfg)
	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err.Error())
	}
	filter := fmt.Sprintf(`Service =="%s"`,name)
	data, err :=  client.Agent().ServicesWithFilter(filter)
	if err != nil{
		panic(err)
	}
	for key,value := range data{
		fmt.Println(key,*value)
	}
}


func main(){
	//err := Register("192.168.16.111",8022,"user_web",[]string{
	//	"danny_shop","user_web",
	//},"user_web_id")
	//if err !=nil{
	//	fmt.Println(err.Error())
	//}
	//AllService("192.168.16.111")

	FilterService("192.168.16.111","user-srv")
}