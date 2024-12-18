// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        (unknown)
// source: icbt/rpc/v1/notification.proto

package rpcv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Notification struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefId         string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId" json:"ref_id,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Created       *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created" json:"created,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Notification) Reset() {
	*x = Notification{}
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_notification_proto_rawDescGZIP(), []int{0}
}

func (x *Notification) GetRefId() string {
	if x != nil {
		return x.RefId
	}
	return ""
}

func (x *Notification) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Notification) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

type NotificationDeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefId         string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId" json:"ref_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NotificationDeleteRequest) Reset() {
	*x = NotificationDeleteRequest{}
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NotificationDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationDeleteRequest) ProtoMessage() {}

func (x *NotificationDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationDeleteRequest.ProtoReflect.Descriptor instead.
func (*NotificationDeleteRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_notification_proto_rawDescGZIP(), []int{1}
}

func (x *NotificationDeleteRequest) GetRefId() string {
	if x != nil {
		return x.RefId
	}
	return ""
}

type NotificationsDeleteAllRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NotificationsDeleteAllRequest) Reset() {
	*x = NotificationsDeleteAllRequest{}
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NotificationsDeleteAllRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationsDeleteAllRequest) ProtoMessage() {}

func (x *NotificationsDeleteAllRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationsDeleteAllRequest.ProtoReflect.Descriptor instead.
func (*NotificationsDeleteAllRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_notification_proto_rawDescGZIP(), []int{2}
}

type NotificationsListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Pagination    *PaginationRequest     `protobuf:"bytes,1,opt,name=pagination" json:"pagination,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NotificationsListRequest) Reset() {
	*x = NotificationsListRequest{}
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NotificationsListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationsListRequest) ProtoMessage() {}

func (x *NotificationsListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationsListRequest.ProtoReflect.Descriptor instead.
func (*NotificationsListRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_notification_proto_rawDescGZIP(), []int{3}
}

func (x *NotificationsListRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

type NotificationsListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Notifications []*Notification        `protobuf:"bytes,1,rep,name=notifications" json:"notifications,omitempty"`
	Pagination    *PaginationResult      `protobuf:"bytes,2,opt,name=pagination" json:"pagination,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NotificationsListResponse) Reset() {
	*x = NotificationsListResponse{}
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NotificationsListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NotificationsListResponse) ProtoMessage() {}

func (x *NotificationsListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_notification_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NotificationsListResponse.ProtoReflect.Descriptor instead.
func (*NotificationsListResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_notification_proto_rawDescGZIP(), []int{4}
}

func (x *NotificationsListResponse) GetNotifications() []*Notification {
	if x != nil {
		return x.Notifications
	}
	return nil
}

func (x *NotificationsListResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.Pagination
	}
	return nil
}

var File_icbt_rpc_v1_notification_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_notification_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0b, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62,
	0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x69, 0x63, 0x62,
	0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61,
	0x69, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x69, 0x63, 0x62, 0x74,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x75, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69,
	0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x66, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x22,
	0x3f, 0x0a, 0x19, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06,
	0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48,
	0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64,
	0x22, 0x1f, 0x0a, 0x1d, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x61, 0x0a, 0x18, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a,
	0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e,
	0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x05, 0xaa, 0x01, 0x02, 0x08, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0xa2, 0x01, 0x0a, 0x19, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x69, 0x63, 0x62, 0x74,
	0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x44, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x05, 0xaa, 0x01, 0x02, 0x08, 0x01, 0x52, 0x0a, 0x70,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0xb1, 0x01, 0x0a, 0x0f, 0x63, 0x6f,
	0x6d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x11, 0x4e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64,
	0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f, 0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69,
	0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49,
	0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a,
	0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x92, 0x03, 0x02, 0x08, 0x02, 0x62, 0x08, 0x65,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
}

var (
	file_icbt_rpc_v1_notification_proto_rawDescOnce sync.Once
	file_icbt_rpc_v1_notification_proto_rawDescData = file_icbt_rpc_v1_notification_proto_rawDesc
)

func file_icbt_rpc_v1_notification_proto_rawDescGZIP() []byte {
	file_icbt_rpc_v1_notification_proto_rawDescOnce.Do(func() {
		file_icbt_rpc_v1_notification_proto_rawDescData = protoimpl.X.CompressGZIP(file_icbt_rpc_v1_notification_proto_rawDescData)
	})
	return file_icbt_rpc_v1_notification_proto_rawDescData
}

var file_icbt_rpc_v1_notification_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_icbt_rpc_v1_notification_proto_goTypes = []any{
	(*Notification)(nil),                  // 0: icbt.rpc.v1.Notification
	(*NotificationDeleteRequest)(nil),     // 1: icbt.rpc.v1.NotificationDeleteRequest
	(*NotificationsDeleteAllRequest)(nil), // 2: icbt.rpc.v1.NotificationsDeleteAllRequest
	(*NotificationsListRequest)(nil),      // 3: icbt.rpc.v1.NotificationsListRequest
	(*NotificationsListResponse)(nil),     // 4: icbt.rpc.v1.NotificationsListResponse
	(*timestamppb.Timestamp)(nil),         // 5: google.protobuf.Timestamp
	(*PaginationRequest)(nil),             // 6: icbt.rpc.v1.PaginationRequest
	(*PaginationResult)(nil),              // 7: icbt.rpc.v1.PaginationResult
}
var file_icbt_rpc_v1_notification_proto_depIdxs = []int32{
	5, // 0: icbt.rpc.v1.Notification.created:type_name -> google.protobuf.Timestamp
	6, // 1: icbt.rpc.v1.NotificationsListRequest.pagination:type_name -> icbt.rpc.v1.PaginationRequest
	0, // 2: icbt.rpc.v1.NotificationsListResponse.notifications:type_name -> icbt.rpc.v1.Notification
	7, // 3: icbt.rpc.v1.NotificationsListResponse.pagination:type_name -> icbt.rpc.v1.PaginationResult
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_notification_proto_init() }
func file_icbt_rpc_v1_notification_proto_init() {
	if File_icbt_rpc_v1_notification_proto != nil {
		return
	}
	file_icbt_rpc_v1_constraints_proto_init()
	file_icbt_rpc_v1_pagination_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_notification_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_notification_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_notification_proto_depIdxs,
		MessageInfos:      file_icbt_rpc_v1_notification_proto_msgTypes,
	}.Build()
	File_icbt_rpc_v1_notification_proto = out.File
	file_icbt_rpc_v1_notification_proto_rawDesc = nil
	file_icbt_rpc_v1_notification_proto_goTypes = nil
	file_icbt_rpc_v1_notification_proto_depIdxs = nil
}
