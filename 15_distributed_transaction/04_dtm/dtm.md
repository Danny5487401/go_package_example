<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [DTM](#dtm)
  - [协议分类](#%E5%8D%8F%E8%AE%AE%E5%88%86%E7%B1%BB)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# DTM
DTM是一款开源的分布式事务管理器，解决跨数据库、跨服务、跨语言栈更新数据的一致性问题

## 协议分类
dtm支持http协议和gRPC协议

- HTTP协议：dtm服务器会监听HTTP端口36789，这里的例子业务会监听HTTP 8081
- gRPC协议：dtm服务器会监听gRPC端口36790，这里的例子业务会监听gRPC 58081