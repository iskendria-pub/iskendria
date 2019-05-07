// Code generated by protoc-gen-go. DO NOT EDIT.
// source: person.proto

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

type StatePerson struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedOn            int64    `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	ModifiedOn           int64    `protobuf:"varint,3,opt,name=modifiedOn,proto3" json:"modifiedOn,omitempty"`
	PublicKey            string   `protobuf:"bytes,4,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	Name                 string   `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Email                string   `protobuf:"bytes,6,opt,name=email,proto3" json:"email,omitempty"`
	IsMajor              bool     `protobuf:"varint,7,opt,name=isMajor,proto3" json:"isMajor,omitempty"`
	IsSigned             bool     `protobuf:"varint,8,opt,name=isSigned,proto3" json:"isSigned,omitempty"`
	Balance              int32    `protobuf:"varint,9,opt,name=balance,proto3" json:"balance,omitempty"`
	BiographyHash        string   `protobuf:"bytes,10,opt,name=biographyHash,proto3" json:"biographyHash,omitempty"`
	BiographyFormat      string   `protobuf:"bytes,11,opt,name=biographyFormat,proto3" json:"biographyFormat,omitempty"`
	Organization         string   `protobuf:"bytes,12,opt,name=organization,proto3" json:"organization,omitempty"`
	Telephone            string   `protobuf:"bytes,13,opt,name=telephone,proto3" json:"telephone,omitempty"`
	Address              string   `protobuf:"bytes,14,opt,name=address,proto3" json:"address,omitempty"`
	PostalCode           string   `protobuf:"bytes,15,opt,name=postalCode,proto3" json:"postalCode,omitempty"`
	Country              string   `protobuf:"bytes,16,opt,name=country,proto3" json:"country,omitempty"`
	ExtraInfo            string   `protobuf:"bytes,17,opt,name=extraInfo,proto3" json:"extraInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatePerson) Reset()         { *m = StatePerson{} }
func (m *StatePerson) String() string { return proto.CompactTextString(m) }
func (*StatePerson) ProtoMessage()    {}
func (*StatePerson) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c9e10cf24b1156d, []int{0}
}

func (m *StatePerson) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatePerson.Unmarshal(m, b)
}
func (m *StatePerson) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatePerson.Marshal(b, m, deterministic)
}
func (m *StatePerson) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatePerson.Merge(m, src)
}
func (m *StatePerson) XXX_Size() int {
	return xxx_messageInfo_StatePerson.Size(m)
}
func (m *StatePerson) XXX_DiscardUnknown() {
	xxx_messageInfo_StatePerson.DiscardUnknown(m)
}

var xxx_messageInfo_StatePerson proto.InternalMessageInfo

func (m *StatePerson) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StatePerson) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *StatePerson) GetModifiedOn() int64 {
	if m != nil {
		return m.ModifiedOn
	}
	return 0
}

func (m *StatePerson) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *StatePerson) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StatePerson) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *StatePerson) GetIsMajor() bool {
	if m != nil {
		return m.IsMajor
	}
	return false
}

func (m *StatePerson) GetIsSigned() bool {
	if m != nil {
		return m.IsSigned
	}
	return false
}

func (m *StatePerson) GetBalance() int32 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *StatePerson) GetBiographyHash() string {
	if m != nil {
		return m.BiographyHash
	}
	return ""
}

func (m *StatePerson) GetBiographyFormat() string {
	if m != nil {
		return m.BiographyFormat
	}
	return ""
}

func (m *StatePerson) GetOrganization() string {
	if m != nil {
		return m.Organization
	}
	return ""
}

func (m *StatePerson) GetTelephone() string {
	if m != nil {
		return m.Telephone
	}
	return ""
}

func (m *StatePerson) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *StatePerson) GetPostalCode() string {
	if m != nil {
		return m.PostalCode
	}
	return ""
}

func (m *StatePerson) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *StatePerson) GetExtraInfo() string {
	if m != nil {
		return m.ExtraInfo
	}
	return ""
}

type CommandPersonCreate struct {
	NewPersonId          string   `protobuf:"bytes,1,opt,name=newPersonId,proto3" json:"newPersonId,omitempty"`
	PublicKey            string   `protobuf:"bytes,2,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email                string   `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandPersonCreate) Reset()         { *m = CommandPersonCreate{} }
func (m *CommandPersonCreate) String() string { return proto.CompactTextString(m) }
func (*CommandPersonCreate) ProtoMessage()    {}
func (*CommandPersonCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c9e10cf24b1156d, []int{1}
}

