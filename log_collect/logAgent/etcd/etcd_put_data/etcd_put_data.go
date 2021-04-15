package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"go_test_project/log_collect/logAgent/utils"
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

	//put 操作
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	value1 := `[{"path":"/Users/python/kafka_test/nginx/nginx.log","topic":"nginx_log"},
				{"path":"/Users/python/kafka_test/redis/redis.log","topic":"redis_log"}]
			`
	//value2 := `[{"path":"/Users/python/kafka_test/nginx/nginx.log","topic":"nginx_log"},
	//			{"path":"/Users/python/kafka_test/redis/redis.log","topic":"redis_log"},
	//			{"path":"/Users/python/kafka_test/mysql/mysql.log","topic":"mysql_log"}]
	//			`
	ip,_ := utils.GetOutboundIP()
	etcdKey := fmt.Sprintf("/logAgent/%s/collect_config",ip)
	//rsp ,err := cli.Put(ctx,"/logAgent/collect_config",value1)
	rsp ,err := cli.Put(ctx,etcdKey,value1)
	cancel()
	if err != nil{
		fmt.Printf("[put ctcd  ] failed:%v\n",err)
	}
	fmt.Printf("[put ctcd  ] success,rsp :%v\n",rsp)

	//// get 操作
	ctx2, cancel := context.WithTimeout(context.Background(),time.Second)

	resp,err := cli.Get(ctx2,etcdKey)
	cancel()
	if err != nil{
		fmt.Printf("[get value] failed:%v\n",err)
		return
	}
	fmt.Println("[get value] success")
	for _,ev := range resp.Kvs{
		fmt.Printf("键值对是%s:%s\n",ev.Key,ev.Value)
	}

}

