// Code generated by protoc-gen-go. DO NOT EDIT.
// source: manuscript.proto

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

type ManuscriptStatus int32

const (
	ManuscriptStatus_init       ManuscriptStatus = 0
	ManuscriptStatus_new        ManuscriptStatus = 1
	ManuscriptStatus_reviewable ManuscriptStatus = 2
	ManuscriptStatus_rejected   ManuscriptStatus = 3
	ManuscriptStatus_published  ManuscriptStatus = 4
	ManuscriptStatus_assigned   ManuscriptStatus = 5
)

var ManuscriptStatus_name = map[int32]string{
	0: "init",
	1: "new",
	2: "reviewable",
	3: "rejected",
	4: "published",
	5: "assigned",
}

var ManuscriptStatus_value = map[string]int32{
	"init":       0,
	"new":        1,
	"reviewable": 2,
	"rejected":   3,
	"published":  4,
	"assigned":   5,
}

func (x ManuscriptStatus) String() string {
	return proto.EnumName(ManuscriptStatus_name, int32(x))
}

func (ManuscriptStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{0}
}

type ManuscriptJudgement int32

const (
	ManuscriptJudgement_judgementRejected ManuscriptJudgement = 0
	ManuscriptJudgement_judgementAccepted ManuscriptJudgement = 1
)

var ManuscriptJudgement_name = map[int32]string{
	0: "judgementRejected",
	1: "judgementAccepted",
}

var ManuscriptJudgement_value = map[string]int32{
	"judgementRejected": 0,
	"judgementAccepted": 1,
}

func (x ManuscriptJudgement) String() string {
	return proto.EnumName(ManuscriptJudgement_name, int32(x))
}

func (ManuscriptJudgement) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{1}
}

type StateManuscript struct {
	Id                   string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedOn            int64            `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	ModifiedOn           int64            `protobuf:"varint,3,opt,name=modifiedOn,proto3" json:"modifiedOn,omitempty"`
	Hash                 string           `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	ThreadId             string           `protobuf:"bytes,5,opt,name=threadId,proto3" json:"threadId,omitempty"`
	VersionNumber        int32            `protobuf:"varint,6,opt,name=versionNumber,proto3" json:"versionNumber,omitempty"`
	CommitMsg            string           `protobuf:"bytes,7,opt,name=commitMsg,proto3" json:"commitMsg,omitempty"`
	Title                string           `protobuf:"bytes,8,opt,name=title,proto3" json:"title,omitempty"`
	Author               []*Author        `protobuf:"bytes,9,rep,name=author,proto3" json:"author,omitempty"`
	Status               ManuscriptStatus `protobuf:"varint,10,opt,name=status,proto3,enum=ManuscriptStatus" json:"status,omitempty"`
	JournalId            string           `protobuf:"bytes,11,opt,name=journalId,proto3" json:"journalId,omitempty"`
	VolumeId             string           `protobuf:"bytes,12,opt,name=volumeId,proto3" json:"volumeId,omitempty"`
	FirstPage            string           `protobuf:"bytes,13,opt,name=firstPage,proto3" json:"firstPage,omitempty"`
	LastPage             string           `protobuf:"bytes,14,opt,name=lastPage,proto3" json:"lastPage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *StateManuscript) Reset()         { *m = StateManuscript{} }
func (m *StateManuscript) String() string { return proto.CompactTextString(m) }
func (*StateManuscript) ProtoMessage()    {}
func (*StateManuscript) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{0}
}

func (m *StateManuscript) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateManuscript.Unmarshal(m, b)
}
func (m *StateManuscript) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateManuscript.Marshal(b, m, deterministic)
}
func (m *StateManuscript) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateManuscript.Merge(m, src)
}
func (m *StateManuscript) XXX_Size() int {
	return xxx_messageInfo_StateManuscript.Size(m)
}
func (m *StateManuscript) XXX_DiscardUnknown() {
	xxx_messageInfo_StateManuscript.DiscardUnknown(m)
}

var xxx_messageInfo_StateManuscript proto.InternalMessageInfo

func (m *StateManuscript) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StateManuscript) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *StateManuscript) GetModifiedOn() int64 {
	if m != nil {
		return m.ModifiedOn
	}
	return 0
}

func (m *StateManuscript) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *StateManuscript) GetThreadId() string {
	if m != nil {
		return m.ThreadId
	}
	return ""
}

func (m *StateManuscript) GetVersionNumber() int32 {
	if m != nil {
		return m.VersionNumber
	}
	return 0
}

func (m *StateManuscript) GetCommitMsg() string {
	if m != nil {
		return m.CommitMsg
	}
	return ""
}

func (m *StateManuscript) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *StateManuscript) GetAuthor() []*Author {
	if m != nil {
		return m.Author
	}
	return nil
}

func (m *StateManuscript) GetStatus() ManuscriptStatus {
	if m != nil {
		return m.Status
	}
	return ManuscriptStatus_init
}

func (m *StateManuscript) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

func (m *StateManuscript) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *StateManuscript) GetFirstPage() string {
	if m != nil {
		return m.FirstPage
	}
	return ""
}

func (m *StateManuscript) GetLastPage() string {
	if m != nil {
		return m.LastPage
	}
	return ""
}

type Author struct {
	AuthorId             string   `protobuf:"bytes,1,opt,name=authorId,proto3" json:"authorId,omitempty"`
	DidSign              bool     `protobuf:"varint,2,opt,name=didSign,proto3" json:"didSign,omitempty"`
	AuthorNumber         int32    `protobuf:"varint,3,opt,name=authorNumber,proto3" json:"authorNumber,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Author) Reset()         { *m = Author{} }
func (m *Author) String() string { return proto.CompactTextString(m) }
func (*Author) ProtoMessage()    {}
func (*Author) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{1}
}

