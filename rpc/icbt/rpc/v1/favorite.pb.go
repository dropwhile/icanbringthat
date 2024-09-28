// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: icbt/rpc/v1/favorite.proto

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

type Favorite struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventRefId string                 `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId,proto3" json:"event_ref_id,omitempty"`
	Created    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *Favorite) Reset() {
	*x = Favorite{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Favorite) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Favorite) ProtoMessage() {}

func (x *Favorite) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Favorite.ProtoReflect.Descriptor instead.
func (*Favorite) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{0}
}

func (x *Favorite) GetEventRefId() string {
	if x != nil {
		return x.EventRefId
	}
	return ""
}

func (x *Favorite) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

type FavoriteAddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventRefId string `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId,proto3" json:"event_ref_id,omitempty"`
}

func (x *FavoriteAddRequest) Reset() {
	*x = FavoriteAddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteAddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteAddRequest) ProtoMessage() {}

func (x *FavoriteAddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteAddRequest.ProtoReflect.Descriptor instead.
func (*FavoriteAddRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{1}
}

func (x *FavoriteAddRequest) GetEventRefId() string {
	if x != nil {
		return x.EventRefId
	}
	return ""
}

type FavoriteAddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Favorite *Favorite `protobuf:"bytes,1,opt,name=favorite,proto3" json:"favorite,omitempty"`
}

func (x *FavoriteAddResponse) Reset() {
	*x = FavoriteAddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteAddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteAddResponse) ProtoMessage() {}

func (x *FavoriteAddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteAddResponse.ProtoReflect.Descriptor instead.
func (*FavoriteAddResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{2}
}

func (x *FavoriteAddResponse) GetFavorite() *Favorite {
	if x != nil {
		return x.Favorite
	}
	return nil
}

type FavoriteRemoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventRefId string `protobuf:"bytes,1,opt,name=event_ref_id,json=eventRefId,proto3" json:"event_ref_id,omitempty"`
}

func (x *FavoriteRemoveRequest) Reset() {
	*x = FavoriteRemoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteRemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteRemoveRequest) ProtoMessage() {}

func (x *FavoriteRemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteRemoveRequest.ProtoReflect.Descriptor instead.
func (*FavoriteRemoveRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{3}
}

func (x *FavoriteRemoveRequest) GetEventRefId() string {
	if x != nil {
		return x.EventRefId
	}
	return ""
}

type FavoriteListEventsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3,oneof" json:"pagination,omitempty"`
	Archived   *bool              `protobuf:"varint,2,opt,name=archived,proto3,oneof" json:"archived,omitempty"`
}

func (x *FavoriteListEventsRequest) Reset() {
	*x = FavoriteListEventsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteListEventsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteListEventsRequest) ProtoMessage() {}

func (x *FavoriteListEventsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteListEventsRequest.ProtoReflect.Descriptor instead.
func (*FavoriteListEventsRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{4}
}

func (x *FavoriteListEventsRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *FavoriteListEventsRequest) GetArchived() bool {
	if x != nil && x.Archived != nil {
		return *x.Archived
	}
	return false
}

type FavoriteListEventsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events     []*Event          `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	Pagination *PaginationResult `protobuf:"bytes,2,opt,name=pagination,proto3,oneof" json:"pagination,omitempty"`
}

func (x *FavoriteListEventsResponse) Reset() {
	*x = FavoriteListEventsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FavoriteListEventsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FavoriteListEventsResponse) ProtoMessage() {}

func (x *FavoriteListEventsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_favorite_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FavoriteListEventsResponse.ProtoReflect.Descriptor instead.
func (*FavoriteListEventsResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_favorite_proto_rawDescGZIP(), []int{5}
}

func (x *FavoriteListEventsResponse) GetEvents() []*Event {
	if x != nil {
		return x.Events
	}
	return nil
}

func (x *FavoriteListEventsResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.Pagination
	}
	return nil
}

var File_icbt_rpc_v1_favorite_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_favorite_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x66, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70,
	0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1c, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a,
	0x08, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x22, 0x43, 0x0a, 0x12, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba,
	0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x48, 0x0a, 0x13, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69,
	0x74, 0x65, 0x41, 0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a,
	0x08, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x08, 0x66, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x22, 0x46, 0x0a, 0x15, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2d, 0x0a, 0x0c, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x0a, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x9d, 0x01, 0x0a, 0x19, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x43, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61,
	0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52,
	0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b,
	0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f,
	0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x22, 0x9b, 0x01, 0x0a, 0x1a, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x06, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x12, 0x42, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0xa8, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69,
	0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c,
	0x65, 0x2f, 0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69, 0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b,
	0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63,
	0x62, 0x74, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74,
	0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52,
	0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_icbt_rpc_v1_favorite_proto_rawDescOnce sync.Once
	file_icbt_rpc_v1_favorite_proto_rawDescData = file_icbt_rpc_v1_favorite_proto_rawDesc
)

func file_icbt_rpc_v1_favorite_proto_rawDescGZIP() []byte {
	file_icbt_rpc_v1_favorite_proto_rawDescOnce.Do(func() {
		file_icbt_rpc_v1_favorite_proto_rawDescData = protoimpl.X.CompressGZIP(file_icbt_rpc_v1_favorite_proto_rawDescData)
	})
	return file_icbt_rpc_v1_favorite_proto_rawDescData
}

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
	if !protoimpl.UnsafeEnabled {
		file_icbt_rpc_v1_favorite_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Favorite); i {
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
		file_icbt_rpc_v1_favorite_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*FavoriteAddRequest); i {
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
		file_icbt_rpc_v1_favorite_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*FavoriteAddResponse); i {
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
		file_icbt_rpc_v1_favorite_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*FavoriteRemoveRequest); i {
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
		file_icbt_rpc_v1_favorite_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*FavoriteListEventsRequest); i {
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
		file_icbt_rpc_v1_favorite_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*FavoriteListEventsResponse); i {
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
	file_icbt_rpc_v1_favorite_proto_msgTypes[4].OneofWrappers = []any{}
	file_icbt_rpc_v1_favorite_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_favorite_proto_rawDesc,
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
	file_icbt_rpc_v1_favorite_proto_rawDesc = nil
	file_icbt_rpc_v1_favorite_proto_goTypes = nil
	file_icbt_rpc_v1_favorite_proto_depIdxs = nil
}
