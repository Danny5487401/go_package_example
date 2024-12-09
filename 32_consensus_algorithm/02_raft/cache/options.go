package cache

import (
	"flag"
)

type Options struct {
	dataDir        string // data directory
	HttpAddress    string // http server address
	RaftTCPAddress string // construct Raft Address
	bootstrap      bool   // start as master or not
	JoinAddress    string // peer address to join
}

func NewOptions() *Options {
	opts := &Options{}

	var httpAddress = flag.String("http", "127.0.0.1:6000", "Http address")
	var raftTCPAddress = flag.String("Raft", "127.0.0.1:7000", "Raft tcp address")
	var node = flag.String("node", "node1", "Raft node name")
	var bootstrap = flag.Bool("bootstrap", false, "start as Raft cluster")
	var joinAddress = flag.String("join", "", "join address for Raft cluster")
	flag.Parse()

	opts.dataDir = "./" + *node
	opts.HttpAddress = *httpAddress
	opts.bootstrap = *bootstrap
	opts.RaftTCPAddress = *raftTCPAddress
	opts.JoinAddress = *joinAddress
	return opts
}