func (m *Author) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Author.Unmarshal(m, b)
}
func (m *Author) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Author.Marshal(b, m, deterministic)
}
func (m *Author) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Author.Merge(m, src)
}
func (m *Author) XXX_Size() int {
	return xxx_messageInfo_Author.Size(m)
}
func (m *Author) XXX_DiscardUnknown() {
	xxx_messageInfo_Author.DiscardUnknown(m)
}

var xxx_messageInfo_Author proto.InternalMessageInfo

func (m *Author) GetAuthorId() string {
	if m != nil {
		return m.AuthorId
	}
	return ""
}

func (m *Author) GetDidSign() bool {
	if m != nil {
		return m.DidSign
	}
	return false
}

func (m *Author) GetAuthorNumber() int32 {
	if m != nil {
		return m.AuthorNumber
	}
	return 0
}

type StateManuscriptThread struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ManuscriptId         []string `protobuf:"bytes,2,rep,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	IsReviewable         bool     `protobuf:"varint,3,opt,name=isReviewable,proto3" json:"isReviewable,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateManuscriptThread) Reset()         { *m = StateManuscriptThread{} }
func (m *StateManuscriptThread) String() string { return proto.CompactTextString(m) }
func (*StateManuscriptThread) ProtoMessage()    {}
func (*StateManuscriptThread) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{2}
}

func (m *StateManuscriptThread) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateManuscriptThread.Unmarshal(m, b)
}
func (m *StateManuscriptThread) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateManuscriptThread.Marshal(b, m, deterministic)
}
func (m *StateManuscriptThread) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateManuscriptThread.Merge(m, src)
}
func (m *StateManuscriptThread) XXX_Size() int {
	return xxx_messageInfo_StateManuscriptThread.Size(m)
}
func (m *StateManuscriptThread) XXX_DiscardUnknown() {
	xxx_messageInfo_StateManuscriptThread.DiscardUnknown(m)
}

var xxx_messageInfo_StateManuscriptThread proto.InternalMessageInfo

func (m *StateManuscriptThread) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StateManuscriptThread) GetManuscriptId() []string {
	if m != nil {
		return m.ManuscriptId
	}
	return nil
}

func (m *StateManuscriptThread) GetIsReviewable() bool {
	if m != nil {
		return m.IsReviewable
	}
	return false
}

