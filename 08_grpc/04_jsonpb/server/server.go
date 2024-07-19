package main

import (
	"github.com/Danny5487401/go_package_example/08_grpc/04_jsonpb/proto"
	"google.golang.org/grpc"

	"context"
	"net"
)

type Member struct {
	proto.UnimplementedMemberServer
}

// 获取用户信息的接口
func (m *Member) GetMember(context.Context, *proto.MemberRequest) (resp *proto.MemberResponse, err error) {
	resp = &proto.MemberResponse{}
	resp.Phone = "15112810201"
	resp.Id = 12
	return resp, nil
}

func main() {
	g := grpc.NewServer()

	proto.RegisterMemberServer(g, &Member{})
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}
