# Protobuf
Protocol buffers 是一种语言无关、平台无关的可扩展机制或者说是数据交换格式，用于序列化结构化数据。与 XML、JSON 相比，Protocol buffers 序列化后的码流更小、速度更快、操作更简单。
## v2 和 v3 主要区别
* 删除原始值字段的字段存在逻辑
* 删除 required 字段
* 删除 optional 字段，默认就是
* 删除 default 字段
* 删除扩展特性，新增 Any 类型来替代它
* 删除 unknown 字段的支持
* 新增 JSON Mapping
* 新增 Map 类型的支持
* 修复 enum 的 unknown 类型
* repeated 默认使用 packed 编码
* 引入了新的语言实现（C＃，JavaScript，Ruby，Objective-C）

## 核心编码原理(包括 Varint 编码、ZigZag编码及 protobuf 特有的 Message Structure 编码结构等)
### 1. Varint编码:protobuf 编码主要依赖于 Varint 编码。
原理：
Varint 是一种紧凑的表示数字的方法。它用一个或多个字节来表示一个数字，值越小的数字使用越少的字节数。这能减少用来表示数字的字节数。

Varint 中的每个字节（最后一个字节除外）都设置了最高有效位（msb），这一位表示还会有更多字节出现。每个字节的低 7 位用于以 7 位组的形式存储数字的二进制补码表示，最低有效组首位

编码方式

1. 转换为二进制表示
2. 每个字节保留后7位，去掉最高位
3. 因为 protobuf 使用的是小端序，所以要将大端序转为小端序
   - 每次从低向高取7位再加上最高有效位(最后一个字节高位补0，其余各字节高位补1)组成编码后的数据。 
4. 最后转成10进制

Note: 最高位为1代表后面7位仍然表示数字，否则为0，后面7位用原码补齐。

如果用不到 1 个字节，那么最高有效位设为 0 ，如下面这个例子，1 用一个字节就可以表示，所以 msb 为 0.
```
0000 0001
```
![](.proto_images/transfer_123456_to_varint.png)
1. 123456用二进制表示为1 11100010 01000000，
2. 每次从低向高取7位再加上最高有效位变成1100 0000 11000100 00000111
3. 所以经过varint编码后123456占用三个字节分别为192 196 7。

解码的过程就是将字节依次取出，去掉最高有效位，因为是小端排序所以先解码的字节要放在低位，
之后解码出来的二进制位继续放在之前已经解码出来的二进制的高位最后转换为10进制数完成varint编码的解码过程。

缺点： 负数需要10个字节显示（因为计算机定义负数的符号位为数字的最高位）。
具体是先将负数是转成了long类型，再进行varint编码，这就是占用10个字节的原因了。

protobuf 采取的解决方式：使用 sint32/sint64 类型表示负数，通过先采用 Zigzag 编码，将正数、负数和0都映射到无符号数，最后再采用varints编码。


编码
```go
const maxVarintBytes = 10 // maximum length of a varint

// 返回Varint类型编码后的字节流
func EncodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	// 下面的编码规则需要详细理解:
	// 1.每个字节的最高位是保留位, 如果是1说明后面的字节还是属于当前数据的,如果是0,那么这是当前数据的最后一个字节数据
	//  看下面代码,因为一个字节最高位是保留位,那么这个字节中只有下面7bits可以保存数据
	//  所以,如果x>127,那么说明这个数据还需大于一个字节保存,所以当前字节最高位是1,看下面的buf[n] = 0x80 | ...
	//  0x80说明将这个字节最高位置为1, 后面的x&0x7F是取得x的低7位数据, 那么0x80 | uint8(x&0x7F)整体的意思就是
	//  这个字节最高位是1表示这不是最后一个字节,后面7为是正式数据! 注意操作下一个字节之前需要将x>>=7
	// 2.看如果x<=127那么说明x现在使用7bits可以表示了,那么最高位没有必要是1,直接是0就ok!所以最后直接是buf[n] = uint8(x)
	//
	// 如果数据大于一个字节(127是一个字节最大数据), 那么继续, 即: 需要在最高位加上1
	for n = 0; x > 127; n++ {
	    // x&0x7F表示取出下7bit数据, 0x80表示在最高位加上1
		buf[n] = 0x80 | uint8(x&0x7F)
		// 右移7位, 继续后面的数据处理
		x >>= 7
	}
	// 最后一个字节数据
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}
```

