package conn

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func GetClient() (cli *clientv3.Client, err error) {
	//  创建连接客户端
	cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{
			"106.14.35.115:2379",
		},
		Username:    "root",
		Password:    "root",
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("[init ctcd client ] failed:%v\n", err)
		return
	}

	fmt.Println("[init etcd client ] success")
	return
}
