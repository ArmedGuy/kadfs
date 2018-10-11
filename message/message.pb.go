// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RPC struct {
	MessageId            int32    `protobuf:"varint,1,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	SenderId             string   `protobuf:"bytes,2,opt,name=SenderId,proto3" json:"SenderId,omitempty"`
	ReceiverId           string   `protobuf:"bytes,3,opt,name=ReceiverId,proto3" json:"ReceiverId,omitempty"`
	RemoteProcedure      string   `protobuf:"bytes,4,opt,name=RemoteProcedure,proto3" json:"RemoteProcedure,omitempty"`
	Length               int32    `protobuf:"varint,5,opt,name=Length,proto3" json:"Length,omitempty"`
	Request              bool     `protobuf:"varint,6,opt,name=Request,proto3" json:"Request,omitempty"`
	SenderAddress        string   `protobuf:"bytes,7,opt,name=SenderAddress,proto3" json:"SenderAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RPC) Reset()         { *m = RPC{} }
func (m *RPC) String() string { return proto.CompactTextString(m) }
func (*RPC) ProtoMessage()    {}
func (*RPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}
func (m *RPC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RPC.Unmarshal(m, b)
}
func (m *RPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RPC.Marshal(b, m, deterministic)
}
func (dst *RPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RPC.Merge(dst, src)
}
func (m *RPC) XXX_Size() int {
	return xxx_messageInfo_RPC.Size(m)
}
func (m *RPC) XXX_DiscardUnknown() {
	xxx_messageInfo_RPC.DiscardUnknown(m)
}

var xxx_messageInfo_RPC proto.InternalMessageInfo

func (m *RPC) GetMessageId() int32 {
	if m != nil {
		return m.MessageId
	}
	return 0
}

func (m *RPC) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

func (m *RPC) GetReceiverId() string {
	if m != nil {
		return m.ReceiverId
	}
	return ""
}

func (m *RPC) GetRemoteProcedure() string {
	if m != nil {
		return m.RemoteProcedure
	}
	return ""
}

func (m *RPC) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *RPC) GetRequest() bool {
	if m != nil {
		return m.Request
	}
	return false
}

func (m *RPC) GetSenderAddress() string {
	if m != nil {
		return m.SenderAddress
	}
	return ""
}

