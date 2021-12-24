// Binary server is an example server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go_grpc_example/08_grpc/01_grpc_helloworld/proto"
	grpcErrProtobuf "go_grpc_example/08_grpc/07_grpc_error/proto"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

var port = flag.Int("port", 50052, "port number")

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	mu    sync.Mutex
	count map[string]int
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Track the number of times the user has been greeted.
	s.count[in.Name]++
	if s.count[in.Name] > 1 {

		// 使用默认自带的error code
		//st := status.New(codes.ResourceExhausted, "Request limit exceeded.")

		// 自定义的error code
		st := status.New(codes.Code(grpcErrProtobuf.Error_RESOURCE_ERR_NOT_FOUND), "Request limit exceeded.")
		// 添加具体描述信息
		ds, err := st.WithDetails(
			&epb.QuotaFailure{
				Violations: []*epb.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "Limit one greeting per person",
				}},
			},
			&grpcErrProtobuf.ErrDetail{
				Key: "hello",
				Msg: "danny",
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	flag.Parse()

	address := fmt.Sprintf(":%v", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{count: make(map[string]int)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
