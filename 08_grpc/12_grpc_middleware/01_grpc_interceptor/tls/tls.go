package grpctls

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func GetTLSCredentialsByCA() credentials.TransportCredentials {
	cert, err := tls.LoadX509KeyPair("08_grpc/12_grpc_middleware/01_grpc_interceptor/server.pem", "08_grpc/12_grpc_middleware/01_grpc_interceptor/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("08_grpc/12_grpc_middleware/01_grpc_interceptor/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	// 尝试解析所传入的 PEM 编码的证书
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	// 构建基于 TLS 的 TransportCredentials 选项
	c := credentials.NewTLS(&tls.Config{
		//tls.Config：Config 结构用于配置 TLS 客户端或服务器
		Certificates: []tls.Certificate{cert}, //Certificates：设置证书链，允许包含一个或多个
		ClientAuth:   tls.RequestClientCert,   //ClientAuth：要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientCAs:    certPool,                //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})
	return c
}
