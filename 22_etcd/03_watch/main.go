package main

import (
	"context"
	"fmt"
	"github.com/Danny5487401/go_package_example/22_etcd/conn"
	"go.etcd.io/etcd/client/v3"
)

func main() {
	//  创建连接客户端
	var cli *clientv3.Client
	var err error
	cli, err = conn.GetClient()
	if err != nil {
		return
	}

	// watch 操作 ，获取key的变化,需要自己关闭退出
	watChan := cli.Watch(context.Background(), "etcd_key", clientv3.WithPrefix())

	for wresp := range watChan {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

}
