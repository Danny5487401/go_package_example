package main

import (
	"fmt"
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
		fmt.Println(err)
	}
	return
}

/**
 * 测试连接
 * @return
 */
func test1() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)

	defer conn.Close()
	var flags int32 = 0
	//flags有4种取值：
	//0:永久，除非手动删除
	//zk.FlagEphemeral = 1:短暂，session断开则改节点也被删除
	//zk.FlagSequence  = 2:会自动在节点后面添加序号
	//3:Ephemeral和Sequence，即，短暂且自动添加序号
	conn.Create("/go_servers", nil, flags, zk.WorldACL(zk.PermAll)) // zk.WorldACL(zk.PermAll)控制访问权限模式

	time.Sleep(20 * time.Second)
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
 * 测试临时节点
 * @return {[type]}
 */
func test2() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)

	defer conn.Close()
	conn.Create("/testadaadsasdsaw", nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	time.Sleep(20 * time.Second)
}

/**
 * 获取所有节点
 */
func test3() {
	zkList := []string{"localhost:2181"}
	conn := getConnect(zkList)

	defer conn.Close()

	children, _, err := conn.Children("/go_servers")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v \n", children)
}

func main() {
	test3()
}
