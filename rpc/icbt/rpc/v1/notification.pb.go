// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: icbt/rpc/v1/notification.proto

package rpcv1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/gofeaturespb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Notification struct {
	state              protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_RefId   string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId"`
	xxx_hidden_Message string                 `protobuf:"bytes,2,opt,name=message"`
	xxx_hidden_Created *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
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

func (x *Notification) GetRefId() string {
	if x != nil {
		return x.xxx_hidden_RefId
	}
	return ""
}

func (x *Notification) GetMessage() string {
	if x != nil {
		return x.xxx_hidden_Message
	}
	return ""
}

func (x *Notification) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.xxx_hidden_Created
	}
	return nil
}

func (x *Notification) SetRefId(v string) {
	x.xxx_hidden_RefId = v
}

func (x *Notification) SetMessage(v string) {
	x.xxx_hidden_Message = v
}

func (x *Notification) SetCreated(v *timestamppb.Timestamp) {
	x.xxx_hidden_Created = v
}

func (x *Notification) HasCreated() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Created != nil
}

func (x *Notification) ClearCreated() {
	x.xxx_hidden_Created = nil
}

type Notification_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	RefId   string
	Message string
	Created *timestamppb.Timestamp
}

func (b0 Notification_builder) Build() *Notification {
	m0 := &Notification{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_RefId = b.RefId
	x.xxx_hidden_Message = b.Message
	x.xxx_hidden_Created = b.Created
	return m0
}

type NotificationDeleteRequest struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_RefId string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
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

func (x *NotificationDeleteRequest) GetRefId() string {
	if x != nil {
		return x.xxx_hidden_RefId
	}
	return ""
}

func (x *NotificationDeleteRequest) SetRefId(v string) {
	x.xxx_hidden_RefId = v
}

type NotificationDeleteRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	RefId string
}

func (b0 NotificationDeleteRequest_builder) Build() *NotificationDeleteRequest {
	m0 := &NotificationDeleteRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_RefId = b.RefId
	return m0
}

type NotificationsDeleteAllRequest struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
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

type NotificationsDeleteAllRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 NotificationsDeleteAllRequest_builder) Build() *NotificationsDeleteAllRequest {
	m0 := &NotificationsDeleteAllRequest{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type NotificationsListRequest struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Pagination *PaginationRequest     `protobuf:"bytes,1,opt,name=pagination"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
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

func (x *NotificationsListRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *NotificationsListRequest) SetPagination(v *PaginationRequest) {
	x.xxx_hidden_Pagination = v
}

func (x *NotificationsListRequest) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *NotificationsListRequest) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

type NotificationsListRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Pagination *PaginationRequest
}

func (b0 NotificationsListRequest_builder) Build() *NotificationsListRequest {
	m0 := &NotificationsListRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Pagination = b.Pagination
	return m0
}

type NotificationsListResponse struct {
	state                    protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Notifications *[]*Notification       `protobuf:"bytes,1,rep,name=notifications"`
	xxx_hidden_Pagination    *PaginationResult      `protobuf:"bytes,2,opt,name=pagination"`
	unknownFields            protoimpl.UnknownFields
	sizeCache                protoimpl.SizeCache
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

func (x *NotificationsListResponse) GetNotifications() []*Notification {
	if x != nil {
		if x.xxx_hidden_Notifications != nil {
			return *x.xxx_hidden_Notifications
		}
	}
	return nil
}

func (x *NotificationsListResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *NotificationsListResponse) SetNotifications(v []*Notification) {
	x.xxx_hidden_Notifications = &v
}

func (x *NotificationsListResponse) SetPagination(v *PaginationResult) {
	x.xxx_hidden_Pagination = v
}

func (x *NotificationsListResponse) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *NotificationsListResponse) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

type NotificationsListResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Notifications []*Notification
	Pagination    *PaginationResult
}

func (b0 NotificationsListResponse_builder) Build() *NotificationsListResponse {
	m0 := &NotificationsListResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Notifications = &b.Notifications
	x.xxx_hidden_Pagination = b.Pagination
	return m0
}

var File_icbt_rpc_v1_notification_proto protoreflect.FileDescriptor

const file_icbt_rpc_v1_notification_proto_rawDesc = "" +
	"\n" +
	"\x1eicbt/rpc/v1/notification.proto\x12\vicbt.rpc.v1\x1a\x1bbuf/validate/validate.proto\x1a!google/protobuf/go_features.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1dicbt/rpc/v1/constraints.proto\x1a\x1cicbt/rpc/v1/pagination.proto\"u\n" +
	"\fNotification\x12\x15\n" +
	"\x06ref_id\x18\x01 \x01(\tR\x05refId\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x124\n" +
	"\acreated\x18\x03 \x01(\v2\x1a.google.protobuf.TimestampR\acreated\"?\n" +
	"\x19NotificationDeleteRequest\x12\"\n" +
	"\x06ref_id\x18\x01 \x01(\tB\v\xbaH\br\x06\x88\u0603\x8b\x02\x01R\x05refId\"\x1f\n" +
	"\x1dNotificationsDeleteAllRequest\"a\n" +
	"\x18NotificationsListRequest\x12E\n" +
	"\n" +
	"pagination\x18\x01 \x01(\v2\x1e.icbt.rpc.v1.PaginationRequestB\x05\xaa\x01\x02\b\x01R\n" +
	"pagination\"\xa2\x01\n" +
	"\x19NotificationsListResponse\x12?\n" +
	"\rnotifications\x18\x01 \x03(\v2\x19.icbt.rpc.v1.NotificationR\rnotifications\x12D\n" +
	"\n" +
	"pagination\x18\x02 \x01(\v2\x1d.icbt.rpc.v1.PaginationResultB\x05\xaa\x01\x02\b\x01R\n" +
	"paginationB\xb6\x01\n" +
	"\x0fcom.icbt.rpc.v1B\x11NotificationProtoP\x01Z8github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1;rpcv1\xa2\x02\x03IRX\xaa\x02\vIcbt.Rpc.V1\xca\x02\vIcbt\\Rpc\\V1\xe2\x02\x17Icbt\\Rpc\\V1\\GPBMetadata\xea\x02\rIcbt::Rpc::V1\x92\x03\a\xd2>\x02\x10\x03\b\x02b\beditionsp\xe8\a"

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
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_icbt_rpc_v1_notification_proto_rawDesc), len(file_icbt_rpc_v1_notification_proto_rawDesc)),
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
	file_icbt_rpc_v1_notification_proto_goTypes = nil
	file_icbt_rpc_v1_notification_proto_depIdxs = nil
}