type StateReview struct {
	Id                   string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedOn            int64     `protobuf:"varint,2,opt,name=createdOn,proto3" json:"createdOn,omitempty"`
	ManuscriptId         string    `protobuf:"bytes,3,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	ReviewAuthorId       string    `protobuf:"bytes,4,opt,name=reviewAuthorId,proto3" json:"reviewAuthorId,omitempty"`
	Hash                 string    `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	Judgement            Judgement `protobuf:"varint,6,opt,name=judgement,proto3,enum=Judgement" json:"judgement,omitempty"`
	IsUsedByEditor       bool      `protobuf:"varint,7,opt,name=isUsedByEditor,proto3" json:"isUsedByEditor,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *StateReview) Reset()         { *m = StateReview{} }
func (m *StateReview) String() string { return proto.CompactTextString(m) }
func (*StateReview) ProtoMessage()    {}
func (*StateReview) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{3}
}

func (m *StateReview) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateReview.Unmarshal(m, b)
}
func (m *StateReview) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateReview.Marshal(b, m, deterministic)
}
func (m *StateReview) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateReview.Merge(m, src)
}
func (m *StateReview) XXX_Size() int {
	return xxx_messageInfo_StateReview.Size(m)
}
func (m *StateReview) XXX_DiscardUnknown() {
	xxx_messageInfo_StateReview.DiscardUnknown(m)
}

var xxx_messageInfo_StateReview proto.InternalMessageInfo

func (m *StateReview) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *StateReview) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *StateReview) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *StateReview) GetReviewAuthorId() string {
	if m != nil {
		return m.ReviewAuthorId
	}
	return ""
}

func (m *StateReview) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *StateReview) GetJudgement() Judgement {
	if m != nil {
		return m.Judgement
	}
	return Judgement_NEGATIVE
}

func (m *StateReview) GetIsUsedByEditor() bool {
	if m != nil {
		return m.IsUsedByEditor
	}
	return false
}

type CommandManuscriptCreate struct {
	ManuscriptId         string   `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	ManuscriptThreadId   string   `protobuf:"bytes,2,opt,name=manuscriptThreadId,proto3" json:"manuscriptThreadId,omitempty"`
	Hash                 string   `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	CommitMsg            string   `protobuf:"bytes,4,opt,name=commitMsg,proto3" json:"commitMsg,omitempty"`
	Title                string   `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	AuthorId             []string `protobuf:"bytes,6,rep,name=authorId,proto3" json:"authorId,omitempty"`
	JournalId            string   `protobuf:"bytes,7,opt,name=journalId,proto3" json:"journalId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandManuscriptCreate) Reset()         { *m = CommandManuscriptCreate{} }
func (m *CommandManuscriptCreate) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptCreate) ProtoMessage()    {}
func (*CommandManuscriptCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{4}
}

func (m *CommandManuscriptCreate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptCreate.Unmarshal(m, b)
}
func (m *CommandManuscriptCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptCreate.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptCreate.Merge(m, src)
}
func (m *CommandManuscriptCreate) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptCreate.Size(m)
}
func (m *CommandManuscriptCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptCreate.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptCreate proto.InternalMessageInfo

func (m *CommandManuscriptCreate) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandManuscriptCreate) GetManuscriptThreadId() string {
	if m != nil {
		return m.ManuscriptThreadId
	}
	return ""
}

func (m *CommandManuscriptCreate) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *CommandManuscriptCreate) GetCommitMsg() string {
	if m != nil {
		return m.CommitMsg
	}
	return ""
}

func (m *CommandManuscriptCreate) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CommandManuscriptCreate) GetAuthorId() []string {
	if m != nil {
		return m.AuthorId
	}
	return nil
}

func (m *CommandManuscriptCreate) GetJournalId() string {
	if m != nil {
		return m.JournalId
	}
	return ""
}

type CommandManuscriptCreateNewVersion struct {
	ManuscriptId         string   `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	PreviousManuscriptId string   `protobuf:"bytes,2,opt,name=previousManuscriptId,proto3" json:"previousManuscriptId,omitempty"`
	Hash                 string   `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	CommitMsg            string   `protobuf:"bytes,4,opt,name=commitMsg,proto3" json:"commitMsg,omitempty"`
	Title                string   `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	AuthorId             []string `protobuf:"bytes,6,rep,name=authorId,proto3" json:"authorId,omitempty"`
	HistoricAuthorId     []string `protobuf:"bytes,7,rep,name=historicAuthorId,proto3" json:"historicAuthorId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandManuscriptCreateNewVersion) Reset()         { *m = CommandManuscriptCreateNewVersion{} }
func (m *CommandManuscriptCreateNewVersion) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptCreateNewVersion) ProtoMessage()    {}
func (*CommandManuscriptCreateNewVersion) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{5}
}

func (m *CommandManuscriptCreateNewVersion) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptCreateNewVersion.Unmarshal(m, b)
}
func (m *CommandManuscriptCreateNewVersion) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptCreateNewVersion.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptCreateNewVersion) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptCreateNewVersion.Merge(m, src)
}
func (m *CommandManuscriptCreateNewVersion) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptCreateNewVersion.Size(m)
}
func (m *CommandManuscriptCreateNewVersion) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptCreateNewVersion.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptCreateNewVersion proto.InternalMessageInfo

func (m *CommandManuscriptCreateNewVersion) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandManuscriptCreateNewVersion) GetPreviousManuscriptId() string {
	if m != nil {
		return m.PreviousManuscriptId
	}
	return ""
}

func (m *CommandManuscriptCreateNewVersion) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *CommandManuscriptCreateNewVersion) GetCommitMsg() string {
	if m != nil {
		return m.CommitMsg
	}
	return ""
}

func (m *CommandManuscriptCreateNewVersion) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CommandManuscriptCreateNewVersion) GetAuthorId() []string {
	if m != nil {
		return m.AuthorId
	}
	return nil
}

func (m *CommandManuscriptCreateNewVersion) GetHistoricAuthorId() []string {
	if m != nil {
		return m.HistoricAuthorId
	}
	return nil
}

type CommandManuscriptAcceptAuthorship struct {
	ManuscriptId         string    `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	Author               []*Author `protobuf:"bytes,2,rep,name=author,proto3" json:"author,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CommandManuscriptAcceptAuthorship) Reset()         { *m = CommandManuscriptAcceptAuthorship{} }
func (m *CommandManuscriptAcceptAuthorship) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptAcceptAuthorship) ProtoMessage()    {}
func (*CommandManuscriptAcceptAuthorship) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{6}
}

func (m *CommandManuscriptAcceptAuthorship) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptAcceptAuthorship.Unmarshal(m, b)
}
func (m *CommandManuscriptAcceptAuthorship) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptAcceptAuthorship.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptAcceptAuthorship) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptAcceptAuthorship.Merge(m, src)
}
func (m *CommandManuscriptAcceptAuthorship) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptAcceptAuthorship.Size(m)
}
func (m *CommandManuscriptAcceptAuthorship) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptAcceptAuthorship.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptAcceptAuthorship proto.InternalMessageInfo

func (m *CommandManuscriptAcceptAuthorship) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandManuscriptAcceptAuthorship) GetAuthor() []*Author {
	if m != nil {
		return m.Author
	}
	return nil
}

type CommandManuscriptAllowReview struct {
	ThreadId             string                 `protobuf:"bytes,1,opt,name=ThreadId,proto3" json:"ThreadId,omitempty"`
	ThreadReference      []*ThreadReferenceItem `protobuf:"bytes,2,rep,name=threadReference,proto3" json:"threadReference,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *CommandManuscriptAllowReview) Reset()         { *m = CommandManuscriptAllowReview{} }
func (m *CommandManuscriptAllowReview) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptAllowReview) ProtoMessage()    {}
func (*CommandManuscriptAllowReview) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{7}
}