func (m *CommandPersonCreate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandPersonCreate.Unmarshal(m, b)
}
func (m *CommandPersonCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandPersonCreate.Marshal(b, m, deterministic)
}
func (m *CommandPersonCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandPersonCreate.Merge(m, src)
}
func (m *CommandPersonCreate) XXX_Size() int {
	return xxx_messageInfo_CommandPersonCreate.Size(m)
}
func (m *CommandPersonCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandPersonCreate.DiscardUnknown(m)
}

var xxx_messageInfo_CommandPersonCreate proto.InternalMessageInfo

func (m *CommandPersonCreate) GetNewPersonId() string {
	if m != nil {
		return m.NewPersonId
	}
	return ""
}

func (m *CommandPersonCreate) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *CommandPersonCreate) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CommandPersonCreate) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type CommandPersonUpdateProperties struct {
	PersonId             string        `protobuf:"bytes,1,opt,name=personId,proto3" json:"personId,omitempty"`
	PublicKeyUpdate      *StringUpdate `protobuf:"bytes,2,opt,name=publicKeyUpdate,proto3" json:"publicKeyUpdate,omitempty"`
	NameUpdate           *StringUpdate `protobuf:"bytes,3,opt,name=nameUpdate,proto3" json:"nameUpdate,omitempty"`
	EmailUpdate          *StringUpdate `protobuf:"bytes,4,opt,name=emailUpdate,proto3" json:"emailUpdate,omitempty"`
	BiographyHashUpdate  *StringUpdate `protobuf:"bytes,5,opt,name=biographyHashUpdate,proto3" json:"biographyHashUpdate,omitempty"`
	OrganizationUpdate   *StringUpdate `protobuf:"bytes,6,opt,name=organizationUpdate,proto3" json:"organizationUpdate,omitempty"`
	TelephoneUpdate      *StringUpdate `protobuf:"bytes,7,opt,name=telephoneUpdate,proto3" json:"telephoneUpdate,omitempty"`
	AddressUpdate        *StringUpdate `protobuf:"bytes,8,opt,name=addressUpdate,proto3" json:"addressUpdate,omitempty"`
	PostalCodeUpdate     *StringUpdate `protobuf:"bytes,9,opt,name=postalCodeUpdate,proto3" json:"postalCodeUpdate,omitempty"`
	CountryUpdate        *StringUpdate `protobuf:"bytes,10,opt,name=countryUpdate,proto3" json:"countryUpdate,omitempty"`
	ExtraInfoUpdate      *StringUpdate `protobuf:"bytes,11,opt,name=extraInfoUpdate,proto3" json:"extraInfoUpdate,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CommandPersonUpdateProperties) Reset()         { *m = CommandPersonUpdateProperties{} }
func (m *CommandPersonUpdateProperties) String() string { return proto.CompactTextString(m) }
func (*CommandPersonUpdateProperties) ProtoMessage()    {}
func (*CommandPersonUpdateProperties) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c9e10cf24b1156d, []int{2}
}

func (m *CommandPersonUpdateProperties) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandPersonUpdateProperties.Unmarshal(m, b)
}
func (m *CommandPersonUpdateProperties) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandPersonUpdateProperties.Marshal(b, m, deterministic)
}
func (m *CommandPersonUpdateProperties) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandPersonUpdateProperties.Merge(m, src)
}
func (m *CommandPersonUpdateProperties) XXX_Size() int {
	return xxx_messageInfo_CommandPersonUpdateProperties.Size(m)
}
func (m *CommandPersonUpdateProperties) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandPersonUpdateProperties.DiscardUnknown(m)
}

var xxx_messageInfo_CommandPersonUpdateProperties proto.InternalMessageInfo

func (m *CommandPersonUpdateProperties) GetPersonId() string {
	if m != nil {
		return m.PersonId
	}
	return ""
}

func (m *CommandPersonUpdateProperties) GetPublicKeyUpdate() *StringUpdate {
	if m != nil {
		return m.PublicKeyUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetNameUpdate() *StringUpdate {
	if m != nil {
		return m.NameUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetEmailUpdate() *StringUpdate {
	if m != nil {
		return m.EmailUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetBiographyHashUpdate() *StringUpdate {
	if m != nil {
		return m.BiographyHashUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetOrganizationUpdate() *StringUpdate {
	if m != nil {
		return m.OrganizationUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetTelephoneUpdate() *StringUpdate {
	if m != nil {
		return m.TelephoneUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetAddressUpdate() *StringUpdate {
	if m != nil {
		return m.AddressUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetPostalCodeUpdate() *StringUpdate {
	if m != nil {
		return m.PostalCodeUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetCountryUpdate() *StringUpdate {
	if m != nil {
		return m.CountryUpdate
	}
	return nil
}

func (m *CommandPersonUpdateProperties) GetExtraInfoUpdate() *StringUpdate {
	if m != nil {
		return m.ExtraInfoUpdate
	}
	return nil
}

type CommandPersonUpdateAuthorization struct {
	PersonId             string     `protobuf:"bytes,1,opt,name=personId,proto3" json:"personId,omitempty"`
	MakeMajor            BoolUpdate `protobuf:"varint,2,opt,name=makeMajor,proto3,enum=BoolUpdate" json:"makeMajor,omitempty"`
	MakeSigned           BoolUpdate `protobuf:"varint,3,opt,name=makeSigned,proto3,enum=BoolUpdate" json:"makeSigned,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CommandPersonUpdateAuthorization) Reset()         { *m = CommandPersonUpdateAuthorization{} }
func (m *CommandPersonUpdateAuthorization) String() string { return proto.CompactTextString(m) }
func (*CommandPersonUpdateAuthorization) ProtoMessage()    {}
func (*CommandPersonUpdateAuthorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c9e10cf24b1156d, []int{3}
}

func (m *CommandPersonUpdateAuthorization) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandPersonUpdateAuthorization.Unmarshal(m, b)
}
func (m *CommandPersonUpdateAuthorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandPersonUpdateAuthorization.Marshal(b, m, deterministic)
}
func (m *CommandPersonUpdateAuthorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandPersonUpdateAuthorization.Merge(m, src)
}
func (m *CommandPersonUpdateAuthorization) XXX_Size() int {
	return xxx_messageInfo_CommandPersonUpdateAuthorization.Size(m)
}
func (m *CommandPersonUpdateAuthorization) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandPersonUpdateAuthorization.DiscardUnknown(m)
}

