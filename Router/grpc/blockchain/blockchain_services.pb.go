// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blockchain_services.proto

package blockchain

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

// The request message for raw transaction creation
type RawTransactionRequest struct {
	AddressFrom          string   `protobuf:"bytes,1,opt,name=addressFrom,proto3" json:"addressFrom,omitempty"`
	SendersPublicKey     string   `protobuf:"bytes,2,opt,name=sendersPublicKey,proto3" json:"sendersPublicKey,omitempty"`
	AddressTo            string   `protobuf:"bytes,3,opt,name=addressTo,proto3" json:"addressTo,omitempty"`
	Amount               string   `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	AssetId              string   `protobuf:"bytes,5,opt,name=assetId,proto3" json:"assetId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawTransactionRequest) Reset()         { *m = RawTransactionRequest{} }
func (m *RawTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*RawTransactionRequest) ProtoMessage()    {}
func (*RawTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{0}
}

func (m *RawTransactionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawTransactionRequest.Unmarshal(m, b)
}
func (m *RawTransactionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawTransactionRequest.Marshal(b, m, deterministic)
}
func (m *RawTransactionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawTransactionRequest.Merge(m, src)
}
func (m *RawTransactionRequest) XXX_Size() int {
	return xxx_messageInfo_RawTransactionRequest.Size(m)
}
func (m *RawTransactionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RawTransactionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RawTransactionRequest proto.InternalMessageInfo

func (m *RawTransactionRequest) GetAddressFrom() string {
	if m != nil {
		return m.AddressFrom
	}
	return ""
}

func (m *RawTransactionRequest) GetSendersPublicKey() string {
	if m != nil {
		return m.SendersPublicKey
	}
	return ""
}

func (m *RawTransactionRequest) GetAddressTo() string {
	if m != nil {
		return m.AddressTo
	}
	return ""
}

func (m *RawTransactionRequest) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *RawTransactionRequest) GetAssetId() string {
	if m != nil {
		return m.AssetId
	}
	return ""
}

// The response message containing raw transaction
type RawTransactionReply struct {
	Tx                   []byte   `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawTransactionReply) Reset()         { *m = RawTransactionReply{} }
func (m *RawTransactionReply) String() string { return proto.CompactTextString(m) }
func (*RawTransactionReply) ProtoMessage()    {}
func (*RawTransactionReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{1}
}

func (m *RawTransactionReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawTransactionReply.Unmarshal(m, b)
}
func (m *RawTransactionReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawTransactionReply.Marshal(b, m, deterministic)
}
func (m *RawTransactionReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawTransactionReply.Merge(m, src)
}
func (m *RawTransactionReply) XXX_Size() int {
	return xxx_messageInfo_RawTransactionReply.Size(m)
}
func (m *RawTransactionReply) XXX_DiscardUnknown() {
	xxx_messageInfo_RawTransactionReply.DiscardUnknown(m)
}

var xxx_messageInfo_RawTransactionReply proto.InternalMessageInfo

func (m *RawTransactionReply) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

// The request message for singing transaction
type SendTransactionRequest struct {
	Tx                   []byte   `protobuf:"bytes,1,opt,name=tx,proto3" json:"tx,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTransactionRequest) Reset()         { *m = SendTransactionRequest{} }
func (m *SendTransactionRequest) String() string { return proto.CompactTextString(m) }
func (*SendTransactionRequest) ProtoMessage()    {}
func (*SendTransactionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{2}
}

func (m *SendTransactionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTransactionRequest.Unmarshal(m, b)
}
func (m *SendTransactionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTransactionRequest.Marshal(b, m, deterministic)
}
func (m *SendTransactionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTransactionRequest.Merge(m, src)
}
func (m *SendTransactionRequest) XXX_Size() int {
	return xxx_messageInfo_SendTransactionRequest.Size(m)
}
func (m *SendTransactionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTransactionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendTransactionRequest proto.InternalMessageInfo

func (m *SendTransactionRequest) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

// The response message containing transaction's id
type SendTransactionReply struct {
	TxId                 string   `protobuf:"bytes,1,opt,name=txId,proto3" json:"txId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendTransactionReply) Reset()         { *m = SendTransactionReply{} }
func (m *SendTransactionReply) String() string { return proto.CompactTextString(m) }
func (*SendTransactionReply) ProtoMessage()    {}
func (*SendTransactionReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{3}
}

func (m *SendTransactionReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendTransactionReply.Unmarshal(m, b)
}
func (m *SendTransactionReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendTransactionReply.Marshal(b, m, deterministic)
}
func (m *SendTransactionReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendTransactionReply.Merge(m, src)
}
func (m *SendTransactionReply) XXX_Size() int {
	return xxx_messageInfo_SendTransactionReply.Size(m)
}
func (m *SendTransactionReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SendTransactionReply.DiscardUnknown(m)
}

var xxx_messageInfo_SendTransactionReply proto.InternalMessageInfo

func (m *SendTransactionReply) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

// The request message for getting transaction's status containing tx id
type GetTransactionStatusRequest struct {
	TxId                 string   `protobuf:"bytes,1,opt,name=txId,proto3" json:"txId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTransactionStatusRequest) Reset()         { *m = GetTransactionStatusRequest{} }
func (m *GetTransactionStatusRequest) String() string { return proto.CompactTextString(m) }
func (*GetTransactionStatusRequest) ProtoMessage()    {}
func (*GetTransactionStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{4}
}

func (m *GetTransactionStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTransactionStatusRequest.Unmarshal(m, b)
}
func (m *GetTransactionStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTransactionStatusRequest.Marshal(b, m, deterministic)
}
func (m *GetTransactionStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTransactionStatusRequest.Merge(m, src)
}
func (m *GetTransactionStatusRequest) XXX_Size() int {
	return xxx_messageInfo_GetTransactionStatusRequest.Size(m)
}
func (m *GetTransactionStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTransactionStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTransactionStatusRequest proto.InternalMessageInfo

func (m *GetTransactionStatusRequest) GetTxId() string {
	if m != nil {
		return m.TxId
	}
	return ""
}

// The response message containing transaction's status: UNKNOWN, PENDING, SUCCESS
type GetTransactionStatusReply struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTransactionStatusReply) Reset()         { *m = GetTransactionStatusReply{} }
func (m *GetTransactionStatusReply) String() string { return proto.CompactTextString(m) }
func (*GetTransactionStatusReply) ProtoMessage()    {}
func (*GetTransactionStatusReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{5}
}

func (m *GetTransactionStatusReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTransactionStatusReply.Unmarshal(m, b)
}
func (m *GetTransactionStatusReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTransactionStatusReply.Marshal(b, m, deterministic)
}
func (m *GetTransactionStatusReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTransactionStatusReply.Merge(m, src)
}
func (m *GetTransactionStatusReply) XXX_Size() int {
	return xxx_messageInfo_GetTransactionStatusReply.Size(m)
}
func (m *GetTransactionStatusReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTransactionStatusReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetTransactionStatusReply proto.InternalMessageInfo

func (m *GetTransactionStatusReply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

// The request message containing task
type AddTaskRequest struct {
	ListenTo             *ListenObject `protobuf:"bytes,1,opt,name=listenTo,proto3" json:"listenTo,omitempty"`
	CallbackType         string        `protobuf:"bytes,2,opt,name=callbackType,proto3" json:"callbackType,omitempty"`
	TaskType             string        `protobuf:"bytes,3,opt,name=taskType,proto3" json:"taskType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AddTaskRequest) Reset()         { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()    {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{6}
}

func (m *AddTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskRequest.Unmarshal(m, b)
}
func (m *AddTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskRequest.Marshal(b, m, deterministic)
}
func (m *AddTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskRequest.Merge(m, src)
}
func (m *AddTaskRequest) XXX_Size() int {
	return xxx_messageInfo_AddTaskRequest.Size(m)
}
func (m *AddTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskRequest proto.InternalMessageInfo

func (m *AddTaskRequest) GetListenTo() *ListenObject {
	if m != nil {
		return m.ListenTo
	}
	return nil
}

func (m *AddTaskRequest) GetCallbackType() string {
	if m != nil {
		return m.CallbackType
	}
	return ""
}

func (m *AddTaskRequest) GetTaskType() string {
	if m != nil {
		return m.TaskType
	}
	return ""
}

type AddTaskResponse struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTaskResponse) Reset()         { *m = AddTaskResponse{} }
func (m *AddTaskResponse) String() string { return proto.CompactTextString(m) }
func (*AddTaskResponse) ProtoMessage()    {}
func (*AddTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{7}
}

func (m *AddTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTaskResponse.Unmarshal(m, b)
}
func (m *AddTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTaskResponse.Marshal(b, m, deterministic)
}
func (m *AddTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTaskResponse.Merge(m, src)
}
func (m *AddTaskResponse) XXX_Size() int {
	return xxx_messageInfo_AddTaskResponse.Size(m)
}
func (m *AddTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddTaskResponse proto.InternalMessageInfo

func (m *AddTaskResponse) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

type RemoveTaskRequest struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveTaskRequest) Reset()         { *m = RemoveTaskRequest{} }
func (m *RemoveTaskRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveTaskRequest) ProtoMessage()    {}
func (*RemoveTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{8}
}

func (m *RemoveTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveTaskRequest.Unmarshal(m, b)
}
func (m *RemoveTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveTaskRequest.Marshal(b, m, deterministic)
}
func (m *RemoveTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveTaskRequest.Merge(m, src)
}
func (m *RemoveTaskRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveTaskRequest.Size(m)
}
func (m *RemoveTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveTaskRequest proto.InternalMessageInfo

func (m *RemoveTaskRequest) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

// ListenObject defines what will service listen to
type ListenObject struct {
	// type can be one of the follow values: TxId or Address
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListenObject) Reset()         { *m = ListenObject{} }
func (m *ListenObject) String() string { return proto.CompactTextString(m) }
func (*ListenObject) ProtoMessage()    {}
func (*ListenObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{9}
}

func (m *ListenObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListenObject.Unmarshal(m, b)
}
func (m *ListenObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListenObject.Marshal(b, m, deterministic)
}
func (m *ListenObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListenObject.Merge(m, src)
}
func (m *ListenObject) XXX_Size() int {
	return xxx_messageInfo_ListenObject.Size(m)
}
func (m *ListenObject) XXX_DiscardUnknown() {
	xxx_messageInfo_ListenObject.DiscardUnknown(m)
}

var xxx_messageInfo_ListenObject proto.InternalMessageInfo

func (m *ListenObject) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ListenObject) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_1bf9757f8350b04c, []int{10}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RawTransactionRequest)(nil), "blockchain.RawTransactionRequest")
	proto.RegisterType((*RawTransactionReply)(nil), "blockchain.RawTransactionReply")
	proto.RegisterType((*SendTransactionRequest)(nil), "blockchain.SendTransactionRequest")
	proto.RegisterType((*SendTransactionReply)(nil), "blockchain.SendTransactionReply")
	proto.RegisterType((*GetTransactionStatusRequest)(nil), "blockchain.GetTransactionStatusRequest")
	proto.RegisterType((*GetTransactionStatusReply)(nil), "blockchain.GetTransactionStatusReply")
	proto.RegisterType((*AddTaskRequest)(nil), "blockchain.AddTaskRequest")
	proto.RegisterType((*AddTaskResponse)(nil), "blockchain.AddTaskResponse")
	proto.RegisterType((*RemoveTaskRequest)(nil), "blockchain.RemoveTaskRequest")
	proto.RegisterType((*ListenObject)(nil), "blockchain.ListenObject")
	proto.RegisterType((*Empty)(nil), "blockchain.Empty")
}

func init() { proto.RegisterFile("blockchain_services.proto", fileDescriptor_1bf9757f8350b04c) }

var fileDescriptor_1bf9757f8350b04c = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x5d, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0x43, 0xf3, 0xd3, 0x69, 0xd4, 0x90, 0x21, 0x44, 0xae, 0x0b, 0x22, 0xac, 0x54, 0xb5,
	0x14, 0x29, 0x12, 0x29, 0x0f, 0x3c, 0xd2, 0x08, 0xa8, 0x2a, 0x90, 0x88, 0xdc, 0xbc, 0xf4, 0x09,
	0x6d, 0xec, 0x11, 0x35, 0x71, 0xbc, 0xc6, 0xbb, 0x09, 0xc9, 0x05, 0x38, 0x00, 0xb7, 0xe0, 0x00,
	0xdc, 0x0f, 0x79, 0xed, 0x38, 0x4e, 0x6a, 0xc2, 0xdb, 0xce, 0x7c, 0xdf, 0xcc, 0x7e, 0x3b, 0xf3,
	0x69, 0xe1, 0x68, 0xec, 0x0b, 0x67, 0xe2, 0xdc, 0x71, 0x2f, 0xf8, 0x22, 0x29, 0x9a, 0x7b, 0x0e,
	0xc9, 0x5e, 0x18, 0x09, 0x25, 0x10, 0xd6, 0x10, 0xfb, 0x63, 0xc0, 0x63, 0x9b, 0xff, 0x18, 0x45,
	0x3c, 0x90, 0xdc, 0x51, 0x9e, 0x08, 0x6c, 0xfa, 0x3e, 0x23, 0xa9, 0xb0, 0x0b, 0x07, 0xdc, 0x75,
	0x23, 0x92, 0xf2, 0x43, 0x24, 0xa6, 0xa6, 0xd1, 0x35, 0xce, 0xf6, 0xed, 0x7c, 0x0a, 0xcf, 0xe1,
	0xa1, 0xa4, 0xc0, 0xa5, 0x48, 0x0e, 0x67, 0x63, 0xdf, 0x73, 0x3e, 0xd2, 0xd2, 0x2c, 0x6b, 0xda,
	0xbd, 0x3c, 0x3e, 0x81, 0xfd, 0xb4, 0x74, 0x24, 0xcc, 0x07, 0x9a, 0xb4, 0x4e, 0x60, 0x07, 0xaa,
	0x7c, 0x2a, 0x66, 0x81, 0x32, 0xf7, 0x34, 0x94, 0x46, 0x68, 0x42, 0x8d, 0x4b, 0x49, 0xea, 0xda,
	0x35, 0x2b, 0x1a, 0x58, 0x85, 0xec, 0x04, 0x1e, 0x6d, 0xcb, 0x0e, 0xfd, 0x25, 0x1e, 0x42, 0x59,
	0x2d, 0xb4, 0xd6, 0x86, 0x5d, 0x56, 0x0b, 0x76, 0x06, 0x9d, 0x1b, 0x0a, 0xdc, 0x82, 0xe7, 0x6d,
	0x33, 0xcf, 0xa1, 0x7d, 0x8f, 0x19, 0x77, 0x44, 0xd8, 0x53, 0x8b, 0x6b, 0x37, 0x7d, 0xbf, 0x3e,
	0xb3, 0x57, 0x70, 0x7c, 0x45, 0x2a, 0x47, 0xbd, 0x51, 0x5c, 0xcd, 0xe4, 0xaa, 0x75, 0x51, 0xc9,
	0x05, 0x1c, 0x15, 0x97, 0xc4, 0x77, 0x74, 0xa0, 0x2a, 0x75, 0x98, 0x96, 0xa4, 0x11, 0xfb, 0x69,
	0xc0, 0xe1, 0xa5, 0xeb, 0x8e, 0xb8, 0x9c, 0xac, 0x7a, 0xbf, 0x86, 0xba, 0xef, 0x49, 0x45, 0xc1,
	0x48, 0x68, 0xf2, 0x41, 0xdf, 0xec, 0xad, 0xd7, 0xd9, 0xfb, 0xa4, 0xb1, 0xcf, 0xe3, 0x6f, 0xe4,
	0x28, 0x3b, 0x63, 0x22, 0x83, 0x86, 0xc3, 0x7d, 0x7f, 0xcc, 0x9d, 0xc9, 0x68, 0x19, 0x52, 0xba,
	0xa5, 0x8d, 0x1c, 0x5a, 0x50, 0x57, 0x5c, 0x26, 0x78, 0xb2, 0xa0, 0x2c, 0x66, 0x2f, 0xa0, 0x99,
	0xe9, 0x90, 0xa1, 0x08, 0x24, 0xc5, 0x9a, 0x63, 0x38, 0x7b, 0x66, 0x1a, 0xb1, 0x97, 0xd0, 0xb2,
	0x69, 0x2a, 0xe6, 0x94, 0x57, 0xfd, 0x2f, 0xf2, 0x1b, 0x68, 0xe4, 0x15, 0xeb, 0xc9, 0xc5, 0xf7,
	0xaf, 0x26, 0x17, 0xeb, 0x6a, 0x43, 0x65, 0xce, 0xfd, 0xd9, 0x4a, 0x74, 0x12, 0xb0, 0x1a, 0x54,
	0xde, 0x4f, 0x43, 0xb5, 0xec, 0xff, 0x2e, 0x43, 0xed, 0xd2, 0xe5, 0xa1, 0xa2, 0x08, 0x6f, 0xa1,
	0x75, 0x45, 0x6a, 0xd3, 0x17, 0xf8, 0x3c, 0x3f, 0x9f, 0x42, 0xab, 0x5b, 0xcf, 0x76, 0x51, 0x42,
	0x7f, 0xc9, 0x4a, 0x78, 0x0b, 0xcd, 0x2d, 0x7b, 0x20, 0xcb, 0x57, 0x15, 0xbb, 0xcc, 0xea, 0xee,
	0xe4, 0x24, 0xad, 0xef, 0xa0, 0x5d, 0x64, 0x0d, 0x3c, 0xcd, 0xd7, 0xee, 0xf0, 0x9b, 0x75, 0xf2,
	0x7f, 0xa2, 0xbe, 0xa9, 0xff, 0xcb, 0x80, 0x7a, 0x32, 0x6f, 0x8a, 0xf0, 0x5d, 0x3c, 0x37, 0xbd,
	0x53, 0xb4, 0xf2, 0x0d, 0x36, 0x0d, 0x67, 0x1d, 0x17, 0x62, 0x89, 0x09, 0x58, 0x09, 0xdf, 0x02,
	0xac, 0xd7, 0x8d, 0x4f, 0x37, 0x06, 0xb9, 0x6d, 0x03, 0xab, 0x95, 0x87, 0xf5, 0xfa, 0x58, 0x69,
	0x70, 0x0a, 0xe8, 0x89, 0xde, 0xd7, 0x28, 0x74, 0x72, 0xe8, 0xa0, 0x39, 0xc8, 0xce, 0xc3, 0xf8,
	0xd3, 0x1a, 0x1a, 0xe3, 0xaa, 0xfe, 0xbd, 0x2e, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0x92, 0xe0,
	0x03, 0x1e, 0xda, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AdapterClient is the client API for Adapter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdapterClient interface {
	// common request to create raw transaction: for main currency or based on this currency assets transfer
	GetRawTransaction(ctx context.Context, in *RawTransactionRequest, opts ...grpc.CallOption) (*RawTransactionReply, error)
	// Send transaction
	SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionReply, error)
	// Get transaction status
	GetTransactionStatus(ctx context.Context, in *GetTransactionStatusRequest, opts ...grpc.CallOption) (*GetTransactionStatusReply, error)
}

type adapterClient struct {
	cc *grpc.ClientConn
}

func NewAdapterClient(cc *grpc.ClientConn) AdapterClient {
	return &adapterClient{cc}
}

func (c *adapterClient) GetRawTransaction(ctx context.Context, in *RawTransactionRequest, opts ...grpc.CallOption) (*RawTransactionReply, error) {
	out := new(RawTransactionReply)
	err := c.cc.Invoke(ctx, "/blockchain.Adapter/GetRawTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adapterClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*SendTransactionReply, error) {
	out := new(SendTransactionReply)
	err := c.cc.Invoke(ctx, "/blockchain.Adapter/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adapterClient) GetTransactionStatus(ctx context.Context, in *GetTransactionStatusRequest, opts ...grpc.CallOption) (*GetTransactionStatusReply, error) {
	out := new(GetTransactionStatusReply)
	err := c.cc.Invoke(ctx, "/blockchain.Adapter/GetTransactionStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdapterServer is the server API for Adapter service.
type AdapterServer interface {
	// common request to create raw transaction: for main currency or based on this currency assets transfer
	GetRawTransaction(context.Context, *RawTransactionRequest) (*RawTransactionReply, error)
	// Send transaction
	SendTransaction(context.Context, *SendTransactionRequest) (*SendTransactionReply, error)
	// Get transaction status
	GetTransactionStatus(context.Context, *GetTransactionStatusRequest) (*GetTransactionStatusReply, error)
}

func RegisterAdapterServer(s *grpc.Server, srv AdapterServer) {
	s.RegisterService(&_Adapter_serviceDesc, srv)
}

func _Adapter_GetRawTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RawTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).GetRawTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.Adapter/GetRawTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).GetRawTransaction(ctx, req.(*RawTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Adapter_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.Adapter/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).SendTransaction(ctx, req.(*SendTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Adapter_GetTransactionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).GetTransactionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.Adapter/GetTransactionStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).GetTransactionStatus(ctx, req.(*GetTransactionStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Adapter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blockchain.Adapter",
	HandlerType: (*AdapterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRawTransaction",
			Handler:    _Adapter_GetRawTransaction_Handler,
		},
		{
			MethodName: "SendTransaction",
			Handler:    _Adapter_SendTransaction_Handler,
		},
		{
			MethodName: "GetTransactionStatus",
			Handler:    _Adapter_GetTransactionStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blockchain_services.proto",
}

// ListenerClient is the client API for Listener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ListenerClient interface {
	// Add task
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error)
	// Remove task
	RemoveTask(ctx context.Context, in *RemoveTaskRequest, opts ...grpc.CallOption) (*Empty, error)
}

type listenerClient struct {
	cc *grpc.ClientConn
}

func NewListenerClient(cc *grpc.ClientConn) ListenerClient {
	return &listenerClient{cc}
}

func (c *listenerClient) AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error) {
	out := new(AddTaskResponse)
	err := c.cc.Invoke(ctx, "/blockchain.Listener/AddTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listenerClient) RemoveTask(ctx context.Context, in *RemoveTaskRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/blockchain.Listener/RemoveTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListenerServer is the server API for Listener service.
type ListenerServer interface {
	// Add task
	AddTask(context.Context, *AddTaskRequest) (*AddTaskResponse, error)
	// Remove task
	RemoveTask(context.Context, *RemoveTaskRequest) (*Empty, error)
}

func RegisterListenerServer(s *grpc.Server, srv ListenerServer) {
	s.RegisterService(&_Listener_serviceDesc, srv)
}

func _Listener_AddTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListenerServer).AddTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.Listener/AddTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListenerServer).AddTask(ctx, req.(*AddTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Listener_RemoveTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListenerServer).RemoveTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.Listener/RemoveTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListenerServer).RemoveTask(ctx, req.(*RemoveTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Listener_serviceDesc = grpc.ServiceDesc{
	ServiceName: "blockchain.Listener",
	HandlerType: (*ListenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTask",
			Handler:    _Listener_AddTask_Handler,
		},
		{
			MethodName: "RemoveTask",
			Handler:    _Listener_RemoveTask_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blockchain_services.proto",
}
