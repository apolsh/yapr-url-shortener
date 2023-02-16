// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.15.8
// source: url_shortener.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PingDBResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsAlive bool `protobuf:"varint,1,opt,name=isAlive,proto3" json:"isAlive,omitempty"`
}

func (x *PingDBResponse) Reset() {
	*x = PingDBResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingDBResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingDBResponse) ProtoMessage() {}

func (x *PingDBResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingDBResponse.ProtoReflect.Descriptor instead.
func (*PingDBResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{0}
}

func (x *PingDBResponse) GetIsAlive() bool {
	if x != nil {
		return x.IsAlive
	}
	return false
}

type GetShortenURLByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlID string `protobuf:"bytes,1,opt,name=urlID,proto3" json:"urlID,omitempty"`
}

func (x *GetShortenURLByIDRequest) Reset() {
	*x = GetShortenURLByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShortenURLByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShortenURLByIDRequest) ProtoMessage() {}

func (x *GetShortenURLByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShortenURLByIDRequest.ProtoReflect.Descriptor instead.
func (*GetShortenURLByIDRequest) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{1}
}

func (x *GetShortenURLByIDRequest) GetUrlID() string {
	if x != nil {
		return x.UrlID
	}
	return ""
}

type GetShortenURLByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalURL string `protobuf:"bytes,1,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *GetShortenURLByIDResponse) Reset() {
	*x = GetShortenURLByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShortenURLByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShortenURLByIDResponse) ProtoMessage() {}

func (x *GetShortenURLByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShortenURLByIDResponse.ProtoReflect.Descriptor instead.
func (*GetShortenURLByIDResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{2}
}

func (x *GetShortenURLByIDResponse) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type URLPair struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortURL    string `protobuf:"bytes,1,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
	OriginalURL string `protobuf:"bytes,2,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *URLPair) Reset() {
	*x = URLPair{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *URLPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*URLPair) ProtoMessage() {}

func (x *URLPair) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use URLPair.ProtoReflect.Descriptor instead.
func (*URLPair) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{3}
}

func (x *URLPair) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

func (x *URLPair) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type GetShortenURLsByUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlPairs []*URLPair `protobuf:"bytes,1,rep,name=urlPairs,proto3" json:"urlPairs,omitempty"`
}

func (x *GetShortenURLsByUserResponse) Reset() {
	*x = GetShortenURLsByUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShortenURLsByUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShortenURLsByUserResponse) ProtoMessage() {}

func (x *GetShortenURLsByUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShortenURLsByUserResponse.ProtoReflect.Descriptor instead.
func (*GetShortenURLsByUserResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{4}
}

func (x *GetShortenURLsByUserResponse) GetUrlPairs() []*URLPair {
	if x != nil {
		return x.UrlPairs
	}
	return nil
}

type SaveShortenURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OriginalURL string `protobuf:"bytes,1,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *SaveShortenURLRequest) Reset() {
	*x = SaveShortenURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLRequest) ProtoMessage() {}

func (x *SaveShortenURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLRequest.ProtoReflect.Descriptor instead.
func (*SaveShortenURLRequest) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{5}
}

func (x *SaveShortenURLRequest) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type SaveShortenURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortenedURL string `protobuf:"bytes,1,opt,name=shortenedURL,proto3" json:"shortenedURL,omitempty"`
}

func (x *SaveShortenURLResponse) Reset() {
	*x = SaveShortenURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLResponse) ProtoMessage() {}

func (x *SaveShortenURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLResponse.ProtoReflect.Descriptor instead.
func (*SaveShortenURLResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{6}
}

func (x *SaveShortenURLResponse) GetShortenedURL() string {
	if x != nil {
		return x.ShortenedURL
	}
	return ""
}

type ShortenInBatchRequestItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationID string `protobuf:"bytes,1,opt,name=correlationID,proto3" json:"correlationID,omitempty"`
	OriginalURL   string `protobuf:"bytes,2,opt,name=originalURL,proto3" json:"originalURL,omitempty"`
}

