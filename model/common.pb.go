// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BoolUpdate int32

const (
	BoolUpdate_UNMODIFIED BoolUpdate = 0
	BoolUpdate_MAKE_FALSE BoolUpdate = 1
	BoolUpdate_MAKE_TRUE  BoolUpdate = 2
)

var BoolUpdate_name = map[int32]string{
	0: "UNMODIFIED",
	1: "MAKE_FALSE",
	2: "MAKE_TRUE",
}

var BoolUpdate_value = map[string]int32{
	"UNMODIFIED": 0,
	"MAKE_FALSE": 1,
	"MAKE_TRUE":  2,
}

func (x BoolUpdate) String() string {
	return proto.EnumName(BoolUpdate_name, int32(x))
}

func (BoolUpdate) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

type Judgement int32

const (
	Judgement_NEGATIVE Judgement = 0
	Judgement_POSITIVE Judgement = 1
)

var Judgement_name = map[int32]string{
	0: "NEGATIVE",
	1: "POSITIVE",
}

var Judgement_value = map[string]int32{
	"NEGATIVE": 0,
	"POSITIVE": 1,
}

func (x Judgement) String() string {
	return proto.EnumName(Judgement_name, int32(x))
}

func (Judgement) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

type IntUpdate struct {
	OldValue             int32    `protobuf:"varint,1,opt,name=oldValue,proto3" json:"oldValue,omitempty"`
	NewValue             int32    `protobuf:"varint,2,opt,name=newValue,proto3" json:"newValue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntUpdate) Reset()         { *m = IntUpdate{} }
func (m *IntUpdate) String() string { return proto.CompactTextString(m) }
func (*IntUpdate) ProtoMessage()    {}
func (*IntUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{0}
}

func (m *IntUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntUpdate.Unmarshal(m, b)
}
func (m *IntUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntUpdate.Marshal(b, m, deterministic)
}
func (m *IntUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntUpdate.Merge(m, src)
}
func (m *IntUpdate) XXX_Size() int {
	return xxx_messageInfo_IntUpdate.Size(m)
}
func (m *IntUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_IntUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_IntUpdate proto.InternalMessageInfo

func (m *IntUpdate) GetOldValue() int32 {
	if m != nil {
		return m.OldValue
	}
	return 0
}

func (m *IntUpdate) GetNewValue() int32 {
	if m != nil {
		return m.NewValue
	}
	return 0
}

type StringUpdate struct {
	OldValue             string   `protobuf:"bytes,1,opt,name=oldValue,proto3" json:"oldValue,omitempty"`
	NewValue             string   `protobuf:"bytes,2,opt,name=newValue,proto3" json:"newValue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringUpdate) Reset()         { *m = StringUpdate{} }
func (m *StringUpdate) String() string { return proto.CompactTextString(m) }
func (*StringUpdate) ProtoMessage()    {}
func (*StringUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_555bd8c177793206, []int{1}
}

func (m *StringUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringUpdate.Unmarshal(m, b)
}
func (m *StringUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringUpdate.Marshal(b, m, deterministic)
}
func (m *StringUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringUpdate.Merge(m, src)
}
func (m *StringUpdate) XXX_Size() int {
	return xxx_messageInfo_StringUpdate.Size(m)
}
func (m *StringUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_StringUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_StringUpdate proto.InternalMessageInfo

func (m *StringUpdate) GetOldValue() string {
	if m != nil {
		return m.OldValue
	}
	return ""
}

func (m *StringUpdate) GetNewValue() string {
	if m != nil {
		return m.NewValue
	}
	return ""
}

func init() {
	proto.RegisterEnum("BoolUpdate", BoolUpdate_name, BoolUpdate_value)
	proto.RegisterEnum("Judgement", Judgement_name, Judgement_value)
	proto.RegisterType((*IntUpdate)(nil), "IntUpdate")
	proto.RegisterType((*StringUpdate)(nil), "StringUpdate")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_555bd8c177793206) }

var fileDescriptor_555bd8c177793206 = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x72, 0xe6, 0xe2, 0xf4, 0xcc, 0x2b, 0x09,
	0x2d, 0x48, 0x49, 0x2c, 0x49, 0x15, 0x92, 0xe2, 0xe2, 0xc8, 0xcf, 0x49, 0x09, 0x4b, 0xcc, 0x29,
	0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x0d, 0x82, 0xf3, 0x41, 0x72, 0x79, 0xa9, 0xe5, 0x10,
	0x39, 0x26, 0x88, 0x1c, 0x8c, 0xaf, 0xe4, 0xc6, 0xc5, 0x13, 0x5c, 0x52, 0x94, 0x99, 0x97, 0x8e,
	0xc3, 0x1c, 0x4e, 0x3c, 0xe6, 0x70, 0x22, 0xcc, 0xd1, 0xb2, 0xe6, 0xe2, 0x72, 0xca, 0xcf, 0xcf,
	0x81, 0x9a, 0xc2, 0xc7, 0xc5, 0x15, 0xea, 0xe7, 0xeb, 0xef, 0xe2, 0xe9, 0xe6, 0xe9, 0xea, 0x22,
	0xc0, 0x00, 0xe2, 0xfb, 0x3a, 0x7a, 0xbb, 0xc6, 0xbb, 0x39, 0xfa, 0x04, 0xbb, 0x0a, 0x30, 0x0a,
	0xf1, 0x72, 0x71, 0x82, 0xf9, 0x21, 0x41, 0xa1, 0xae, 0x02, 0x4c, 0x5a, 0xea, 0x5c, 0x9c, 0x5e,
	0xa5, 0x29, 0xe9, 0xa9, 0xb9, 0xa9, 0x79, 0x25, 0x42, 0x3c, 0x5c, 0x1c, 0x7e, 0xae, 0xee, 0x8e,
	0x21, 0x9e, 0x61, 0xae, 0x02, 0x0c, 0x20, 0x5e, 0x80, 0x7f, 0xb0, 0x27, 0x98, 0xc7, 0xe8, 0xc4,
	0x1e, 0xc5, 0x9a, 0x9b, 0x9f, 0x92, 0x9a, 0x93, 0xc4, 0x06, 0x0e, 0x02, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x5d, 0xee, 0xb6, 0x03, 0x12, 0x01, 0x00, 0x00,
}