// Code generated by protoc-gen-go. DO NOT EDIT.
// source: journal.proto

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

type EditorState int32

const (
	EditorState_editorProposed EditorState = 0
	EditorState_editorAccepted EditorState = 1
)

var EditorState_name = map[int32]string{
	0: "editorProposed",
	1: "editorAccepted",
}

var EditorState_value = map[string]int32{
	"editorProposed": 0,
	"editorAccepted": 1,
}

func (x EditorState) String() string {
	return proto.EnumName(EditorState_name, int32(x))
}

func (EditorState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{0}
}

type StateJournal struct {
	Id                   string        `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedOn            int64         `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	ModifiedOn           int64         `protobuf:"varint,3,opt,name=modifiedOn,proto3" json:"modifiedOn,omitempty"`
	Title                string        `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	IsSigned             bool          `protobuf:"varint,5,opt,name=isSigned,proto3" json:"isSigned,omitempty"`
	DescriptionHash      string        `protobuf:"bytes,6,opt,name=descriptionHash,proto3" json:"descriptionHash,omitempty"`
	EditorInfo           []*EditorInfo `protobuf:"bytes,7,rep,name=editorInfo,proto3" json:"editorInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *StateJournal) Reset()         { *m = StateJournal{} }
func (m *StateJournal) String() string { return proto.CompactTextString(m) }
func (*StateJournal) ProtoMessage()    {}
func (*StateJournal) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{0}
}

func (m *StateJournal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateJournal.Unmarshal(m, b)
}
func (m *StateJournal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateJournal.Marshal(b, m, deterministic)
}
func (m *StateJournal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateJournal.Merge(m, src)
}
func (m *StateJournal) XXX_Size() int {
	return xxx_messageInfo_StateJournal.Size(m)
}
func (m *StateJournal) XXX_DiscardUnknown() {
	xxx_messageInfo_StateJournal.DiscardUnknown(m)
}

var xxx_messageInfo_StateJournal proto.InternalMessageInfo

func (m *StateJournal) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StateJournal) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *StateJournal) GetModifiedOn() int64 {
	if m != nil {
		return m.ModifiedOn
	}
	return 0
}

func (m *StateJournal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *StateJournal) GetIsSigned() bool {
	if m != nil {
		return m.IsSigned
	}
	return false
}

func (m *StateJournal) GetDescriptionHash() string {
	if m != nil {
		return m.DescriptionHash
	}
	return ""
}

func (m *StateJournal) GetEditorInfo() []*EditorInfo {
	if m != nil {
		return m.EditorInfo
	}
	return nil
}

