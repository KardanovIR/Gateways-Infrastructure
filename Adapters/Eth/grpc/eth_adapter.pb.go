// Code generated by protoc-gen-go. DO NOT EDIT.
// source: eth_adapter.proto

package ethAdapter

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type GasPriceRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GasPriceRequest) Reset()         { *m = GasPriceRequest{} }
func (m *GasPriceRequest) String() string { return proto.CompactTextString(m) }
func (*GasPriceRequest) ProtoMessage()    {}
func (*GasPriceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc5170b1544d27bc, []int{0}
}

func (m *GasPriceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GasPriceRequest.Unmarshal(m, b)
}
func (m *GasPriceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GasPriceRequest.Marshal(b, m, deterministic)
}
func (m *GasPriceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPriceRequest.Merge(m, src)
}
func (m *GasPriceRequest) XXX_Size() int {
	return xxx_messageInfo_GasPriceRequest.Size(m)
}
func (m *GasPriceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPriceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GasPriceRequest proto.InternalMessageInfo

// The response message containing the greetings
type GasPriceReply struct {
	GasPrice             string   `protobuf:"bytes,1,opt,name=gasPrice,proto3" json:"gasPrice,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GasPriceReply) Reset()         { *m = GasPriceReply{} }
func (m *GasPriceReply) String() string { return proto.CompactTextString(m) }
func (*GasPriceReply) ProtoMessage()    {}
func (*GasPriceReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc5170b1544d27bc, []int{1}
}

func (m *GasPriceReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GasPriceReply.Unmarshal(m, b)
}
func (m *GasPriceReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GasPriceReply.Marshal(b, m, deterministic)
}
func (m *GasPriceReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPriceReply.Merge(m, src)
}
func (m *GasPriceReply) XXX_Size() int {
	return xxx_messageInfo_GasPriceReply.Size(m)
}
func (m *GasPriceReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPriceReply.DiscardUnknown(m)
}

var xxx_messageInfo_GasPriceReply proto.InternalMessageInfo

func (m *GasPriceReply) GetGasPrice() string {
	if m != nil {
		return m.GasPrice
	}
	return ""
}

func init() {
	proto.RegisterType((*GasPriceRequest)(nil), "ethAdapter.GasPriceRequest")
	proto.RegisterType((*GasPriceReply)(nil), "ethAdapter.GasPriceReply")
}

func init() { proto.RegisterFile("eth_adapter.proto", fileDescriptor_cc5170b1544d27bc) }

var fileDescriptor_cc5170b1544d27bc = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0x2d, 0xc9, 0x70,
	0x4c, 0x49, 0x2c, 0x28, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0x88,
	0x28, 0x09, 0x72, 0xf1, 0xbb, 0x27, 0x16, 0x07, 0x14, 0x65, 0x26, 0xa7, 0x06, 0xa5, 0x16, 0x96,
	0xa6, 0x16, 0x97, 0x28, 0x69, 0x73, 0xf1, 0x22, 0x84, 0x0a, 0x72, 0x2a, 0x85, 0xa4, 0xb8, 0x38,
	0xd2, 0xa1, 0x02, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x70, 0xbe, 0x91, 0x1f, 0x17, 0x9b,
	0x73, 0x7e, 0x6e, 0x6e, 0x7e, 0x9e, 0x90, 0x0b, 0x17, 0x07, 0x4c, 0x9b, 0x90, 0xb4, 0x1e, 0x92,
	0xa5, 0x68, 0xe6, 0x4b, 0x49, 0x62, 0x97, 0x2c, 0xc8, 0xa9, 0x54, 0x62, 0x70, 0xd2, 0xe0, 0x12,
	0xce, 0xcc, 0xd7, 0x4b, 0x2f, 0x2a, 0x48, 0xd6, 0x4b, 0x84, 0x2a, 0x49, 0x2d, 0xc9, 0x70, 0xe2,
	0x77, 0x85, 0x6b, 0x09, 0x00, 0xf9, 0x21, 0x80, 0x31, 0x89, 0x0d, 0xec, 0x19, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x65, 0xcd, 0x07, 0xa3, 0xe0, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CommonClient is the client API for Common service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommonClient interface {
	// Sends a greeting
	GasPrice(ctx context.Context, in *GasPriceRequest, opts ...grpc.CallOption) (*GasPriceReply, error)
}

type commonClient struct {
	cc *grpc.ClientConn
}

func NewCommonClient(cc *grpc.ClientConn) CommonClient {
	return &commonClient{cc}
}

func (c *commonClient) GasPrice(ctx context.Context, in *GasPriceRequest, opts ...grpc.CallOption) (*GasPriceReply, error) {
	out := new(GasPriceReply)
	err := c.cc.Invoke(ctx, "/ethAdapter.Common/GasPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommonServer is the server API for Common service.
type CommonServer interface {
	// Sends a greeting
	GasPrice(context.Context, *GasPriceRequest) (*GasPriceReply, error)
}

func RegisterCommonServer(s *grpc.Server, srv CommonServer) {
	s.RegisterService(&_Common_serviceDesc, srv)
}

func _Common_GasPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GasPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommonServer).GasPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethAdapter.Common/GasPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommonServer).GasPrice(ctx, req.(*GasPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Common_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethAdapter.Common",
	HandlerType: (*CommonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GasPrice",
			Handler:    _Common_GasPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "eth_adapter.proto",
}