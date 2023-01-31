# Grpc

![](.grpc_images/grpc_layer.png)

最底层为TCP或Unix Socket协议，在此之上是HTTP/2协议的实现，然后在HTTP/2协议之上又构建了针对Go语言的gRPC核心库。应用程序通过gRPC插件生产的Stub代码和gRPC核心库通信，也可以直接和gRPC核心库通信。

## grpc 流程
![](.grpc_images/simple_grpc_process.png)

在 gRPC 中, 可以将 gRPC 的流程大致分为两个阶段, 分别是 RPC 连接阶段, 以及 RPC 交互阶段

- 在 RPC 连接阶段, client 和 server 之间建立起来 TCP 连接, 并且由于 gRPC 底层依赖于 HTTP2, 因此 client 和 server 还需要协调 frame 的相关设置, 例如 frame 的大小, 滑动窗口的大小等等.
- 在 RPC 交互阶段, client 将数据发送给 server, 并等待 server 执行指定 method 之后返回结果.
。

## grpc分类
### 1. unary
![](.grpc_images/unary.png)

### 2. client streaming
![](.grpc_images/client_streaming.png)    

### 3. server streaming
![](.grpc_images/server_streaming.png)

### 4. bidi streaming   
![](.grpc_images/bidi_streaming.png)

## 拦截器
![](.grpc_images/intercepter.png)


## grpc调优
* GRPC默认的参数对于传输大数据块来说不够友好，我们需要进行特定参数的调优。

* MaxSendMsgSizeGRPC最大允许发送的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。

* MaxRecvMsgSizeGRPC最大允许接收的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。

* InitialWindowSize基于Stream的滑动窗口，类似于TCP的滑动窗口，用来做流控，默认64KiB，吞吐量上不去，Client和Server我们调到1GiB。

* InitialConnWindowSize基于Connection的滑动窗口，默认16 * 64KiB，吞吐量上不去，Client和Server我们也都调到1GiB。

* KeepAliveTime每隔KeepAliveTime时间，发送PING帧测量最小往返时间，确定空闲连接是否仍然有效，我们设置为10S。

* KeepAliveTimeout超过KeepAliveTimeout，关闭连接，我们设置为3S。

* PermitWithoutStream如果为true，当连接空闲时仍然发送PING帧监测，如果为false，则不发送忽略。我们设置为true


## 参考链接
1. [官方 example](https://github.com/grpc/grpc-go/tree/master/examples/features)