package main

import (
	"context"
	"go_package_example/22_etcd/conn"
	"log"

	// 注意是这，不是github.com/coreos/etcd v3.3.17+incompatible
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	//  创建连接客户端
	var cli *clientv3.Client
	var err error
	cli, err = conn.GetClient()
	if err != nil {
		return
	}
	//创建一个session，并根据业务情况设置锁的ttl
	s, _ := concurrency.NewSession(cli, concurrency.WithTTL(3))
	defer s.Close()
	//初始化一个锁的实例，并进行加锁解锁操作。
	mu := concurrency.NewMutex(s, "mutex-prefix")
	if err := mu.Lock(context.TODO()); err != nil {
		log.Fatal("m lock err: ", err)
	}
	//do something
	if err := mu.Unlock(context.TODO()); err != nil {
		log.Fatal("m unlock err: ", err)
	}

}
