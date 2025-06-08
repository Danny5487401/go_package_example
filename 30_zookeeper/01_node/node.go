package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

/**
 * 获取一个zk连接
 * @return {[type]}
 */
func getConnect(zkList []string) (conn *zk.Conn) {
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	if err != nil {
		log.Fatalf("connect err : %v", err)
	}
	return
}

func createNode() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)

	defer conn.Close()
	var flags int32 = zk.FlagPersistent
	//flags有4种取值：
	//zk.FlagPersistent=0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//FlagEphemeralSequential=3 ，即，短暂且自动添加序号
	conn.Create("/go_servers", nil, flags, zk.WorldACL(zk.PermAll)) // zk.WorldACL(zk.PermAll)控制访问权限模式

}

/*
删改与增不同在于其函数中的version参数,其中version是用于 CAS支持
func (c *Conn) Set(path string, data []byte, version int32) (*Stat, error)
func (c *Conn) Delete(path string, version int32) error

demo：
if err = conn.Delete(migrateLockPath, -1); err != nil {
    log.Error("conn.Delete(\"%s\") error(%v)", migrateLockPath, err)
}
*/

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
	listNode()
}
