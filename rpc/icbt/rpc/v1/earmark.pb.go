// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: icbt/rpc/v1/earmark.proto

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

type Earmark struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefId          string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId,proto3" json:"ref_id,omitempty"`
	EventItemRefId string                 `protobuf:"bytes,2,opt,name=event_item_ref_id,json=eventItemRefId,proto3" json:"event_item_ref_id,omitempty"`
	Note           string                 `protobuf:"bytes,3,opt,name=note,proto3" json:"note,omitempty"`
	Owner          string                 `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Created        *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created,proto3" json:"created,omitempty"`
}

func (x *Earmark) Reset() {
	*x = Earmark{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Earmark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Earmark) ProtoMessage() {}

func (x *Earmark) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Earmark.ProtoReflect.Descriptor instead.
func (*Earmark) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{0}
}

func (x *Earmark) GetRefId() string {
	if x != nil {
		return x.RefId
	}
	return ""
}

func (x *Earmark) GetEventItemRefId() string {
	if x != nil {
		return x.EventItemRefId
	}
	return ""
}

func (x *Earmark) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *Earmark) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Earmark) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.Created
	}
	return nil
}

type EarmarkCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventItemRefId string `protobuf:"bytes,1,opt,name=event_item_ref_id,json=eventItemRefId,proto3" json:"event_item_ref_id,omitempty"`
	Note           string `protobuf:"bytes,2,opt,name=note,proto3" json:"note,omitempty"` // required, but can be empty
}

func (x *EarmarkCreateRequest) Reset() {
	*x = EarmarkCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarkCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkCreateRequest) ProtoMessage() {}

func (x *EarmarkCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarkCreateRequest.ProtoReflect.Descriptor instead.
func (*EarmarkCreateRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{1}
}

func (x *EarmarkCreateRequest) GetEventItemRefId() string {
	if x != nil {
		return x.EventItemRefId
	}
	return ""
}

func (x *EarmarkCreateRequest) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

type EarmarkCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Earmark *Earmark `protobuf:"bytes,1,opt,name=earmark,proto3" json:"earmark,omitempty"`
}

func (x *EarmarkCreateResponse) Reset() {
	*x = EarmarkCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarkCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkCreateResponse) ProtoMessage() {}

func (x *EarmarkCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarkCreateResponse.ProtoReflect.Descriptor instead.
func (*EarmarkCreateResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{2}
}

func (x *EarmarkCreateResponse) GetEarmark() *Earmark {
	if x != nil {
		return x.Earmark
	}
	return nil
}

type EarmarkRemoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefId string `protobuf:"bytes,1,opt,name=ref_id,json=refId,proto3" json:"ref_id,omitempty"`
}

func (x *EarmarkRemoveRequest) Reset() {
	*x = EarmarkRemoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarkRemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkRemoveRequest) ProtoMessage() {}

func (x *EarmarkRemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarkRemoveRequest.ProtoReflect.Descriptor instead.
func (*EarmarkRemoveRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{3}
}

func (x *EarmarkRemoveRequest) GetRefId() string {
	if x != nil {
		return x.RefId
	}
	return ""
}

type EarmarkGetDetailsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefId string `protobuf:"bytes,1,opt,name=ref_id,json=refId,proto3" json:"ref_id,omitempty"`
}

func (x *EarmarkGetDetailsRequest) Reset() {
	*x = EarmarkGetDetailsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarkGetDetailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkGetDetailsRequest) ProtoMessage() {}

func (x *EarmarkGetDetailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarkGetDetailsRequest.ProtoReflect.Descriptor instead.
func (*EarmarkGetDetailsRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{4}
}

func (x *EarmarkGetDetailsRequest) GetRefId() string {
	if x != nil {
		return x.RefId
	}
	return ""
}

type EarmarkGetDetailsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Earmark    *Earmark `protobuf:"bytes,1,opt,name=earmark,proto3" json:"earmark,omitempty"`
	EventRefId string   `protobuf:"bytes,2,opt,name=event_ref_id,json=eventRefId,proto3" json:"event_ref_id,omitempty"`
}

