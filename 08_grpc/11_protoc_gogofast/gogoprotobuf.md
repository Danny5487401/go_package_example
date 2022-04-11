# gogo/protobuf
生成 Golang 代码的插件 protoc-gen-go」，这个插件其实是 golang 官方提供的 一个 Protobuf api 实现。

而gogo/protobuf是基于 golang/protobuf 的一个增强版实现。

gogo 库基于官方库开发，增加了很多的功能，包括：

- 快速的序列化和反序列化
- 更规范的 Go 数据结构
- goprotobuf 兼容
- 可选择的产生一些辅助方法，减少使用中的代码输入
- 可以选择产生测试代码和 benchmark 代码
- 其它序列化格式

目前很多知名的项目都在使用该库，如 etcd、k8s、tidb、docker swarmkit 等

## 使用
gogo 库目前有三种生成代码的方式
1. gofast: 速度优先，但此方式不支持其它 gogoprotobuf 的扩展选项
```shell script
$ go get github.com/gogo/protobuf/protoc-gen-gofast
$ protoc --gofast_out=. myproto.proto
```
2. gogofast、gogofaster、gogoslick: 更快的速度、会生成更多的代码。
   - gogofast类似gofast，但是会引入 gogoprotobuf 库。
   - gogofaster类似gogofast，但是不会产生XXX_unrecognized类的指针字段，可以减少垃圾回收时间。
   - gogoslick类似gogofaster，但是会增加一些额外的string、gostring和equal method等。
```shell
$ go get github.com/gogo/protobuf/proto
$ go get github.com/gogo/protobuf/{binary} //protoc-gen-gogofast、protoc-gen-gogofaster 、protoc-gen-gogoslick 
$ go get github.com/gogo/protobuf/gogoproto
$ protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --{binary}_out=. myproto.proto // 这里的{binary}不包含「protoc-gen」前缀
```

3. protoc-gen-gogo: 最快的速度，最多的可定制化
   - 可以通过扩展选项高度定制序列化。
```shell
$ go get github.com/gogo/protobuf/proto
$ go get github.com/gogo/protobuf/jsonpb
$ go get github.com/gogo/protobuf/protoc-gen-gogo
$ go get github.com/gogo/protobuf/gogoproto
```
