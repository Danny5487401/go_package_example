package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	proto2 "github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"os"

	"github.com/Danny5487401/go_package_example/08_grpc/04_jsonpb/proto"
)

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
	// 方式一:打印结果
	fmt.Println("优化前直接给前端", r)
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
	fmt.Println("优化后发送给前端:", string(jsonCnt))
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
