// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: icbt/rpc/v1/service.proto

package rpcv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_icbt_rpc_v1_service_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_service_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x17, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x69, 0x63, 0x62, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x76, 0x31, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xdd, 0x0d, 0x0a, 0x0e, 0x49, 0x63, 0x62, 0x74, 0x52, 0x70,
	0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x0d, 0x45, 0x61, 0x72, 0x6d,
	0x61, 0x72, 0x6b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x21, 0x2e, 0x69, 0x63, 0x62, 0x74,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x69,
	0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61,
	0x72, 0x6b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x62, 0x0a, 0x11, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x25, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x69,
	0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61,
	0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x21, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x53, 0x0a, 0x0c, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x20, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x21, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x46, 0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1f,
	0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4d, 0x0a, 0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1e, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x47,
	0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x23, 0x2e, 0x69, 0x63, 0x62, 0x74,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x47, 0x65, 0x74,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24,
	0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x59, 0x0a, 0x0e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x22, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x62, 0x0a, 0x11, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x61, 0x72, 0x6d,
	0x61, 0x72, 0x6b, 0x73, 0x12, 0x25, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x61, 0x72, 0x6d,
	0x61, 0x72, 0x6b, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c,
	0x69, 0x73, 0x74, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x20, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x23, 0x2e, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4e, 0x0a, 0x0f, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x23, 0x2e, 0x69, 0x63, 0x62, 0x74,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x50, 0x0a, 0x0b, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x41, 0x64, 0x64, 0x12, 0x1f, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x41, 0x64, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0e, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x22, 0x2e, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x65, 0x0a, 0x12, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x26, 0x2e, 0x69,
	0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a,
	0x12, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x12, 0x26, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x5c, 0x0a, 0x16, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x12, 0x2a, 0x2e,
	0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41,
	0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x62, 0x0a, 0x11, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x25, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26, 0x2e,
	0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xa7, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f,
	0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69, 0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70,
	0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70,
	0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74,
	0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52,
	0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_icbt_rpc_v1_service_proto_goTypes = []any{
	(*EarmarkCreateRequest)(nil),          // 0: icbt.rpc.v1.EarmarkCreateRequest
	(*EarmarkGetDetailsRequest)(nil),      // 1: icbt.rpc.v1.EarmarkGetDetailsRequest
	(*EarmarkRemoveRequest)(nil),          // 2: icbt.rpc.v1.EarmarkRemoveRequest
	(*EarmarksListRequest)(nil),           // 3: icbt.rpc.v1.EarmarksListRequest
	(*EventCreateRequest)(nil),            // 4: icbt.rpc.v1.EventCreateRequest
	(*EventUpdateRequest)(nil),            // 5: icbt.rpc.v1.EventUpdateRequest
	(*EventDeleteRequest)(nil),            // 6: icbt.rpc.v1.EventDeleteRequest
	(*EventsListRequest)(nil),             // 7: icbt.rpc.v1.EventsListRequest
	(*EventGetDetailsRequest)(nil),        // 8: icbt.rpc.v1.EventGetDetailsRequest
	(*EventListItemsRequest)(nil),         // 9: icbt.rpc.v1.EventListItemsRequest
	(*EventListEarmarksRequest)(nil),      // 10: icbt.rpc.v1.EventListEarmarksRequest
	(*EventAddItemRequest)(nil),           // 11: icbt.rpc.v1.EventAddItemRequest
	(*EventUpdateItemRequest)(nil),        // 12: icbt.rpc.v1.EventUpdateItemRequest
	(*EventRemoveItemRequest)(nil),        // 13: icbt.rpc.v1.EventRemoveItemRequest
	(*FavoriteAddRequest)(nil),            // 14: icbt.rpc.v1.FavoriteAddRequest
	(*FavoriteRemoveRequest)(nil),         // 15: icbt.rpc.v1.FavoriteRemoveRequest
	(*FavoriteListEventsRequest)(nil),     // 16: icbt.rpc.v1.FavoriteListEventsRequest
	(*NotificationDeleteRequest)(nil),     // 17: icbt.rpc.v1.NotificationDeleteRequest
	(*NotificationsDeleteAllRequest)(nil), // 18: icbt.rpc.v1.NotificationsDeleteAllRequest
	(*NotificationsListRequest)(nil),      // 19: icbt.rpc.v1.NotificationsListRequest
	(*EarmarkCreateResponse)(nil),         // 20: icbt.rpc.v1.EarmarkCreateResponse
	(*EarmarkGetDetailsResponse)(nil),     // 21: icbt.rpc.v1.EarmarkGetDetailsResponse
	(*emptypb.Empty)(nil),                 // 22: google.protobuf.Empty
	(*EarmarksListResponse)(nil),          // 23: icbt.rpc.v1.EarmarksListResponse
	(*EventCreateResponse)(nil),           // 24: icbt.rpc.v1.EventCreateResponse
	(*EventsListResponse)(nil),            // 25: icbt.rpc.v1.EventsListResponse
	(*EventGetDetailsResponse)(nil),       // 26: icbt.rpc.v1.EventGetDetailsResponse
	(*EventListItemsResponse)(nil),        // 27: icbt.rpc.v1.EventListItemsResponse
	(*EventListEarmarksResponse)(nil),     // 28: icbt.rpc.v1.EventListEarmarksResponse
	(*EventAddItemResponse)(nil),          // 29: icbt.rpc.v1.EventAddItemResponse
	(*EventUpdateItemResponse)(nil),       // 30: icbt.rpc.v1.EventUpdateItemResponse
	(*FavoriteAddResponse)(nil),           // 31: icbt.rpc.v1.FavoriteAddResponse
	(*FavoriteListEventsResponse)(nil),    // 32: icbt.rpc.v1.FavoriteListEventsResponse
	(*NotificationsListResponse)(nil),     // 33: icbt.rpc.v1.NotificationsListResponse
}
var file_icbt_rpc_v1_service_proto_depIdxs = []int32{
	0,  // 0: icbt.rpc.v1.IcbtRpcService.EarmarkCreate:input_type -> icbt.rpc.v1.EarmarkCreateRequest
	1,  // 1: icbt.rpc.v1.IcbtRpcService.EarmarkGetDetails:input_type -> icbt.rpc.v1.EarmarkGetDetailsRequest
	2,  // 2: icbt.rpc.v1.IcbtRpcService.EarmarkRemove:input_type -> icbt.rpc.v1.EarmarkRemoveRequest
	3,  // 3: icbt.rpc.v1.IcbtRpcService.EarmarksList:input_type -> icbt.rpc.v1.EarmarksListRequest
	4,  // 4: icbt.rpc.v1.IcbtRpcService.EventCreate:input_type -> icbt.rpc.v1.EventCreateRequest
	5,  // 5: icbt.rpc.v1.IcbtRpcService.EventUpdate:input_type -> icbt.rpc.v1.EventUpdateRequest
	6,  // 6: icbt.rpc.v1.IcbtRpcService.EventDelete:input_type -> icbt.rpc.v1.EventDeleteRequest
	7,  // 7: icbt.rpc.v1.IcbtRpcService.EventsList:input_type -> icbt.rpc.v1.EventsListRequest
	8,  // 8: icbt.rpc.v1.IcbtRpcService.EventGetDetails:input_type -> icbt.rpc.v1.EventGetDetailsRequest
	9,  // 9: icbt.rpc.v1.IcbtRpcService.EventListItems:input_type -> icbt.rpc.v1.EventListItemsRequest
	10, // 10: icbt.rpc.v1.IcbtRpcService.EventListEarmarks:input_type -> icbt.rpc.v1.EventListEarmarksRequest
	11, // 11: icbt.rpc.v1.IcbtRpcService.EventAddItem:input_type -> icbt.rpc.v1.EventAddItemRequest
	12, // 12: icbt.rpc.v1.IcbtRpcService.EventUpdateItem:input_type -> icbt.rpc.v1.EventUpdateItemRequest
	13, // 13: icbt.rpc.v1.IcbtRpcService.EventRemoveItem:input_type -> icbt.rpc.v1.EventRemoveItemRequest
	14, // 14: icbt.rpc.v1.IcbtRpcService.FavoriteAdd:input_type -> icbt.rpc.v1.FavoriteAddRequest
	15, // 15: icbt.rpc.v1.IcbtRpcService.FavoriteRemove:input_type -> icbt.rpc.v1.FavoriteRemoveRequest
	16, // 16: icbt.rpc.v1.IcbtRpcService.FavoriteListEvents:input_type -> icbt.rpc.v1.FavoriteListEventsRequest
	17, // 17: icbt.rpc.v1.IcbtRpcService.NotificationDelete:input_type -> icbt.rpc.v1.NotificationDeleteRequest
	18, // 18: icbt.rpc.v1.IcbtRpcService.NotificationsDeleteAll:input_type -> icbt.rpc.v1.NotificationsDeleteAllRequest
	19, // 19: icbt.rpc.v1.IcbtRpcService.NotificationsList:input_type -> icbt.rpc.v1.NotificationsListRequest
	20, // 20: icbt.rpc.v1.IcbtRpcService.EarmarkCreate:output_type -> icbt.rpc.v1.EarmarkCreateResponse
	21, // 21: icbt.rpc.v1.IcbtRpcService.EarmarkGetDetails:output_type -> icbt.rpc.v1.EarmarkGetDetailsResponse
	22, // 22: icbt.rpc.v1.IcbtRpcService.EarmarkRemove:output_type -> google.protobuf.Empty
	23, // 23: icbt.rpc.v1.IcbtRpcService.EarmarksList:output_type -> icbt.rpc.v1.EarmarksListResponse
	24, // 24: icbt.rpc.v1.IcbtRpcService.EventCreate:output_type -> icbt.rpc.v1.EventCreateResponse
	22, // 25: icbt.rpc.v1.IcbtRpcService.EventUpdate:output_type -> google.protobuf.Empty
	22, // 26: icbt.rpc.v1.IcbtRpcService.EventDelete:output_type -> google.protobuf.Empty
	25, // 27: icbt.rpc.v1.IcbtRpcService.EventsList:output_type -> icbt.rpc.v1.EventsListResponse
	26, // 28: icbt.rpc.v1.IcbtRpcService.EventGetDetails:output_type -> icbt.rpc.v1.EventGetDetailsResponse
	27, // 29: icbt.rpc.v1.IcbtRpcService.EventListItems:output_type -> icbt.rpc.v1.EventListItemsResponse
	28, // 30: icbt.rpc.v1.IcbtRpcService.EventListEarmarks:output_type -> icbt.rpc.v1.EventListEarmarksResponse
	29, // 31: icbt.rpc.v1.IcbtRpcService.EventAddItem:output_type -> icbt.rpc.v1.EventAddItemResponse
	30, // 32: icbt.rpc.v1.IcbtRpcService.EventUpdateItem:output_type -> icbt.rpc.v1.EventUpdateItemResponse
	22, // 33: icbt.rpc.v1.IcbtRpcService.EventRemoveItem:output_type -> google.protobuf.Empty
	31, // 34: icbt.rpc.v1.IcbtRpcService.FavoriteAdd:output_type -> icbt.rpc.v1.FavoriteAddResponse
	22, // 35: icbt.rpc.v1.IcbtRpcService.FavoriteRemove:output_type -> google.protobuf.Empty
	32, // 36: icbt.rpc.v1.IcbtRpcService.FavoriteListEvents:output_type -> icbt.rpc.v1.FavoriteListEventsResponse
	22, // 37: icbt.rpc.v1.IcbtRpcService.NotificationDelete:output_type -> google.protobuf.Empty
	22, // 38: icbt.rpc.v1.IcbtRpcService.NotificationsDeleteAll:output_type -> google.protobuf.Empty
	33, // 39: icbt.rpc.v1.IcbtRpcService.NotificationsList:output_type -> icbt.rpc.v1.NotificationsListResponse
	20, // [20:40] is the sub-list for method output_type
	0,  // [0:20] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_service_proto_init() }
func file_icbt_rpc_v1_service_proto_init() {
	if File_icbt_rpc_v1_service_proto != nil {
		return
	}
	file_icbt_rpc_v1_earmark_proto_init()
	file_icbt_rpc_v1_event_proto_init()
	file_icbt_rpc_v1_favorite_proto_init()
	file_icbt_rpc_v1_notification_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_icbt_rpc_v1_service_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_service_proto_depIdxs,
	}.Build()
	File_icbt_rpc_v1_service_proto = out.File
	file_icbt_rpc_v1_service_proto_rawDesc = nil
	file_icbt_rpc_v1_service_proto_goTypes = nil
	file_icbt_rpc_v1_service_proto_depIdxs = nil
}