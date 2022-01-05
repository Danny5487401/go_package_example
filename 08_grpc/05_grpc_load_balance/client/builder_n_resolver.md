# Builder接口和Resolver接口
gRPC已提供了简单的负载均衡策略（如：Round Robin），我们只需实现它提供的Builder和Resolver接口，就能完成gRPC客户端负载均衡。

Builder接口：创建一个resolver（本文称之服务发现），用于监视名称解析更新。
```go
    type Builder interface {
    Build(target Target, cc ClientConn, opts BuildOption) (Resolver, error)//为给定目标创建一个新的resolver，当调用grpc.Dial()时执行
    Scheme() string //返回此resolver支持的方案
}

```

Resolver接口：监视指定目标的更新，包括地址更新和服务配置更新。
```go
type Resolver interface {
    ResolveNow(ResolveNowOption) // 被 gRPC 调用，以尝试再次解析目标名称。只用于提示，可忽略该方法。
    Close()
}
```
