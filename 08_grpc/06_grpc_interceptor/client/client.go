package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go_grpc_example/08_grpc/06_grpc_interceptor/proto"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
)

func main() {
	// 方法一 拦截器定义
	//interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//	start := time.Now()
	//	err := invoker(ctx, method, req, reply, cc, opts...)
	//	fmt.Printf("耗时%s\n", time.Since(start))
	//	return err
	//}
	//opt := grpc.WithUnaryInterceptor(interceptor)
	//conn,err := grpc.Dial("127.0.0.1:9000",grpc.WithInsecure(),opt)  //grpc.WithInsecure(),opt两个相同

	// 方法二
	//var opts []grpc.DialOption
	//opts = append(opts, grpc.WithInsecure())
	//opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	//conn, err := grpc.Dial("127.0.0.1:9000", opts...)

	cert, err := tls.LoadX509KeyPair("08_grpc/06_grpc_interceptor/client.pem", "08_grpc/06_grpc_interceptor/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("08_grpc/06_grpc_interceptor/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "go-grpc-example", //使用证书生成的Common Name 字段
		RootCAs:      certPool,
	})
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithTransportCredentials(c))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := proto.NewGreeterClient(conn)
	r, err := client.SayHello(context.Background(), &proto.HelloRequest{
		Name: "danny",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}

/*
在 Client 中绝大部分与 Server 一致，不同点的地方是，在 Client 请求 Server 端时，Client 端会使用根证书和 ServerName 去对 Server 端进行校验

简单流程大致如下：

	Client 通过请求得到 Server 端的证书
	使用 CA 认证的根证书对 Server 端的证书进行可靠性、有效性等校验
	校验 ServerName 是否可用、有效
注意点:
	golang 1.15+版本上，用 gRPC通过TLS实现数据传输加密时，会报错证书的问题
原因:
	因为我们用的证书，并没有开启SAN扩展（默认是没有开启SAN扩展）所生成的，导致客户端和服务端无法建立连接
什么是 SAN
	SAN(Subject Alternative Name) 是 SSL 标准 x509 中定义的一个扩展。使用了 SAN 字段的 SSL 证书，可以扩展此证书支持的域名，
	使得一个证书可以支持多个不同域名的解析。
*/
