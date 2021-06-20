package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main()  {
	//  创建连接客户端
	cli,err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			"81.68.197.3:2379",
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil{
		fmt.Printf("[init ctcd client ] failed:%v\n",err)
		return
	}

	fmt.Println("[init etcd client ] success")
	defer cli.Close()

	// put 操作
	//ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	//rsp ,err := cli.Put(ctx,"etcd_danny_key","etcd_value")
	//cancel()
	//if err != nil{
	//	fmt.Printf("[put ctcd  ] failed:%v\n",err)
	//}
	//fmt.Printf("[put ctcd  ] success,rsp :%v\n",rsp)

	// get 操作
	ctx2, cancel := context.WithTimeout(context.Background(),time.Second)
	//resp,err := cli.Get(ctx2,"etcd_danny_key", clientv3.WithPrefix() )  //  clientv3.WithPrefix() 业务前缀
	resp,err := cli.Get(ctx2,"etcd_key")
	cancel()
	if err != nil{
		fmt.Printf("[get value] failed:%v\n",err)
		return
	}
	fmt.Println("[get value] success")
	for _,ev := range resp.Kvs{
		fmt.Printf("键值对是%s:%s\n",ev.Key,ev.Value)
	}

	// watch 操作 ，获取key的变化
	watChan := cli.Watch(context.Background(),"etcd_key")
	for wResp := range watChan {
		for _,ev := range wResp.Events{
			fmt.Printf("变化后的Type:%s Key:%s Value:%s\n",ev.Type,ev.Kv.Key,ev.Kv.Value)
		}
	}







}
