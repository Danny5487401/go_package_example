package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"go_package_example/08_grpc/17_fieldmask/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
}

func (s *UserServer) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*emptypb.Empty, error) {
	// naming func(string) string) (Mask, error): 接收的naming参数本质上是一个将字段掩码字段名映射到 Go 结构中使用的名称的函数，它必须根据你的实际需求实现。
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(in.FieldMask, strcase.ToCamel)
	var userDst = make(map[string]interface{})
	// 将数据读取到map[string]interface{}
	// fieldmask-utils支持读取到结构体等
	fieldmask_utils.StructToMap(mask, in.User, userDst)
	// do update with userDst
	fmt.Printf("userDst:%#v\n", userDst)
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//1。初始化
	s := grpc.NewServer()
	// 2.注册服务 service放在 m map[string]*service 中
	proto.RegisterUserServiceServer(s, &UserServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