func (m *CommandManuscriptAllowReview) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptAllowReview.Unmarshal(m, b)
}
func (m *CommandManuscriptAllowReview) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptAllowReview.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptAllowReview) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptAllowReview.Merge(m, src)
}
func (m *CommandManuscriptAllowReview) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptAllowReview.Size(m)
}
func (m *CommandManuscriptAllowReview) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptAllowReview.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptAllowReview proto.InternalMessageInfo

func (m *CommandManuscriptAllowReview) GetThreadId() string {
	if m != nil {
		return m.ThreadId
	}
	return ""
}

func (m *CommandManuscriptAllowReview) GetThreadReference() []*ThreadReferenceItem {
	if m != nil {
		return m.ThreadReference
	}
	return nil
}

type ThreadReferenceItem struct {
	ManuscriptId         string           `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	ManuscriptStatus     ManuscriptStatus `protobuf:"varint,2,opt,name=manuscriptStatus,proto3,enum=ManuscriptStatus" json:"manuscriptStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ThreadReferenceItem) Reset()         { *m = ThreadReferenceItem{} }
func (m *ThreadReferenceItem) String() string { return proto.CompactTextString(m) }
func (*ThreadReferenceItem) ProtoMessage()    {}
func (*ThreadReferenceItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{8}
}

func (m *ThreadReferenceItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ThreadReferenceItem.Unmarshal(m, b)
}
func (m *ThreadReferenceItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ThreadReferenceItem.Marshal(b, m, deterministic)
}
func (m *ThreadReferenceItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThreadReferenceItem.Merge(m, src)
}
func (m *ThreadReferenceItem) XXX_Size() int {
	return xxx_messageInfo_ThreadReferenceItem.Size(m)
}
func (m *ThreadReferenceItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ThreadReferenceItem.DiscardUnknown(m)
}

var xxx_messageInfo_ThreadReferenceItem proto.InternalMessageInfo

func (m *ThreadReferenceItem) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *ThreadReferenceItem) GetManuscriptStatus() ManuscriptStatus {
	if m != nil {
		return m.ManuscriptStatus
	}
	return ManuscriptStatus_init
}

