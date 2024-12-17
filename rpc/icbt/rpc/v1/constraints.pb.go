// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: icbt/rpc/v1/constraints.proto

package rpcv1

import (
	validate "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var file_icbt_rpc_v1_constraints_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*validate.StringRules)(nil),
		ExtensionType: (*bool)(nil),
		Field:         70000001,
		Name:          "icbt.rpc.v1.refid",
		Tag:           "varint,70000001,opt,name=refid",
		Filename:      "icbt/rpc/v1/constraints.proto",
	},
}

// Extension fields to validate.StringRules.
var (
	// optional bool refid = 70000001;
	E_Refid = &file_icbt_rpc_v1_constraints_proto_extTypes[0]
)

var File_icbt_rpc_v1_constraints_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_constraints_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f,
	0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75,
	0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3a, 0x83, 0x01, 0x0a, 0x05, 0x72, 0x65,
	0x66, 0x69, 0x64, 0x12, 0x19, 0x2e, 0x62, 0x75, 0x66, 0x2e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x81,
	0xbb, 0xb0, 0x21, 0x20, 0x01, 0x28, 0x08, 0x42, 0x4f, 0xc2, 0x48, 0x4c, 0x0a, 0x4a, 0x0a, 0x0c,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x72, 0x65, 0x66, 0x69, 0x64, 0x12, 0x17, 0x6d, 0x75,
	0x73, 0x74, 0x20, 0x62, 0x65, 0x20, 0x69, 0x6e, 0x20, 0x72, 0x65, 0x66, 0x69, 0x64, 0x20, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x1a, 0x21, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x73, 0x28, 0x27, 0x5e, 0x5b, 0x30, 0x2d, 0x39, 0x61, 0x2d, 0x7a, 0x41, 0x2d, 0x5a,
	0x5d, 0x7b, 0x32, 0x36, 0x7d, 0x24, 0x27, 0x29, 0x52, 0x05, 0x72, 0x65, 0x66, 0x69, 0x64, 0x42,
	0xab, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x42, 0x10, 0x43, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f, 0x69, 0x63,
	0x61, 0x6e, 0x62, 0x72, 0x69, 0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76,
	0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x2e, 0x52,
	0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d,
	0x49, 0x63, 0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x08, 0x65,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
}

var file_icbt_rpc_v1_constraints_proto_goTypes = []any{
	(*validate.StringRules)(nil), // 0: buf.validate.StringRules
}
var file_icbt_rpc_v1_constraints_proto_depIdxs = []int32{
	0, // 0: icbt.rpc.v1.refid:extendee -> buf.validate.StringRules
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	0, // [0:1] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_constraints_proto_init() }
func file_icbt_rpc_v1_constraints_proto_init() {
	if File_icbt_rpc_v1_constraints_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_constraints_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 1,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_constraints_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_constraints_proto_depIdxs,
		ExtensionInfos:    file_icbt_rpc_v1_constraints_proto_extTypes,
	}.Build()
	File_icbt_rpc_v1_constraints_proto = out.File
	file_icbt_rpc_v1_constraints_proto_rawDesc = nil
	file_icbt_rpc_v1_constraints_proto_goTypes = nil
	file_icbt_rpc_v1_constraints_proto_depIdxs = nil
}
