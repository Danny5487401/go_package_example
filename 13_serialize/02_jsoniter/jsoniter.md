<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [第三方包 json-iterator/go](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%8C%85-json-iteratorgo)
  - [优化点](#%E4%BC%98%E5%8C%96%E7%82%B9)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



# 第三方包 json-iterator/go

是一个非常优秀的go json解析库，完全兼容官方的json解析库。

从性能上，jsoniter 能够比众多大神联合开发的官方库性能还快的主要原因，一个是尽量减少不必要的内存复制，另一个是减少 reflect 的使用——同一类型的对象，jsoniter 只调用 reflect 解析一次之后即缓存下来。
不过随着 go 版本的迭代，原生 json 库的性能也越来越高，jsonter 的性能优势也越来越窄.

json-iterator还提供了很多其他方便的功能，如开放的序列化/反序列化配置、Extension、FieldEncoder/FieldDecoder、懒解析Any对象等等增强功能，应对不同使用场景下的json编码和解析，满足各种复杂的需求



## 优化点


1，单次扫描：所有解析都是在字节数组流中直接在一次传递中完成的。readInt或readString一次完成，并没有做json的token切分，直接读取字符，转换成目标类型，readFloat或readDouble都以这种方式实现。避免重复扫描的同时，也最大限度避免了内存的申请和释放。

2，它不解析令牌，然后分支。相反，它是先将目标需要绑定的golang对象类型和对应的解析器解析出来，并缓存。然后遍历json串的时候，对取出来的每个key，结合json当前上下文，去map里取对应的解析器，去解析并绑定值。

3，对于不需要解析的字段，会跳过它所有的嵌套对象，因为匹配不到解析器，避免不必要的解析。跳过整个对象时，我们不关心嵌套字段名称

4，绑定到对象不使用反射api。而是取出原始指针interface{}，然后转换为正确的指针类型以设置值。例如：*((*int)(ptr)) = iter.ReadInt()

5，尽量避免map的分配和寻址，对于小于等于10个字段的结构体，通过计算key的hash的方式，分配每个字段的结构体和对应的解析函数，这样解析到key的时候，直接通过hash值的匹配，避免了字符串匹配和map的分配，以及匹配