type CommandWriteReview struct {
	ReviewId             string    `protobuf:"bytes,1,opt,name=reviewId,proto3" json:"reviewId,omitempty"`
	ManuscriptId         string    `protobuf:"bytes,2,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	Hash                 string    `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Judgement            Judgement `protobuf:"varint,4,opt,name=judgement,proto3,enum=Judgement" json:"judgement,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CommandWriteReview) Reset()         { *m = CommandWriteReview{} }
func (m *CommandWriteReview) String() string { return proto.CompactTextString(m) }
func (*CommandWriteReview) ProtoMessage()    {}
func (*CommandWriteReview) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{9}
}

func (m *CommandWriteReview) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandWriteReview.Unmarshal(m, b)
}
func (m *CommandWriteReview) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandWriteReview.Marshal(b, m, deterministic)
}
func (m *CommandWriteReview) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandWriteReview.Merge(m, src)
}
func (m *CommandWriteReview) XXX_Size() int {
	return xxx_messageInfo_CommandWriteReview.Size(m)
}
func (m *CommandWriteReview) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandWriteReview.DiscardUnknown(m)
}

var xxx_messageInfo_CommandWriteReview proto.InternalMessageInfo

func (m *CommandWriteReview) GetReviewId() string {
	if m != nil {
		return m.ReviewId
	}
	return ""
}

func (m *CommandWriteReview) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandWriteReview) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *CommandWriteReview) GetJudgement() Judgement {
	if m != nil {
		return m.Judgement
	}
	return Judgement_NEGATIVE
}

type CommandManuscriptJudge struct {
	ManuscriptId         string              `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	ReviewId             []string            `protobuf:"bytes,2,rep,name=reviewId,proto3" json:"reviewId,omitempty"`
	Judgement            ManuscriptJudgement `protobuf:"varint,3,opt,name=judgement,proto3,enum=ManuscriptJudgement" json:"judgement,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CommandManuscriptJudge) Reset()         { *m = CommandManuscriptJudge{} }
func (m *CommandManuscriptJudge) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptJudge) ProtoMessage()    {}
func (*CommandManuscriptJudge) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{10}
}

func (m *CommandManuscriptJudge) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptJudge.Unmarshal(m, b)
}
func (m *CommandManuscriptJudge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptJudge.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptJudge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptJudge.Merge(m, src)
}
func (m *CommandManuscriptJudge) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptJudge.Size(m)
}
func (m *CommandManuscriptJudge) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptJudge.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptJudge proto.InternalMessageInfo

func (m *CommandManuscriptJudge) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandManuscriptJudge) GetReviewId() []string {
	if m != nil {
		return m.ReviewId
	}
	return nil
}

func (m *CommandManuscriptJudge) GetJudgement() ManuscriptJudgement {
	if m != nil {
		return m.Judgement
	}
	return ManuscriptJudgement_judgementRejected
}

type CommandManuscriptAssign struct {
	ManuscriptId         string   `protobuf:"bytes,1,opt,name=manuscriptId,proto3" json:"manuscriptId,omitempty"`
	VolumeId             string   `protobuf:"bytes,2,opt,name=volumeId,proto3" json:"volumeId,omitempty"`
	FirstPage            string   `protobuf:"bytes,3,opt,name=firstPage,proto3" json:"firstPage,omitempty"`
	LastPage             string   `protobuf:"bytes,4,opt,name=lastPage,proto3" json:"lastPage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CommandManuscriptAssign) Reset()         { *m = CommandManuscriptAssign{} }
func (m *CommandManuscriptAssign) String() string { return proto.CompactTextString(m) }
func (*CommandManuscriptAssign) ProtoMessage()    {}
func (*CommandManuscriptAssign) Descriptor() ([]byte, []int) {
	return fileDescriptor_fb127795525a7311, []int{11}
}

func (m *CommandManuscriptAssign) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommandManuscriptAssign.Unmarshal(m, b)
}
func (m *CommandManuscriptAssign) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommandManuscriptAssign.Marshal(b, m, deterministic)
}
func (m *CommandManuscriptAssign) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommandManuscriptAssign.Merge(m, src)
}
func (m *CommandManuscriptAssign) XXX_Size() int {
	return xxx_messageInfo_CommandManuscriptAssign.Size(m)
}
func (m *CommandManuscriptAssign) XXX_DiscardUnknown() {
	xxx_messageInfo_CommandManuscriptAssign.DiscardUnknown(m)
}

var xxx_messageInfo_CommandManuscriptAssign proto.InternalMessageInfo

