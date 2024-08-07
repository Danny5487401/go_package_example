<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [jsonpb](#jsonpb)
  - [使用](#%E4%BD%BF%E7%94%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# jsonpb




## 使用 
```go

type MemberResponse struct {
	Id    int32  `json "id"`
	Phone string `json "phone"`
	Age   int8   `json "age"`
}
// 返回结果	Id:12  Phone:"15112810201"


```

问题：
因为server未返回age字段，所以没有age。在某些情况下对前端也是不太友好的，尤其是APP客户端，更需要明确的json响应字段结构

解决方式
1. 直接修改经过protoc生成的member.pb.go文件代码，删除掉不希望被忽略的字段tag标签中的omitempty即可，
   但是*.pb.go一般我们不建议去修改它，而且我们会经常去调整grpc微服务协议中的方法或者字段内容，这样每次protoc之后，
   都需要我们去修改，这显然是不太现实的，因此就有了第二种办法；
2. 通过grpc官方库中的jsonpb来实现,官方在它的设定中有一个结构体用来实现protoc buffer转换为JSON结构，并可以根据字段来配置转换的要求，
   
```go
// Marshaler is a configurable object for converting between
   // protocol buffer objects and a JSON representation for them.
   type Marshaler struct {
   // 是否将枚举值设定为整数，而不是字符串类型.
   EnumsAsInts bool
   // 是否将字段值为空的渲染到JSON结构中
   EmitDefaults bool
   //缩进每个级别的字符串
   Indent string
   //是否使用原生的proto协议中的字段
   OrigName bool
   }
```