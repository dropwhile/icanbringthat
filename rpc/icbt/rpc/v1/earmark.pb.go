// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: icbt/rpc/v1/earmark.proto

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

type Earmark struct {
	state                     protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_RefId          string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId" json:"ref_id,omitempty"`
	xxx_hidden_EventItemRefId string                 `protobuf:"bytes,2,opt,name=event_item_ref_id,json=eventItemRefId" json:"event_item_ref_id,omitempty"`
	xxx_hidden_Note           string                 `protobuf:"bytes,3,opt,name=note" json:"note,omitempty"`
	xxx_hidden_Owner          string                 `protobuf:"bytes,4,opt,name=owner" json:"owner,omitempty"`
	xxx_hidden_Created        *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created" json:"created,omitempty"`
	unknownFields             protoimpl.UnknownFields
	sizeCache                 protoimpl.SizeCache
}

func (x *Earmark) Reset() {
	*x = Earmark{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Earmark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Earmark) ProtoMessage() {}

func (x *Earmark) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *Earmark) GetRefId() string {
	if x != nil {
		return x.xxx_hidden_RefId
	}
	return ""
}

func (x *Earmark) GetEventItemRefId() string {
	if x != nil {
		return x.xxx_hidden_EventItemRefId
	}
	return ""
}

func (x *Earmark) GetNote() string {
	if x != nil {
		return x.xxx_hidden_Note
	}
	return ""
}

func (x *Earmark) GetOwner() string {
	if x != nil {
		return x.xxx_hidden_Owner
	}
	return ""
}

func (x *Earmark) GetCreated() *timestamppb.Timestamp {
	if x != nil {
		return x.xxx_hidden_Created
	}
	return nil
}

func (x *Earmark) SetRefId(v string) {
	x.xxx_hidden_RefId = v
}

func (x *Earmark) SetEventItemRefId(v string) {
	x.xxx_hidden_EventItemRefId = v
}

func (x *Earmark) SetNote(v string) {
	x.xxx_hidden_Note = v
}

func (x *Earmark) SetOwner(v string) {
	x.xxx_hidden_Owner = v
}

func (x *Earmark) SetCreated(v *timestamppb.Timestamp) {
	x.xxx_hidden_Created = v
}

func (x *Earmark) HasCreated() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Created != nil
}

func (x *Earmark) ClearCreated() {
	x.xxx_hidden_Created = nil
}

type Earmark_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	RefId          string
	EventItemRefId string
	Note           string
	Owner          string
	Created        *timestamppb.Timestamp
}