解码
```go
func DecodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
    // 下面这个分成三步走:
		// 1: b & 0x7F 获取下7bits有效数据
		// 2: (b & 0x7F) << shift 由于是小端序, 所以每次处理一个Byte数据, 都需要向高位移动7bits
		// 3: 将数据x和当前的这个字节数据 | 在一起
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}
```

## 使用
### 基本定义
```protobuf
option go_package = "{out_path};out_go_package"; // 前一个参数用于指定生成文件的位置，后一个参数指定生成的 .go 文件的 package
package import; // 表示当前 protobuf 文件属于 import包，这个package不是 Go 语言中的那个package
```

### 1. 引入其他proto文件
```shell
pwd 
# /Users/xiaxin/Desktop/go_grpc_example
cd 08_grpc
```
目录结构   
![](.proto_images/dir_proto.png)

Note: Goland proto插件展示问题，需要手动添加路径，不添加也不影响(这是插件问题)  
![](.proto_images/goland_proto_display_problem.png)   
解决方式:解决后   
![](.proto_images/goland_protobuf_plugin.png)
![](.proto_images/goland_protobuf_display_fix.png)

### 生成protobuf
```makefile
.PHONY: proto
proto:
	protoc --proto_path=. --go_out=. ./proto/dir_import/*.proto
```
解释：
1) --proto_path =.  指定在当前目录(go_grpc_example/08_grpc)寻找 import 的文件
```protobuf
// 08_grpc/proto/dir_import/computer.proto
import "proto/dir_import/component.proto";
```
所以最终会去找 go_grpc_example/08_grpc/proto/dir_import/component.proto

2）–go_out=.
指定将生成文件放在当前目录( go_grpc_example/08_grpc)，同时因为 proto 文件中也指定了目录为protobuf/import,具体如下：
```protobuf
option go_package = "proto/dir_import;proto";
```
所以最终生成目录为--go_out+go_package= go_grpc_example/08_grpc/proto/dir_import

Note:  可以通过参数 --go_opt=paths=source_relative 来指定使用绝对路径，从而忽略掉 proto 文件中的 go_package 路径，直接生成在 –go_out 指定的路径

