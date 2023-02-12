# Grpc
准确来说gRPC设计上是分层的，底层支持不同的协议，目前gRPC支持：

- [gRPC over HTTP2](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
- [gRPC Web](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md)

但是大多数情况下，讨论都是基于gRPC over HTTP2。
![](.grpc_images/grpc_layer.png)

最底层为TCP或Unix Socket协议，在此之上是HTTP/2协议的实现，然后在HTTP/2协议之上又构建了针对Go语言的gRPC核心库。应用程序通过gRPC插件生产的Stub代码和gRPC核心库通信，也可以直接和gRPC核心库通信。

## grpc 流程
![](.grpc_images/simple_grpc_process.png)


在 gRPC 中, 可以将 gRPC 的流程大致分为两个阶段, 分别是 RPC 连接阶段, 以及 RPC 交互阶段

- 在 RPC 连接阶段, client 和 server 之间建立起来 TCP 连接, 并且由于 gRPC 底层依赖于 HTTP2, 因此 client 和 server 还需要协调 frame 的相关设置, 例如 frame 的大小, 滑动窗口的大小等等.
- 在 RPC 交互阶段, client 将数据发送给 server, 并等待 server 执行指定 method 之后返回结果。


### Server 流程

在 RPC 连接阶段: server 在完成一些初始化的配置之后, 开始监听某个 TCP 端口. 在和某个 client 建立了 TCP 连接之后完成 HTTP2 settings frame 的交互.


在 RPC 交互阶段:

1. Server 等待 client 发来的 header frame, 从而创建出一个 stream 对象来管理整个交互流程. 根据 header frame 中的信息, server 知道 client 请求的是哪个 service 的哪个 method.
2. Server 接收到 client 发来的 data frame, 并执行 method.
3. Server 将执行是否成功等信息方法放在 header frame 中发送给 client.
4. Server 将 method 执行的结果 (返回值) 放在 data frame 中发送给 client


### client 大概流程

在 RPC 连接阶段, client 接收到一个目标地址 (string) 和一系列的 DialOptions, 然后

1. 配置连接参数, interceptor 等, 启动 resolver.
2. Resolver 根据目标地址获取 server 的地址列表 (比如一个 DNS name 可能会指向多个 server ip, dnsResovler 是 gRPC 内置的 resolver 之一). 启动 balancer.
3. Balancer 根据平衡策略, 从诸多 server 地址中选择一个或多个建立 TCP 连接.
4. client 在 TCP 连接建立完成之后, 等待 server 发来的 HTTP2 Settings frame, 并调整自身的 HTTP2 相关配置. 随后向 server 发送 HTTP2 Settings frame.


在 RPC 交互阶段,某个 local method 被调用之后,

1. Client 创建一个 stream 对象用来管理整个交互流程.
2. Client 将 service name, method name 等信息放到 header frame 中并发送给 server.
3. Client 将 method 的参数信息放到 data frame 中并发送给 server.
4. Client 等待 server 传回的 header frame 和 data frame. 一次 RPC call 的 result status 会被包含在 header frame 中, 而 method 的返回值被包含在 data frame 中


### wireshark 抓包结果
[案例代码](./01_grpc_helloworld)

![](.grpc_images/wireshark_config2.png)
点击 "Protobuf search paths" 旁的 Edit 按钮。

在弹出的对话框中，添加 proto 文件所在的路径。记得勾选右边的 "Load all files" 选项。

![](.grpc_images/wireshark_config1.png)
为了避免捕获无关流量，开始前，设定 Capture Filter 为 "tcp port 50051"，50051 为 gRPC 服务端的地址。

把 50051 端口的 TCP 消息解码成 HTTP2 协议信息，把 50051 端口的 TCP 消息解码成 HTTP2 协议信息

![](.grpc_images/wireshark_config.png)

![](.grpc_images/wireshark_info.png)

#### 1. 客户端请求
gGRPC把元数据放到HTTP/2 Headers里
![](.grpc_images/client_hello_req.png)

请求的参数在DATA frame里.
![](.grpc_images/client_hello_req2.png)

```protobuf
HEADERS (flags = END_HEADERS)
:method = POST
:scheme = http
:path = /google.pubsub.v2.PublisherService/CreateTopic
:authority = pubsub.googleapis.com
grpc-timeout = 1S
content-type = application/grpc+proto
grpc-encoding = gzip
authorization = Bearer y235.wef315yfh138vh31hv93hv8h3v

DATA (flags = END_STREAM)
<Length-Prefixed Message>

```


#### 2. 服务端返回
![](.grpc_images/server_hello_response.png)
```
HEADERS (flags = END_HEADERS)
:status = 200
grpc-encoding = gzip
content-type = application/grpc+proto

DATA
<Length-Prefixed Message>

HEADERS (flags = END_STREAM, END_HEADERS)
grpc-status = 0 # OK
trace-proto-bin = jher831yy13JHy3hc
```

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

GRPC默认的参数对于传输大数据块来说不够友好，我们需要进行特定参数的调优。

```go
// /Users/python/go/pkg/mod/google.golang.org/grpc@v1.47.0/clientconn.go
const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultClientMaxSendMessageSize    = math.MaxInt32
)
```
* MaxRecvMsgSizeGRPC最大允许接收的字节数，默认4MiB，如果超过了GRPC会报错。Client和Server我们都调到4GiB。

* InitialWindowSize基于Stream的滑动窗口，类似于TCP的滑动窗口，用来做流控，默认64KiB，吞吐量上不去，Client和Server我们调到1GiB。
```go
const (
	// The default value of flow control window size in HTTP2 spec.
	defaultWindowSize = 65535
)
```

* InitialConnWindowSize基于Connection的滑动窗口，默认16 * 64KiB，吞吐量上不去，Client和Server我们也都调到1GiB。

* KeepAliveTime每隔KeepAliveTime时间，发送PING帧测量最小往返时间，确定空闲连接是否仍然有效，我们设置为10s。

* KeepAliveTimeout超过KeepAliveTimeout，关闭连接，我们设置为3S。

* PermitWithoutStream如果为true，当连接空闲时仍然发送PING帧监测，如果为false，则不发送忽略。我们设置为true


## 参考链接
1. [官方 example 源码使用](https://github.com/grpc/grpc-go/tree/master/examples/features)
2. [wireshark 抓取grpc ](https://www.imlc.me/how-to-inspect-grpc-with-wireshark/)
3. [官方 grpc 设计准则](https://grpc.io/blog/principles/)
4. [官方 grpc 文档](https://github.com/grpc/grpc/tree/master/doc)