package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	proto2 "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"os"

	"go_grpc_example/08_grpc/04_jsonpb/proto"
)

/*
结构体
type MemberResponse struct {
	Id    int32  `json "id"`
	Phone string `json "phone"`
	Age   int8   `json "age"`
}
返回结果	Id:12  Phone:"15112810201"
问题
	因为server未返回age字段，所以没有age。在某些情况下对前端也是不太友好的，尤其是APP客户端，更需要明确的json响应字段结构
解决方式
1. 直接修改经过protoc生成的member.pb.go文件代码，删除掉不希望被忽略的字段tag标签中的omitempty即可，
	但是*.pb.go一般我们不建议去修改它，而且我们会经常去调整grpc微服务协议中的方法或者字段内容，这样每次protoc之后，
	都需要我们去修改，这显然是不太现实的，因此就有了第二种办法；
2. 通过grpc官方库中的jsonpb来实现,官方在它的设定中有一个结构体用来实现protoc buffer转换为JSON结构，并可以根据字段来配置转换的要求，
	结构体如下
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
*/

var jsonpbMarshaler *jsonpb.Marshaler

func main() {
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewMemberClient(conn)

	r, err := c.GetMember(context.Background(), &proto.MemberRequest{
		Id: 1,
	})
	if err != nil {
		panic(err)
	}
	//方式一:打印结果
	//fmt.Println(r)
	// 方式二：Marshaler使用
	MarshalData(r)

	UnmarshalPbData()

}

func MarshalData(r proto2.Message) {

	var (
		_buffer bytes.Buffer
	)
	// 记住初始化
	jsonpbMarshaler = &jsonpb.Marshaler{
		EnumsAsInts:  true, // 将枚举值设定为整数，而不是字符串类型.
		EmitDefaults: true, // 将字段值为空的渲染到JSON结构中
		OrigName:     true, // 使用原生的proto协议中的字段
	}
	//调用此方法实现转换
	jsonpbMarshaler.Marshal(&_buffer, r)
	jsonCnt := _buffer.Bytes()
	//发送给前端: {"Id":12,"Phone":"15112810201","Age":0,"data":null}  注意当data字段为空指针，会返回null给前端
	fmt.Println("发送给前端:", string(jsonCnt))
	// 会把空数值打印处理，方便返回给前端
}

func UnmarshalPbData() {
	// Unmarshaler使用:转换成pb对象
	// school没有这字段，Age为0
	mockData := `{"Id":12,"Phone":"15112810201","Age":0,"school":"Beijing"}`
	bufferObj := new(bytes.Buffer)
	bufferObj.WriteString(mockData)

	memberRsp := proto.MemberResponse{}
	unmarshaler := jsonpb.Unmarshaler{
		AllowUnknownFields: true, //允许忽略未知字段，如school字段不存在
	}
	if err := unmarshaler.Unmarshal(bufferObj, &memberRsp); err != nil {
		fmt.Println("jsonpb UnmarshalString fail: ", err)
		os.Exit(0)
	}
	fmt.Printf("member info pb反序列化: %+v", memberRsp.String())
}
