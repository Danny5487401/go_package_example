//package etcd
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/coreos/etcd/clientv3"
//
//	//"github.com/coreos/etcd/clientv3"
//	"time"
//)
//
//var (
//	cli *clientv3.Client
//)
//
//// 返回需要收集的信息条目
//type LogEntries struct {
//	Path  string `json:"path"`  //日志存放的路径
//	Topic string `json:"topic"` //日志在kafka发送的topic
//}
//
//// 初始化etcd
//func Init(addr string, timeout time.Duration) (err error) {
//	//  创建连接客户端
//	cli, err = clientv3.New(clientv3.Config{
//		Endpoints: []string{
//			addr,
//		},
//		DialTimeout: timeout,
//	})
//	if err != nil {
//		fmt.Printf("[init ctcd client ] failed:%v\n", err)
//		return
//	}
//	fmt.Println("[init etcd client ] success")
//	return
//}
//
//// 从etcd中获取配置项
//func GetConf(key string) (logEntries []*LogEntries, err error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	resp, err := cli.Get(ctx, key)
//	cancel()
//	if err != nil {
//		fmt.Printf("Get [%s] failed:%v\n", key, err)
//		return nil, err
//	}
//	fmt.Printf("Get [%s] success", key)
//	for _, ev := range resp.Kvs {
//		fmt.Printf("键值对是%s:%s\n", ev.Key, ev.Value)
//		err = json.Unmarshal(ev.Value, &logEntries)
//		if err != nil {
//			fmt.Printf("[Unmarshal]failed,%v", err)
//			return
//		}
//	}
//	return
//
//}
//func WatchConfChange(key string, ch chan<- []*LogEntries) {
//	watChan := cli.Watch(context.Background(), key)
//	for wResp := range watChan {
//		for _, ev := range wResp.Events {
//			fmt.Printf("变化后的Type:%s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
//			// 不需要判断操作的类型 delete,put
//			var newConf []*LogEntries
//			//if ev.Type == clientv3.EventTypeDelete{
//			//	// 如果是删除操作,手动传递一个空的[]*LogEntries
//			//
//			//}else {
//			//	err := json.Unmarshal(ev.Kv.Value,&newConf)
//			//	if err != nil{
//			//		fmt.Printf("unmarshal failed,err:%v\n",err)
//			//		continue
//			//	}
//			//}
//			// 以下等价写法
//			if ev.Type != clientv3.EventTypeDelete {
//				err := json.Unmarshal(ev.Kv.Value, &newConf)
//				if err != nil {
//					fmt.Printf("unmarshal failed,err:%v\n", err)
//					continue
//				}
//			}
//
//			fmt.Printf("unmarshal new conf success:%v\n", newConf)
//			// 通知 tailTask
//			ch <- newConf
//		}
//	}
//}