func (x *ShortenInBatchRequestItem) Reset() {
	*x = ShortenInBatchRequestItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenInBatchRequestItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenInBatchRequestItem) ProtoMessage() {}

func (x *ShortenInBatchRequestItem) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenInBatchRequestItem.ProtoReflect.Descriptor instead.
func (*ShortenInBatchRequestItem) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{7}
}

func (x *ShortenInBatchRequestItem) GetCorrelationID() string {
	if x != nil {
		return x.CorrelationID
	}
	return ""
}

func (x *ShortenInBatchRequestItem) GetOriginalURL() string {
	if x != nil {
		return x.OriginalURL
	}
	return ""
}

type ShortenInBatchResponseItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationID string `protobuf:"bytes,1,opt,name=correlationID,proto3" json:"correlationID,omitempty"`
	ShortURL      string `protobuf:"bytes,2,opt,name=shortURL,proto3" json:"shortURL,omitempty"`
}

func (x *ShortenInBatchResponseItem) Reset() {
	*x = ShortenInBatchResponseItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenInBatchResponseItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenInBatchResponseItem) ProtoMessage() {}

func (x *ShortenInBatchResponseItem) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenInBatchResponseItem.ProtoReflect.Descriptor instead.
func (*ShortenInBatchResponseItem) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{8}
}

func (x *ShortenInBatchResponseItem) GetCorrelationID() string {
	if x != nil {
		return x.CorrelationID
	}
	return ""
}

func (x *ShortenInBatchResponseItem) GetShortURL() string {
	if x != nil {
		return x.ShortURL
	}
	return ""
}

type SaveShortenURLsInBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ShortenInBatchRequestItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SaveShortenURLsInBatchRequest) Reset() {
	*x = SaveShortenURLsInBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLsInBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLsInBatchRequest) ProtoMessage() {}

func (x *SaveShortenURLsInBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLsInBatchRequest.ProtoReflect.Descriptor instead.
func (*SaveShortenURLsInBatchRequest) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{9}
}

func (x *SaveShortenURLsInBatchRequest) GetItems() []*ShortenInBatchRequestItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type SaveShortenURLsInBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*ShortenInBatchResponseItem `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *SaveShortenURLsInBatchResponse) Reset() {
	*x = SaveShortenURLsInBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveShortenURLsInBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveShortenURLsInBatchResponse) ProtoMessage() {}

func (x *SaveShortenURLsInBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveShortenURLsInBatchResponse.ProtoReflect.Descriptor instead.
func (*SaveShortenURLsInBatchResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{10}
}

func (x *SaveShortenURLsInBatchResponse) GetItems() []*ShortenInBatchResponseItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type DeleteShortenURLsInBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []string `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *DeleteShortenURLsInBatchRequest) Reset() {
	*x = DeleteShortenURLsInBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteShortenURLsInBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteShortenURLsInBatchRequest) ProtoMessage() {}

func (x *DeleteShortenURLsInBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteShortenURLsInBatchRequest.ProtoReflect.Descriptor instead.
func (*DeleteShortenURLsInBatchRequest) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{11}
}

func (x *DeleteShortenURLsInBatchRequest) GetItems() []string {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetAppStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	URLs  int64 `protobuf:"varint,1,opt,name=URLs,proto3" json:"URLs,omitempty"`
	Users int64 `protobuf:"varint,2,opt,name=Users,proto3" json:"Users,omitempty"`
}

func (x *GetAppStatsResponse) Reset() {
	*x = GetAppStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_url_shortener_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAppStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAppStatsResponse) ProtoMessage() {}

