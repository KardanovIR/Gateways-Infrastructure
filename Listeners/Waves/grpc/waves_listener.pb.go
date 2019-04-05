// Code generated by protoc-gen-go. DO NOT EDIT.
// source: waves_listener.proto

package wavesListener

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

// The request message containing task.
type AddTaskRequest struct {
	ListenTo             *ListenObject `protobuf:"bytes,1,opt,name=listenTo,proto3" json:"listenTo,omitempty"`
	CallbackType         string        `protobuf:"bytes,2,opt,name=callbackType,proto3" json:"callbackType,omitempty"`
	CallbackUrl          string        `protobuf:"bytes,3,opt,name=callbackUrl,proto3" json:"callbackUrl,omitempty"`
	TaskType             string        `protobuf:"bytes,4,opt,name=taskType,proto3" json:"taskType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *AddTaskRequest) Reset()         { *m = AddTaskRequest{} }
func (m *AddTaskRequest) String() string { return proto.CompactTextString(m) }
func (*AddTaskRequest) ProtoMessage()    {}
func (*AddTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f600a1e4d89480e4, []int{0}
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

func (m *AddTaskRequest) GetCallbackUrl() string {
	if m != nil {
		return m.CallbackUrl
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
	return fileDescriptor_f600a1e4d89480e4, []int{1}
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
	return fileDescriptor_f600a1e4d89480e4, []int{2}
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

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_f600a1e4d89480e4, []int{3}
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
	return fileDescriptor_f600a1e4d89480e4, []int{4}
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

func init() {
	proto.RegisterType((*AddTaskRequest)(nil), "wavesListener.AddTaskRequest")
	proto.RegisterType((*AddTaskResponse)(nil), "wavesListener.AddTaskResponse")
	proto.RegisterType((*RemoveTaskRequest)(nil), "wavesListener.RemoveTaskRequest")
	proto.RegisterType((*Empty)(nil), "wavesListener.Empty")
	proto.RegisterType((*ListenObject)(nil), "wavesListener.ListenObject")
}

func init() { proto.RegisterFile("waves_listener.proto", fileDescriptor_f600a1e4d89480e4) }

var fileDescriptor_f600a1e4d89480e4 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0xc6, 0xbb, 0xef, 0xdb, 0xbf, 0xd3, 0xaa, 0x38, 0x94, 0x52, 0x2a, 0x4a, 0xd8, 0x53, 0x45,
	0xd8, 0x43, 0x3d, 0xe8, 0xd5, 0x82, 0x82, 0x22, 0x58, 0x42, 0xc5, 0xa3, 0x6c, 0x93, 0x45, 0x6a,
	0xd3, 0x6e, 0xcc, 0x6e, 0x23, 0xfd, 0x32, 0x7e, 0x01, 0xbf, 0xa4, 0x64, 0x37, 0x89, 0x49, 0xa4,
	0xb7, 0x9d, 0x67, 0x9e, 0x59, 0x9e, 0xdf, 0x30, 0xd0, 0xff, 0xe4, 0xb1, 0x50, 0xaf, 0xc1, 0x52,
	0x69, 0xb1, 0x11, 0x11, 0x0b, 0x23, 0xa9, 0x25, 0x1e, 0x18, 0xf5, 0x31, 0x15, 0xe9, 0x37, 0x81,
	0xc3, 0x1b, 0xdf, 0x9f, 0x73, 0xb5, 0x72, 0xc5, 0xc7, 0x56, 0x28, 0x8d, 0x57, 0xd0, 0xb6, 0x33,
	0x73, 0x39, 0x24, 0x0e, 0x19, 0x77, 0x27, 0x27, 0xac, 0x34, 0xc4, 0xec, 0xe3, 0x69, 0xf1, 0x2e,
	0x3c, 0xed, 0xe6, 0x66, 0xa4, 0xd0, 0xf3, 0x78, 0x10, 0x2c, 0xb8, 0xb7, 0x9a, 0xef, 0x42, 0x31,
	0xfc, 0xe7, 0x90, 0x71, 0xc7, 0x2d, 0x69, 0xe8, 0x40, 0x37, 0xab, 0x9f, 0xa3, 0x60, 0xf8, 0xdf,
	0x58, 0x8a, 0x12, 0x8e, 0xa0, 0xad, 0xb9, 0xb2, 0x3f, 0xd4, 0x4d, 0x3b, 0xaf, 0xe9, 0x39, 0x1c,
	0xe5, 0x61, 0x55, 0x28, 0x37, 0x4a, 0xe0, 0x00, 0x9a, 0x49, 0xfb, 0xde, 0x37, 0x59, 0x3b, 0x6e,
	0x5a, 0xd1, 0x0b, 0x38, 0x76, 0xc5, 0x5a, 0xc6, 0xa2, 0x88, 0xb6, 0xcf, 0xdc, 0x82, 0xc6, 0xed,
	0x3a, 0xd4, 0x3b, 0x7a, 0x0d, 0xbd, 0x22, 0x1c, 0x22, 0xd4, 0x75, 0x12, 0xc4, 0xda, 0xcd, 0x1b,
	0xfb, 0xd0, 0x88, 0x79, 0xb0, 0xcd, 0xf8, 0x6c, 0x31, 0xf9, 0x22, 0xd0, 0xce, 0x16, 0x84, 0x0f,
	0xd0, 0x4a, 0x73, 0xe2, 0x69, 0x65, 0x77, 0xe5, 0x65, 0x8f, 0xce, 0xf6, 0xb5, 0x2d, 0x1e, 0xad,
	0xe1, 0x1d, 0xc0, 0x2f, 0x08, 0x3a, 0x15, 0xff, 0x1f, 0xc6, 0x51, 0xbf, 0xe2, 0xb0, 0x60, 0xb5,
	0x29, 0x83, 0xc1, 0x52, 0xb2, 0xb7, 0x28, 0xf4, 0x58, 0x7e, 0x12, 0xc6, 0x39, 0xc5, 0x97, 0xe2,
	0xc0, 0x2c, 0x39, 0x93, 0x19, 0x59, 0x34, 0xcd, 0xbd, 0x5c, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x7a, 0xee, 0x10, 0x63, 0x47, 0x02, 0x00, 0x00,
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
	// Sends a greeting
	AddTask(ctx context.Context, in *AddTaskRequest, opts ...grpc.CallOption) (*AddTaskResponse, error)
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
	err := c.cc.Invoke(ctx, "/wavesListener.Listener/AddTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listenerClient) RemoveTask(ctx context.Context, in *RemoveTaskRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/wavesListener.Listener/RemoveTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListenerServer is the server API for Listener service.
type ListenerServer interface {
	// Sends a greeting
	AddTask(context.Context, *AddTaskRequest) (*AddTaskResponse, error)
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
		FullMethod: "/wavesListener.Listener/AddTask",
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
		FullMethod: "/wavesListener.Listener/RemoveTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListenerServer).RemoveTask(ctx, req.(*RemoveTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Listener_serviceDesc = grpc.ServiceDesc{
	ServiceName: "wavesListener.Listener",
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
	Metadata: "waves_listener.proto",
}
