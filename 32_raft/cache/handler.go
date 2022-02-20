package cache

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"net/http"
	"time"
)

func (h *HttpServer) doGet(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	key := vars.Get("key")
	if key == "" {
		h.log.Println("doGet() error, get nil key")
		fmt.Fprint(w, "")
		return
	}
	// 由于leader对set操作返回的时候，follower可能还没有apply数据，所以从follower的get查询可能返回旧数据或者空数据。
	// 如果要保证能从follower查询到的一定是最新的数据还需要很多额外的工作，即做到linearizable read
	// 参考文章： https://aphyr.com/posts/316-call-me-maybe-etcd-and-consul
	ret := h.ctx.St.Cm.Get(key)
	fmt.Fprintf(w, "%s\n", ret)
}

// doSet saves data to cache, only raft master node provides this api
func (h *HttpServer) doSet(w http.ResponseWriter, r *http.Request) {

	//对有些场景而言，应用程序需要感知leader状态，比如对stcache而言，理论上只有leader才能处理set请求来写数据，follower应该只能处理get请求查询数据。
	// 为了模拟说明这个情况，我们在stcache里面我们设置一个写标志位，当本节点是leader的时候标识位置true，可以处理set请求，否则标识位为false，不能处理set请求
	if !h.checkWritePermission() {
		fmt.Fprint(w, "write method not allowed\n")
		return
	}
	vars := r.URL.Query()

	key := vars.Get("key")
	value := vars.Get("value")
	if key == "" || value == "" {
		h.log.Println("doSet() error, get nil key or nil value")
		fmt.Fprint(w, "param error\n")
		return
	}

	event := LogEntryData{Key: key, Value: value}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		h.log.Printf("json.Marshal failed, err:%v", err)
		fmt.Fprint(w, "internal error\n")
		return
	}

	applyFuture := h.ctx.St.Raft.Raft.Apply(eventBytes, 5*time.Second)
	if err := applyFuture.Error(); err != nil {
		h.log.Printf("raft.Apply failed:%v", err)
		fmt.Fprint(w, "internal error\n")
		return
	}

	fmt.Fprintf(w, "ok\n")
}

// doJoin handles joining cluster request
func (h *HttpServer) doJoin(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()

	peerAddress := vars.Get("peerAddress")
	if peerAddress == "" {
		h.log.Println("invalid PeerAddress")
		fmt.Fprint(w, "invalid peerAddress\n")
		return
	}
	// r把这个节点加入到集群即可。申请加入的节点会进入follower状态，这以后集群节点之间就可以正常通信，leader也会把数据同步给follower。
	addPeerFuture := h.ctx.St.Raft.Raft.AddVoter(raft.ServerID(peerAddress), raft.ServerAddress(peerAddress), 0, 0)
	if err := addPeerFuture.Error(); err != nil {
		h.log.Printf("Error joining peer to raft, peerAddress:%s, err:%v, code:%d", peerAddress, err, http.StatusInternalServerError)
		fmt.Fprint(w, "internal error\n")
		return
	}
	fmt.Fprint(w, "ok")
}