type EditorInfo struct {
	EditorId             string      `protobuf:"bytes,1,opt,name=editorId,proto3" json:"editorId,omitempty"`
	EditorState          EditorState `protobuf:"varint,2,opt,name=editorState,proto3,enum=EditorState" json:"editorState,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *EditorInfo) Reset()         { *m = EditorInfo{} }
func (m *EditorInfo) String() string { return proto.CompactTextString(m) }
func (*EditorInfo) ProtoMessage()    {}
func (*EditorInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{1}
}

func (m *EditorInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EditorInfo.Unmarshal(m, b)
}
func (m *EditorInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EditorInfo.Marshal(b, m, deterministic)
}
func (m *EditorInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EditorInfo.Merge(m, src)
}
func (m *EditorInfo) XXX_Size() int {
	return xxx_messageInfo_EditorInfo.Size(m)
}
func (m *EditorInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EditorInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EditorInfo proto.InternalMessageInfo

func (m *EditorInfo) GetEditorId() string {
	if m != nil {
		return m.EditorId
	}
	return ""
}

func (m *EditorInfo) GetEditorState() EditorState {
	if m != nil {
		return m.EditorState
	}
	return EditorState_editorProposed
}

type CommandJournalCreate struct {
	JournalId            string   `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	Title                string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	DescriptionHash      string   `protobuf:"bytes,3,opt,name=descriptionHash,proto3" json:"descriptionHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandJournalCreate) Reset()         { *m = CommandJournalCreate{} }
func (m *CommandJournalCreate) String() string { return proto.CompactTextString(m) }
func (*CommandJournalCreate) ProtoMessage()    {}
func (*CommandJournalCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{2}
}

func (m *CommandJournalCreate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalCreate.Unmarshal(m, b)
}
func (m *CommandJournalCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalCreate.Marshal(b, m, deterministic)
}
func (m *CommandJournalCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalCreate.Merge(m, src)
}
func (m *CommandJournalCreate) XXX_Size() int {
	return xxx_messageInfo_CommandJournalCreate.Size(m)
}
func (m *CommandJournalCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalCreate.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalCreate proto.InternalMessageInfo

func (m *CommandJournalCreate) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *CommandJournalCreate) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CommandJournalCreate) GetDescriptionHash() string {
	if m != nil {
		return m.DescriptionHash
	}
	return ""
}

type CommandJournalUpdateProperties struct {
	JournalId             string        `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	TitleUpdate           *StringUpdate `protobuf:"bytes,2,opt,name=titleUpdate,proto3" json:"titleUpdate,omitempty"`
	DescriptionHashUpdate *StringUpdate `protobuf:"bytes,3,opt,name=descriptionHashUpdate,proto3" json:"descriptionHashUpdate,omitempty"`
	XXX_NoUnkeyedLiteral  struct{}      `json:"-"`
	XXX_unrecognized      []byte        `json:"-"`
	XXX_sizecache         int32         `json:"-"`
}

func (m *CommandJournalUpdateProperties) Reset()         { *m = CommandJournalUpdateProperties{} }
func (m *CommandJournalUpdateProperties) String() string { return proto.CompactTextString(m) }
func (*CommandJournalUpdateProperties) ProtoMessage()    {}
func (*CommandJournalUpdateProperties) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{3}
}

func (m *CommandJournalUpdateProperties) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalUpdateProperties.Unmarshal(m, b)
}
func (m *CommandJournalUpdateProperties) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalUpdateProperties.Marshal(b, m, deterministic)
}
func (m *CommandJournalUpdateProperties) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalUpdateProperties.Merge(m, src)
}
func (m *CommandJournalUpdateProperties) XXX_Size() int {
	return xxx_messageInfo_CommandJournalUpdateProperties.Size(m)
}
func (m *CommandJournalUpdateProperties) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalUpdateProperties.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalUpdateProperties proto.InternalMessageInfo

func (m *CommandJournalUpdateProperties) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *CommandJournalUpdateProperties) GetTitleUpdate() *StringUpdate {
	if m != nil {
		return m.TitleUpdate
	}
	return nil
}

func (m *CommandJournalUpdateProperties) GetDescriptionHashUpdate() *StringUpdate {
	if m != nil {
		return m.DescriptionHashUpdate
	}
	return nil
}

type CommandJournalUpdateAuthorization struct {
	JournalId            string   `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	MakeSigned           bool     `protobuf:"varint,2,opt,name=makeSigned,proto3" json:"makeSigned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandJournalUpdateAuthorization) Reset()         { *m = CommandJournalUpdateAuthorization{} }
func (m *CommandJournalUpdateAuthorization) String() string { return proto.CompactTextString(m) }
func (*CommandJournalUpdateAuthorization) ProtoMessage()    {}
func (*CommandJournalUpdateAuthorization) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{4}
}

func (m *CommandJournalUpdateAuthorization) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalUpdateAuthorization.Unmarshal(m, b)
}
func (m *CommandJournalUpdateAuthorization) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalUpdateAuthorization.Marshal(b, m, deterministic)
}
func (m *CommandJournalUpdateAuthorization) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalUpdateAuthorization.Merge(m, src)
}
func (m *CommandJournalUpdateAuthorization) XXX_Size() int {
	return xxx_messageInfo_CommandJournalUpdateAuthorization.Size(m)
}
func (m *CommandJournalUpdateAuthorization) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalUpdateAuthorization.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalUpdateAuthorization proto.InternalMessageInfo