var xxx_messageInfo_CommandPersonUpdateAuthorization proto.InternalMessageInfo

func (m *CommandPersonUpdateAuthorization) GetPersonId() string {
	if m != nil {
		return m.PersonId
	}
	return ""
}

func (m *CommandPersonUpdateAuthorization) GetMakeMajor() BoolUpdate {
	if m != nil {
		return m.MakeMajor
	}
	return BoolUpdate_UNMODIFIED
}

func (m *CommandPersonUpdateAuthorization) GetMakeSigned() BoolUpdate {
	if m != nil {
		return m.MakeSigned
	}
	return BoolUpdate_UNMODIFIED
}

type CommandPersonUpdateBalanceIncrement struct {
	PersonId             string   `protobuf:"bytes,1,opt,name=personId,proto3" json:"personId,omitempty"`
	BalanceIncrement     int32    `protobuf:"varint,2,opt,name=balanceIncrement,proto3" json:"balanceIncrement,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandPersonUpdateBalanceIncrement) Reset()         { *m = CommandPersonUpdateBalanceIncrement{} }
func (m *CommandPersonUpdateBalanceIncrement) String() string { return proto.CompactTextString(m) }
func (*CommandPersonUpdateBalanceIncrement) ProtoMessage()    {}
func (*CommandPersonUpdateBalanceIncrement) Descriptor() ([]byte, []int) {
	return fileDescriptor_4c9e10cf24b1156d, []int{4}
}

func (m *CommandPersonUpdateBalanceIncrement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandPersonUpdateBalanceIncrement.Unmarshal(m, b)
}
func (m *CommandPersonUpdateBalanceIncrement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandPersonUpdateBalanceIncrement.Marshal(b, m, deterministic)
}
func (m *CommandPersonUpdateBalanceIncrement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandPersonUpdateBalanceIncrement.Merge(m, src)
}
func (m *CommandPersonUpdateBalanceIncrement) XXX_Size() int {
	return xxx_messageInfo_CommandPersonUpdateBalanceIncrement.Size(m)
}
func (m *CommandPersonUpdateBalanceIncrement) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandPersonUpdateBalanceIncrement.DiscardUnknown(m)
}

var xxx_messageInfo_CommandPersonUpdateBalanceIncrement proto.InternalMessageInfo

func (m *CommandPersonUpdateBalanceIncrement) GetPersonId() string {
	if m != nil {
		return m.PersonId
	}
	return ""
}

func (m *CommandPersonUpdateBalanceIncrement) GetBalanceIncrement() int32 {
	if m != nil {
		return m.BalanceIncrement
	}
	return 0
}

func init() {
	proto.RegisterType((*StatePerson)(nil), "StatePerson")
	proto.RegisterType((*CommandPersonCreate)(nil), "CommandPersonCreate")
	proto.RegisterType((*CommandPersonUpdateProperties)(nil), "CommandPersonUpdateProperties")
	proto.RegisterType((*CommandPersonUpdateAuthorization)(nil), "CommandPersonUpdateAuthorization")
	proto.RegisterType((*CommandPersonUpdateBalanceIncrement)(nil), "CommandPersonUpdateBalanceIncrement")
}

func init() { proto.RegisterFile("person.proto", fileDescriptor_4c9e10cf24b1156d) }

var fileDescriptor_4c9e10cf24b1156d = []byte{
	// 614 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0x56, 0x9b, 0x66, 0x6d, 0x5f, 0xda, 0x6d, 0x78, 0x1c, 0xac, 0x09, 0x50, 0x14, 0x38, 0x14,
	0x10, 0x43, 0xda, 0x0e, 0x88, 0x03, 0x42, 0x6c, 0x12, 0x62, 0x42, 0x08, 0x94, 0x89, 0x0b, 0x37,
	0x37, 0xf1, 0x5a, 0x43, 0x6c, 0x47, 0x8e, 0x27, 0x18, 0x17, 0xfe, 0x05, 0x7f, 0x92, 0x1f, 0xc0,
	0x15, 0xd9, 0x71, 0xd2, 0xa4, 0x33, 0xbb, 0xf5, 0x7d, 0xef, 0xfb, 0xfa, 0xbe, 0xbc, 0xf7, 0x35,
	0x85, 0x59, 0x49, 0x55, 0x25, 0xc5, 0x51, 0xa9, 0xa4, 0x96, 0x87, 0xb3, 0x4c, 0x72, 0xde, 0x54,
	0xc9, 0xdf, 0x00, 0xa2, 0x0b, 0x4d, 0x34, 0xfd, 0x64, 0x39, 0x68, 0x17, 0x86, 0x2c, 0xc7, 0x83,
	0x78, 0xb0, 0x98, 0xa6, 0x43, 0x96, 0xa3, 0x7b, 0x30, 0xcd, 0x14, 0x25, 0x9a, 0xe6, 0x1f, 0x05,
	0x1e, 0xc6, 0x83, 0x45, 0x90, 0x6e, 0x00, 0xf4, 0x00, 0x80, 0xcb, 0x9c, 0x5d, 0x32, 0xdb, 0x0e,
	0x6c, 0xbb, 0x83, 0x18, 0x75, 0x79, 0xb5, 0x2c, 0x58, 0xf6, 0x9e, 0x5e, 0xe3, 0x91, 0xfd, 0xd2,
	0x0d, 0x80, 0x10, 0x8c, 0x04, 0xe1, 0x14, 0x87, 0xb6, 0x61, 0x3f, 0xa3, 0xbb, 0x10, 0x52, 0x4e,
	0x58, 0x81, 0x77, 0x2c, 0x58, 0x17, 0x08, 0xc3, 0x98, 0x55, 0x1f, 0xc8, 0x57, 0xa9, 0xf0, 0x38,
	0x1e, 0x2c, 0x26, 0x69, 0x53, 0xa2, 0x43, 0x98, 0xb0, 0xea, 0x82, 0xad, 0x04, 0xcd, 0xf1, 0xc4,
	0xb6, 0xda, 0xda, 0xa8, 0x96, 0xa4, 0x20, 0x22, 0xa3, 0x78, 0x1a, 0x0f, 0x16, 0x61, 0xda, 0x94,
	0xe8, 0x11, 0xcc, 0x97, 0x4c, 0xae, 0x14, 0x29, 0xd7, 0xd7, 0xef, 0x48, 0xb5, 0xc6, 0x60, 0xa7,
	0xf5, 0x41, 0xb4, 0x80, 0xbd, 0x16, 0x78, 0x2b, 0x15, 0x27, 0x1a, 0x47, 0x96, 0xb7, 0x0d, 0xa3,
	0x04, 0x66, 0x52, 0xad, 0x88, 0x60, 0x3f, 0x89, 0x66, 0x52, 0xe0, 0x99, 0xa5, 0xf5, 0x30, 0xb3,
	0x0b, 0x4d, 0x0b, 0x5a, 0xae, 0xa5, 0xa0, 0x78, 0x5e, 0xef, 0xa2, 0x05, 0x8c, 0x57, 0x92, 0xe7,
	0x8a, 0x56, 0x15, 0xde, 0xb5, 0xbd, 0xa6, 0x34, 0x3b, 0x2e, 0x65, 0xa5, 0x49, 0x71, 0x26, 0x73,
	0x8a, 0xf7, 0x6c, 0xb3, 0x83, 0x18, 0x65, 0x26, 0xaf, 0x84, 0x56, 0xd7, 0x78, 0xbf, 0x56, 0xba,
	0xd2, 0x4c, 0xa4, 0x3f, 0xb4, 0x22, 0xe7, 0xe2, 0x52, 0xe2, 0x3b, 0xf5, 0xc4, 0x16, 0x48, 0x7e,
	0xc1, 0xc1, 0x99, 0xe4, 0x9c, 0x88, 0xbc, 0x3e, 0xfd, 0x99, 0xbd, 0x2a, 0x8a, 0x21, 0x12, 0xf4,
	0x7b, 0x0d, 0x9d, 0x37, 0x49, 0xe8, 0x42, 0xfd, 0xa3, 0x0e, 0xff, 0x77, 0xd4, 0xc0, 0x77, 0xd4,
	0x51, 0xe7, 0xa8, 0xc9, 0x9f, 0x11, 0xdc, 0xef, 0x39, 0xf8, 0x5c, 0xe6, 0x26, 0x88, 0x4a, 0x96,
	0x54, 0x69, 0x46, 0x2b, 0x73, 0xdc, 0xb2, 0x6f, 0xa4, 0xad, 0xd1, 0x0b, 0xd8, 0x6b, 0x87, 0xd6,
	0x42, 0xeb, 0x25, 0x3a, 0x9e, 0x1f, 0x5d, 0x68, 0xc5, 0xc4, 0xaa, 0x06, 0xd3, 0x6d, 0x16, 0x7a,
	0x06, 0x60, 0x4c, 0x39, 0x4d, 0xe0, 0xd3, 0x74, 0x08, 0xe8, 0x39, 0x44, 0xd6, 0xae, 0xe3, 0x8f,
	0x7c, 0xfc, 0x2e, 0x03, 0xbd, 0x86, 0x83, 0x5e, 0x8c, 0x9c, 0x30, 0xf4, 0x09, 0x7d, 0x4c, 0xf4,
	0x0a, 0x50, 0x37, 0x38, 0x4e, 0xbf, 0xe3, 0xd3, 0x7b, 0x88, 0x66, 0x31, 0x6d, 0xac, 0x9c, 0x76,
	0xec, 0x5d, 0xcc, 0x16, 0x0b, 0x9d, 0xc0, 0xdc, 0x65, 0xce, 0xc9, 0x26, 0x3e, 0x59, 0x9f, 0x83,
	0x5e, 0xc2, 0xfe, 0x26, 0x8b, 0x4e, 0x37, 0xf5, 0xe9, 0x6e, 0xd0, 0xcc, 0x3c, 0x97, 0x54, 0xa7,
	0x03, 0xef, 0xbc, 0x1e, 0xc7, 0x3c, 0x5d, 0x1b, 0x61, 0x27, 0x8b, 0xbc, 0x4f, 0xb7, 0xc5, 0x4a,
	0x7e, 0x0f, 0x20, 0xf6, 0xa4, 0xed, 0xcd, 0x95, 0x5e, 0x4b, 0xd5, 0xfc, 0x46, 0x6f, 0x0b, 0xdc,
	0x63, 0x98, 0x72, 0xf2, 0x8d, 0xd6, 0x6f, 0x21, 0x13, 0xb5, 0xdd, 0xe3, 0xe8, 0xe8, 0x54, 0x4a,
	0x77, 0xf7, 0x74, 0xd3, 0x45, 0x4f, 0x01, 0x4c, 0xe1, 0x5e, 0x4b, 0xc1, 0x4d, 0x6e, 0xa7, 0x9d,
	0x70, 0x78, 0xe8, 0xf1, 0x75, 0x5a, 0xbf, 0xa9, 0xce, 0x45, 0xa6, 0x28, 0xa7, 0x42, 0xdf, 0x6a,
	0xed, 0x09, 0xec, 0x2f, 0xb7, 0xf8, 0xd6, 0x61, 0x98, 0xde, 0xc0, 0x4f, 0xc7, 0x5f, 0x42, 0x2e,
	0x73, 0x5a, 0x2c, 0x77, 0xec, 0x1f, 0xc0, 0xc9, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x11, 0xa1,
	0x77, 0xc1, 0x1e, 0x06, 0x00, 0x00,
}