package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"net"

	"google.golang.org/grpc"

	"go_test_project/14_metadata_test/proto"
)

type Server struct {

}

func (s *Server)SayHello(ctx context.Context,request *proto.HelloRequest)  (*proto.HelloReply,error){
	md,ok := metadata.FromIncomingContext(ctx)
	if !ok{
		fmt.Println("get metadata errors")
	}
	/*
	timestamp [2021-03-20 15:10:53]
	:authority [127.0.0.1:9000]
	content-type [application/grpc]
	*/
	if nameSlice,ok := md["name"];ok{
		fmt.Println(nameSlice)
		for _,v := range nameSlice{
			fmt.Println(v)
		}
	}
	//for key, val := range md{
	//	fmt.Println(key,val)
	//}
	return &proto.HelloReply{
		Message:"hello," + request.Name,
	},nil
}

func main()  {
	g := grpc.NewServer()

	proto.RegisterGreeterServer(g,&Server{})
	lis,err := net.Listen("tcp","127.0.0.1:9000")
	if err != nil{
		panic("faild to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}
