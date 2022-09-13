package main

import (
	"context"
	"fmt"
	"go_package_example/01_consul/02_grpc/api"
	pb "go_package_example/08_grpc/01_grpc_helloworld/proto"
	grpctls "go_package_example/08_grpc/12_grpc_middleware/01_grpc_interceptor/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	address = ":3333"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func RegisterToConsul() {
	api.RegisterService("tencent.danny.games:8500", &api.ConsulService{
		Name: "helloworld",
		Tag:  []string{"grpc"},
		IP:   "8.tcp.ngrok.io", //grpc的调用地址
		Port: 19416,
	})
}

//health
type HealthImpl struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	fmt.Print("health checking\n")
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	opts := []grpc.ServerOption{
		grpc.Creds(grpctls.GetTLSCredentialsByCA()),
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(s, &server{})
	grpc_health_v1.RegisterHealthServer(s, &HealthImpl{})
	RegisterToConsul()
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	onExit()
}

// 监听信号
func onExit() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)

	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("Program Exit...", s)

			//api.DeregisterService("tencent.danny.games:8500", &api.ConsulService{
			//	Name: "helloworld",
			//	Tag:  []string{"grpc"},
			//	IP:   "abcdefgh.vaiwan.cn", //grpc的调用地址
			//	Port: 8888,
			//})
			return

		default:
			fmt.Println("other signal", s)
		}
	}

}
