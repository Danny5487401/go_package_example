package main

import (
	"context"
	"fmt"
	"go_package_example/22_etcd/conn"

	// 注意是这，不是github.com/coreos/etcd v3.3.17+incompatible
	"go.etcd.io/etcd/client/v3"
	"time"
)

// 与grpc版本冲突
func main() {
	//  创建连接客户端
	var cli *clientv3.Client
	var err error
	cli, err = conn.GetClient()
	if err != nil {
		return
	}

	// put 操作
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//rsp, err := cli.Put(ctx, "etcd_danny_key", "etcd_value")
	//cancel()
	//if err != nil {
	//	fmt.Printf("[put etcd  ] failed:%v\n", err)
	//}
	//fmt.Printf("[put etcd  ] success,rsp :%v\n", rsp)

	// get 操作
	ctx2, cancel := context.WithTimeout(context.Background(), time.Second*3) // 注意超时时间
	//resp,err := cli.Get(ctx2,"etcd_danny_key", clientv3.WithPrefix() )  //  clientv3.WithPrefix() 业务前缀
	resp, err := cli.Get(ctx2, "name")
	cancel()
	if err != nil {
		fmt.Printf("[get value] failed:%v\n", err)
		return
	}
	fmt.Println("[get value] success")
	for _, ev := range resp.Kvs {
		fmt.Printf("键值对是%s:%s,创建版本%v,修改的次数%v，修改版本%v\n", ev.Key, ev.Value, ev.CreateRevision, ev.Version, ev.ModRevision)
	}

	//// watch 操作 ，获取key的变化
	//watChan := cli.Watch(context.Background(), "etcd_key")
	//for wResp := range watChan {
	//	for _, ev := range wResp.Events {
	//		fmt.Printf("变化后的Type:%s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
	//	}
	//}

}