func (m *CommandJournalUpdateAuthorization) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *CommandJournalUpdateAuthorization) GetMakeSigned() bool {
	if m != nil {
		return m.MakeSigned
	}
	return false
}

type CommandJournalEditorResign struct {
	JournalId            string   `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandJournalEditorResign) Reset()         { *m = CommandJournalEditorResign{} }
func (m *CommandJournalEditorResign) String() string { return proto.CompactTextString(m) }
func (*CommandJournalEditorResign) ProtoMessage()    {}
func (*CommandJournalEditorResign) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{5}
}

func (m *CommandJournalEditorResign) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalEditorResign.Unmarshal(m, b)
}
func (m *CommandJournalEditorResign) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalEditorResign.Marshal(b, m, deterministic)
}
func (m *CommandJournalEditorResign) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalEditorResign.Merge(m, src)
}
func (m *CommandJournalEditorResign) XXX_Size() int {
	return xxx_messageInfo_CommandJournalEditorResign.Size(m)
}
func (m *CommandJournalEditorResign) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalEditorResign.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalEditorResign proto.InternalMessageInfo

func (m *CommandJournalEditorResign) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

type CommandJournalEditorInvite struct {
	JournalId            string   `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	InvitedEditorId      string   `protobuf:"bytes,2,opt,name=invitedEditorId,proto3" json:"invitedEditorId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandJournalEditorInvite) Reset()         { *m = CommandJournalEditorInvite{} }
func (m *CommandJournalEditorInvite) String() string { return proto.CompactTextString(m) }
func (*CommandJournalEditorInvite) ProtoMessage()    {}
func (*CommandJournalEditorInvite) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{6}
}

func (m *CommandJournalEditorInvite) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalEditorInvite.Unmarshal(m, b)
}
func (m *CommandJournalEditorInvite) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalEditorInvite.Marshal(b, m, deterministic)
}
func (m *CommandJournalEditorInvite) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalEditorInvite.Merge(m, src)
}
func (m *CommandJournalEditorInvite) XXX_Size() int {
	return xxx_messageInfo_CommandJournalEditorInvite.Size(m)
}
func (m *CommandJournalEditorInvite) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalEditorInvite.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalEditorInvite proto.InternalMessageInfo

func (m *CommandJournalEditorInvite) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *CommandJournalEditorInvite) GetInvitedEditorId() string {
	if m != nil {
		return m.InvitedEditorId
	}
	return ""
}

type CommandJournalEditorAcceptDuty struct {
	JournalId            string   `protobuf:"bytes,1,opt,name=journalId,proto3" json:"journalId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandJournalEditorAcceptDuty) Reset()         { *m = CommandJournalEditorAcceptDuty{} }
func (m *CommandJournalEditorAcceptDuty) String() string { return proto.CompactTextString(m) }
func (*CommandJournalEditorAcceptDuty) ProtoMessage()    {}
func (*CommandJournalEditorAcceptDuty) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{7}
}

func (m *CommandJournalEditorAcceptDuty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandJournalEditorAcceptDuty.Unmarshal(m, b)
}
func (m *CommandJournalEditorAcceptDuty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandJournalEditorAcceptDuty.Marshal(b, m, deterministic)
}
func (m *CommandJournalEditorAcceptDuty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandJournalEditorAcceptDuty.Merge(m, src)
}
func (m *CommandJournalEditorAcceptDuty) XXX_Size() int {
	return xxx_messageInfo_CommandJournalEditorAcceptDuty.Size(m)
}
func (m *CommandJournalEditorAcceptDuty) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandJournalEditorAcceptDuty.DiscardUnknown(m)
}

var xxx_messageInfo_CommandJournalEditorAcceptDuty proto.InternalMessageInfo

func (m *CommandJournalEditorAcceptDuty) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

