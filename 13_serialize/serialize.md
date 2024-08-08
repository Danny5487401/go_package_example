<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [序列化Marshal和反序列化UnMarshal](#%E5%BA%8F%E5%88%97%E5%8C%96marshal%E5%92%8C%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96unmarshal)
  - [常见序列化协议](#%E5%B8%B8%E8%A7%81%E5%BA%8F%E5%88%97%E5%8C%96%E5%8D%8F%E8%AE%AE)
    - [xml（Extensible Markup Language）](#xmlextensible-markup-language)
    - [JSON(JavaScript Object Notation, JS 对象标记)](#jsonjavascript-object-notation-js-%E5%AF%B9%E8%B1%A1%E6%A0%87%E8%AE%B0)
    - [Thrift](#thrift)
    - [Avro](#avro)
    - [Protobuf](#protobuf)
  - [json 协议](#json-%E5%8D%8F%E8%AE%AE)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 序列化Marshal和反序列化UnMarshal

- 序列化（编码）是将对象序列化为二进制形式（字节数组），主要用于网络传输、数据持久化等；
- 而反序列化（解码）则是将从网络、磁盘等读取的字节数组还原成原始对象，主要用于网络传输对象的解码，以便完成远程调用。

影响序列化性能的关键因素：序列化后的码流大小（网络带宽的占用）、序列化的性能（CPU资源占用）；是否支持跨语言（异构系统的对接和开发语言切换）。

## 常见序列化协议

### xml（Extensible Markup Language）
- 优点：人机可读性好，可指定元素或特性的名称。
- 缺点：序列化数据只包含数据本身以及类的结构，不包括类型标识和程序集信息；只能序列化公共属性和字段；不能序列化方法；文件庞大，文件格式复杂，传输占带宽。

适用场景：当做配置文件存储数据，实时数据转换。

### JSON(JavaScript Object Notation, JS 对象标记) 
是一种轻量级的数据交换格式，
- 优点：兼容性高、数据格式比较简单，易于读写、序列化后数据较小，可扩展性好，兼容性好、与XML相比，其协议比较简单，解析速度比较快。
- 缺点：数据的描述性比XML差、不适合性能要求为ms级别的情况、额外空间开销比较大。

适用场景（可替代ＸＭＬ）：跨防火墙访问、可调式性要求高、基于Web browser的Ajax请求、传输数据量相对小，实时性要求相对低（例如秒级别）的服务。

### Thrift
不仅是序列化协议，还是一个RPC框架
- 优点：序列化后的体积小, 速度快、支持多种语言和丰富的数据类型、对于数据字段的增删具有较强的兼容性、支持二进制压缩编码。
- 缺点：使用者较少、跨防火墙访问时，不安全、不具有可读性，调试代码时相对困难、不能与其他传输层协议共同使用（例如HTTP）、无法支持向持久层直接读写数据，即不适合做数据持久化序列化协议。

适用场景：分布式系统的RPC解决方案

### Avro
Hadoop的一个子项目，解决了JSON的冗长和没有IDL的问题。

- 优点：支持丰富的数据类型、简单的动态语言结合功能、具有自我描述属性、提高了数据解析速度、快速可压缩的二进制数据形式、可以实现远程过程调用RPC、支持跨编程语言实现。
- 缺点：对于习惯于静态类型语言的用户不直观。

适用场景：在Hadoop中做Hive、Pig和MapReduce的持久化数据格式。

### Protobuf
将数据结构以.proto文件进行描述，通过代码生成工具可以生成对应数据结构的POJO对象和Protobuf相关的方法和属性
- 优点：序列化后码流小，性能高、结构化数据存储格式（XML JSON等）、通过标识字段的顺序，可以实现协议的前向兼容、结构化的文档更容易管理和维护。
- 缺点：需要依赖于工具生成代码、支持的语言相对较少，官方只支持Java 、C++ 、python,但是可以扩展。

适用场景：对性能要求高的RPC调用、具有良好的跨防火墙的访问属性、适合应用层对象的持久化






## json 协议


协议结构包括要素
- 对象（Object）：由一对大括号{}包围，内部是零个或多个键值对，每个键值对由冒号:分隔，键（key）是一个字符串，值（value）可以是字符串、数字、布尔值、对象、数组或null。
- 数组（Array）：由一对方括号[]包围，内部是零个或多个值，值可以是字符串、数字、布尔值、对象、数组或null，多个值之间用逗号,分隔。
- 字符串（String）：由双引号""包围的Unicode字符序列，可以包含任意字符，使用转义字符\来表示特殊字符。
- 数字（Number）：整数或浮点数。
- 布尔值（Boolean）：true或false。
- null：表示空值。



JSON语法规则
- 数据由键值对组成，键和值之间使用冒号（:）分隔。
- 键必须是字符串，使用双引号（"）括起来。
- 值可以是字符串、数字、布尔值、数组、对象或null。
- 多个键值对之间使用逗号（,）分隔。
- 对象使用花括号（{}）表示，键值对之间没有顺序。
- 数组使用方括号（[]）表示，值之间使用逗号分隔。

```json

{
  "name": "John",
  "age": 30,
  "isStudent": true,
  "address": {
    "street": "123 Main St",
    "city": "New York"
  },
  "hobbies": ["reading", "music", "sports"],
  "scores": [98, 85, 92, 76],
  "isMarried": null
}

```



## 参考


