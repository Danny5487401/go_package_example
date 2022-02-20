package cache

import (
	"log"
	"net/http"
	"sync/atomic"
)

const (
	EnableWriteTrue  = int32(1)
	EnableWriteFalse = int32(0)
)

type HttpServer struct {
	ctx         *StCachedContext
	log         *log.Logger
	Mux         *http.ServeMux
	enableWrite int32
}

func NewHttpServer(ctx *StCachedContext, log *log.Logger) *HttpServer {
	mux := http.NewServeMux()
	s := &HttpServer{
		ctx:         ctx,
		log:         log,
		Mux:         mux,
		enableWrite: EnableWriteFalse,
	}

	mux.HandleFunc("/set", s.doSet)
	mux.HandleFunc("/get", s.doGet)

	// 加入集群
	mux.HandleFunc("/join", s.doJoin)
	return s
}

func (h *HttpServer) checkWritePermission() bool {
	return atomic.LoadInt32(&h.enableWrite) == EnableWriteTrue
}

func (h *HttpServer) SetWriteFlag(flag bool) {
	if flag {
		atomic.StoreInt32(&h.enableWrite, EnableWriteTrue)
	} else {
		atomic.StoreInt32(&h.enableWrite, EnableWriteFalse)
	}
}
