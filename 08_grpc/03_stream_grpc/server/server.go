package main

import (
	"fmt"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/Danny5487401/go_package_example/08_grpc/03_stream_grpc/proto"
)

const PORT = ":50051"

type server struct {
	proto.UnimplementedGreeterServer
}

// GetStream 服务端 推送流
func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		res.Send(&proto.StreamResData{Data: fmt.Sprintf("服务端发送给客户端连续数据%v", time.Now().Unix())})
		time.Sleep(1 * time.Second)
		if i > 10 {
			break
		}
	}
	return nil
}

// PutStream 客户端 推送流
func (s *server) PutStream(cliStr proto.Greeter_PutStreamServer) error {

	ctx := cliStr.Context()
	peerAddr := "0.0.0.0"
	if peerInfo, ok := peer.FromContext(ctx); ok {
		peerAddr = peerInfo.Addr.String()
	}
	fmt.Printf("客户端 推送流：获取客户端地址信息%v \n", peerAddr)
	for {
		if tem, err := cliStr.Recv(); err == nil {
			log.Println(tem)
		} else {
			log.Println("break, err :", err)
			break
		}
	}

	return nil
}

// AllStream 客户端服务端 双向流
func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error {

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			data, err := allStr.Recv()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("服务端接收数据", data.GetData())

		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 5; i++ {
			allStr.Send(&proto.StreamResData{Data: fmt.Sprintf("服务端发送连续数据:%v", i)})
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil

}

func main() {
	//监听端口
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		return
	}
	//创建一个grpc 服务器
	s := grpc.NewServer()
	//注册事件
	proto.RegisterGreeterServer(s, &server{})
	//处理链接
	s.Serve(lis)
}
