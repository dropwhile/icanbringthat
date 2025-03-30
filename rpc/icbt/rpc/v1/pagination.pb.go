// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: icbt/rpc/v1/pagination.proto

package rpcv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/gofeaturespb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PaginationRequest struct {
	state             protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Limit  uint32                 `protobuf:"varint,1,opt,name=limit"`
	xxx_hidden_Offset uint32                 `protobuf:"varint,2,opt,name=offset"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *PaginationRequest) Reset() {
	*x = PaginationRequest{}
	mi := &file_icbt_rpc_v1_pagination_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaginationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationRequest) ProtoMessage() {}

func (x *PaginationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_pagination_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PaginationRequest) GetLimit() uint32 {
	if x != nil {
		return x.xxx_hidden_Limit
	}
	return 0
}

func (x *PaginationRequest) GetOffset() uint32 {
	if x != nil {
		return x.xxx_hidden_Offset
	}
	return 0
}

func (x *PaginationRequest) SetLimit(v uint32) {
	x.xxx_hidden_Limit = v
}

func (x *PaginationRequest) SetOffset(v uint32) {
	x.xxx_hidden_Offset = v
}

type PaginationRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Limit  uint32
	Offset uint32
}

func (b0 PaginationRequest_builder) Build() *PaginationRequest {
	m0 := &PaginationRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Limit = b.Limit
	x.xxx_hidden_Offset = b.Offset
	return m0
}

type PaginationResult struct {
	state             protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Limit  uint32                 `protobuf:"varint,1,opt,name=limit"`
	xxx_hidden_Offset uint32                 `protobuf:"varint,2,opt,name=offset"`
	xxx_hidden_Count  uint32                 `protobuf:"varint,3,opt,name=count"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *PaginationResult) Reset() {
	*x = PaginationResult{}
	mi := &file_icbt_rpc_v1_pagination_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaginationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationResult) ProtoMessage() {}

func (x *PaginationResult) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_pagination_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *PaginationResult) GetLimit() uint32 {
	if x != nil {
		return x.xxx_hidden_Limit
	}
	return 0
}

func (x *PaginationResult) GetOffset() uint32 {
	if x != nil {
		return x.xxx_hidden_Offset
	}
	return 0
}

func (x *PaginationResult) GetCount() uint32 {
	if x != nil {
		return x.xxx_hidden_Count
	}
	return 0
}

func (x *PaginationResult) SetLimit(v uint32) {
	x.xxx_hidden_Limit = v
}

func (x *PaginationResult) SetOffset(v uint32) {
	x.xxx_hidden_Offset = v
}

func (x *PaginationResult) SetCount(v uint32) {
	x.xxx_hidden_Count = v
}

type PaginationResult_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Limit  uint32
	Offset uint32
	Count  uint32
}

func (b0 PaginationResult_builder) Build() *PaginationResult {
	m0 := &PaginationResult{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Limit = b.Limit
	x.xxx_hidden_Offset = b.Offset
	x.xxx_hidden_Count = b.Count
	return m0
}

var File_icbt_rpc_v1_pagination_proto protoreflect.FileDescriptor

const file_icbt_rpc_v1_pagination_proto_rawDesc = "" +
	"\n" +
	"\x1cicbt/rpc/v1/pagination.proto\x12\vicbt.rpc.v1\x1a\x1bbuf/validate/validate.proto\x1a!google/protobuf/go_features.proto\"J\n" +
	"\x11PaginationRequest\x12\x1d\n" +
	"\x05limit\x18\x01 \x01(\rB\a\xbaH\x04*\x02 \x00R\x05limit\x12\x16\n" +
	"\x06offset\x18\x02 \x01(\rR\x06offset\"V\n" +
	"\x10PaginationResult\x12\x14\n" +
	"\x05limit\x18\x01 \x01(\rR\x05limit\x12\x16\n" +
	"\x06offset\x18\x02 \x01(\rR\x06offset\x12\x14\n" +
	"\x05count\x18\x03 \x01(\rR\x05countB\xb4\x01\n" +
	"\x0fcom.icbt.rpc.v1B\x0fPaginationProtoP\x01Z8github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1;rpcv1\xa2\x02\x03IRX\xaa\x02\vIcbt.Rpc.V1\xca\x02\vIcbt\\Rpc\\V1\xe2\x02\x17Icbt\\Rpc\\V1\\GPBMetadata\xea\x02\rIcbt::Rpc::V1\x92\x03\a\xd2>\x02\x10\x03\b\x02b\beditionsp\xe8\a"

var file_icbt_rpc_v1_pagination_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_icbt_rpc_v1_pagination_proto_goTypes = []any{
	(*PaginationRequest)(nil), // 0: icbt.rpc.v1.PaginationRequest
	(*PaginationResult)(nil),  // 1: icbt.rpc.v1.PaginationResult
}
var file_icbt_rpc_v1_pagination_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_pagination_proto_init() }
func file_icbt_rpc_v1_pagination_proto_init() {
	if File_icbt_rpc_v1_pagination_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_icbt_rpc_v1_pagination_proto_rawDesc), len(file_icbt_rpc_v1_pagination_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_pagination_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_pagination_proto_depIdxs,
		MessageInfos:      file_icbt_rpc_v1_pagination_proto_msgTypes,
	}.Build()
	File_icbt_rpc_v1_pagination_proto = out.File
	file_icbt_rpc_v1_pagination_proto_goTypes = nil
	file_icbt_rpc_v1_pagination_proto_depIdxs = nil
}
