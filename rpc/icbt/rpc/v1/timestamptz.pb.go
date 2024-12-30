// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: icbt/rpc/v1/timestamptz.proto

package rpcv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/gofeaturespb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TimestampTZ struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Ts *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=ts" json:"ts,omitempty"`
	xxx_hidden_Tz string                 `protobuf:"bytes,2,opt,name=tz" json:"tz,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TimestampTZ) Reset() {
	*x = TimestampTZ{}
	mi := &file_icbt_rpc_v1_timestamptz_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TimestampTZ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimestampTZ) ProtoMessage() {}

func (x *TimestampTZ) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_timestamptz_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *TimestampTZ) GetTs() *timestamppb.Timestamp {
	if x != nil {
		return x.xxx_hidden_Ts
	}
	return nil
}

func (x *TimestampTZ) GetTz() string {
	if x != nil {
		return x.xxx_hidden_Tz
	}
	return ""
}

func (x *TimestampTZ) SetTs(v *timestamppb.Timestamp) {
	x.xxx_hidden_Ts = v
}

func (x *TimestampTZ) SetTz(v string) {
	x.xxx_hidden_Tz = v
}

func (x *TimestampTZ) HasTs() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Ts != nil
}

func (x *TimestampTZ) ClearTs() {
	x.xxx_hidden_Ts = nil
}

type TimestampTZ_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	// timestamp
	Ts *timestamppb.Timestamp
	// timezone
	Tz string
}

func (b0 TimestampTZ_builder) Build() *TimestampTZ {
	m0 := &TimestampTZ{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Ts = b.Ts
	x.xxx_hidden_Tz = b.Tz
	return m0
}

var File_icbt_rpc_v1_timestamptz_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_timestamptz_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x5f, 0x66, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a,
	0x0b, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x54, 0x5a, 0x12, 0x2a, 0x0a, 0x02,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x02, 0x74, 0x7a, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0d, 0xba, 0x48, 0x0a, 0x72, 0x08, 0x10, 0x03, 0x92, 0x02, 0x03,
	0x55, 0x54, 0x43, 0x52, 0x02, 0x74, 0x7a, 0x42, 0xb5, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e,
	0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x10, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70,
	0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f, 0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69, 0x6e, 0x67, 0x74,
	0x68, 0x61, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa,
	0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b,
	0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x49, 0x63,
	0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70,
	0x63, 0x3a, 0x3a, 0x56, 0x31, 0x92, 0x03, 0x07, 0xd2, 0x3e, 0x02, 0x10, 0x03, 0x08, 0x02, 0x62,
	0x08, 0x65, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
}

var file_icbt_rpc_v1_timestamptz_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_icbt_rpc_v1_timestamptz_proto_goTypes = []any{
	(*TimestampTZ)(nil),           // 0: icbt.rpc.v1.TimestampTZ
	(*timestamppb.Timestamp)(nil), // 1: google.protobuf.Timestamp
}
var file_icbt_rpc_v1_timestamptz_proto_depIdxs = []int32{
	1, // 0: icbt.rpc.v1.TimestampTZ.ts:type_name -> google.protobuf.Timestamp
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_timestamptz_proto_init() }
func file_icbt_rpc_v1_timestamptz_proto_init() {
	if File_icbt_rpc_v1_timestamptz_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_timestamptz_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_timestamptz_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_timestamptz_proto_depIdxs,
		MessageInfos:      file_icbt_rpc_v1_timestamptz_proto_msgTypes,
	}.Build()
	File_icbt_rpc_v1_timestamptz_proto = out.File
	file_icbt_rpc_v1_timestamptz_proto_rawDesc = nil
	file_icbt_rpc_v1_timestamptz_proto_goTypes = nil
	file_icbt_rpc_v1_timestamptz_proto_depIdxs = nil
}
