package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go_package_example/22_etcd/conn"
	"time"
)

func main() {
	var (
		client       *clientv3.Client
		err          error
		kv           clientv3.KV
		keepResp     *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
	)

	client, err = conn.GetClient()
	if err != nil {
		return
	}

	//创建租约
	lease := clientv3.NewLease(client)
	//判断是否有问题
	if leaseRes, err := lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("授权一个租约失败", err)
		return
	} else {
		//得到租约id
		leaseId := leaseRes.ID

		//定义一个上下文使得租约5秒过期
		ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)

		//自动续租（底层会每次讲租约信息扔到 <-chan *clientv3.LeaseKeepAliveResponse 这个管道中）,这里ttl是10s
		if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
			fmt.Println(err)
			return
		}
		//启动一个新的协程来select这个管道
		go func() {
			for {
				select {
				case keepResp = <-keepRespChan:
					if keepResp == nil {
						fmt.Println("租约失效了")
						goto END //失效跳出循环
					} else {
						//每秒收到一次应答
						fmt.Printf("收到租约应答%v,ttl是%v\n", keepResp.ID, keepResp.TTL)
					}

				}
			}
		END:
		}()
		//得到操作键值对的kv
		kv = clientv3.NewKV(client)
		//进行写操作
		if putResp, err := kv.Put(context.TODO(), "/cron/lock/job1", "danny", clientv3.WithLease(leaseId) /*高速etcd这个key对应的租约*/); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Printf("写入成功,版本是%v,raft的term是%v\n", putResp.Header.Revision, putResp.Header.RaftTerm /*这东西你可以理解为每次操作的id*/)
		}
	}
	//每2秒监听这个key的租约是否过期
	for {
		var getResp *clientv3.GetResponse
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}

		fmt.Println("kv没过期", getResp.Kvs)
		time.Sleep(2 * time.Second)

	}
}
