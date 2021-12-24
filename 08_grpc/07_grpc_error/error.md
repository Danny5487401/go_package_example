# grpc错误

## errors包
```go
package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func wrapNewPointerError() error {
	// 	Go1.13版本为fmt.Errorf函数新加了一个%w占位符用来生成一个可以包裹Error的Wrapping Error。
	return fmt.Errorf("wrap err0:%w", fmt.Errorf("i am a error0"))
}

func wrapConstantPointerError() error {
	return fmt.Errorf("wrap err1:%w", constantErr)
}

var constantErr = fmt.Errorf("i am a error1 ")

func main() {
	fmt.Println("第一个结果", errors.Is(wrapNewPointerError(), fmt.Errorf("i am a error0"))) // false
	fmt.Println("第二个结果", errors.Is(wrapConstantPointerError(), constantErr))            //true
}

```

## gRPC网络传输的Error
![](.error_images/error_transfer_in_grpc.png)     
我们客户端在获取到gRPC的error的时候，是否可以使用上文说的官方errors.Is进行判断呢。
如果我们直接使用该方法，通过判断error地址是否相等，是无法做到的。原因是因为我们在使用gRPC的时候，在远程调用过程中，客户端获取的服务端返回的error，在tcp传递的时候实际上是一串文本。
客户端拿到这个文本，是要将其反序列化转换为error，在这个反序列化的过程中，其实是new了一个新的error地址，这样就无法判断error地址是否相等。

为了更好的解释gRPC网络传输的error，以下描述了整个error的处理流程
- 客户端通过invoker方法将请求发送到服务端。

- 服务端通过processUnaryRPC方法，获取到用户代码的error信息。

- 服务端通过status.FromError方法，将error转化为status.Status。

- 服务端通过WriteStatus方法将status.Status里的数据，写入到grpc-status、grpc-message、grpc-status-details-bin的header头里。

- 客户端通过网络获取到这些header头，使用strconv.ParseInt解析到grpc-status信息、decodeGrpcMessage解析到grpc-message信息、decodeGRPCStatusDetails解析为grpc-status-details-bin信息。

- 客户端通过a.Status().Err()获取到用户代码的错误。

为了方便理解，我们抓个包，看下error具体的报文情况。
![](.error_images/error_packets_in_grpc.png)


## grpc status包
对外暴露的方法，首先看返回err
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/status/status.go
package status

import (
	"context"
	"fmt"

	spb "google.golang.org/genproto/googleapis/rpc/status"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
)


// Status references google.golang.org/grpc/internal/status. It represents an
// RPC status code, message, and details.  It is immutable and should be
// created with New, Newf, or FromProto.
// https://godoc.org/google.golang.org/grpc/internal/status
// 不可改变的
type Status = status.Status

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *Status {
	return status.New(c, msg)
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c codes.Code, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

// Error returns an error representing c and msg.  If c is OK, returns nil.
func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}
```
内部包
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/grpc@v1.32.0/internal/status/status.go
type Status struct {
	s *spb.Status
}

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *Status {
	return &Status{s: &spb.Status{Code: int32(c), Message: msg}}
}
```
内部定义的pb包:status结构体
```go
// /Users/xiaxin/go/pkg/mod/google.golang.org/genproto@v0.0.0-20210729151513-df9385d47c1b/googleapis/rpc/status/status.pb.go
type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// A developer-facing error message, which should be in English. Any
	// user-facing error message should be localized and sent in the
	// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// A list of messages that carry the error details.  There is a common set of
	// message types for APIs to use.
	Details []*anypb.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
}
```
