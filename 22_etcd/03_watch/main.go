package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go_package_example/22_etcd/conn"
)

func main() {
	//  创建连接客户端
	var cli *clientv3.Client
	var err error
	cli, err = conn.GetClient()
	if err != nil {
		return
	}

	// watch 操作 ，获取key的变化
	watChan := cli.Watch(context.Background(), "etcd_key")
	for wResp := range watChan {
		for _, ev := range wResp.Events {
			fmt.Printf("变化后的Type:%s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}
