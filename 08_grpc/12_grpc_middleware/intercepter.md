<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [截取器](#%E6%88%AA%E5%8F%96%E5%99%A8)
  - [使用](#%E4%BD%BF%E7%94%A8)
  - [grpc源码](#grpc%E6%BA%90%E7%A0%81)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 截取器
![](../.grpc_images/intercepter.png)

gRPC中的grpc.UnaryInterceptor和grpc.StreamInterceptor分别对普通方法和流方法提供了截取器的支持。我们这里简单介绍普通方法的截取器用法。

## 使用
```go
// 日志拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}
```

函数的ctx和req参数就是每个普通的RPC方法的前两个参数。
第三个info参数表示当前是对应的那个gRPC方法，第四个handler参数对应当前的gRPC方法函数。

```go
grpc.UnaryInterceptor(LoggingInterceptor)
```
不过gRPC框架中只能为每个服务设置一个截取器，因此所有的截取工作只能在一个函数中完成。开源的grpc-ecosystem项目中的go-grpc-middleware包已经基于gRPC对截取器实现了链式截取器的支持。

## grpc源码
```go
// /Users/python/go/pkg/mod/google.golang.org/grpc@v1.45.0/interceptor.go
```
分类
- 服务端拦截
- 客户端拦截