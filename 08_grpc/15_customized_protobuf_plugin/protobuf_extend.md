<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [protobuf插件](#protobuf%E6%8F%92%E4%BB%B6)
  - [自定义选项](#%E8%87%AA%E5%AE%9A%E4%B9%89%E9%80%89%E9%A1%B9)
  - [场景描述](#%E5%9C%BA%E6%99%AF%E6%8F%8F%E8%BF%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

extend google.protobuf.FieldOptions {
  ...
}
```

proto3 定义
```protobuf

syntax = "proto3";

package main;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
    string default_string = 50000;
    int32 default_int = 50001;
}

message Message {
    string name = 1 [(default_string) = "gopher"];
    int32 age = 2[(default_int) = 10];
}
```
其中成员后面的方括号内部的就是扩展语法。重新生成Go语言代码，里面会包含扩展选项相关的元信息
```go
var E_DefaultString = &proto.ExtensionDesc{
    ExtendedType:  (*descriptor.FieldOptions)(nil),
    ExtensionType: (*string)(nil),
    Field:         50000,
    Name:          "main.default_string",
    Tag:           "bytes,50000,opt,name=default_string,json=defaultString",
    Filename:      "helloworld.proto",
}

var E_DefaultInt = &proto.ExtensionDesc{
    ExtendedType:  (*descriptor.FieldOptions)(nil),
    ExtensionType: (*int32)(nil),
    Field:         50001,
    Name:          "main.default_int",
    Tag:           "varint,50001,opt,name=default_int,json=defaultInt",
    Filename:      "helloworld.proto",
}
```
