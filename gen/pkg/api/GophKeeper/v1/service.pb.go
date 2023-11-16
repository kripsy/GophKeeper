// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: pkg/api/GophKeeper/v1/service.proto

package GophKeeper

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type AuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *AuthRequest) Reset() {
	*x = AuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRequest) ProtoMessage() {}

func (x *AuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRequest.ProtoReflect.Descriptor instead.
func (*AuthRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *AuthRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AuthRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AuthResponse) Reset() {
	*x = AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthResponse) ProtoMessage() {}

func (x *AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthResponse.ProtoReflect.Descriptor instead.
func (*AuthResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *AuthResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type MultipartUploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content  []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Guid     string `protobuf:"bytes,3,opt,name=guid,proto3" json:"guid,omitempty"`
	Hash     string `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *MultipartUploadFileRequest) Reset() {
	*x = MultipartUploadFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MultipartUploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MultipartUploadFileRequest) ProtoMessage() {}

func (x *MultipartUploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MultipartUploadFileRequest.ProtoReflect.Descriptor instead.
func (*MultipartUploadFileRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *MultipartUploadFileRequest) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *MultipartUploadFileRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *MultipartUploadFileRequest) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

func (x *MultipartUploadFileRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type MultipartUploadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=fileId,proto3" json:"fileId,omitempty"`
}

func (x *MultipartUploadFileResponse) Reset() {
	*x = MultipartUploadFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MultipartUploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MultipartUploadFileResponse) ProtoMessage() {}

func (x *MultipartUploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MultipartUploadFileResponse.ProtoReflect.Descriptor instead.
func (*MultipartUploadFileResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *MultipartUploadFileResponse) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

type MultipartDownloadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Hash     string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Guid     string `protobuf:"bytes,3,opt,name=guid,proto3" json:"guid,omitempty"`
}

func (x *MultipartDownloadFileRequest) Reset() {
	*x = MultipartDownloadFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MultipartDownloadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MultipartDownloadFileRequest) ProtoMessage() {}

func (x *MultipartDownloadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MultipartDownloadFileRequest.ProtoReflect.Descriptor instead.
func (*MultipartDownloadFileRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *MultipartDownloadFileRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *MultipartDownloadFileRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *MultipartDownloadFileRequest) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

type MultipartDownloadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content  []byte `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	Hash     string `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Guid     string `protobuf:"bytes,4,opt,name=guid,proto3" json:"guid,omitempty"`
}

func (x *MultipartDownloadFileResponse) Reset() {
	*x = MultipartDownloadFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MultipartDownloadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MultipartDownloadFileResponse) ProtoMessage() {}

func (x *MultipartDownloadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MultipartDownloadFileResponse.ProtoReflect.Descriptor instead.
func (*MultipartDownloadFileResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *MultipartDownloadFileResponse) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *MultipartDownloadFileResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *MultipartDownloadFileResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *MultipartDownloadFileResponse) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

type BlockStoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid     string `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
	IsFinish bool   `protobuf:"varint,2,opt,name=is_finish,json=isFinish,proto3" json:"is_finish,omitempty"`
}

func (x *BlockStoreRequest) Reset() {
	*x = BlockStoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockStoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockStoreRequest) ProtoMessage() {}

func (x *BlockStoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockStoreRequest.ProtoReflect.Descriptor instead.
func (*BlockStoreRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *BlockStoreRequest) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

func (x *BlockStoreRequest) GetIsFinish() bool {
	if x != nil {
		return x.IsFinish
	}
	return false
}

type BlockStoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid string `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
}

func (x *BlockStoreResponse) Reset() {
	*x = BlockStoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BlockStoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BlockStoreResponse) ProtoMessage() {}

func (x *BlockStoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BlockStoreResponse.ProtoReflect.Descriptor instead.
func (*BlockStoreResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{7}
}

func (x *BlockStoreResponse) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

type ApplyChangesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid string `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
}

func (x *ApplyChangesRequest) Reset() {
	*x = ApplyChangesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplyChangesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyChangesRequest) ProtoMessage() {}

func (x *ApplyChangesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyChangesRequest.ProtoReflect.Descriptor instead.
func (*ApplyChangesRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{8}
}

func (x *ApplyChangesRequest) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

type ApplyChangesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid string `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
}

func (x *ApplyChangesResponse) Reset() {
	*x = ApplyChangesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApplyChangesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApplyChangesResponse) ProtoMessage() {}

func (x *ApplyChangesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_GophKeeper_v1_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApplyChangesResponse.ProtoReflect.Descriptor instead.
func (*ApplyChangesResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP(), []int{9}
}

func (x *ApplyChangesResponse) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

var File_pkg_api_GophKeeper_v1_service_proto protoreflect.FileDescriptor

var file_pkg_api_GophKeeper_v1_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x47, 0x6f, 0x70, 0x68, 0x4b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67,
	0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x57, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48,
	0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x2d, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1d, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x96,
	0x01, 0x0a, 0x1a, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0x35, 0x0a, 0x1b, 0x4d, 0x75, 0x6c, 0x74, 0x69,
	0x70, 0x61, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x7e,
	0x0a, 0x1c, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24,
	0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x12, 0x1b, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x22, 0x99,
	0x01, 0x0a, 0x1d, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x44, 0x6f, 0x77, 0x6e,
	0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba,
	0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x1b, 0x0a,
	0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x22, 0x4d, 0x0a, 0x11, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba,
	0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09,
	0x69, 0x73, 0x5f, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x08, 0x69, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x22, 0x31, 0x0a, 0x12, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1b, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba,
	0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x22, 0x32, 0x0a, 0x13,
	0x41, 0x70, 0x70, 0x6c, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64,
	0x22, 0x33, 0x0a, 0x14, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52,
	0x04, 0x67, 0x75, 0x69, 0x64, 0x32, 0xc9, 0x05, 0x0a, 0x11, 0x47, 0x6f, 0x70, 0x68, 0x4b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x08, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x22, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x6b,
	0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x50, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x22, 0x2e, 0x70, 0x6b, 0x67, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x7e, 0x0a, 0x13, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x31, 0x2e, 0x70, 0x6b, 0x67, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x32, 0x2e, 0x70,
	0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x55, 0x70,
	0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x28, 0x01, 0x12, 0x65, 0x0a, 0x0a, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x28, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x70, 0x6b, 0x67,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x36, 0x0a, 0x04, 0x50, 0x69, 0x6e,
	0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x12, 0x84, 0x01, 0x0a, 0x15, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x44,
	0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x33, 0x2e, 0x70, 0x6b,
	0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61, 0x72, 0x74, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x34, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b,
	0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x61,
	0x72, 0x74, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x67, 0x0a, 0x0c, 0x41, 0x70, 0x70, 0x6c,
	0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x2a, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x70, 0x70, 0x6c, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x67,
	0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70,
	0x6c, 0x79, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6b, 0x72, 0x69, 0x70, 0x73, 0x79, 0x2f, 0x47, 0x6f, 0x70, 0x68, 0x4b, 0x65, 0x65, 0x70, 0x65,
	0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x47, 0x6f, 0x70, 0x68, 0x4b, 0x65,
	0x65, 0x70, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_api_GophKeeper_v1_service_proto_rawDescOnce sync.Once
	file_pkg_api_GophKeeper_v1_service_proto_rawDescData = file_pkg_api_GophKeeper_v1_service_proto_rawDesc
)

func file_pkg_api_GophKeeper_v1_service_proto_rawDescGZIP() []byte {
	file_pkg_api_GophKeeper_v1_service_proto_rawDescOnce.Do(func() {
		file_pkg_api_GophKeeper_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_GophKeeper_v1_service_proto_rawDescData)
	})
	return file_pkg_api_GophKeeper_v1_service_proto_rawDescData
}