func (m *CommandManuscriptAssign) GetManuscriptId() string {
	if m != nil {
		return m.ManuscriptId
	}
	return ""
}

func (m *CommandManuscriptAssign) GetVolumeId() string {
	if m != nil {
		return m.VolumeId
	}
	return ""
}

func (m *CommandManuscriptAssign) GetFirstPage() string {
	if m != nil {
		return m.FirstPage
	}
	return ""
}

func (m *CommandManuscriptAssign) GetLastPage() string {
	if m != nil {
		return m.LastPage
	}
	return ""
}

func init() {
	proto.RegisterEnum("ManuscriptStatus", ManuscriptStatus_name, ManuscriptStatus_value)
	proto.RegisterEnum("ManuscriptJudgement", ManuscriptJudgement_name, ManuscriptJudgement_value)
	proto.RegisterType((*StateManuscript)(nil), "StateManuscript")
	proto.RegisterType((*Author)(nil), "Author")
	proto.RegisterType((*StateManuscriptThread)(nil), "StateManuscriptThread")
	proto.RegisterType((*StateReview)(nil), "StateReview")
	proto.RegisterType((*CommandManuscriptCreate)(nil), "CommandManuscriptCreate")
	proto.RegisterType((*CommandManuscriptCreateNewVersion)(nil), "CommandManuscriptCreateNewVersion")
	proto.RegisterType((*CommandManuscriptAcceptAuthorship)(nil), "CommandManuscriptAcceptAuthorship")
	proto.RegisterType((*CommandManuscriptAllowReview)(nil), "CommandManuscriptAllowReview")
	proto.RegisterType((*ThreadReferenceItem)(nil), "ThreadReferenceItem")
	proto.RegisterType((*CommandWriteReview)(nil), "CommandWriteReview")
	proto.RegisterType((*CommandManuscriptJudge)(nil), "CommandManuscriptJudge")
	proto.RegisterType((*CommandManuscriptAssign)(nil), "CommandManuscriptAssign")
}

func init() { proto.RegisterFile("manuscript.proto", fileDescriptor_fb127795525a7311) }