type FindValueRequest struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindValueRequest) Reset()         { *m = FindValueRequest{} }
func (m *FindValueRequest) String() string { return proto.CompactTextString(m) }
func (*FindValueRequest) ProtoMessage()    {}
func (*FindValueRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}
func (m *FindValueRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindValueRequest.Unmarshal(m, b)
}
func (m *FindValueRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindValueRequest.Marshal(b, m, deterministic)
}
func (dst *FindValueRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindValueRequest.Merge(dst, src)
}
func (m *FindValueRequest) XXX_Size() int {
	return xxx_messageInfo_FindValueRequest.Size(m)
}
func (m *FindValueRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindValueRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindValueRequest proto.InternalMessageInfo

func (m *FindValueRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

//
// A FindValueResponse can either have the file or K number of contacts.
// This is indicated by the hasData flag
//
type FindValueResponse struct {
	HasData              bool       `protobuf:"varint,1,opt,name=HasData,proto3" json:"HasData,omitempty"`
	Data                 []byte     `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	Contacts             []*Contact `protobuf:"bytes,3,rep,name=Contacts,proto3" json:"Contacts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *FindValueResponse) Reset()         { *m = FindValueResponse{} }
func (m *FindValueResponse) String() string { return proto.CompactTextString(m) }
func (*FindValueResponse) ProtoMessage()    {}
func (*FindValueResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}
func (m *FindValueResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindValueResponse.Unmarshal(m, b)
}
func (m *FindValueResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindValueResponse.Marshal(b, m, deterministic)
}
func (dst *FindValueResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindValueResponse.Merge(dst, src)
}
func (m *FindValueResponse) XXX_Size() int {
	return xxx_messageInfo_FindValueResponse.Size(m)
}
func (m *FindValueResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindValueResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindValueResponse proto.InternalMessageInfo

func (m *FindValueResponse) GetHasData() bool {
	if m != nil {
		return m.HasData
	}
	return false
}

func (m *FindValueResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FindValueResponse) GetContacts() []*Contact {
	if m != nil {
		return m.Contacts
	}
	return nil
}

type DeleteValueRequest struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteValueRequest) Reset()         { *m = DeleteValueRequest{} }
func (m *DeleteValueRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteValueRequest) ProtoMessage()    {}
func (*DeleteValueRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}
func (m *DeleteValueRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteValueRequest.Unmarshal(m, b)
}
func (m *DeleteValueRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteValueRequest.Marshal(b, m, deterministic)
}
func (dst *DeleteValueRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteValueRequest.Merge(dst, src)
}
func (m *DeleteValueRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteValueRequest.Size(m)
}
func (m *DeleteValueRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteValueRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteValueRequest proto.InternalMessageInfo

func (m *DeleteValueRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type DeleteValueResponse struct {
	Deleted              bool     `protobuf:"varint,1,opt,name=Deleted,proto3" json:"Deleted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteValueResponse) Reset()         { *m = DeleteValueResponse{} }
func (m *DeleteValueResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteValueResponse) ProtoMessage()    {}
func (*DeleteValueResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}
func (m *DeleteValueResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteValueResponse.Unmarshal(m, b)
}
func (m *DeleteValueResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteValueResponse.Marshal(b, m, deterministic)
}
func (dst *DeleteValueResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteValueResponse.Merge(dst, src)
}
func (m *DeleteValueResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteValueResponse.Size(m)
}
func (m *DeleteValueResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteValueResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteValueResponse proto.InternalMessageInfo

func (m *DeleteValueResponse) GetDeleted() bool {
	if m != nil {
		return m.Deleted
	}
	return false
}

type SendDataMessage struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Hash                 string   `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendDataMessage) Reset()         { *m = SendDataMessage{} }
func (m *SendDataMessage) String() string { return proto.CompactTextString(m) }
func (*SendDataMessage) ProtoMessage()    {}
func (*SendDataMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}
func (m *SendDataMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendDataMessage.Unmarshal(m, b)
}
func (m *SendDataMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendDataMessage.Marshal(b, m, deterministic)
}
func (dst *SendDataMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendDataMessage.Merge(dst, src)
}
func (m *SendDataMessage) XXX_Size() int {
	return xxx_messageInfo_SendDataMessage.Size(m)
}
func (m *SendDataMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SendDataMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SendDataMessage proto.InternalMessageInfo

func (m *SendDataMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SendDataMessage) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type FindNodeRequest struct {
	TargetID             string   `protobuf:"bytes,1,opt,name=TargetID,proto3" json:"TargetID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindNodeRequest) Reset()         { *m = FindNodeRequest{} }
func (m *FindNodeRequest) String() string { return proto.CompactTextString(m) }
func (*FindNodeRequest) ProtoMessage()    {}
func (*FindNodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{6}
}
func (m *FindNodeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindNodeRequest.Unmarshal(m, b)
}
func (m *FindNodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindNodeRequest.Marshal(b, m, deterministic)
}
func (dst *FindNodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindNodeRequest.Merge(dst, src)
}
func (m *FindNodeRequest) XXX_Size() int {
	return xxx_messageInfo_FindNodeRequest.Size(m)
}
func (m *FindNodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FindNodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FindNodeRequest proto.InternalMessageInfo

func (m *FindNodeRequest) GetTargetID() string {
	if m != nil {
		return m.TargetID
	}
	return ""
}

type FindNodeResponse struct {
	Contacts             []*Contact `protobuf:"bytes,1,rep,name=contacts,proto3" json:"contacts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *FindNodeResponse) Reset()         { *m = FindNodeResponse{} }
func (m *FindNodeResponse) String() string { return proto.CompactTextString(m) }
func (*FindNodeResponse) ProtoMessage()    {}
func (*FindNodeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{7}
}
func (m *FindNodeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindNodeResponse.Unmarshal(m, b)
}
func (m *FindNodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindNodeResponse.Marshal(b, m, deterministic)
}
func (dst *FindNodeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindNodeResponse.Merge(dst, src)
}
func (m *FindNodeResponse) XXX_Size() int {
	return xxx_messageInfo_FindNodeResponse.Size(m)
}
func (m *FindNodeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_FindNodeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_FindNodeResponse proto.InternalMessageInfo

func (m *FindNodeResponse) GetContacts() []*Contact {
	if m != nil {
		return m.Contacts
	}
	return nil
}

type Contact struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Contact) Reset()         { *m = Contact{} }
func (m *Contact) String() string { return proto.CompactTextString(m) }
func (*Contact) ProtoMessage()    {}
func (*Contact) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{8}
}
func (m *Contact) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Contact.Unmarshal(m, b)
}
func (m *Contact) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Contact.Marshal(b, m, deterministic)
}
func (dst *Contact) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Contact.Merge(dst, src)
}
func (m *Contact) XXX_Size() int {
	return xxx_messageInfo_Contact.Size(m)
}
func (m *Contact) XXX_DiscardUnknown() {
	xxx_messageInfo_Contact.DiscardUnknown(m)
}

var xxx_messageInfo_Contact proto.InternalMessageInfo

func (m *Contact) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Contact) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*RPC)(nil), "message.RPC")
	proto.RegisterType((*FindValueRequest)(nil), "message.FindValueRequest")
	proto.RegisterType((*FindValueResponse)(nil), "message.FindValueResponse")
	proto.RegisterType((*DeleteValueRequest)(nil), "message.DeleteValueRequest")
	proto.RegisterType((*DeleteValueResponse)(nil), "message.DeleteValueResponse")
	proto.RegisterType((*SendDataMessage)(nil), "message.SendDataMessage")
	proto.RegisterType((*FindNodeRequest)(nil), "message.FindNodeRequest")
	proto.RegisterType((*FindNodeResponse)(nil), "message.FindNodeResponse")
	proto.RegisterType((*Contact)(nil), "message.Contact")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x51, 0x6b, 0xdb, 0x30,
	0x14, 0x85, 0xb1, 0x9d, 0xc4, 0xce, 0xdd, 0xb2, 0x64, 0x1a, 0x0c, 0x11, 0xc6, 0x30, 0x66, 0x0c,
	0x3f, 0x64, 0x19, 0x2c, 0x4f, 0x7b, 0xdb, 0x88, 0x19, 0x31, 0xac, 0x25, 0xa8, 0xa5, 0xef, 0xaa,
	0x75, 0x49, 0x02, 0x89, 0x95, 0x5a, 0x4a, 0xff, 0x6f, 0xff, 0x49, 0xb1, 0x2c, 0x39, 0x4e, 0xa1,
	0xf4, 0x4d, 0xe7, 0xe8, 0x58, 0xf7, 0xd3, 0x91, 0x61, 0x74, 0x40, 0xa5, 0xf8, 0x06, 0xe7, 0xc7,
	0x4a, 0x6a, 0x49, 0x42, 0x2b, 0x93, 0x27, 0x0f, 0x02, 0xb6, 0x5e, 0x92, 0x2f, 0x30, 0xbc, 0x6a,
	0xac, 0x5c, 0x50, 0x2f, 0xf6, 0xd2, 0x3e, 0x3b, 0x1b, 0x64, 0x0a, 0xd1, 0x0d, 0x96, 0x02, 0xab,
	0x5c, 0x50, 0x3f, 0xf6, 0xd2, 0x21, 0x6b, 0x35, 0xf9, 0x0a, 0xc0, 0xb0, 0xc0, 0xdd, 0xa3, 0xd9,
	0x0d, 0xcc, 0x6e, 0xc7, 0x21, 0x29, 0x8c, 0x19, 0x1e, 0xa4, 0xc6, 0x75, 0x25, 0x0b, 0x14, 0xa7,
	0x0a, 0x69, 0xcf, 0x84, 0x5e, 0xda, 0xe4, 0x33, 0x0c, 0xfe, 0x63, 0xb9, 0xd1, 0x5b, 0xda, 0x37,
	0x00, 0x56, 0x11, 0x0a, 0x21, 0xc3, 0x87, 0x13, 0x2a, 0x4d, 0x07, 0xb1, 0x97, 0x46, 0xcc, 0x49,
	0xf2, 0x0d, 0x46, 0x0d, 0xc7, 0x5f, 0x21, 0x2a, 0x54, 0x8a, 0x86, 0xe6, 0xe4, 0x4b, 0x33, 0xf9,
	0x0e, 0x93, 0x7f, 0xbb, 0x52, 0xdc, 0xf1, 0xfd, 0x09, 0xdd, 0x97, 0x04, 0x7a, 0x2b, 0xae, 0xb6,
	0xe6, 0xaa, 0x43, 0x66, 0xd6, 0x89, 0x84, 0x8f, 0x9d, 0x9c, 0x3a, 0xca, 0x52, 0x61, 0x3d, 0x7c,
	0xc5, 0x55, 0xc6, 0x35, 0x37, 0xd9, 0x88, 0x39, 0x59, 0x1f, 0x61, 0xec, 0xba, 0x90, 0xf7, 0xcc,
	0xac, 0xc9, 0x0c, 0xa2, 0xa5, 0x2c, 0x35, 0x2f, 0xb4, 0xa2, 0x41, 0x1c, 0xa4, 0xef, 0x7e, 0x4d,
	0xe6, 0xae, 0x79, 0xbb, 0xc1, 0xda, 0x44, 0x92, 0x02, 0xc9, 0x70, 0x8f, 0x1a, 0xdf, 0x44, 0xfb,
	0x09, 0x9f, 0x2e, 0x92, 0x67, 0xb8, 0xc6, 0x16, 0x0e, 0xce, 0xca, 0xe4, 0x37, 0x8c, 0xeb, 0x12,
	0x6a, 0x28, 0xfb, 0x8c, 0x2d, 0xaf, 0xd7, 0xe1, 0x75, 0xb3, 0xfc, 0xce, 0xac, 0x1f, 0x30, 0xae,
	0x6b, 0xb8, 0x96, 0xa2, 0x45, 0x9a, 0x42, 0x74, 0xcb, 0xab, 0x0d, 0xea, 0x3c, 0xb3, 0x58, 0xad,
	0x4e, 0xfe, 0x34, 0xed, 0x36, 0x71, 0xcb, 0x35, 0x83, 0xa8, 0x70, 0x35, 0x78, 0xaf, 0xd5, 0xe0,
	0x12, 0xc9, 0x02, 0x42, 0x6b, 0x92, 0x0f, 0xe0, 0xb7, 0x23, 0xfc, 0x3c, 0xab, 0x2f, 0xe8, 0x9e,
	0xb6, 0x41, 0x74, 0xf2, 0x7e, 0x60, 0x7e, 0xe4, 0xc5, 0x73, 0x00, 0x00, 0x00, 0xff, 0xff, 0xca,
	0x51, 0x41, 0xe2, 0xd9, 0x02, 0x00, 0x00,
}
