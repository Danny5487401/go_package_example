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

package main

import (
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/util"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/Danny5487401/go_package_example/04_nacos/service_center/v2/server"
)

func main() {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("tencent.danny.games", 8848, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId("nacos_dev"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithRotateTime("1h"),
		constant.WithMaxAge(3),
		constant.WithLogLevel("debug"),
	)

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8890,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	//DeRegister with ip,port,serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	//Note:ip=10.0.0.10,port=8848 should belong to the cluster of DEFAULT and the group of DEFAULT_GROUP.
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.10",
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	GroupName:   "group-a",
	//	Cluster:     "cluster-a",
	//	Ephemeral:   true, //it must be true
	//})

	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})

	time.Sleep(1 * time.Second)

	//Get service with serviceName
	//ClusterName=DEFAULT, GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_GetService(client, vo.GetServiceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectAllInstance
	//GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectAllInstances(client, vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectInstances(client, vo.SelectInstancesParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//SelectOneHealthyInstance return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	server.ExampleServiceClient_SelectOneHealthyInstance(client, vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
	})

	//Subscribe key=serviceName+groupName+cluster
	//Note:We call add multiple SubscribeCallback with the same key.
	param := &vo.SubscribeParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",
		Clusters:    []string{"cluster-a"},
		SubscribeCallback: func(services []model.Instance, err error) {
			fmt.Printf("callback111 return services:%s \n\n", util.ToJsonString(services))
		},
	}
	server.ExampleServiceClient_Subscribe(client, param)

	server.ExampleServiceClient_RegisterServiceInstance(client, vo.RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		GroupName:   "group-a",
		ClusterName: "cluster-a",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "beijing"},
	})
	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	//Now we just unsubscribe callback1, and callback2 will still receive change event
	server.ExampleServiceClient_UnSubscribe(client, param)
	//server.ExampleServiceClient_DeRegisterServiceInstance(client, vo.DeregisterInstanceParam{
	//	Ip:          "10.0.0.112",
	//	Ephemeral:   true,
	//	Port:        8848,
	//	ServiceName: "demo.go",
	//	Cluster:     "cluster-b",
	//})
	//wait for client pull change from server
	time.Sleep(3 * time.Second)

	//GeAllService will get the list of service name
	//NameSpace default value is public.If the client set the namespaceId, NameSpace will use it.
	//GroupName default value is DEFAULT_GROUP
	server.ExampleServiceClient_GetAllService(client, vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 10,
	})
}