func (b0 Earmark_builder) Build() *Earmark {
	m0 := &Earmark{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_RefId = b.RefId
	x.xxx_hidden_EventItemRefId = b.EventItemRefId
	x.xxx_hidden_Note = b.Note
	x.xxx_hidden_Owner = b.Owner
	x.xxx_hidden_Created = b.Created
	return m0
}

type EarmarkCreateRequest struct {
	state                     protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_EventItemRefId string                 `protobuf:"bytes,1,opt,name=event_item_ref_id,json=eventItemRefId" json:"event_item_ref_id,omitempty"`
	xxx_hidden_Note           string                 `protobuf:"bytes,2,opt,name=note" json:"note,omitempty"`
	unknownFields             protoimpl.UnknownFields
	sizeCache                 protoimpl.SizeCache
}

func (x *EarmarkCreateRequest) Reset() {
	*x = EarmarkCreateRequest{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarkCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkCreateRequest) ProtoMessage() {}

func (x *EarmarkCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarkCreateRequest) GetEventItemRefId() string {
	if x != nil {
		return x.xxx_hidden_EventItemRefId
	}
	return ""
}

func (x *EarmarkCreateRequest) GetNote() string {
	if x != nil {
		return x.xxx_hidden_Note
	}
	return ""
}

func (x *EarmarkCreateRequest) SetEventItemRefId(v string) {
	x.xxx_hidden_EventItemRefId = v
}

func (x *EarmarkCreateRequest) SetNote(v string) {
	x.xxx_hidden_Note = v
}

type EarmarkCreateRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	EventItemRefId string
	Note           string
}

func (b0 EarmarkCreateRequest_builder) Build() *EarmarkCreateRequest {
	m0 := &EarmarkCreateRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_EventItemRefId = b.EventItemRefId
	x.xxx_hidden_Note = b.Note
	return m0
}

type EarmarkCreateResponse struct {
	state              protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Earmark *Earmark               `protobuf:"bytes,1,opt,name=earmark" json:"earmark,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *EarmarkCreateResponse) Reset() {
	*x = EarmarkCreateResponse{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarkCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkCreateResponse) ProtoMessage() {}

func (x *EarmarkCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarkCreateResponse) GetEarmark() *Earmark {
	if x != nil {
		return x.xxx_hidden_Earmark
	}
	return nil
}

func (x *EarmarkCreateResponse) SetEarmark(v *Earmark) {
	x.xxx_hidden_Earmark = v
}

func (x *EarmarkCreateResponse) HasEarmark() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Earmark != nil
}

func (x *EarmarkCreateResponse) ClearEarmark() {
	x.xxx_hidden_Earmark = nil
}

type EarmarkCreateResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Earmark *Earmark
}

func (b0 EarmarkCreateResponse_builder) Build() *EarmarkCreateResponse {
	m0 := &EarmarkCreateResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Earmark = b.Earmark
	return m0
}

type EarmarkRemoveRequest struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_RefId string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId" json:"ref_id,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *EarmarkRemoveRequest) Reset() {
	*x = EarmarkRemoveRequest{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarkRemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkRemoveRequest) ProtoMessage() {}

func (x *EarmarkRemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarkRemoveRequest) GetRefId() string {
	if x != nil {
		return x.xxx_hidden_RefId
	}
	return ""
}

func (x *EarmarkRemoveRequest) SetRefId(v string) {
	x.xxx_hidden_RefId = v
}

type EarmarkRemoveRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	RefId string
}

func (b0 EarmarkRemoveRequest_builder) Build() *EarmarkRemoveRequest {
	m0 := &EarmarkRemoveRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_RefId = b.RefId
	return m0
}

type EarmarkGetDetailsRequest struct {
	state            protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_RefId string                 `protobuf:"bytes,1,opt,name=ref_id,json=refId" json:"ref_id,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *EarmarkGetDetailsRequest) Reset() {
	*x = EarmarkGetDetailsRequest{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarkGetDetailsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkGetDetailsRequest) ProtoMessage() {}

func (x *EarmarkGetDetailsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarkGetDetailsRequest) GetRefId() string {
	if x != nil {
		return x.xxx_hidden_RefId
	}
	return ""
}

func (x *EarmarkGetDetailsRequest) SetRefId(v string) {
	x.xxx_hidden_RefId = v
}

type EarmarkGetDetailsRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	RefId string
}

func (b0 EarmarkGetDetailsRequest_builder) Build() *EarmarkGetDetailsRequest {
	m0 := &EarmarkGetDetailsRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_RefId = b.RefId
	return m0
}

type EarmarkGetDetailsResponse struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Earmark    *Earmark               `protobuf:"bytes,1,opt,name=earmark" json:"earmark,omitempty"`
	xxx_hidden_EventRefId string                 `protobuf:"bytes,2,opt,name=event_ref_id,json=eventRefId" json:"event_ref_id,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *EarmarkGetDetailsResponse) Reset() {
	*x = EarmarkGetDetailsResponse{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarkGetDetailsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarkGetDetailsResponse) ProtoMessage() {}

func (x *EarmarkGetDetailsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarkGetDetailsResponse) GetEarmark() *Earmark {
	if x != nil {
		return x.xxx_hidden_Earmark
	}
	return nil
}

func (x *EarmarkGetDetailsResponse) GetEventRefId() string {
	if x != nil {
		return x.xxx_hidden_EventRefId
	}
	return ""
}

func (x *EarmarkGetDetailsResponse) SetEarmark(v *Earmark) {
	x.xxx_hidden_Earmark = v
}

func (x *EarmarkGetDetailsResponse) SetEventRefId(v string) {
	x.xxx_hidden_EventRefId = v
}

func (x *EarmarkGetDetailsResponse) HasEarmark() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Earmark != nil
}

func (x *EarmarkGetDetailsResponse) ClearEarmark() {
	x.xxx_hidden_Earmark = nil
}

type EarmarkGetDetailsResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Earmark    *Earmark
	EventRefId string
}

func (b0 EarmarkGetDetailsResponse_builder) Build() *EarmarkGetDetailsResponse {
	m0 := &EarmarkGetDetailsResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Earmark = b.Earmark
	x.xxx_hidden_EventRefId = b.EventRefId
	return m0
}

type EarmarksListRequest struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Pagination  *PaginationRequest     `protobuf:"bytes,1,opt,name=pagination" json:"pagination,omitempty"`
	xxx_hidden_Archived    bool                   `protobuf:"varint,2,opt,name=archived" json:"archived,omitempty"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *EarmarksListRequest) Reset() {
	*x = EarmarksListRequest{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarksListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarksListRequest) ProtoMessage() {}

func (x *EarmarksListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarksListRequest) GetPagination() *PaginationRequest {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *EarmarksListRequest) GetArchived() bool {
	if x != nil {
		return x.xxx_hidden_Archived
	}
	return false
}

func (x *EarmarksListRequest) SetPagination(v *PaginationRequest) {
	x.xxx_hidden_Pagination = v
}

func (x *EarmarksListRequest) SetArchived(v bool) {
	x.xxx_hidden_Archived = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 1, 2)
}

