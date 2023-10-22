<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [RPC 方法做自定义认证](#rpc-%E6%96%B9%E6%B3%95%E5%81%9A%E8%87%AA%E5%AE%9A%E4%B9%89%E8%AE%A4%E8%AF%81)
  - [gRPC 中默认定义 PerRPCCredential](#grpc-%E4%B8%AD%E9%BB%98%E8%AE%A4%E5%AE%9A%E4%B9%89-perrpccredential)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

#  RPC 方法做自定义认证

##  gRPC 中默认定义 PerRPCCredential
通常来说，认证信息是需要每次都携带，但如果需要单次携带 metadata，可以使用 metadata.NewOutgoingContext 方法来创建一个携带 metadata 的 context。

```go
type PerRPCCredentials interface {
    GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) // 获取当前请求认证所需的元数据（metadata）
    RequireTransportSecurity() bool // 是否需要基于 TLS 认证进行安全传输
}
```
gRPC 默认提供用于自定义认证的接口，它的作用是将所需的安全认证信息添加到每个 RPC 方法的上下文中。

