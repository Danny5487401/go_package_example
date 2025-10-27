<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [SSE（Server-Sent Events）](#sseserver-sent-events)
  - [SSE的主要特点](#sse%E7%9A%84%E4%B8%BB%E8%A6%81%E7%89%B9%E7%82%B9)
  - [格式](#%E6%A0%BC%E5%BC%8F)
  - [性能优化](#%E6%80%A7%E8%83%BD%E4%BC%98%E5%8C%96)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# SSE（Server-Sent Events）


SSE（Server-Sent Events）是一种服务器推送技术，允许服务器通过HTTP连接向客户端发送实时更新。
与WebSocket不同，SSE是单向通信机制，只能从服务器向客户端推送数据，不支持客户端向服务器发送数据。


## SSE的主要特点

- 基于HTTP协议：不需要特殊的协议支持，使用标准的HTTP连接
- 自动重连：浏览器内置支持断线重连机制
- 事件Id和类型：支持消息Id和事件类型，便于客户端处理不同类型的事件
- 纯文本传输：使用UTF-8编码的文本数据
- 单向通信：只能服务器向客户端推送数据


## 格式

协议内容放在http返回的body里，每次返回一个Event信息。每个Event里可以包含5个属性：

* id, id用于表示Event的序号，客户端通过序号实现断线重连功能。需要重连的时候，客户端在HTTP的header里加一个Last-Event-ID字段，把最后接收到的id传给服务端。服务端实现了重连功能，就能继续传Last-Event-ID之后的消息给客户端。
* event, event表示自定义事件类型，客户端通过该字段区分不同消息。
* data, data表示返回的业务数据，如果数据很长可以分成多行返回
* retry, retry表示重连的间隔，以毫秒为单位。
* :（注释消息）

SSE的数据格式非常简单，每条消息由一个或多个字段组成，每个字段由字段名、冒号和字段值组成，以换行符分隔：
```html
field: value\n
```

一个完整的SSE消息示例：
```html
id: 1\n
event: update\n
data: {"message": "Hello, World!"}\n\n
```
其中，双换行符（\n\n）表示一条消息的结束。

## 性能优化
在高并发场景下，SSE服务的性能优化至关重要：

连接管理：SSE使用长连接，服务器需处理大量并发连接，可以设置最大连接数或使用连接池。
数据压缩：对于大数据量推送，使用HTTP压缩（如gzip）减少带宽占用。
事件合并：将多个小事件合并为一个大事件，降低推送频率。