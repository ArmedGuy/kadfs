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
	SenderIp             string   `protobuf:"bytes,7,opt,name=SenderIp,proto3" json:"SenderIp,omitempty"`
	SenderPort           string   `protobuf:"bytes,8,opt,name=SenderPort,proto3" json:"SenderPort,omitempty"`
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

func (m *RPC) GetSenderIp() string {
	if m != nil {
		return m.SenderIp
	}
	return ""
}

func (m *RPC) GetSenderPort() string {
	if m != nil {
		return m.SenderPort
	}
	return ""
}

type FindDataMessage struct {
	Hash                 string   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FindDataMessage) Reset()         { *m = FindDataMessage{} }
func (m *FindDataMessage) String() string { return proto.CompactTextString(m) }
func (*FindDataMessage) ProtoMessage()    {}
func (*FindDataMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}
func (m *FindDataMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FindDataMessage.Unmarshal(m, b)
}
func (m *FindDataMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FindDataMessage.Marshal(b, m, deterministic)
}
func (dst *FindDataMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FindDataMessage.Merge(dst, src)
}
func (m *FindDataMessage) XXX_Size() int {
	return xxx_messageInfo_FindDataMessage.Size(m)
}
func (m *FindDataMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_FindDataMessage.DiscardUnknown(m)
}

var xxx_messageInfo_FindDataMessage proto.InternalMessageInfo

func (m *FindDataMessage) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type SendDataMessage struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendDataMessage) Reset()         { *m = SendDataMessage{} }
func (m *SendDataMessage) String() string { return proto.CompactTextString(m) }
func (*SendDataMessage) ProtoMessage()    {}
func (*SendDataMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
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

func init() {
	proto.RegisterType((*RPC)(nil), "message.RPC")
	proto.RegisterType((*FindDataMessage)(nil), "message.FindDataMessage")
	proto.RegisterType((*SendDataMessage)(nil), "message.SendDataMessage")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 231 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0xc6, 0x59, 0xdb, 0xe6, 0xcf, 0xa0, 0x14, 0xe6, 0x20, 0x83, 0x88, 0x84, 0x82, 0x90, 0x93,
	0x17, 0x1f, 0x41, 0x11, 0x0a, 0x0a, 0x61, 0x7c, 0x82, 0xb5, 0x3b, 0x34, 0x3d, 0x34, 0x1b, 0x77,
	0xb7, 0xbe, 0xba, 0x57, 0xc9, 0x24, 0xda, 0xda, 0xdb, 0x7c, 0xbf, 0xf9, 0x16, 0x7e, 0xb3, 0x70,
	0xb5, 0x97, 0x18, 0xed, 0x56, 0x1e, 0xfa, 0xe0, 0x93, 0xc7, 0x7c, 0x8a, 0xab, 0x6f, 0x03, 0x33,
	0x6e, 0x9e, 0xf0, 0x16, 0xca, 0xb7, 0x11, 0xad, 0x1d, 0x99, 0xca, 0xd4, 0x0b, 0x3e, 0x02, 0xbc,
	0x81, 0xe2, 0x5d, 0x3a, 0x27, 0x61, 0xed, 0xe8, 0xa2, 0x32, 0x75, 0xc9, 0x7f, 0x19, 0xef, 0x00,
	0x58, 0x36, 0xb2, 0xfb, 0xd2, 0xed, 0x4c, 0xb7, 0x27, 0x04, 0x6b, 0x58, 0xb2, 0xec, 0x7d, 0x92,
	0x26, 0xf8, 0x8d, 0xb8, 0x43, 0x10, 0x9a, 0x6b, 0xe9, 0x1c, 0xe3, 0x35, 0x64, 0xaf, 0xd2, 0x6d,
	0x53, 0x4b, 0x0b, 0x15, 0x98, 0x12, 0x12, 0xe4, 0x2c, 0x9f, 0x07, 0x89, 0x89, 0xb2, 0xca, 0xd4,
	0x05, 0xff, 0xc6, 0x13, 0xaf, 0x9e, 0xf2, 0x7f, 0x5e, 0xfd, 0xe0, 0x35, 0xce, 0x8d, 0x0f, 0x89,
	0x8a, 0xd1, 0xeb, 0x48, 0x56, 0xf7, 0xb0, 0x7c, 0xd9, 0x75, 0xee, 0xd9, 0x26, 0x3b, 0x1d, 0x8a,
	0x08, 0xf3, 0xd6, 0xc6, 0x56, 0xef, 0x2f, 0x59, 0xe7, 0xa1, 0x36, 0x3c, 0x3a, 0xab, 0x39, 0x9b,
	0xac, 0xd6, 0x2e, 0x59, 0xe7, 0x8f, 0x4c, 0xff, 0xf5, 0xf1, 0x27, 0x00, 0x00, 0xff, 0xff, 0x17,
	0x09, 0xc1, 0xc7, 0x68, 0x01, 0x00, 0x00,
}
