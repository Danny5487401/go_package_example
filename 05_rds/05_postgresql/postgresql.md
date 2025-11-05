<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [PostgreSQL](#postgresql)
  - [环境搭建](#%E7%8E%AF%E5%A2%83%E6%90%AD%E5%BB%BA)
  - [LISTEN 和 NOTIFY 机制](#listen-%E5%92%8C-notify-%E6%9C%BA%E5%88%B6)
  - [流复制](#%E6%B5%81%E5%A4%8D%E5%88%B6)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# PostgreSQL

## 环境搭建
```shell
docker run --rm --name my-postgres -v postgre-data:/var/lib/postgresql/data -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e LANG=C.UTF-8 -d postgres:16
```



## LISTEN 和 NOTIFY 机制


PostgreSQL 的 LISTEN 和 NOTIFY 是一种内置的消息通知系统，允许应用程序订阅数据库事件，并在事件发生时接收通知。


## 流复制
PostgreSQL 支持 COPY 操作，COPY 操作通过流复制协议（Streaming Replication Protocol）实现。COPY 命令允许在服务器之间进行高速批量数据传输，有三种流复制模式：


COPY-IN 模式 : 数据从客户端传输到服务器端。

COPY-OUT 模式 : 数据从服务器端传输到客户端。

COPY-BOTH 模式 : 服务器端和客户端数据可以双向传输。


## 参考
- https://www.postgresql.org/docs/17/index.html
- [掌握 PostgreSQL 的 LISTEN 和 NOTIFY 机制：实时数据库通知的艺术](https://blog.csdn.net/2401_85761762/article/details/139885992)