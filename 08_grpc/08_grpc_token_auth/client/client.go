package main

import (
	"context"
	"fmt"
	"github.com/Danny5487401/go_package_example/08_grpc/08_grpc_token_auth/proto"

	"google.golang.org/grpc"
)

type CustomCredential struct {
	User     string
	Password string
}

func (c CustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"appid": c.User, "appkey": c.Password}, nil
}

func (c CustomCredential) RequireTransportSecurity() bool {
	// 不需要基于 TLS 认证进行安全传输
	// 在真实的环境中建议必须要求底层启用安全的链接，否则认证信息有泄露和被篡改的风险。
	return false
}

func main() {
	// 方法一
	//interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error{
	//	start := time.Now()
	//	md := metadata.Pairs("appid","123456",
	//		"appkey","i am a key")
	//	ctx = metadata.NewOutgoingContext(context.Background(), md)
	//
	//	err := invoker(ctx,method,req,reply,cc,opts...)
	//	fmt.Printf("耗时%s\n",time.Since(start))
	//	return err
	//}
	//var opts []grpc.DialOption
	//opts = append(opts,grpc.WithInsecure())
	//opts = append(opts,grpc.WithUnaryInterceptor(interceptor))
	// 以上可以换另外一种简单方法

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithPerRPCCredentials(CustomCredential{
		User:     "name",
		Password: "danny",
	}))

	conn, err := grpc.Dial("127.0.0.1:9000", opts...)

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Name: "danny",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
