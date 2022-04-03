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
- 服务端推送
#### 定义
![](.grpc_images/definition.png)

#### 二进制分帧层binary frame
![](.grpc_images/binary_frame.png)
binary frame在应用层和TCP层中间.
- 优先级控制 
- 流量控制
- 服务端推送

#### 多路复用multiplex(实现无序传输 )
![](.grpc_images/multi_routes.png)
![](.grpc_images/multi_routes2.png)
一个tcp上面多个stream，stream只是个抽象概念。
- 并行交错地发送多个请求，请求之间互不影响。
- 并行交错地发送多个响应，响应之间互不干扰。
- 使用一个连接并行发送多个请求和响应。

#### message消息
具有业务含义 

#### frame帧
![](.http_images/frame.png) 
最小传输单元
- 前三行，标准头部，9个字节
- 第四行，payload
flag介绍
![](.http_images/flag_frame.png)

##### frame types类型(10种，主要2种data和headers)
![](.grpc_images/frame_type.png)

1. Magic:
Magic 帧的主要作用是建立 HTTP/2 请求的前言。在 HTTP/2 中，要求两端都要发送一个连接前言，作为对所使用协议的最终确认，并确定 HTTP/2 连接的初始设置，客户端和服务端各自发送不同的连接前言。

Magic 帧是客户端的前言之一，内容为 PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n，以确定启用 HTTP/2 连接。
2. Settings连接级参数
![](.http_images/settings_frame.png)
SETTINGS 帧的主要作用是设置这一个连接的参数，作用域是整个连接而并非单一的流,比如最大并发量。

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



### http3.0
![](.grpc_images/http3.png)
