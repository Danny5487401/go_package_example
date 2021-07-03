package main

import (
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/grpclb"
	"sync"

	"go_grpc_example/08_grpc/23_stream_grpc/proto"

	"context"
	"log"
	"time"
)

const (
	ADDRESS = "localhost:50051"
)

func main() {
	//通过grpc 库 建立一个连接
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	//通过刚刚的连接 生成一个client对象。
	c := proto.NewGreeterClient(conn)
	//调用服务端推送流
	reqStreamData := &proto.StreamReqData{Data: "客户端发送单次数据aaa"}
	res, _ := c.GetStream(context.Background(), reqStreamData)
	for {
		aa, err := res.Recv()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("接收到服务端数据", aa)
	}
	//客户端 推送 流
	putRes, _ := c.PutStream(context.Background())
	i := 1
	for {
		i++
		putRes.Send(&proto.StreamReqData{Data: "客户端发送连续数据ss"})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}
	//服务端 客户端 双向流
	allStr, _ := c.AllStream(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			data, err := allStr.Recv()
			if err != nil {
				break
			}
			log.Println("接收到服务端数据", data)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 5; i++ {
			allStr.Send(&proto.StreamReqData{Data: "客户端发送连续数据ssss"})
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	wg.Wait()

}
