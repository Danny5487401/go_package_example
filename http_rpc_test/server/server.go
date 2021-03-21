package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {

}

func (s *HelloService) Hello(request string, reply *string) error{
	*reply = "hello, " + request
	return nil
}

func main(){
	// 1.实例化一个server
	_ = rpc.RegisterName("HelloService", &HelloService{})
	http.HandleFunc("/jsonrpc",func(w http.ResponseWriter, r *http.Request){
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			Writer:w,
			ReadCloser:r.Body,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	http.ListenAndServe(":1234",nil)
}
