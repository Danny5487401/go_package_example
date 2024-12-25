package main

import (
	"fmt"
	"github.com/Danny5487401/go_package_example/32_consensus_algorithm/02_raft/cache"

	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	st := &cache.StCached{
		Opts: cache.NewOptions(),
		Log:  log.New(os.Stderr, "stCached: ", log.Ldate|log.Ltime),
		Cm:   cache.NewCacheManager(),
	}
	ctx := &cache.StCachedContext{St: st}

	var l net.Listener
	var err error
	l, err = net.Listen("tcp", st.Opts.HttpAddress)
	if err != nil {
		st.Log.Fatal(fmt.Sprintf("listen %s failed: %s", st.Opts.HttpAddress, err))
	}
	st.Log.Printf("http server listen:%s", l.Addr())

	// http业务服务器
	logger := log.New(os.Stderr, "httpserver: ", log.Ldate|log.Ltime)
	httpServer := cache.NewHttpServer(ctx, logger)
	st.Hs = httpServer
	go func() {
		http.Serve(l, httpServer.Mux)
	}()

	// raft 节点启动
	raftNode, err := cache.NewRaftNode(st.Opts, ctx)
	if err != nil {
		st.Log.Fatal(fmt.Sprintf("new raft node failed:%v", err))
	}
	st.Raft = raftNode
	if st.Opts.JoinAddress != "" {
		err = cache.JoinRaftCluster(st.Opts)
		if err != nil {
			st.Log.Fatal(fmt.Sprintf("join raft cluster failed:%v", err))
		}
	}

	// monitor leadership
	// 模拟选举
	for {
		select {
		case leader := <-st.Raft.LeaderNotifyCh:
			if leader {
				st.Log.Println("become leader, enable write api")
				st.Hs.SetWriteFlag(true)
			} else {
				st.Log.Println("become follower, close write api")
				st.Hs.SetWriteFlag(false)
			}

		}
	}
}

// 运行 cd 32_raft
// go build -o main sql_squirrel_test.go
// ./main -bootstrap true -node node1
// 测试数据
// 写： curl "http://127.0.0.1:6000/set?key=name&value=danny"
// 读： curl "http://127.0.0.1:6000/get?key=name"
// ./main -join 127.0.0.1:6000 -Raft 127.0.0.1:7001 -http 127.0.0.1:6001 -node node2
// ./main -join 127.0.0.1:6000 -Raft 127.0.0.1:7002 -http 127.0.0.1:6002 -node node3
// 读： curl "http://127.0.0.1:6001/get?key=name"
// 读： curl "http://127.0.0.1:6002/get?key=name"
