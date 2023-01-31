# http介绍

## 基本概念
### RTT(Round Trip Time)
![](.grpc_images/RTT.png)

解释：Round trip time for a packet

公式： rtt = recvtime - sendtime

### tcp相关概念
1. timewait过程-->过多连接会timewait等待很久,端口占用
![](.grpc_images/time_wait.png)
![](.grpc_images/port_exhaustion.png)
一个tcp连接由四元组构成，即只有一个原地址和端口，也只有一个目的地址和端口,所以两个域肯定两个tcp连接。

2. tcp拥塞控制-->tcp慢启动,探测网络环境
   ![](.grpc_images/tcp_slow_start.png)


## http的三个版本介绍

### http1.0
![](.grpc_images/http1.1.png)

* Head-of-Line Blocking (No Pipelining)
* 一个连接同一时间只能处理一个请求，是串行，因为是无状态，不像tcp有序列号，http2会建立映射关系
* 如果当前请求阻塞，那么该连接就无法复用

### http1.1
![](.grpc_images/http_pipeline.png)

* 未完全解决head of line blocking，第一个阻塞，后面数据可能丢失(有管线化)
* fifo原则, 需要等待最后的响应
* 多数http proxy不支持
* 多数浏览器默认关闭 http1.1 pipeline

### http2.0
![](.grpc_images/http1.1VShttp2.0.png)

![](.grpc_images/definition.png)

优点：
- 多路复用，基于stream模型
- header压缩（hpack)
![](.http_images/header_packed_info.png)
![](.http_images/header_packed_info2.png)
![](.http_images/header_packed_info3.png)
![](.http_images/header_packed_info4.png)
  - 动态表：第二个请求只会发送与第一个请求不一样的内容
  - 静态表：常用的缺省表
- 流控
- 优先级在flag标志位上
- 服务端推送(server push)
  - Server Push指的是服务端主动向客户端推送数据，相当于对客户端的一次请求，服务端可以主动返回多次结果。
  这个功能打破了严格的请求---响应的语义，对客户端和服务端双方通信的互动上，开启了一个崭新的可能性。
  但是这个推送跟websocket中的推送功能不是一回事，Server Push的存在不是为了解决websocket推送的这种需求。
  

#### 1. 二进制分帧层binary frame
![](.grpc_images/binary_frame.png)
binary frame在应用层和TCP层中间.
- 优先级控制 
- 流量控制
- 服务端推送

##### message消息
具有业务含义，类似Request/Response消息，每个消息包含一个或多个帧

##### frame帧
![](.http_images/frame.png)
HTTP/2协议里通信的最小单位，每个帧有自己的格式，不同类型的帧负责传输不同的消息
- Length: 表示Frame Payload的大小，是一个24-bit的整型，表明Frame Payload的大小不应该超过2^24-1字节，但其实payload默认的大小是不超过2^14字节，可以通过SETTING Frame来设置SETTINGS_MAX_FRAME_SIZE修改允许的Payload大小。

- Type: 表示Frame的类型,目前定义了0-9共10种类型。
  ![](.http_images/flag_frame.png)
- Flags: 为一些特定类型的Frame预留的标志位，比如Header, Data, Setting, Ping等，都会用到。

- R: 1-bit的保留位，目前没用，值必须为0

- Stream Identifier: Steam的id标识，表明id的范围只能为0到2^31-1之间，其中0用来传输控制信息，比如Setting, Ping；客户端发起的Stream id 必须为奇数，服务端发起的Stream id必须为偶数；并且每次建立新Stream的时候，id必须比上一次的建立的Stream的id大；当在一个连接里，如果无限建立Stream，最后id大于2^31时，必须从新建立TCP连接，来发送请求。如果是服务端的Stream id超过上限，需要对客户端发送一个GOWAY的Frame来强制客户端重新发起连接。
  总结，前三行：标准头部，9个字节；第四行：payload


#### 2. 多路复用multiplex(实现无序传输 )
![](.grpc_images/multi_routes.png)
![](.grpc_images/multi_routes2.png)
一个tcp上面多个stream，stream只是个抽象概念，是建立链接后的一个双向字节流，用来传输消息，每次传输的是一个或多个帧。
- 并行交错地发送多个请求，请求之间互不影响。
- 并行交错地发送多个响应，响应之间互不干扰。
- 使用一个连接并行发送多个请求和响应。
- 每个stream都有唯一的id标识和一些优先级信息，客户端发起的stream的id为单数，服务端发起的stream id为偶数

#### 3. serverpush
![](.http_images/serverPush.png)
在HTTP/1.x里，为了展示这个页面，客户端会先发起一次 GET /index.html 的请求，拿到返回结果进行分析后，再发起两个资源的请求，一共是三次请求, 并且有串行的请求存在。

在HTTP/2里，当客户端发起 GET /index.html的请求后，如果服务端进行了Server Push的支持，那么会直接把客户端需要的/index.html和另外两份文件资源一起返回，避免了串行和多次请求的发送
  
#### frame types类型(10种，主要2种data和headers)
![](.grpc_images/frame_type.png)

1. Magic:
Magic 帧的主要作用是建立 HTTP/2 请求的前言。在 HTTP/2 中，要求两端都要发送一个连接前言，作为对所使用协议的最终确认，并确定 HTTP/2 连接的初始设置，客户端和服务端各自发送不同的连接前言。

Magic 帧是客户端的前言之一，内容为 PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n，以确定启用 HTTP/2 连接。
2. Settings连接级参数
![](.http_images/settings_frame.png)
SETTINGS 帧的主要作用是设置这一个连接的参数，作用域是整个连接而并非单一的流,比如最大并发量。
ETTINGS帧必须在连接开始时由通信双方发送，并且可以在任何其他时间由任一端点在连接的生命周期内发送。SETTINGS帧必须在id为0的stream上进行发送，不能通过其他stream发送；SETTINGS影响的是整个TCP链接，而不是某个stream；在SETTINGS设置出现错误时，必须当做connection error重置整个链接。

3. Headers frame

    HEADERS 帧的主要作用是存储和传播 HTTP 的标头信息。我们关注到 HEADERS 里有一些眼熟的信息，分别如下：
    ![](.http_images/header_frame.png)
    method：POST
    scheme：http
    path：/proto.SearchService/Search
    authority：:10001
    content-type：application/grpc
    user-agent：grpc-go/1.20.0-dev

4. Data frame(数据帧)
DATA 帧的主要作用是装填主体信息
![](.grpc_images/data_frame.png)
5. PING/PONG
   主要作用是判断当前连接是否仍然可用，也常用于计算往返时间。

6. WINDOW_UPDATE流量控制
   主要作用是管理和流的窗口控制。

7. GOAWAY停止

    用于关闭连接，GOAWAY允许端点优雅地停止接受新流，同时仍然完成先前建立的流的处理。
    这个就厉害了，当服务端需要维护时，发送一个GOAWAY的Frame给客户端，那么发送之前的Stream都正常处理了，发送GOAWAY后，客户端会新启用一个链接，继续刚才未完成的Stream发送。

### http3.0
![](.grpc_images/http3.png)
