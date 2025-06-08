package main

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

func main() {
	servers := []string{"127.0.0.1:2181"} // 替换为你的 ZooKeeper 实例的地址
	conn, _, err := zk.Connect(servers, time.Second*10)
	if err != nil {
		log.Fatalf("Failed to connect to ZooKeeper: %v", err)
	}
	defer conn.Close()

	lockPath := "/distributed_lock"
	lock := zk.NewLock(conn, lockPath, zk.WorldACL(zk.PermAll))
	err = lock.Lock()
	if err != nil {
		log.Fatalf("Failed to acquire lock: %v", err)

	}
	log.Println("Create lock")
	lock.Unlock()
	log.Println("release lock")
}
