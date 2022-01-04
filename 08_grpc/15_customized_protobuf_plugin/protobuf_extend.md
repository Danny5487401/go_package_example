# protobuf插件

## 自定义选项
在 proto3 中，常见的实现插件的方式是使用 自定义选项，也就是 extend 标签，其中支持的 extend Options 有：

* MethodOptions
* ServiceOptions
* EnumOptions
* EnumValueOptions
* MessageOptions
* FieldOptions
* FileOptions
* OneofOptions
* ExtensionRangeOptions


## 场景描述
我们有很多的拦截器，其中不同的 service 可能会使用一个或多个拦截器，不同的 method 也可能会使用一个或多个拦截器，在 helloworld.proto 中

- service Greeter{} 支持登录令牌验证
- rpc SayHello1() 支持 IP 白名单限制和记录日志
- rpc SayHello2() 支持禁止记录日志
```protobuf
// helloworld.proto

service Greeter {
  rpc SayHello1 (HelloRequest) returns (HelloReply) {}
  rpc SayHello2 (HelloRequest) returns (HelloReply) {}
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

需要用到 MethodOptions 和 ServiceOptions 选项，通过名字大概也能猜到 MethodOptions 是定义方法选项的，ServiceOptions 是定义服务选项的。
```protobuf
extend google.protobuf.MethodOptions {
  ...
}

extend google.protobuf.ServiceOptions {
  ...
}
```