3）./protobuf/import/*.proto 

指定编译 import 目录下的所有 proto 文件，由于有文件的引入所以需要一起编译才能生效。

Note: 当然也可以一个一个编译，只要把相关文件都编译好即可。


## protobuf优化
![](.proto_images/proto_optimize.png)

### wiretype     
![](.proto_images/wire_type.png)

## 工具
### protoc
![](.proto_images/protoc_process.png)   
protoc是protobuf文件（.proto）的编译器，可以借助这个工具把 .proto 文件转译成各种编程语言对应的源码，包含数据类型定义、调用接口等。

通过查看protoc的源码（参见github库）可以知道，protoc在设计上把protobuf和不同的语言解耦了，底层用c++来实现protobuf结构的存储，然后通过插件的形式来生成不同语言的源码。可以把protoc的编译过程分成简单的两个步骤

1. 解析.proto文件，转译成protobuf的原生数据结构在内存中保存；    

2. 把protobuf相关的数据结构传递给相应语言的编译插件，由插件负责根据接收到的protobuf原生结构渲染输出特定语言的模板

Note:包含的插件有 csharp、java、js、objectivec、php、python、ruby等多种,不包含go.

### protoc-gen-go
![](.proto_images/protoc_gen_go_files.png)   
原生protoc并不包含Go版本的插件,protoc-gen-go是protobuf编译插件系列中的Go版本。
由于protoc-gen-go是Go写的，所以安装它变得很简单，只需要运行 go get -u github.com/golang/protobuf/protoc-gen-go

#### protoc-gen-go 源码目录分析
main包

- doc.go 主要是说明。
- link_grpc.go 显式引用protoc-gen-go/grpc包，触发grpc的init函数。
- main.go 代码不到50行，初始generator，并调用generator相应的方输出protobuf的Go语言文件。

- generator.go 包含了大部分由protobuf原生结构到Go语言文件的渲染方法，其中 func (g *Generator) P(str ...interface{}) 这个方法会把渲染输出到generator的output（generator匿名嵌套了bytes.Buffer，因此有Buffer的方法）。
name_test.go 测试，主要包含generator中名称相关方法的测试。

- grpc.go 与generator相似，但是包含了很多生成grpc相关方法的方法，比如渲染转译protobuf中定义的rpc方法（在generator中不包含，其默认不转译service的定义）
descriptor 包含protobuf的描述文件（.proto文件及其对应的Go编译文件），其中proto文件来自于proto库

- plugin 包含plugin的描述文件（.proto文件及其对应的Go编译文件），其中proto文件来自于proto库



### protoc-gen-god的替代版本:gogoprotobuf

在go中使用protobuf，有两个可选用的包goprotobuf（go官方出品）和gogoprotobuf。gogoprotobuf完全兼容google protobuf，
它生成的代码质量和编解码性能均比goprotobuf高一些。
主要是它在goprotobuf之上extend了一些option。这些option也是有级别区分的，有的option只能修饰field，有的可以修饰enum，有的可以修饰message，有的是修饰package（即对整个文件都有效)

gogoprotobuf有两个插件可以使用

protoc-gen-gogo：和protoc-gen-go生成的文件差不多，性能也几乎一样(稍微快一点点)
protoc-gen-gofast：生成的文件更复杂，性能也更高(快5-7倍)

```shell
#安装 the protoc-gen-gofast binary
go get github.com/gogo/protobuf/protoc-gen-gofast
#生成
protoc --gofast_out=. myproto.proto
```
## 生成方式
参考scripts脚本

## 生成的protobuf.pb.go源码分析
```go
//message接口
type Message = protoiface.MessageV1
type MessageV1 interface {
    Reset()
    String() string
    ProtoMessage()
}
```
proto编译成的Go结构体都是符合Message接口的，从Marshal可知Go结构体有3种序列化方式：
```go
func Marshal(pb Message) ([]byte, error) {
	if m, ok := pb.(newMarshaler); ok {
		siz := m.XXX_Size()
		b := make([]byte, 0, siz)
		return m.XXX_Marshal(b, false)
	}
	if m, ok := pb.(Marshaler); ok {
		// If the message can marshal itself, let it do it, for compatibility.
		// NOTE: This is not efficient.
		return m.Marshal()
	}
	// in case somehow we didn't generate the wrapper
	if pb == nil {
		return nil, ErrNil
	}
	var info InternalMessageInfo
	siz := info.Size(pb)
	b := make([]byte, 0, siz)
	return info.Marshal(b, pb, false)
}
//newMarshaler接口
type newMarshaler interface {
    XXX_Size() int
    XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
}
//Marshaler接口
type Marshaler interface {
    Marshal() ([]byte, error)
}
```

1. pb Message满足newMarshaler接口，则调用XXX_Marshal()进行序列化。   
2. pb满足Marshaler接口，则调用Marshal()进行序列化，这种方式适合某类型自定义序列化规则的情况。   
3. 否则，使用默认的序列化方式，创建一个Warpper，利用wrapper对pb进行序列化，后面会介绍方式1实际就是使用方式3。