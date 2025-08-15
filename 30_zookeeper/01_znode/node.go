package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

/**
 * 获取一个zk连接
 */
func getConnect(zkList []string) (conn *zk.Conn) {
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	if err != nil {
		log.Fatalf("connect err: %v", err)
	}
	return
}

func getOrCreateNode() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)
	defer conn.Close()

	path := "/go_servers"
	data, stat, err := conn.Get(path)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, zk.ErrNoNode) { // 节点不存在
			var flags int32 = zk.FlagPersistent
			//flags有4种取值：
			//zk.FlagPersistent=0:永久，除非手动删除
			//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
			//zk.FlagSequence  = 2:会自动在节点后面添加序号
			//FlagEphemeralSequential=3 ，即，短暂且自动添加序号
			data := []byte("data1")
			conn.Create(path, data, flags, zk.WorldACL(zk.PermAll)) // zk.WorldACL(zk.PermAll)控制访问权限模式

		}
		return
	}
	fmt.Printf("path:%s,data:%s,stat:%+v", path, string(data), stat)

}

/**
 * 获取所有节点
 */
func listNode() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)

	defer conn.Close()

	children, stat, err := conn.Children("/go_servers")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("stat: %+v \n", stat)
	fmt.Printf("children: %v \n", children)
}

func main() {
	getOrCreateNode()
	listNode()
}