func (x *GetAppStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_url_shortener_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAppStatsResponse.ProtoReflect.Descriptor instead.
func (*GetAppStatsResponse) Descriptor() ([]byte, []int) {
	return file_url_shortener_proto_rawDescGZIP(), []int{12}
}

func (x *GetAppStatsResponse) GetURLs() int64 {
	if x != nil {
		return x.URLs
	}
	return 0
}

func (x *GetAppStatsResponse) GetUsers() int64 {
	if x != nil {
		return x.Users
	}
	return 0
}

var File_url_shortener_proto protoreflect.FileDescriptor

var file_url_shortener_proto_rawDesc = []byte{
	0x0a, 0x13, 0x75, 0x72, 0x6c, 0x5f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a, 0x0e, 0x50, 0x69, 0x6e,
	0x67, 0x44, 0x42, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x69,
	0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73,
	0x41, 0x6c, 0x69, 0x76, 0x65, 0x22, 0x30, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x72, 0x6c, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x75, 0x72, 0x6c, 0x49, 0x44, 0x22, 0x3d, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c,
	0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x22, 0x47, 0x0a, 0x07, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x69,
	0x72, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x20, 0x0a,
	0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x22,
	0x4a, 0x0a, 0x1c, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c,
	0x73, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2a, 0x0a, 0x08, 0x75, 0x72, 0x6c, 0x50, 0x61, 0x69, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x52, 0x4c, 0x50, 0x61, 0x69,
	0x72, 0x52, 0x08, 0x75, 0x72, 0x6c, 0x50, 0x61, 0x69, 0x72, 0x73, 0x22, 0x39, 0x0a, 0x15, 0x53,
	0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c,
	0x55, 0x52, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x22, 0x3c, 0x0a, 0x16, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x22, 0x0a, 0x0c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x64, 0x55, 0x52, 0x4c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65,
	0x64, 0x55, 0x52, 0x4c, 0x22, 0x63, 0x0a, 0x19, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x49,
	0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x72,
	0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x22, 0x5e, 0x0a, 0x1a, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x55, 0x52, 0x4c, 0x22, 0x57, 0x0a, 0x1d, 0x53, 0x61, 0x76,
	0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x22, 0x59, 0x0a, 0x1e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x55, 0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x37, 0x0a,
	0x1f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52,
	0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x3f, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x55, 0x52, 0x4c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x55, 0x52, 0x4c,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x32, 0xc9, 0x04, 0x0a, 0x0c, 0x55, 0x52, 0x4c, 0x53,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x06, 0x50, 0x69, 0x6e, 0x67,
	0x44, 0x42, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x44, 0x42, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x56, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55,
	0x52, 0x4c, 0x42, 0x79, 0x49, 0x44, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x49, 0x44,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x42, 0x79, 0x49,
	0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x14, 0x47, 0x65, 0x74,
	0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x42, 0x79, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73,
	0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d,
	0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c,
	0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x0b, 0x47, 0x65, 0x74, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74,
	0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x65, 0x0a, 0x16, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55,
	0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x24, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52,
	0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x25, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x53, 0x68, 0x6f,
	0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x55, 0x52, 0x4c, 0x73, 0x49, 0x6e, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x70, 0x6f, 0x6c, 0x73, 0x68, 0x2f, 0x79, 0x61, 0x70, 0x72, 0x2d, 0x75, 0x72,
	0x6c, 0x2d, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_url_shortener_proto_rawDescOnce sync.Once
	file_url_shortener_proto_rawDescData = file_url_shortener_proto_rawDesc
)

func file_url_shortener_proto_rawDescGZIP() []byte {
	file_url_shortener_proto_rawDescOnce.Do(func() {
		file_url_shortener_proto_rawDescData = protoimpl.X.CompressGZIP(file_url_shortener_proto_rawDescData)
	})
	return file_url_shortener_proto_rawDescData
}

var file_url_shortener_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_url_shortener_proto_goTypes = []interface{}{
	(*PingDBResponse)(nil),                  // 0: proto.PingDBResponse
	(*GetShortenURLByIDRequest)(nil),        // 1: proto.GetShortenURLByIDRequest
	(*GetShortenURLByIDResponse)(nil),       // 2: proto.GetShortenURLByIDResponse
	(*URLPair)(nil),                         // 3: proto.URLPair
	(*GetShortenURLsByUserResponse)(nil),    // 4: proto.GetShortenURLsByUserResponse
	(*SaveShortenURLRequest)(nil),           // 5: proto.SaveShortenURLRequest
	(*SaveShortenURLResponse)(nil),          // 6: proto.SaveShortenURLResponse
	(*ShortenInBatchRequestItem)(nil),       // 7: proto.ShortenInBatchRequestItem
	(*ShortenInBatchResponseItem)(nil),      // 8: proto.ShortenInBatchResponseItem
	(*SaveShortenURLsInBatchRequest)(nil),   // 9: proto.SaveShortenURLsInBatchRequest
	(*SaveShortenURLsInBatchResponse)(nil),  // 10: proto.SaveShortenURLsInBatchResponse
	(*DeleteShortenURLsInBatchRequest)(nil), // 11: proto.DeleteShortenURLsInBatchRequest
	(*GetAppStatsResponse)(nil),             // 12: proto.GetAppStatsResponse
	(*emptypb.Empty)(nil),                   // 13: google.protobuf.Empty
}
var file_url_shortener_proto_depIdxs = []int32{
	3,  // 0: proto.GetShortenURLsByUserResponse.urlPairs:type_name -> proto.URLPair
	7,  // 1: proto.SaveShortenURLsInBatchRequest.items:type_name -> proto.ShortenInBatchRequestItem
	8,  // 2: proto.SaveShortenURLsInBatchResponse.items:type_name -> proto.ShortenInBatchResponseItem
	13, // 3: proto.URLShortener.PingDB:input_type -> google.protobuf.Empty
	1,  // 4: proto.URLShortener.GetShortenURLByID:input_type -> proto.GetShortenURLByIDRequest
	13, // 5: proto.URLShortener.GetShortenURLsByUser:input_type -> google.protobuf.Empty
	5,  // 6: proto.URLShortener.SaveShortenURL:input_type -> proto.SaveShortenURLRequest
	13, // 7: proto.URLShortener.GetAppStats:input_type -> google.protobuf.Empty
	9,  // 8: proto.URLShortener.SaveShortenURLsInBatch:input_type -> proto.SaveShortenURLsInBatchRequest
	11, // 9: proto.URLShortener.DeleteShortenURLsInBatch:input_type -> proto.DeleteShortenURLsInBatchRequest
	0,  // 10: proto.URLShortener.PingDB:output_type -> proto.PingDBResponse
	2,  // 11: proto.URLShortener.GetShortenURLByID:output_type -> proto.GetShortenURLByIDResponse
	4,  // 12: proto.URLShortener.GetShortenURLsByUser:output_type -> proto.GetShortenURLsByUserResponse
	6,  // 13: proto.URLShortener.SaveShortenURL:output_type -> proto.SaveShortenURLResponse
	12, // 14: proto.URLShortener.GetAppStats:output_type -> proto.GetAppStatsResponse
	10, // 15: proto.URLShortener.SaveShortenURLsInBatch:output_type -> proto.SaveShortenURLsInBatchResponse
	13, // 16: proto.URLShortener.DeleteShortenURLsInBatch:output_type -> google.protobuf.Empty
	10, // [10:17] is the sub-list for method output_type
	3,  // [3:10] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_url_shortener_proto_init() }
func file_url_shortener_proto_init() {
	if File_url_shortener_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_url_shortener_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingDBResponse); i {
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
		file_url_shortener_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShortenURLByIDRequest); i {
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
		file_url_shortener_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShortenURLByIDResponse); i {
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
		file_url_shortener_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*URLPair); i {
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
		file_url_shortener_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShortenURLsByUserResponse); i {
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
		file_url_shortener_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveShortenURLRequest); i {
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
		file_url_shortener_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveShortenURLResponse); i {
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
		file_url_shortener_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenInBatchRequestItem); i {
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
		file_url_shortener_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShortenInBatchResponseItem); i {
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
		file_url_shortener_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveShortenURLsInBatchRequest); i {
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
		file_url_shortener_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveShortenURLsInBatchResponse); i {
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
		file_url_shortener_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteShortenURLsInBatchRequest); i {
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
		file_url_shortener_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAppStatsResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_url_shortener_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_url_shortener_proto_goTypes,
		DependencyIndexes: file_url_shortener_proto_depIdxs,
		MessageInfos:      file_url_shortener_proto_msgTypes,
	}.Build()
	File_url_shortener_proto = out.File
	file_url_shortener_proto_rawDesc = nil
	file_url_shortener_proto_goTypes = nil
	file_url_shortener_proto_depIdxs = nil
}
