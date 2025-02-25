// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: icbt/rpc/v1/favorite.proto

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

type Favorite struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_EventRefId string                 `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId"`
	xxx_hidden_Created    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *Favorite) Reset() {
	*x = Favorite{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Favorite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Favorite) ProtoMessage() {}

func (x *Favorite) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Favorite) GetEventRefId() string {
	if x != nil {
		return x.xxx_hidden_EventRefId
	}
	return ""
}

func (x *Favorite) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.xxx_hidden_Created
	}
	return nil
}

func (x *Favorite) SetEventRefId(v string) {
	x.xxx_hidden_EventRefId = v
}

func (x *Favorite) SetCreated(v *timestamppb.Timestamp) {
	x.xxx_hidden_Created = v
}

func (x *Favorite) HasCreated() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Created != nil
}

func (x *Favorite) ClearCreated() {
	x.xxx_hidden_Created = nil
}

type Favorite_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	EventRefId string
	Created    *timestamppb.Timestamp
}

func (b0 Favorite_builder) Build() *Favorite {
	m0 := &Favorite{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_EventRefId = b.EventRefId
	x.xxx_hidden_Created = b.Created
	return m0
}

type FavoriteAddRequest struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_EventRefId string                 `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *FavoriteAddRequest) Reset() {
	*x = FavoriteAddRequest{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FavoriteAddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteAddRequest) ProtoMessage() {}

func (x *FavoriteAddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FavoriteAddRequest) GetEventRefId() string {
	if x != nil {
		return x.xxx_hidden_EventRefId
	}
	return ""
}

func (x *FavoriteAddRequest) SetEventRefId(v string) {
	x.xxx_hidden_EventRefId = v
}

type FavoriteAddRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	EventRefId string
}

func (b0 FavoriteAddRequest_builder) Build() *FavoriteAddRequest {
	m0 := &FavoriteAddRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_EventRefId = b.EventRefId
	return m0
}

type FavoriteAddResponse struct {
	state               protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Favorite *Favorite              `protobuf:"bytes,1,opt,name=favorite"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *FavoriteAddResponse) Reset() {
	*x = FavoriteAddResponse{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FavoriteAddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteAddResponse) ProtoMessage() {}

func (x *FavoriteAddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FavoriteAddResponse) GetFavorite() *Favorite {
	if x != nil {
		return x.xxx_hidden_Favorite
	}
	return nil
}

func (x *FavoriteAddResponse) SetFavorite(v *Favorite) {
	x.xxx_hidden_Favorite = v
}

func (x *FavoriteAddResponse) HasFavorite() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Favorite != nil
}

func (x *FavoriteAddResponse) ClearFavorite() {
	x.xxx_hidden_Favorite = nil
}

type FavoriteAddResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Favorite *Favorite
}

func (b0 FavoriteAddResponse_builder) Build() *FavoriteAddResponse {
	m0 := &FavoriteAddResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Favorite = b.Favorite
	return m0
}

type FavoriteRemoveRequest struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_EventRefId string                 `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *FavoriteRemoveRequest) Reset() {
	*x = FavoriteRemoveRequest{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FavoriteRemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteRemoveRequest) ProtoMessage() {}

func (x *FavoriteRemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FavoriteRemoveRequest) GetEventRefId() string {
	if x != nil {
		return x.xxx_hidden_EventRefId
	}
	return ""
}

func (x *FavoriteRemoveRequest) SetEventRefId(v string) {
	x.xxx_hidden_EventRefId = v
}

type FavoriteRemoveRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	EventRefId string
}

func (b0 FavoriteRemoveRequest_builder) Build() *FavoriteRemoveRequest {
	m0 := &FavoriteRemoveRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_EventRefId = b.EventRefId
	return m0
}

type FavoriteListEventsRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Pagination  *PaginationRequest     `protobuf:"bytes,1,opt,name=pagination"`
	xxx_hidden_Archived    bool                   `protobuf:"varint,2,opt,name=archived"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *FavoriteListEventsRequest) Reset() {
	*x = FavoriteListEventsRequest{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FavoriteListEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteListEventsRequest) ProtoMessage() {}

func (x *FavoriteListEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FavoriteListEventsRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *FavoriteListEventsRequest) GetArchived() bool {
	if x != nil {
		return x.xxx_hidden_Archived
	}
	return false
}

func (x *FavoriteListEventsRequest) SetPagination(v *PaginationRequest) {
	x.xxx_hidden_Pagination = v
}

func (x *FavoriteListEventsRequest) SetArchived(v bool) {
	x.xxx_hidden_Archived = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 1, 2)
}

func (x *FavoriteListEventsRequest) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *FavoriteListEventsRequest) HasArchived() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 1)
}

func (x *FavoriteListEventsRequest) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

func (x *FavoriteListEventsRequest) ClearArchived() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 1)
	x.xxx_hidden_Archived = false
}

type FavoriteListEventsRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Pagination *PaginationRequest
	Archived   *bool
}

func (b0 FavoriteListEventsRequest_builder) Build() *FavoriteListEventsRequest {
	m0 := &FavoriteListEventsRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Pagination = b.Pagination
	if b.Archived != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 1, 2)
		x.xxx_hidden_Archived = *b.Archived
	}
	return m0
}

