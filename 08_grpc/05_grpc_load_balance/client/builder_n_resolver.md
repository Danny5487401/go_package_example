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
    ResolveNow(ResolveNowOption) // 被 gRPC 调用，以尝试再次解析目标名称。只用于提示，可忽略该方法。 需要并发安全的
    Close()
}
```

## 解析过程
![](.builder_n_resolver_images/builder_n_resolver.png)
1. SchemeBuilder将自身实例注册到resolver包的map中； 
2. grpc.Dial/DialContext时使用特定形式的target参数
3. 对target解析后，根据target.Scheme到resolver包的map中查找Scheme对应的Buider；
4. 调用Buider的Build方法
5. Build方法构建出SchemeResolver实例；
6. 后续由SchemeResolver实例监视service instance变更状态并在有变更的时候更新ClientConnection
7. 当address被作为target的实参传入grpc.DialContext后，它会被grpcutil.ParseTarget解析为一个resolver.Target结构体
```go
type Target struct {
	Scheme    string
	Authority string
	Endpoint  string
}
```

解析的方法
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/internal/grpcutil/target.go
// ParseTarget splits target into a resolver.Target struct containing scheme,
// authority and endpoint.
//
// If target is not a valid scheme://authority/endpoint, it returns {Endpoint:
// target}.
func ParseTarget(target string) (ret resolver.Target) {
	var ok bool
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, ok = split2(ret.Endpoint, "/")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	return ret
}
```

gRPC会根据Target.Scheme的值到resolver包中的builder map中查找是否有对应的Resolver Builder实例。
到目前为止gRPC内置的的resolver Builder都无法匹配该Scheme值。
```go
	resolverBuilder := cc.getResolver(cc.parsedTarget.Scheme)
```
```go
func (cc *ClientConn) getResolver(scheme string) resolver.Builder {
	for _, rb := range cc.dopts.resolvers {
		if scheme == rb.Scheme() {
			return rb
		}
	}
    return resolver.Get(scheme)
}
```
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/resolver/resolver.go
var (
    // m is a map from scheme to resolver builder.
    m = make(map[string]Builder)
    // defaultScheme is the default scheme to use.
    defaultScheme = "passthrough"
)
func Get(scheme string) Builder {
	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}
```