func (x *EarmarkGetDetailsResponse) Reset() {
	*x = EarmarkGetDetailsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarkGetDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkGetDetailsResponse) ProtoMessage() {}

func (x *EarmarkGetDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarkGetDetailsResponse.ProtoReflect.Descriptor instead.
func (*EarmarkGetDetailsResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{5}
}

func (x *EarmarkGetDetailsResponse) GetEarmark() *Earmark {
	if x != nil {
		return x.Earmark
	}
	return nil
}

func (x *EarmarkGetDetailsResponse) GetEventRefId() string {
	if x != nil {
		return x.EventRefId
	}
	return ""
}

type EarmarksListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pagination *PaginationRequest `protobuf:"bytes,1,opt,name=pagination,proto3,oneof" json:"pagination,omitempty"`
	Archived   *bool              `protobuf:"varint,2,opt,name=archived,proto3,oneof" json:"archived,omitempty"`
}

func (x *EarmarksListRequest) Reset() {
	*x = EarmarksListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarksListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarksListRequest) ProtoMessage() {}

func (x *EarmarksListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarksListRequest.ProtoReflect.Descriptor instead.
func (*EarmarksListRequest) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{6}
}

func (x *EarmarksListRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.Pagination
	}
	return nil
}

func (x *EarmarksListRequest) GetArchived() bool {
	if x != nil && x.Archived != nil {
		return *x.Archived
	}
	return false
}

type EarmarksListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Earmarks   []*Earmark        `protobuf:"bytes,1,rep,name=earmarks,proto3" json:"earmarks,omitempty"`
	Pagination *PaginationResult `protobuf:"bytes,2,opt,name=pagination,proto3,oneof" json:"pagination,omitempty"`
}

func (x *EarmarksListResponse) Reset() {
	*x = EarmarksListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EarmarksListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarksListResponse) ProtoMessage() {}

func (x *EarmarksListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EarmarksListResponse.ProtoReflect.Descriptor instead.
func (*EarmarksListResponse) Descriptor() ([]byte, []int) {
	return file_icbt_rpc_v1_earmark_proto_rawDescGZIP(), []int{7}
}

func (x *EarmarksListResponse) GetEarmarks() []*Earmark {
	if x != nil {
		return x.Earmarks
	}
	return nil
}

func (x *EarmarksListResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.Pagination
	}
	return nil
}

var File_icbt_rpc_v1_earmark_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_earmark_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x07, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12,
	0x15, 0x0a, 0x06, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x11, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x74, 0x65, 0x6d, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x66, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x6f, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x34, 0x0a, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x22, 0x62, 0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x11, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83, 0x8b, 0x02,
	0x01, 0x52, 0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x66, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x6f, 0x74, 0x65, 0x22, 0x47, 0x0a, 0x15, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e,
	0x0a, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x22, 0x3a,
	0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83,
	0x8b, 0x02, 0x01, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x22, 0x3e, 0x0a, 0x18, 0x45, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8, 0x83,
	0x8b, 0x02, 0x01, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x22, 0x6d, 0x0a, 0x19, 0x45, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61,
	0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x07,
	0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x20, 0x0a, 0x0c, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x97, 0x01, 0x0a, 0x13, 0x45, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x43, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76,
	0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x01, 0x52, 0x08, 0x61, 0x72, 0x63, 0x68,
	0x69, 0x76, 0x65, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x72, 0x63, 0x68, 0x69,
	0x76, 0x65, 0x64, 0x22, 0x9b, 0x01, 0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73,
	0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x08,
	0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72,
	0x6d, 0x61, 0x72, 0x6b, 0x52, 0x08, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x12, 0x42,
	0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x88,
	0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0xa7, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x64, 0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f, 0x69, 0x63, 0x61, 0x6e,
	0x62, 0x72, 0x69, 0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x63,
	0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x49, 0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x2e, 0x52, 0x70, 0x63,
	0x2e, 0x56, 0x31, 0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49, 0x63,
	0x62, 0x74, 0x3a, 0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_icbt_rpc_v1_earmark_proto_rawDescOnce sync.Once
	file_icbt_rpc_v1_earmark_proto_rawDescData = file_icbt_rpc_v1_earmark_proto_rawDesc
)