func (x *EarmarksListRequest) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *EarmarksListRequest) HasArchived() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 1)
}

func (x *EarmarksListRequest) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

func (x *EarmarksListRequest) ClearArchived() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 1)
	x.xxx_hidden_Archived = false
}

type EarmarksListRequest_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Pagination *PaginationRequest
	Archived   *bool
}

func (b0 EarmarksListRequest_builder) Build() *EarmarksListRequest {
	m0 := &EarmarksListRequest{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Pagination = b.Pagination
	if b.Archived != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 1, 2)
		x.xxx_hidden_Archived = *b.Archived
	}
	return m0
}

type EarmarksListResponse struct {
	state                 protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Earmarks   *[]*Earmark            `protobuf:"bytes,1,rep,name=earmarks" json:"earmarks,omitempty"`
	xxx_hidden_Pagination *PaginationResult      `protobuf:"bytes,2,opt,name=pagination" json:"pagination,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *EarmarksListResponse) Reset() {
	*x = EarmarksListResponse{}
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EarmarksListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EarmarksListResponse) ProtoMessage() {}

func (x *EarmarksListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_icbt_rpc_v1_earmark_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *EarmarksListResponse) GetEarmarks() []*Earmark {
	if x != nil {
		if x.xxx_hidden_Earmarks != nil {
			return *x.xxx_hidden_Earmarks
		}
	}
	return nil
}

func (x *EarmarksListResponse) GetPagination() *PaginationResult {
	if x != nil {
		return x.xxx_hidden_Pagination
	}
	return nil
}

func (x *EarmarksListResponse) SetEarmarks(v []*Earmark) {
	x.xxx_hidden_Earmarks = &v
}

func (x *EarmarksListResponse) SetPagination(v *PaginationResult) {
	x.xxx_hidden_Pagination = v
}

func (x *EarmarksListResponse) HasPagination() bool {
	if x == nil {
		return false
	}
	return x.xxx_hidden_Pagination != nil
}

func (x *EarmarksListResponse) ClearPagination() {
	x.xxx_hidden_Pagination = nil
}

type EarmarksListResponse_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Earmarks   []*Earmark
	Pagination *PaginationResult
}

func (b0 EarmarksListResponse_builder) Build() *EarmarksListResponse {
	m0 := &EarmarksListResponse{}
	b, x := &b0, m0
	_, _ = b, x
	x.xxx_hidden_Earmarks = &b.Earmarks
	x.xxx_hidden_Pagination = b.Pagination
	return m0
}

var File_icbt_rpc_v1_earmark_proto protoreflect.FileDescriptor

var file_icbt_rpc_v1_earmark_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x69, 0x63, 0x62,
	0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x6f, 0x5f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x69, 0x63, 0x62, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e,
	0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x69, 0x63, 0x62, 0x74, 0x2f, 0x72,
	0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xab, 0x01, 0x0a, 0x07, 0x45, 0x61, 0x72, 0x6d, 0x61,
	0x72, 0x6b, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x11, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x66, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65,
	0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x34,
	0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x22, 0x62, 0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x11,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06, 0x88, 0xd8,
	0x83, 0x8b, 0x02, 0x01, 0x52, 0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x66, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x22, 0x47, 0x0a, 0x15, 0x45, 0x61, 0x72, 0x6d,
	0x61, 0x72, 0x6b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2e, 0x0a, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72,
	0x6b, 0x22, 0x3a, 0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06, 0x72, 0x65, 0x66,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06,
	0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x22, 0x3e, 0x0a,
	0x18, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x06, 0x72, 0x65, 0x66,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xba, 0x48, 0x08, 0x72, 0x06,
	0x88, 0xd8, 0x83, 0x8b, 0x02, 0x01, 0x52, 0x05, 0x72, 0x65, 0x66, 0x49, 0x64, 0x22, 0x6d, 0x0a,
	0x19, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x65, 0x61,
	0x72, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x63,
	0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72,
	0x6b, 0x52, 0x07, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x20, 0x0a, 0x0c, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x66, 0x49, 0x64, 0x22, 0x7f, 0x0a, 0x13,
	0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x05, 0xaa, 0x01, 0x02, 0x08, 0x01, 0x52, 0x0a,
	0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x08, 0x61, 0x72,
	0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x42, 0x05, 0xaa, 0x01,
	0x02, 0x08, 0x01, 0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x22, 0x8e, 0x01,
	0x0a, 0x14, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x65, 0x61, 0x72, 0x6d, 0x61, 0x72,
	0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x52, 0x08,
	0x65, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x73, 0x12, 0x44, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69,
	0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x05, 0xaa, 0x01, 0x02,
	0x08, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0xb1,
	0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x63, 0x62, 0x74, 0x2e, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x42, 0x0c, 0x45, 0x61, 0x72, 0x6d, 0x61, 0x72, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64,
	0x72, 0x6f, 0x70, 0x77, 0x68, 0x69, 0x6c, 0x65, 0x2f, 0x69, 0x63, 0x61, 0x6e, 0x62, 0x72, 0x69,
	0x6e, 0x67, 0x74, 0x68, 0x61, 0x74, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x69, 0x63, 0x62, 0x74, 0x2f,
	0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x70, 0x63, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x49,
	0x52, 0x58, 0xaa, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x2e, 0x52, 0x70, 0x63, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x0b, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0xe2, 0x02,
	0x17, 0x49, 0x63, 0x62, 0x74, 0x5c, 0x52, 0x70, 0x63, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0d, 0x49, 0x63, 0x62, 0x74, 0x3a,
	0x3a, 0x52, 0x70, 0x63, 0x3a, 0x3a, 0x56, 0x31, 0x92, 0x03, 0x07, 0xd2, 0x3e, 0x02, 0x10, 0x03,
	0x08, 0x02, 0x62, 0x08, 0x65, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x70, 0xe8, 0x07,
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
