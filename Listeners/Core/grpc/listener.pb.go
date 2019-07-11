// Code generated by protoc-gen-go. DO NOT EDIT.
// source: listener.proto

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

// The request message containing task
type AddTaskRequest struct {
	ListenTo             *ListenObject `protobuf:"bytes,1,opt,name=listenTo,proto3" json:"listenTo,omitempty"`
	CallbackType         string        `protobuf:"bytes,2,opt,name=callbackType,proto3" json:"callbackType,omitempty"`
	TaskType             string        `protobuf:"bytes,3,opt,name=taskType,proto3" json:"taskType,omitempty"`
	ProcessId            string        `protobuf:"bytes,4,opt,name=processId,proto3" json:"processId,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AddTaskRequest) Reset()         { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()    {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f75aade3a9f7de9c, []int{0}
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

func (m *AddTaskRequest) GetProcessId() string {
	if m != nil {
		return m.ProcessId
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
	return fileDescriptor_f75aade3a9f7de9c, []int{1}
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
	return fileDescriptor_f75aade3a9f7de9c, []int{2}
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
	return fileDescriptor_f75aade3a9f7de9c, []int{3}
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
	return fileDescriptor_f75aade3a9f7de9c, []int{4}
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
	proto.RegisterType((*AddTaskRequest)(nil), "blockchain.AddTaskRequest")
	proto.RegisterType((*AddTaskResponse)(nil), "blockchain.AddTaskResponse")
	proto.RegisterType((*RemoveTaskRequest)(nil), "blockchain.RemoveTaskRequest")
	proto.RegisterType((*ListenObject)(nil), "blockchain.ListenObject")
	proto.RegisterType((*Empty)(nil), "blockchain.Empty")
}

func init() { proto.RegisterFile("listener.proto", fileDescriptor_f75aade3a9f7de9c) }

var fileDescriptor_f75aade3a9f7de9c = []byte{
	// 310 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x4a, 0x03, 0x31,
	0x10, 0xc6, 0x1b, 0xed, 0xdf, 0xb1, 0x56, 0x3b, 0x88, 0x2c, 0xab, 0x42, 0x89, 0x97, 0x8a, 0xb0,
	0x87, 0xea, 0xc1, 0xa3, 0x16, 0x3d, 0x14, 0x0a, 0x96, 0xa5, 0x2f, 0x90, 0xcd, 0x06, 0xad, 0xdd,
	0x36, 0x71, 0x93, 0x16, 0xfa, 0x1a, 0x3e, 0x83, 0x0f, 0x2a, 0x9b, 0xfd, 0x5b, 0xd4, 0x5b, 0x66,
	0xbe, 0x6f, 0xe0, 0xfb, 0x4d, 0x06, 0x7a, 0xd1, 0x42, 0x1b, 0xb1, 0x16, 0xb1, 0xa7, 0x62, 0x69,
	0x24, 0x42, 0x10, 0x49, 0xbe, 0xe4, 0xef, 0x6c, 0xb1, 0xa6, 0xdf, 0x04, 0x7a, 0x4f, 0x61, 0x38,
	0x67, 0x7a, 0xe9, 0x8b, 0xcf, 0x8d, 0xd0, 0x06, 0xef, 0xa1, 0x9d, 0x0e, 0xcc, 0xa5, 0x43, 0x06,
	0x64, 0x78, 0x34, 0x72, 0xbc, 0x72, 0xc2, 0x9b, 0x5a, 0xed, 0x35, 0xf8, 0x10, 0xdc, 0xf8, 0x85,
	0x13, 0x29, 0x74, 0x39, 0x8b, 0xa2, 0x80, 0xf1, 0xe5, 0x7c, 0xa7, 0x84, 0x73, 0x30, 0x20, 0xc3,
	0x8e, 0xbf, 0xd7, 0x43, 0x17, 0xda, 0x86, 0xe9, 0x54, 0x3f, 0xb4, 0x7a, 0x51, 0xe3, 0x25, 0x74,
	0x54, 0x2c, 0xb9, 0xd0, 0x7a, 0x12, 0x3a, 0x75, 0x2b, 0x96, 0x0d, 0x7a, 0x03, 0x27, 0x45, 0x4a,
	0xad, 0xe4, 0x5a, 0x0b, 0x3c, 0x87, 0x66, 0x32, 0x3c, 0x09, 0x6d, 0xc8, 0x8e, 0x9f, 0x55, 0xf4,
	0x16, 0xfa, 0xbe, 0x58, 0xc9, 0xad, 0xa8, 0x32, 0xfd, 0x67, 0x7e, 0x80, 0x6e, 0x95, 0x07, 0x11,
	0xea, 0x26, 0x49, 0x97, 0xba, 0xec, 0x1b, 0xcf, 0xa0, 0xb1, 0x65, 0xd1, 0x26, 0x47, 0x4a, 0x0b,
	0xda, 0x82, 0xc6, 0xcb, 0x4a, 0x99, 0xdd, 0xe8, 0x8b, 0x40, 0x7b, 0x9a, 0x2d, 0x18, 0x9f, 0xa1,
	0x95, 0xe5, 0x44, 0xb7, 0xba, 0xb4, 0xfd, 0x15, 0xbb, 0x17, 0x7f, 0x6a, 0x29, 0x18, 0xad, 0xe1,
	0x23, 0x40, 0x89, 0x80, 0x57, 0x55, 0xf3, 0x2f, 0x34, 0xb7, 0x5f, 0x95, 0x6d, 0x24, 0x5a, 0x1b,
	0x5f, 0xc3, 0xe9, 0x42, 0x7a, 0x6f, 0xb1, 0xe2, 0x5e, 0xfe, 0xf9, 0xe3, 0xe3, 0x3c, 0xe5, 0x2c,
	0xb9, 0x82, 0x19, 0x09, 0x9a, 0xf6, 0x1c, 0xee, 0x7e, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb2, 0xac,
	0x8c, 0x6a, 0x20, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

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
	Metadata: "listener.proto",
}