package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	cli *clientv3.Client
)

// 返回需要收集的信息条目
type LogEntries struct {
	Path string `json:"path"`  //日志存放的路径
	Topic string `json:"topic"`  //日志在kafka发送的topic
}

// 初始化etcd
func Init(addr string,timeout time.Duration) (err error){
	//  创建连接客户端
	cli,err = clientv3.New(clientv3.Config{
		Endpoints: []string{
			addr,
		},
		DialTimeout: timeout,
	})
	if err != nil{
		fmt.Printf("[init ctcd client ] failed:%v\n",err)
		return
	}
	fmt.Println("[init etcd client ] success")
	return
}

// 从etcd中获取配置项
func GetConf(key string)(logEntries []*LogEntries,err error) {
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	resp,err := cli.Get(ctx,key)
	cancel()
	if err != nil{
		fmt.Printf("Get [%s] failed:%v\n",key,err)
		return nil,err
	}
	fmt.Printf("Get [%s] success",key)
	for _,ev := range resp.Kvs{
		fmt.Printf("键值对是%s:%s\n",ev.Key,ev.Value)
		err = json.Unmarshal(ev.Value,&logEntries)
		if err != nil {
			fmt.Printf("[Unmarshal]failed,%v",err)
			return
		}
	}
	return

}