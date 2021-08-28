package main

import (
	"context"
	"net"

	"go_grpc_example/08_grpc/07_grpc_error/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return nil, status.Errorf(codes.AlreadyExists, "账号已存在%s", request.Name)
	//return &proto.HelloReply{
	//	Message:"hello," + request.Name,
	//},nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	_ = g.Serve(lis)
}

/*
gRPC的错误码，原代码见链接，我们大概了解其原因即可：
	OK 正常
	Canceled 客户端取消
	Unknown 未知
	InvalidArgument 未知参数
	DeadlineExceeded 超时
	NotFound 未找到资源
	AlreadyExists 资源已经创建
	PermissionDenied 权限不足
	ResourceExhausted 资源耗尽
	FailedPrecondition 前置条件不满足
	Aborted 异常退出
	OutOfRange 超出范围
	Unimplemented 未实现方法
	Internal 内部问题
	Unavailable 不可用状态
	DataLoss 数据丢失
	Unauthenticated 未认证
*/