var file_pkg_api_GophKeeper_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_pkg_api_GophKeeper_v1_service_proto_goTypes = []interface{}{
	(*AuthRequest)(nil),                   // 0: pkg.api.gophkeeper.v1.AuthRequest
	(*AuthResponse)(nil),                  // 1: pkg.api.gophkeeper.v1.AuthResponse
	(*MultipartUploadFileRequest)(nil),    // 2: pkg.api.gophkeeper.v1.MultipartUploadFileRequest
	(*MultipartUploadFileResponse)(nil),   // 3: pkg.api.gophkeeper.v1.MultipartUploadFileResponse
	(*MultipartDownloadFileRequest)(nil),  // 4: pkg.api.gophkeeper.v1.MultipartDownloadFileRequest
	(*MultipartDownloadFileResponse)(nil), // 5: pkg.api.gophkeeper.v1.MultipartDownloadFileResponse
	(*BlockStoreRequest)(nil),             // 6: pkg.api.gophkeeper.v1.BlockStoreRequest
	(*BlockStoreResponse)(nil),            // 7: pkg.api.gophkeeper.v1.BlockStoreResponse
	(*ApplyChangesRequest)(nil),           // 8: pkg.api.gophkeeper.v1.ApplyChangesRequest
	(*ApplyChangesResponse)(nil),          // 9: pkg.api.gophkeeper.v1.ApplyChangesResponse
	(*emptypb.Empty)(nil),                 // 10: google.protobuf.Empty
}
var file_pkg_api_GophKeeper_v1_service_proto_depIdxs = []int32{
	0,  // 0: pkg.api.gophkeeper.v1.GophKeeperService.Register:input_type -> pkg.api.gophkeeper.v1.AuthRequest
	0,  // 1: pkg.api.gophkeeper.v1.GophKeeperService.Login:input_type -> pkg.api.gophkeeper.v1.AuthRequest
	2,  // 2: pkg.api.gophkeeper.v1.GophKeeperService.MultipartUploadFile:input_type -> pkg.api.gophkeeper.v1.MultipartUploadFileRequest
	6,  // 3: pkg.api.gophkeeper.v1.GophKeeperService.BlockStore:input_type -> pkg.api.gophkeeper.v1.BlockStoreRequest
	10, // 4: pkg.api.gophkeeper.v1.GophKeeperService.Ping:input_type -> google.protobuf.Empty
	4,  // 5: pkg.api.gophkeeper.v1.GophKeeperService.MultipartDownloadFile:input_type -> pkg.api.gophkeeper.v1.MultipartDownloadFileRequest
	8,  // 6: pkg.api.gophkeeper.v1.GophKeeperService.ApplyChanges:input_type -> pkg.api.gophkeeper.v1.ApplyChangesRequest
	1,  // 7: pkg.api.gophkeeper.v1.GophKeeperService.Register:output_type -> pkg.api.gophkeeper.v1.AuthResponse
	1,  // 8: pkg.api.gophkeeper.v1.GophKeeperService.Login:output_type -> pkg.api.gophkeeper.v1.AuthResponse
	3,  // 9: pkg.api.gophkeeper.v1.GophKeeperService.MultipartUploadFile:output_type -> pkg.api.gophkeeper.v1.MultipartUploadFileResponse
	7,  // 10: pkg.api.gophkeeper.v1.GophKeeperService.BlockStore:output_type -> pkg.api.gophkeeper.v1.BlockStoreResponse
	10, // 11: pkg.api.gophkeeper.v1.GophKeeperService.Ping:output_type -> google.protobuf.Empty
	5,  // 12: pkg.api.gophkeeper.v1.GophKeeperService.MultipartDownloadFile:output_type -> pkg.api.gophkeeper.v1.MultipartDownloadFileResponse
	9,  // 13: pkg.api.gophkeeper.v1.GophKeeperService.ApplyChanges:output_type -> pkg.api.gophkeeper.v1.ApplyChangesResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_api_GophKeeper_v1_service_proto_init() }
func file_pkg_api_GophKeeper_v1_service_proto_init() {
	if File_pkg_api_GophKeeper_v1_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRequest); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthResponse); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MultipartUploadFileRequest); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MultipartUploadFileResponse); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MultipartDownloadFileRequest); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MultipartDownloadFileResponse); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockStoreRequest); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BlockStoreResponse); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplyChangesRequest); i {
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
		file_pkg_api_GophKeeper_v1_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApplyChangesResponse); i {
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
			RawDescriptor: file_pkg_api_GophKeeper_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_api_GophKeeper_v1_service_proto_goTypes,
		DependencyIndexes: file_pkg_api_GophKeeper_v1_service_proto_depIdxs,
		MessageInfos:      file_pkg_api_GophKeeper_v1_service_proto_msgTypes,
	}.Build()
	File_pkg_api_GophKeeper_v1_service_proto = out.File
	file_pkg_api_GophKeeper_v1_service_proto_rawDesc = nil
	file_pkg_api_GophKeeper_v1_service_proto_goTypes = nil
	file_pkg_api_GophKeeper_v1_service_proto_depIdxs = nil
}