var fileDescriptor_fb127795525a7311 = []byte{
	// 839 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0xcd, 0x6e, 0xeb, 0x44,
	0x14, 0xbe, 0xb6, 0xf3, 0xe3, 0x9c, 0xa6, 0xb9, 0xbe, 0x73, 0x73, 0xc1, 0x8a, 0x2a, 0x08, 0x16,
	0x42, 0xa1, 0x8b, 0x2c, 0xc2, 0x1a, 0xa4, 0xb4, 0x62, 0x11, 0xa4, 0x16, 0xe4, 0x16, 0x90, 0x58,
	0xe1, 0x78, 0xa6, 0xf1, 0x54, 0xb6, 0x27, 0xf2, 0x8c, 0x1b, 0xe0, 0x15, 0x60, 0x05, 0x4f, 0xc1,
	0xc3, 0x21, 0x5e, 0x01, 0x79, 0xc6, 0xff, 0x71, 0xab, 0xb0, 0xb9, 0xbb, 0x9c, 0xef, 0x9c, 0xf1,
	0xf9, 0xfb, 0xce, 0x39, 0x01, 0x2b, 0xf2, 0xe2, 0x94, 0xfb, 0x09, 0xdd, 0x8b, 0xe5, 0x3e, 0x61,
	0x82, 0xcd, 0xc6, 0x3e, 0x8b, 0x22, 0x16, 0x2b, 0xc9, 0xf9, 0xdb, 0x80, 0xd7, 0x77, 0xc2, 0x13,
	0xe4, 0xa6, 0xb4, 0x43, 0x13, 0xd0, 0x29, 0xb6, 0xb5, 0xb9, 0xb6, 0x18, 0xb9, 0x3a, 0xc5, 0xe8,
	0x02, 0x46, 0x7e, 0x42, 0x3c, 0x41, 0xf0, 0xb7, 0xb1, 0xad, 0xcf, 0xb5, 0x85, 0xe1, 0x56, 0x00,
	0xfa, 0x08, 0x20, 0x62, 0x98, 0x3e, 0x50, 0xa9, 0x36, 0xa4, 0xba, 0x86, 0x20, 0x04, 0xbd, 0xc0,
	0xe3, 0x81, 0xdd, 0x93, 0xdf, 0x93, 0xbf, 0xd1, 0x0c, 0x4c, 0x11, 0x24, 0xc4, 0xc3, 0x1b, 0x6c,
	0xf7, 0x25, 0x5e, 0xca, 0xe8, 0x53, 0x38, 0x7f, 0x22, 0x09, 0xa7, 0x2c, 0xbe, 0x4d, 0xa3, 0x2d,
	0x49, 0xec, 0xc1, 0x5c, 0x5b, 0xf4, 0xdd, 0x26, 0x28, 0x63, 0x62, 0x51, 0x44, 0xc5, 0x0d, 0xdf,
	0xd9, 0x43, 0xf9, 0x89, 0x0a, 0x40, 0x53, 0xe8, 0x0b, 0x2a, 0x42, 0x62, 0x9b, 0x52, 0xa3, 0x04,
	0xf4, 0x31, 0x0c, 0xbc, 0x54, 0x04, 0x2c, 0xb1, 0x47, 0x73, 0x63, 0x71, 0xb6, 0x1a, 0x2e, 0xd7,
	0x52, 0x74, 0x73, 0x18, 0x7d, 0x0e, 0x03, 0x2e, 0x3c, 0x91, 0x72, 0x1b, 0xe6, 0xda, 0x62, 0xb2,
	0x7a, 0xb3, 0xac, 0xaa, 0x72, 0x27, 0x15, 0x6e, 0x6e, 0x90, 0xf9, 0x7f, 0x64, 0x69, 0x12, 0x7b,
	0xe1, 0x06, 0xdb, 0x67, 0xca, 0x7f, 0x09, 0x64, 0xf9, 0x3d, 0xb1, 0x30, 0x8d, 0xc8, 0x06, 0xdb,
	0x63, 0x95, 0x5f, 0x21, 0x67, 0x2f, 0x1f, 0x68, 0xc2, 0xc5, 0x77, 0xde, 0x8e, 0xd8, 0xe7, 0xea,
	0x65, 0x09, 0x64, 0x2f, 0x43, 0x2f, 0x57, 0x4e, 0xd4, 0xcb, 0x42, 0x76, 0xb6, 0x30, 0x50, 0x01,
	0x67, 0x56, 0x2a, 0xe4, 0x4d, 0xd1, 0xa7, 0x52, 0x46, 0x36, 0x0c, 0x31, 0xc5, 0x77, 0x74, 0xa7,
	0x7a, 0x65, 0xba, 0x85, 0x88, 0x1c, 0x18, 0x2b, 0xab, 0xbc, 0xb0, 0x86, 0x2c, 0x6c, 0x03, 0x73,
	0x18, 0xbc, 0x6b, 0xd1, 0xe1, 0x5e, 0x36, 0xe6, 0x88, 0x14, 0x0e, 0x8c, 0x2b, 0x6a, 0x6d, 0xb0,
	0xad, 0xcf, 0x8d, 0xc5, 0xc8, 0x6d, 0x60, 0x99, 0x0d, 0xe5, 0x2e, 0x79, 0xa2, 0xe4, 0xe0, 0x6d,
	0x43, 0x22, 0x1d, 0x9a, 0x6e, 0x03, 0x73, 0xfe, 0xd1, 0xe0, 0x4c, 0x7a, 0x54, 0xd8, 0xff, 0x24,
	0x5f, 0x3b, 0x0a, 0x43, 0xbe, 0x6b, 0x46, 0xf1, 0x19, 0x4c, 0x12, 0xf9, 0xed, 0x75, 0x51, 0x32,
	0x45, 0xc5, 0x16, 0x5a, 0x12, 0xb5, 0x5f, 0x23, 0xea, 0x02, 0x46, 0x8f, 0x29, 0xde, 0x91, 0x88,
	0xc4, 0x42, 0x12, 0x71, 0xb2, 0x82, 0xe5, 0x37, 0x05, 0xe2, 0x56, 0xca, 0xcc, 0x0b, 0xe5, 0xdf,
	0x73, 0x82, 0xaf, 0x7e, 0xfd, 0x1a, 0x53, 0xc1, 0x12, 0xc9, 0x4a, 0xd3, 0x6d, 0xa1, 0xce, 0xbf,
	0x1a, 0x7c, 0x78, 0xcd, 0xa2, 0xc8, 0x8b, 0x71, 0x55, 0xe3, 0x6b, 0x99, 0xd0, 0x51, 0x36, 0x5a,
	0x47, 0x36, 0x4b, 0x40, 0x51, 0xab, 0x37, 0xb2, 0xfa, 0x99, 0x65, 0x87, 0xa6, 0xcc, 0xca, 0xa8,
	0x65, 0xd5, 0x18, 0x9e, 0xde, 0xb3, 0xc3, 0xd3, 0xaf, 0x0f, 0x4f, 0x9d, 0x72, 0x03, 0xd9, 0xeb,
	0x8a, 0x72, 0x8d, 0x61, 0x18, 0xb6, 0x86, 0xc1, 0xf9, 0x5d, 0x87, 0x4f, 0x9e, 0xc9, 0xf8, 0x96,
	0x1c, 0x7e, 0x50, 0x63, 0x7d, 0x52, 0xee, 0x2b, 0x98, 0xee, 0xb3, 0xa6, 0xb1, 0x94, 0xdf, 0x34,
	0xb9, 0x97, 0xd9, 0x76, 0xea, 0xde, 0x4b, 0xfe, 0x97, 0x60, 0x05, 0x94, 0x0b, 0x96, 0x50, 0xbf,
	0xe4, 0xd8, 0x50, 0xda, 0x1c, 0xe1, 0x4e, 0xd0, 0x51, 0x8c, 0xb5, 0xef, 0x93, 0xbd, 0x50, 0x26,
	0x3c, 0xa0, 0xfb, 0x93, 0x8a, 0x51, 0x6d, 0x33, 0xbd, 0x73, 0x9b, 0x39, 0xbf, 0xc1, 0xc5, 0xb1,
	0xa7, 0x30, 0x64, 0x87, 0x7c, 0xd2, 0x66, 0x60, 0x96, 0xfc, 0xc9, 0x97, 0x48, 0xc9, 0x9a, 0xaf,
	0xe0, 0xb5, 0x5a, 0xc8, 0x2e, 0x79, 0x20, 0x09, 0x89, 0x7d, 0x92, 0x7b, 0x99, 0x2e, 0xef, 0x9b,
	0xf8, 0x46, 0x90, 0xc8, 0x6d, 0x1b, 0x3b, 0xbf, 0xc0, 0xdb, 0x0e, 0xbb, 0x93, 0xf2, 0xfa, 0xb2,
	0x7e, 0xb3, 0xd4, 0xd6, 0x95, 0x0d, 0xee, 0x5c, 0xc7, 0x47, 0xa6, 0xce, 0x9f, 0x1a, 0xa0, 0x3c,
	0xed, 0x1f, 0x13, 0x5a, 0xae, 0x95, 0x19, 0x98, 0x6a, 0xdc, 0xab, 0x64, 0x0b, 0xb9, 0x63, 0x95,
	0x1d, 0x47, 0xd5, 0x45, 0xa3, 0xc6, 0x72, 0xe8, 0xbd, 0xb0, 0x1c, 0x9c, 0x3f, 0x34, 0xf8, 0xe0,
	0xa8, 0x17, 0xd2, 0xf2, 0xa4, 0x92, 0xd4, 0x83, 0x57, 0x7b, 0xb6, 0x0a, 0x7e, 0x55, 0x0f, 0xc2,
	0x90, 0x41, 0x4c, 0x97, 0x2d, 0x27, 0xed, 0x70, 0xfe, 0xea, 0xda, 0x41, 0x6b, 0xce, 0xf3, 0x23,
	0x71, 0x4a, 0x3c, 0xe5, 0x79, 0xd3, 0x5f, 0x3a, 0x6f, 0xc6, 0x4b, 0xe7, 0xad, 0xd7, 0x3c, 0x6f,
	0x97, 0x3f, 0x83, 0xd5, 0xee, 0x2f, 0x32, 0xa1, 0x47, 0x63, 0x2a, 0xac, 0x57, 0x68, 0x08, 0x46,
	0x4c, 0x0e, 0x96, 0x86, 0x26, 0x00, 0x49, 0x79, 0x3e, 0x2c, 0x1d, 0x8d, 0xb3, 0xe2, 0x3c, 0x12,
	0x5f, 0x10, 0x6c, 0x19, 0xe8, 0x1c, 0x46, 0xfb, 0x74, 0x1b, 0x52, 0x1e, 0x10, 0x6c, 0xf5, 0x32,
	0xa5, 0x27, 0xf3, 0x22, 0xd8, 0xea, 0x5f, 0x5e, 0xc3, 0xdb, 0x8e, 0xca, 0xa0, 0x77, 0xf0, 0xa6,
	0xac, 0x8d, 0x5b, 0x7c, 0xea, 0x55, 0x03, 0x56, 0x13, 0x4a, 0xb0, 0xa5, 0x5d, 0x0d, 0x7f, 0xea,
	0x47, 0x0c, 0x93, 0x70, 0x3b, 0x90, 0xff, 0xa0, 0xbe, 0xf8, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x3f,
	0x44, 0x5e, 0xb7, 0x63, 0x09, 0x00, 0x00,
}
