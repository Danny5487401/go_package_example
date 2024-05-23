/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// 服务器测试： nacos 2.0.3
// 客户端版本: Nacos-Go-Client:v1.0.7

package main

import (
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/Danny5487401/go_package_example/04_nacos/service_center/v1/server"
)

func main() {
	sc := []constant.ServerConfig{
		{
			IpAddr: "tencent.danny.games",
			Port:   8848,
		},
	}
	//更优雅的配置
	_ = []constant.ServerConfig{
		*constant.NewServerConfig("tencent.danny.games", 8848),
	}

	cc := constant.ClientConfig{
		NamespaceId:          "nacos_dev", //namespace id
		TimeoutMs:            5000,        //htto请求超时时间，单位ms毫秒
		NotLoadCacheAtStart:  true,        // 在启动式不读取本地缓存数据
		UpdateCacheWhenEmpty: true,        // 当服务列表为空时是否更新本地缓存
		LogDir:               "/tmp/nacos/log",
		CacheDir:             "/tmp/nacos/cache",
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "debug",
		BeatInterval:         3000, // 心跳间隔时间，单位毫秒(仅在serviceClient中有效),默认default value is 5000ms
		//ListenInterval: 3000,  //废弃(仅在configClient中有效)
		UpdateThreadNum: 2, //更新服务的线程数目,默认20
	}
	//更优雅的配置
	_ = *constant.NewClientConfig(
		constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)

	// a more graceful way to create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	// 1。 注册实例
	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//Register with cluster name
	//GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.11",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different cluster
	//GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.12",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	//Register different group
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.13",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.14",
		Port:        8848,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "cluster-b",
		GroupName:   "group-b",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	// 2。 注销实例
	//DeRegister with ip,port,serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	//Note:ip=10.0.0.10,port=8848 should belong to the cluster of DEFAULT and the group of DEFAULT_GROUP.
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.10",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Ephemeral:   true, //it must be true
	//})
	//
	////DeRegister with ip,port,serviceName,cluster
	////GroupName=DEFAULT_GROUP
	////Note:ip=10.0.0.10,port=8848,cluster=cluster-a should belong to the group of DEFAULT_GROUP.
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.11",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Cluster:     "cluster-a",
	//	Ephemeral:   true, //it must be true
	//})
	//
	////DeRegister with ip,port,serviceName,cluster,group
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.14",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Cluster:     "cluster-b",
	//	GroupName:   "group-b",
	//	Ephemeral:   true, //it must be true
	//})

	// 3。获取服务
	//Get service with serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
	})
	//Get service with serviceName and cluster
	//GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})
	//Get service with serviceName ,group
	//ClusterName=DEFAULT
	server.ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectAllInstance return all instances,include healthy=false,enable=false,weight<=0
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-a", "cluster-b"},
	})

	//SelectAllInstance
	//ClusterName=DEFAULT
	server.ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectInstances(client, vo.SelectInstancesParam{
		ServiceName: "demo.go",
	})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	// WRR(Weighted Round Robin)，加权轮训调度算法
	server.ExampleServiceClient_SelectOneHealthyInstance(client, vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
	})

	// 4. 服务监听
	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo.SubscribeParam{
		ServiceName: "demo.go",
		Clusters:    []string{"cluster-b"},
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("callback111 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	server.ExampleServiceClient_Subscribe(client, param)
	//param2 := &vo.SubscribeParam{
	//	ServiceName: "demo.go",
	//	Clusters:    []string{"cluster-b"},
	//	SubscribeCallback: func(services []model.SubscribeService, err error) {
	//		fmt.Printf("callback222 return services:%s \n\n", util.ToJsonString(services))
	//	},
	//}
	//server.ExampleServiceClient_Subscribe(client, param2)
	//server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
	//	Ip:          "10.0.0.112",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Weight:      10,
	//	ClusterName: "cluster-b",
	//	Enable:      true,
	//	Healthy:     true,
	//	Ephemeral:   true,
	//})
	//wait for client pull change from server
	// 监听保证10秒钟
	time.Sleep(10 * time.Second)
	//
	////Now we just unsubscribe callback1, and callback2 will still receive change event
	//server.ExampleServiceClient_UnSubscribe(client, param)
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.112",
	//	Ephemeral:   true,
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Cluster:     "cluster-b",
	//})
	////wait for client pull change from server
	//time.Sleep(10 * time.Second)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	server.ExampleServiceClient_GetAllService(client, vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})

	server.ExampleServiceClient_GetAllService(client, vo.GetAllServiceInfoParam{
		NameSpace: "nacos_dev",
		PageNo:    1,
		PageSize:  10,
	})
}
