// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: 03_stream_grpc/proto/stream.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//stream请求结构
type StreamReqData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StreamReqData) Reset() {
	*x = StreamReqData{}
	if protoimpl.UnsafeEnabled {
		mi := &file__03_stream_grpc_proto_stream_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamReqData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamReqData) ProtoMessage() {}

func (x *StreamReqData) ProtoReflect() protoreflect.Message {
	mi := &file__03_stream_grpc_proto_stream_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamReqData.ProtoReflect.Descriptor instead.
func (*StreamReqData) Descriptor() ([]byte, []int) {
	return file__03_stream_grpc_proto_stream_proto_rawDescGZIP(), []int{0}
}

func (x *StreamReqData) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

//stream返回结构
type StreamResData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *StreamResData) Reset() {
	*x = StreamResData{}
	if protoimpl.UnsafeEnabled {
		mi := &file__03_stream_grpc_proto_stream_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StreamResData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamResData) ProtoMessage() {}

func (x *StreamResData) ProtoReflect() protoreflect.Message {
	mi := &file__03_stream_grpc_proto_stream_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamResData.ProtoReflect.Descriptor instead.
func (*StreamResData) Descriptor() ([]byte, []int) {
	return file__03_stream_grpc_proto_stream_proto_rawDescGZIP(), []int{1}
}

func (x *StreamResData) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File__03_stream_grpc_proto_stream_proto protoreflect.FileDescriptor

var file__03_stream_grpc_proto_stream_proto_rawDesc = []byte{
	0x0a, 0x21, 0x30, 0x33, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x23, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0x9e, 0x01,
	0x0a, 0x07, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x09, 0x47, 0x65, 0x74,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52,
	0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x30, 0x01, 0x12, 0x2f, 0x0a, 0x09, 0x50, 0x75,
	0x74, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x28, 0x01, 0x12, 0x31, 0x0a, 0x09, 0x41,
	0x6c, 0x6c, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x12, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x1a, 0x0e, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x1c,
	0x5a, 0x1a, 0x30, 0x33, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file__03_stream_grpc_proto_stream_proto_rawDescOnce sync.Once
	file__03_stream_grpc_proto_stream_proto_rawDescData = file__03_stream_grpc_proto_stream_proto_rawDesc
)

func file__03_stream_grpc_proto_stream_proto_rawDescGZIP() []byte {
	file__03_stream_grpc_proto_stream_proto_rawDescOnce.Do(func() {
		file__03_stream_grpc_proto_stream_proto_rawDescData = protoimpl.X.CompressGZIP(file__03_stream_grpc_proto_stream_proto_rawDescData)
	})
	return file__03_stream_grpc_proto_stream_proto_rawDescData
}

var file__03_stream_grpc_proto_stream_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file__03_stream_grpc_proto_stream_proto_goTypes = []interface{}{
	(*StreamReqData)(nil), // 0: StreamReqData
	(*StreamResData)(nil), // 1: StreamResData
}
var file__03_stream_grpc_proto_stream_proto_depIdxs = []int32{
	0, // 0: Greeter.GetStream:input_type -> StreamReqData
	0, // 1: Greeter.PutStream:input_type -> StreamReqData
	0, // 2: Greeter.AllStream:input_type -> StreamReqData
	1, // 3: Greeter.GetStream:output_type -> StreamResData
	1, // 4: Greeter.PutStream:output_type -> StreamResData
	1, // 5: Greeter.AllStream:output_type -> StreamResData
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file__03_stream_grpc_proto_stream_proto_init() }
func file__03_stream_grpc_proto_stream_proto_init() {
	if File__03_stream_grpc_proto_stream_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file__03_stream_grpc_proto_stream_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamReqData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file__03_stream_grpc_proto_stream_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StreamResData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file__03_stream_grpc_proto_stream_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file__03_stream_grpc_proto_stream_proto_goTypes,
		DependencyIndexes: file__03_stream_grpc_proto_stream_proto_depIdxs,
		MessageInfos:      file__03_stream_grpc_proto_stream_proto_msgTypes,
	}.Build()
	File__03_stream_grpc_proto_stream_proto = out.File
	file__03_stream_grpc_proto_stream_proto_rawDesc = nil
	file__03_stream_grpc_proto_stream_proto_goTypes = nil
	file__03_stream_grpc_proto_stream_proto_depIdxs = nil
}