type StateVolume struct {
	Id                     string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedOn              int64    `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	JournalId              string   `protobuf:"bytes,3,opt,name=journalId,proto3" json:"journalId,omitempty"`
	Issue                  string   `protobuf:"bytes,4,opt,name=issue,proto3" json:"issue,omitempty"`
	LogicalPublicationTime int64    `protobuf:"varint,5,opt,name=logicalPublicationTime,proto3" json:"logicalPublicationTime,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *StateVolume) Reset()         { *m = StateVolume{} }
func (m *StateVolume) String() string { return proto.CompactTextString(m) }
func (*StateVolume) ProtoMessage()    {}
func (*StateVolume) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{8}
}

func (m *StateVolume) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateVolume.Unmarshal(m, b)
}
func (m *StateVolume) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateVolume.Marshal(b, m, deterministic)
}
func (m *StateVolume) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateVolume.Merge(m, src)
}
func (m *StateVolume) XXX_Size() int {
	return xxx_messageInfo_StateVolume.Size(m)
}
func (m *StateVolume) XXX_DiscardUnknown() {
	xxx_messageInfo_StateVolume.DiscardUnknown(m)
}

var xxx_messageInfo_StateVolume proto.InternalMessageInfo

func (m *StateVolume) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StateVolume) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *StateVolume) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *StateVolume) GetIssue() string {
	if m != nil {
		return m.Issue
	}
	return ""
}

func (m *StateVolume) GetLogicalPublicationTime() int64 {
	if m != nil {
		return m.LogicalPublicationTime
	}
	return 0
}

type CommandVolumeCreate struct {
	VolumeId               string   `protobuf:"bytes,1,opt,name=volumeId,proto3" json:"volumeId,omitempty"`
	JournalId              string   `protobuf:"bytes,2,opt,name=journalId,proto3" json:"journalId,omitempty"`
	Issue                  string   `protobuf:"bytes,3,opt,name=issue,proto3" json:"issue,omitempty"`
	LogicalPublicationTime int64    `protobuf:"varint,4,opt,name=logicalPublicationTime,proto3" json:"logicalPublicationTime,omitempty"`
	XXX_NoUnkeyedLiteral   struct{} `json:"-"`
	XXX_unrecognized       []byte   `json:"-"`
	XXX_sizecache          int32    `json:"-"`
}

func (m *CommandVolumeCreate) Reset()         { *m = CommandVolumeCreate{} }
func (m *CommandVolumeCreate) String() string { return proto.CompactTextString(m) }
func (*CommandVolumeCreate) ProtoMessage()    {}
func (*CommandVolumeCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_04fd98cceb1b9191, []int{9}
}

func (m *CommandVolumeCreate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandVolumeCreate.Unmarshal(m, b)
}
func (m *CommandVolumeCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandVolumeCreate.Marshal(b, m, deterministic)
}
func (m *CommandVolumeCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandVolumeCreate.Merge(m, src)
}
func (m *CommandVolumeCreate) XXX_Size() int {
	return xxx_messageInfo_CommandVolumeCreate.Size(m)
}
func (m *CommandVolumeCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandVolumeCreate.DiscardUnknown(m)
}

var xxx_messageInfo_CommandVolumeCreate proto.InternalMessageInfo

func (m *CommandVolumeCreate) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *CommandVolumeCreate) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *CommandVolumeCreate) GetIssue() string {
	if m != nil {
		return m.Issue
	}
	return ""
}

func (m *CommandVolumeCreate) GetLogicalPublicationTime() int64 {
	if m != nil {
		return m.LogicalPublicationTime
	}
	return 0
}

