package main

import (
	"fmt"
	"go_grpc_example/32_raft/cache"

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

	logger := log.New(os.Stderr, "httpserver: ", log.Ldate|log.Ltime)
	httpServer := cache.NewHttpServer(ctx, logger)
	st.Hs = httpServer
	go func() {
		http.Serve(l, httpServer.Mux)
	}()

	raft, err := cache.NewRaftNode(st.Opts, ctx)
	if err != nil {
		st.Log.Fatal(fmt.Sprintf("new raft node failed:%v", err))
	}
	st.Raft = raft

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
