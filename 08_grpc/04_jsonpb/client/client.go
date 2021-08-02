package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"go_grpc_example/08_grpc/04_jsonpb/proto"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewMemberClient(conn)
	r, err := c.GetMember(context.Background(), &proto.MemberRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