func file_icbt_rpc_v1_earmark_proto_rawDescGZIP() []byte {
	file_icbt_rpc_v1_earmark_proto_rawDescOnce.Do(func() {
		file_icbt_rpc_v1_earmark_proto_rawDescData = protoimpl.X.CompressGZIP(file_icbt_rpc_v1_earmark_proto_rawDescData)
	})
	return file_icbt_rpc_v1_earmark_proto_rawDescData
}

var file_icbt_rpc_v1_earmark_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_icbt_rpc_v1_earmark_proto_goTypes = []any{
	(*Earmark)(nil),                   // 0: icbt.rpc.v1.Earmark
	(*EarmarkCreateRequest)(nil),      // 1: icbt.rpc.v1.EarmarkCreateRequest
	(*EarmarkCreateResponse)(nil),     // 2: icbt.rpc.v1.EarmarkCreateResponse
	(*EarmarkRemoveRequest)(nil),      // 3: icbt.rpc.v1.EarmarkRemoveRequest
	(*EarmarkGetDetailsRequest)(nil),  // 4: icbt.rpc.v1.EarmarkGetDetailsRequest
	(*EarmarkGetDetailsResponse)(nil), // 5: icbt.rpc.v1.EarmarkGetDetailsResponse
	(*EarmarksListRequest)(nil),       // 6: icbt.rpc.v1.EarmarksListRequest
	(*EarmarksListResponse)(nil),      // 7: icbt.rpc.v1.EarmarksListResponse
	(*timestamppb.Timestamp)(nil),     // 8: google.protobuf.Timestamp
	(*PaginationRequest)(nil),         // 9: icbt.rpc.v1.PaginationRequest
	(*PaginationResult)(nil),          // 10: icbt.rpc.v1.PaginationResult
}
var file_icbt_rpc_v1_earmark_proto_depIdxs = []int32{
	8,  // 0: icbt.rpc.v1.Earmark.created:type_name -> google.protobuf.Timestamp
	0,  // 1: icbt.rpc.v1.EarmarkCreateResponse.earmark:type_name -> icbt.rpc.v1.Earmark
	0,  // 2: icbt.rpc.v1.EarmarkGetDetailsResponse.earmark:type_name -> icbt.rpc.v1.Earmark
	9,  // 3: icbt.rpc.v1.EarmarksListRequest.pagination:type_name -> icbt.rpc.v1.PaginationRequest
	0,  // 4: icbt.rpc.v1.EarmarksListResponse.earmarks:type_name -> icbt.rpc.v1.Earmark
	10, // 5: icbt.rpc.v1.EarmarksListResponse.pagination:type_name -> icbt.rpc.v1.PaginationResult
	6,  // [6:6] is the sub-list for method output_type
	6,  // [6:6] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_icbt_rpc_v1_earmark_proto_init() }
func file_icbt_rpc_v1_earmark_proto_init() {
	if File_icbt_rpc_v1_earmark_proto != nil {
		return
	}
	file_icbt_rpc_v1_constraints_proto_init()
	file_icbt_rpc_v1_pagination_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_icbt_rpc_v1_earmark_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Earmark); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarkCreateRequest); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarkCreateResponse); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarkRemoveRequest); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarkGetDetailsRequest); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarkGetDetailsResponse); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarksListRequest); i {
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
		file_icbt_rpc_v1_earmark_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*EarmarksListResponse); i {
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
	file_icbt_rpc_v1_earmark_proto_msgTypes[6].OneofWrappers = []any{}
	file_icbt_rpc_v1_earmark_proto_msgTypes[7].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_icbt_rpc_v1_earmark_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_icbt_rpc_v1_earmark_proto_goTypes,
		DependencyIndexes: file_icbt_rpc_v1_earmark_proto_depIdxs,
		MessageInfos:      file_icbt_rpc_v1_earmark_proto_msgTypes,
	}.Build()
	File_icbt_rpc_v1_earmark_proto = out.File
	file_icbt_rpc_v1_earmark_proto_rawDesc = nil
	file_icbt_rpc_v1_earmark_proto_goTypes = nil
	file_icbt_rpc_v1_earmark_proto_depIdxs = nil
}
