package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net"
	"runtime/debug"

	"go_grpc_example/08_grpc/06_grpc_interceptor/proto"
	"google.golang.org/grpc"
)

// 拦截器
//在 gRPC 中，大类可分为两种 RPC 方法，与拦截器的对应关系是：
//普通方法：一元拦截器（grpc.UnaryInterceptor）
//流方法：流拦截器（grpc.StreamInterceptor）

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello," + request.Name,
	}, nil
}

func main() {
	// 1。定义单个拦截器
	//ctx context.Context：请求上下文
	//req interface{}：RPC 方法的请求参数
	//info *UnaryServerInfo：RPC 方法的所有信息
	//handler UnaryHandler：RPC 方法本身
	//interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//	// 继续处理请求
	//	fmt.Println("接收到新请求,开始时间")
	//	res, err := handler(ctx, req)
	//	fmt.Println("请求处理完成，结束时间")
	//	return res, err
	//}
	//opt := grpc.UnaryInterceptor(interceptor)
	//g := grpc.NewServer(opt)

	// 2。定义多个拦截器
	c := GetTLSCredentialsByCA()
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)

	proto.RegisterGreeterServer(server, &Server{})
	lis, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("启动失败")
		return
	}
	_ = server.Serve(lis)
}

func GetTLSCredentialsByCA() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("08_grpc/06_grpc_interceptor/server.pem", "08_grpc/06_grpc_interceptor/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("08_grpc/06_grpc_interceptor/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	//尝试解析所传入的 PEM 编码的证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	//构建基于 TLS 的 TransportCredentials 选项
	c := credentials.NewTLS(&tls.Config{
		//tls.Config：Config 结构用于配置 TLS 客户端或服务器
		Certificates: []tls.Certificate{cert}, //Certificates：设置证书链，允许包含一个或多个
		ClientAuth:   tls.RequestClientCert,   //ClientAuth：要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientCAs:    certPool,                //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})
	return c
}

//日志拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

// 异常保护
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}

/*
注意：
	可以发现 gRPC 本身居然只能设置一个拦截器，难道所有的逻辑都只能写在一起？
解决
	关于这一点，你可以放心。采用开源项目 go-grpc-middleware 就可以解决这个问题
*/
