// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: 04_jsonpb/proto/member.proto

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

type MemberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *MemberRequest) Reset() {
	*x = MemberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file__04_jsonpb_proto_member_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberRequest) ProtoMessage() {}

func (x *MemberRequest) ProtoReflect() protoreflect.Message {
	mi := &file__04_jsonpb_proto_member_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberRequest.ProtoReflect.Descriptor instead.
func (*MemberRequest) Descriptor() ([]byte, []int) {
	return file__04_jsonpb_proto_member_proto_rawDescGZIP(), []int{0}
}

func (x *MemberRequest) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type MemberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Phone string  `protobuf:"bytes,2,opt,name=Phone,proto3" json:"Phone,omitempty"`
	Age   int32   `protobuf:"varint,3,opt,name=Age,proto3" json:"Age,omitempty"`
	Data  *Detail `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *MemberResponse) Reset() {
	*x = MemberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file__04_jsonpb_proto_member_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemberResponse) ProtoMessage() {}

func (x *MemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file__04_jsonpb_proto_member_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemberResponse.ProtoReflect.Descriptor instead.
func (*MemberResponse) Descriptor() ([]byte, []int) {
	return file__04_jsonpb_proto_member_proto_rawDescGZIP(), []int{1}
}

func (x *MemberResponse) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MemberResponse) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *MemberResponse) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *MemberResponse) GetData() *Detail {
	if x != nil {
		return x.Data
	}
	return nil
}

type Detail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	School int32 `protobuf:"varint,1,opt,name=School,proto3" json:"School,omitempty"`
}

func (x *Detail) Reset() {
	*x = Detail{}
	if protoimpl.UnsafeEnabled {
		mi := &file__04_jsonpb_proto_member_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Detail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Detail) ProtoMessage() {}

func (x *Detail) ProtoReflect() protoreflect.Message {
	mi := &file__04_jsonpb_proto_member_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Detail.ProtoReflect.Descriptor instead.
func (*Detail) Descriptor() ([]byte, []int) {
	return file__04_jsonpb_proto_member_proto_rawDescGZIP(), []int{2}
}

func (x *Detail) GetSchool() int32 {
	if x != nil {
		return x.School
	}
	return 0
}

var File__04_jsonpb_proto_member_proto protoreflect.FileDescriptor

var file__04_jsonpb_proto_member_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x30, 0x34, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1f,
	0x0a, 0x0d, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x22,
	0x65, 0x0a, 0x0e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x50, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x41, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x41, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x20, 0x0a, 0x06, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x12, 0x16, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x53, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x32, 0x36, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x2c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x0e, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x0f, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x17, 0x5a, 0x15, 0x30, 0x34, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x70, 0x62, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file__04_jsonpb_proto_member_proto_rawDescOnce sync.Once
	file__04_jsonpb_proto_member_proto_rawDescData = file__04_jsonpb_proto_member_proto_rawDesc
)

func file__04_jsonpb_proto_member_proto_rawDescGZIP() []byte {
	file__04_jsonpb_proto_member_proto_rawDescOnce.Do(func() {
		file__04_jsonpb_proto_member_proto_rawDescData = protoimpl.X.CompressGZIP(file__04_jsonpb_proto_member_proto_rawDescData)
	})
	return file__04_jsonpb_proto_member_proto_rawDescData
}

var file__04_jsonpb_proto_member_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file__04_jsonpb_proto_member_proto_goTypes = []interface{}{
	(*MemberRequest)(nil),  // 0: MemberRequest
	(*MemberResponse)(nil), // 1: MemberResponse
	(*Detail)(nil),         // 2: Detail
}
var file__04_jsonpb_proto_member_proto_depIdxs = []int32{
	2, // 0: MemberResponse.data:type_name -> Detail
	0, // 1: Member.GetMember:input_type -> MemberRequest
	1, // 2: Member.GetMember:output_type -> MemberResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file__04_jsonpb_proto_member_proto_init() }
func file__04_jsonpb_proto_member_proto_init() {
	if File__04_jsonpb_proto_member_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file__04_jsonpb_proto_member_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberRequest); i {
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
		file__04_jsonpb_proto_member_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemberResponse); i {
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
		file__04_jsonpb_proto_member_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Detail); i {
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
			RawDescriptor: file__04_jsonpb_proto_member_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file__04_jsonpb_proto_member_proto_goTypes,
		DependencyIndexes: file__04_jsonpb_proto_member_proto_depIdxs,
		MessageInfos:      file__04_jsonpb_proto_member_proto_msgTypes,
	}.Build()
	File__04_jsonpb_proto_member_proto = out.File
	file__04_jsonpb_proto_member_proto_rawDesc = nil
	file__04_jsonpb_proto_member_proto_goTypes = nil
	file__04_jsonpb_proto_member_proto_depIdxs = nil
}
