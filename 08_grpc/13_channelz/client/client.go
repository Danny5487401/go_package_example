// Binary client is an example client.
package main

import (
	"context"
	"flag"
	"google.golang.org/grpc/channelz/service"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:50052", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	//	/***** Set up the server serving channelz service. *****/
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(lis)
	defer s.Stop()

	/***** Initialize manual resolver and Dial *****/
	r := manual.NewBuilderWithScheme("whatever")
	// Set up a connection to the server.
	conn, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithResolvers(r), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Manually provide resolved addresses for the target.
	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: ":10001"}, {Addr: ":10002"}, {Addr: ":10003"}}})

	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.

	/***** Make 100 SayHello RPCs *****/
	for i := 0; i < 100; i++ {
		// Setting a 150ms timeout on the RPC.
		// 客户端设置： 客户端将使100个SayHello RPC到达指定的目标，并使用循环策略负载均衡工作负载。每个呼叫都有150毫秒的超时时间。记录RPC响应和错误是为了进行调试
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Printf("could not greet: %v", err)
		} else {
			log.Printf("Greeting: %s", r.Message)
		}
	}

	/***** Wait for user exiting the program *****/
	// Unless you exit the program (e.g. CTRL+C), channelz data will be available for querying.
	// Users can take time to examine and learn about the info provided by channelz.
	select {}
}