type FavoriteListEventsResponse struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Events     *[]*Event              `protobuf:"bytes,1,rep,name=events"`
	xxx_hidden_Pagination *PaginationResult      `protobuf:"bytes,2,opt,name=pagination"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *FavoriteListEventsResponse) Reset() {
	*x = FavoriteListEventsResponse{}
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FavoriteListEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteListEventsResponse) ProtoMessage() {}

func (x *FavoriteListEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FavoriteListEventsResponse) GetEvents() []*Event {
	if x != nil {
		if x.xxx_hidden_Events != nil {
			return *x.xxx_hidden_Events
		}
	}
	return nil
}

func (x *FavoriteListEventsResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *FavoriteListEventsResponse) SetEvents(v []*Event) {
	x.xxx_hidden_Events = &v
}

func (x *FavoriteListEventsResponse) SetPagination(v *PaginationResult) {
	x.xxx_hidden_Pagination = v
}

func (x *FavoriteListEventsResponse) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *FavoriteListEventsResponse) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

type FavoriteListEventsResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Events     []*Event
	Pagination *PaginationResult
}

func (b0 FavoriteListEventsResponse_builder) Build() *FavoriteListEventsResponse {
	m0 := &FavoriteListEventsResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Events = &b.Events
	x.xxx_hidden_Pagination = b.Pagination
	return m0
}

var File_icbt_rpc_v1_favorite_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_favorite_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x5f, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x69, 0x63, 0x62, 0x74,
	0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69,
	0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x69, 0x63, 0x62, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f,
	0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x62, 0x0a, 0x08, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0c,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x12, 0x34,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x22, 0x43, 0x0a, 0x12, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x0c, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x0a, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x13, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x31, 0x0a, 0x08, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x08, 0x66, 0x61, 0x76, 0x6f, 0x72,
	0x69, 0x74, 0x65, 0x22, 0x46, 0x0a, 0x15, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52,
	0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x0c,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x85, 0x01, 0x0a, 0x19,
	0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x0a, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x05, 0xaa,
	0x01, 0x02, 0x08, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x21, 0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x42, 0x05, 0xaa, 0x01, 0x02, 0x08, 0x01, 0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69,
	0x76, 0x65, 0x64, 0x22, 0x8e, 0x01, 0x0a, 0x1a, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x44,
	0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x42, 0x05, 0xaa, 0x01, 0x02, 0x08, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x42, 0xb2, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f,
	0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69, 0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70,
	0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70,
	0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74,
	0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52,
	0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63,
	0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea,
	0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x92,
	0x03, 0x07, 0xd2, 0x3e, 0x02, 0x10, 0x03, 0x08, 0x02, 0x62, 0x08, 0x65, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
})

var file_icbt_rpc_v1_favorite_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_icbt_rpc_v1_favorite_proto_goTypes = []any{
	(*Favorite)(nil),                   // 0: icbt.rpc.v1.Favorite
	(*FavoriteAddRequest)(nil),         // 1: icbt.rpc.v1.FavoriteAddRequest
	(*FavoriteAddResponse)(nil),        // 2: icbt.rpc.v1.FavoriteAddResponse
	(*FavoriteRemoveRequest)(nil),      // 3: icbt.rpc.v1.FavoriteRemoveRequest
	(*FavoriteListEventsRequest)(nil),  // 4: icbt.rpc.v1.FavoriteListEventsRequest
	(*FavoriteListEventsResponse)(nil), // 5: icbt.rpc.v1.FavoriteListEventsResponse
	(*timestamppb.Timestamp)(nil),      // 6: google.protobuf.Timestamp
	(*PaginationRequest)(nil),          // 7: icbt.rpc.v1.PaginationRequest
	(*Event)(nil),                      // 8: icbt.rpc.v1.Event
	(*PaginationResult)(nil),           // 9: icbt.rpc.v1.PaginationResult
}
var file_icbt_rpc_v1_favorite_proto_depIdxs = []int32{
	6, // 0: icbt.rpc.v1.Favorite.created:type_name -> google.protobuf.Timestamp
	0, // 1: icbt.rpc.v1.FavoriteAddResponse.favorite:type_name -> icbt.rpc.v1.Favorite
	7, // 2: icbt.rpc.v1.FavoriteListEventsRequest.pagination:type_name -> icbt.rpc.v1.PaginationRequest
	8, // 3: icbt.rpc.v1.FavoriteListEventsResponse.events:type_name -> icbt.rpc.v1.Event
	9, // 4: icbt.rpc.v1.FavoriteListEventsResponse.pagination:type_name -> icbt.rpc.v1.PaginationResult
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_favorite_proto_init() }
func file_icbt_rpc_v1_favorite_proto_init() {
	if File_icbt_rpc_v1_favorite_proto != nil {
		return
	}
	file_icbt_rpc_v1_constraints_proto_init()
	file_icbt_rpc_v1_event_proto_init()
	file_icbt_rpc_v1_pagination_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_icbt_rpc_v1_favorite_proto_rawDesc), len(file_icbt_rpc_v1_favorite_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_favorite_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_favorite_proto_depIdxs,
		MessageInfos:      file_icbt_rpc_v1_favorite_proto_msgTypes,
	}.Build()
	File_icbt_rpc_v1_favorite_proto = out.File
	file_icbt_rpc_v1_favorite_proto_goTypes = nil
	file_icbt_rpc_v1_favorite_proto_depIdxs = nil
}