func init() {
	proto.RegisterEnum("EditorState", EditorState_name, EditorState_value)
	proto.RegisterType((*StateJournal)(nil), "StateJournal")
	proto.RegisterType((*EditorInfo)(nil), "EditorInfo")
	proto.RegisterType((*CommandJournalCreate)(nil), "CommandJournalCreate")
	proto.RegisterType((*CommandJournalUpdateProperties)(nil), "CommandJournalUpdateProperties")
	proto.RegisterType((*CommandJournalUpdateAuthorization)(nil), "CommandJournalUpdateAuthorization")
	proto.RegisterType((*CommandJournalEditorResign)(nil), "CommandJournalEditorResign")
	proto.RegisterType((*CommandJournalEditorInvite)(nil), "CommandJournalEditorInvite")
	proto.RegisterType((*CommandJournalEditorAcceptDuty)(nil), "CommandJournalEditorAcceptDuty")
	proto.RegisterType((*StateVolume)(nil), "StateVolume")
	proto.RegisterType((*CommandVolumeCreate)(nil), "CommandVolumeCreate")
}

func init() { proto.RegisterFile("journal.proto", fileDescriptor_04fd98cceb1b9191) }

var fileDescriptor_04fd98cceb1b9191 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x76, 0xd3, 0xa6, 0xe3, 0x34, 0xad, 0x96, 0x82, 0xac, 0x08, 0x55, 0xc1, 0x27, 0x0b,
	0xa4, 0x20, 0x05, 0xc1, 0x81, 0x03, 0x52, 0x09, 0x95, 0x28, 0x17, 0x2a, 0x07, 0x10, 0xe2, 0xb6,
	0xf5, 0x6e, 0xd3, 0x01, 0xdb, 0x6b, 0xad, 0xd7, 0x95, 0xe0, 0x5d, 0x38, 0xf3, 0x02, 0xbc, 0x12,
	0xef, 0x81, 0xbc, 0xeb, 0xd8, 0x4e, 0xe4, 0x62, 0xb8, 0x79, 0xbe, 0x99, 0xd9, 0xf9, 0xe6, 0xe7,
	0x33, 0x1c, 0x7c, 0x11, 0x85, 0x4c, 0x69, 0x3c, 0xcb, 0xa4, 0x50, 0x62, 0x32, 0x8a, 0x44, 0x92,
	0x88, 0xd4, 0x58, 0xfe, 0x6f, 0x0b, 0x46, 0x4b, 0x45, 0x15, 0x7f, 0x6b, 0x82, 0xc8, 0x18, 0x6c,
	0x64, 0x9e, 0x35, 0xb5, 0x82, 0xfd, 0xd0, 0x46, 0x46, 0x1e, 0xc0, 0x7e, 0x24, 0x39, 0x55, 0x9c,
	0xbd, 0x4b, 0x3d, 0x7b, 0x6a, 0x05, 0x4e, 0xd8, 0x00, 0xe4, 0x04, 0x20, 0x11, 0x0c, 0xaf, 0x50,
	0xbb, 0x1d, 0xed, 0x6e, 0x21, 0xe4, 0x18, 0x06, 0x0a, 0x55, 0xcc, 0xbd, 0x1d, 0xfd, 0xa0, 0x31,
	0xc8, 0x04, 0x86, 0x98, 0x2f, 0x71, 0x95, 0x72, 0xe6, 0x0d, 0xa6, 0x56, 0x30, 0x0c, 0x6b, 0x9b,
	0x04, 0x70, 0xc8, 0x78, 0x1e, 0x49, 0xcc, 0x14, 0x8a, 0xf4, 0x0d, 0xcd, 0xaf, 0xbd, 0x5d, 0x9d,
	0xbb, 0x0d, 0x93, 0xc7, 0x00, 0x9c, 0xa1, 0x12, 0xf2, 0x3c, 0xbd, 0x12, 0xde, 0xde, 0xd4, 0x09,
	0xdc, 0xb9, 0x3b, 0x3b, 0xab, 0xa1, 0xb0, 0xe5, 0xf6, 0x3f, 0x01, 0x34, 0x9e, 0x92, 0x40, 0xe5,
	0x5b, 0xb7, 0x5a, 0xdb, 0x64, 0x06, 0xae, 0xf9, 0xd6, 0x63, 0xd1, 0x2d, 0x8f, 0xe7, 0xa3, 0xea,
	0x5d, 0x8d, 0x85, 0xed, 0x00, 0x5f, 0xc1, 0xf1, 0x42, 0x24, 0x09, 0x4d, 0x59, 0x35, 0xc2, 0x85,
	0x9e, 0x4e, 0x39, 0xb8, 0x6a, 0xf0, 0x75, 0x91, 0x06, 0x68, 0x06, 0x63, 0xb7, 0x07, 0xd3, 0xd1,
	0xbc, 0xd3, 0xd9, 0xbc, 0xff, 0xcb, 0x82, 0x93, 0xcd, 0xb2, 0x1f, 0x32, 0x46, 0x15, 0xbf, 0x90,
	0x22, 0xe3, 0x52, 0x21, 0xcf, 0x7b, 0x08, 0x3c, 0x01, 0x57, 0xd7, 0x34, 0x69, 0x9a, 0x86, 0x3b,
	0x3f, 0x98, 0x2d, 0x95, 0xc4, 0x74, 0x65, 0xc0, 0xb0, 0x1d, 0x41, 0x16, 0x70, 0x6f, 0x8b, 0x44,
	0x95, 0xea, 0x74, 0xa5, 0x76, 0xc7, 0xfa, 0x14, 0x1e, 0x76, 0xb1, 0x3e, 0x2d, 0xd4, 0xb5, 0x90,
	0xf8, 0x9d, 0x96, 0xe1, 0x3d, 0xc4, 0xcb, 0x93, 0xa3, 0x5f, 0x79, 0x75, 0x3e, 0xb6, 0x3e, 0x9f,
	0x16, 0xe2, 0xbf, 0x80, 0xc9, 0x66, 0x09, 0xb3, 0xb9, 0x90, 0xe7, 0xb8, 0xea, 0x79, 0xdb, 0x67,
	0xdd, 0xb9, 0xe7, 0xe9, 0x0d, 0xf6, 0x6e, 0x34, 0x80, 0x43, 0xd4, 0x71, 0xec, 0x6c, 0x7d, 0x5a,
	0x66, 0xb7, 0xdb, 0xb0, 0xff, 0x72, 0x7b, 0x75, 0xc6, 0x73, 0x1a, 0x45, 0x3c, 0x53, 0xaf, 0x0b,
	0xf5, 0xad, 0x87, 0xe5, 0x4f, 0x0b, 0x5c, 0x7d, 0x7b, 0x1f, 0x45, 0x5c, 0x24, 0xfc, 0x3f, 0x25,
	0xbb, 0xf1, 0xb6, 0xd3, 0x71, 0x97, 0x98, 0xe7, 0x45, 0x2d, 0x58, 0x6d, 0x90, 0xe7, 0x70, 0x3f,
	0x16, 0x2b, 0x8c, 0x68, 0x7c, 0x51, 0x5c, 0xc6, 0x18, 0xe9, 0x3d, 0xbd, 0xc7, 0x84, 0x6b, 0xf9,
	0x3a, 0xe1, 0x2d, 0x5e, 0xff, 0x87, 0x05, 0x77, 0xab, 0x56, 0x0d, 0xd7, 0x4a, 0x1b, 0x13, 0x18,
	0xde, 0x68, 0xbb, 0xd1, 0xdf, 0xda, 0xde, 0xe4, 0x67, 0xdf, 0xca, 0xcf, 0xf9, 0x37, 0x7e, 0x3b,
	0x7f, 0xe3, 0xf7, 0xe8, 0x19, 0xb8, 0x2d, 0x5d, 0x13, 0x02, 0x63, 0xa3, 0xec, 0x52, 0x45, 0x22,
	0xe7, 0xec, 0xe8, 0x4e, 0x83, 0x99, 0xf5, 0x70, 0x76, 0x64, 0xbd, 0xda, 0xfb, 0x3c, 0x48, 0x04,
	0xe3, 0xf1, 0xe5, 0xae, 0xfe, 0x89, 0x3e, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0x0a, 0xff, 0x20,
	0x9d, 0x63, 0x05, 0x00, 0x00,
}
