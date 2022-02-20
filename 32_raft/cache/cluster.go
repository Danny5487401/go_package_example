package cache

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type RaftNodeInfo struct {
	Raft           *raft.Raft
	fsm            *FSM
	LeaderNotifyCh chan bool
}

func newRaftTransport(opts *Options) (*raft.NetworkTransport, error) {
	address, err := net.ResolveTCPAddr("tcp", opts.RaftTCPAddress)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	return transport, nil
}

type lockedBytesBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBytesBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBytesBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func NewRaftNode(opts *Options, ctx *StCachedContext) (*RaftNodeInfo, error) {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(opts.RaftTCPAddress)
	var logbuf lockedBytesBuffer
	raftConfig.Logger = hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Info,
		Output: &logbuf,
	})

	// 需要两个条件同时满足才会生成和保存一次快照，默认config里面配置的条件比较高
	// SnapshotInterval指每间隔多久生成一次快照
	raftConfig.SnapshotInterval = 20 * time.Second
	// SnapshotThreshold 每commit多少log entry后生成一次快照
	raftConfig.SnapshotThreshold = 2

	// 当leader状态变化时会往这个chan写数据，写入的变更消息能够缓存在channel里面，应用程序能够通过它获取到最新的状态变化。
	leaderNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = leaderNotifyCh

	// 集群内部节点之间的通信渠道两种方式来实现，
	// 一种是通过TCPTransport，基于tcp，可以跨机器跨网络通信；
	// 另一种是InmemTransport，不走网络，在内存里面通过channel来通信。
	transport, err := newRaftTransport(opts)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(opts.dataDir, 0700); err != nil {
		return nil, err
	}

	fsm := &FSM{
		Ctx: ctx,
		Log: log.New(os.Stderr, "FSM: ", log.Ldate|log.Ltime),
	}
	snapshotStore, err := raft.NewFileSnapshotStore(opts.dataDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(opts.dataDir, "Raft-Log.bolt"))
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(opts.dataDir, "Raft-stable.bolt"))
	if err != nil {
		return nil, err
	}

	raftNode, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	if opts.bootstrap {
		// 集群最开始的时候只有一个节点，我们让第一个节点通过bootstrap的方式启动，它启动后成为leader。
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		raftNode.BootstrapCluster(configuration)
	}

	return &RaftNodeInfo{Raft: raftNode, fsm: fsm, LeaderNotifyCh: leaderNotifyCh}, nil
}

// joinRaftCluster joins a node to Raft cluster
func JoinRaftCluster(opts *Options) error {
	// 后续的节点启动的时候需要加入集群，启动的时候指定第一个节点的地址，并发送请求加入集群，这里我们定义成直接通过http请求

	url := fmt.Sprintf("http://%s/join?peerAddress=%s", opts.JoinAddress, opts.RaftTCPAddress)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "ok" {
		return errors.New(fmt.Sprintf("Error joining cluster: %s", body))
	}

	return nil
}
