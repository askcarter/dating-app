// Code generated by protoc-gen-go.
// source: messages.proto
// DO NOT EDIT!

/*
Package helloworld is a generated protocol buffer package.

option go_package = "datingapp/pb";

It is generated from these files:
	messages.proto

It has these top-level messages:
	HelloRequest
	HelloReply
	User
	DebugListUsersRequest
	DebugListUsersResponse
	ChatRequest
	ChatResponse
*/
package helloworld

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
type HelloRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type HelloReply struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type User struct {
	ID          int64   `protobuf:"varint,1,opt,name=iD" json:"iD,omitempty"`
	DisplayName string  `protobuf:"bytes,2,opt,name=displayName" json:"displayName,omitempty"`
	Matches     []int64 `protobuf:"varint,3,rep,packed,name=matches" json:"matches,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *User) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *User) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *User) GetMatches() []int64 {
	if m != nil {
		return m.Matches
	}
	return nil
}

type DebugListUsersRequest struct {
}

func (m *DebugListUsersRequest) Reset()                    { *m = DebugListUsersRequest{} }
func (m *DebugListUsersRequest) String() string            { return proto.CompactTextString(m) }
func (*DebugListUsersRequest) ProtoMessage()               {}
func (*DebugListUsersRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type DebugListUsersResponse struct {
	Users []*User `protobuf:"bytes,1,rep,name=users" json:"users,omitempty"`
}

func (m *DebugListUsersResponse) Reset()                    { *m = DebugListUsersResponse{} }
func (m *DebugListUsersResponse) String() string            { return proto.CompactTextString(m) }
func (*DebugListUsersResponse) ProtoMessage()               {}
func (*DebugListUsersResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *DebugListUsersResponse) GetUsers() []*User {
	if m != nil {
		return m.Users
	}
	return nil
}

type ChatRequest struct {
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
}

func (m *ChatRequest) Reset()                    { *m = ChatRequest{} }
func (m *ChatRequest) String() string            { return proto.CompactTextString(m) }
func (*ChatRequest) ProtoMessage()               {}
func (*ChatRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ChatRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type ChatResponse struct {
}

func (m *ChatResponse) Reset()                    { *m = ChatResponse{} }
func (m *ChatResponse) String() string            { return proto.CompactTextString(m) }
func (*ChatResponse) ProtoMessage()               {}
func (*ChatResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
	proto.RegisterType((*User)(nil), "helloworld.User")
	proto.RegisterType((*DebugListUsersRequest)(nil), "helloworld.DebugListUsersRequest")
	proto.RegisterType((*DebugListUsersResponse)(nil), "helloworld.DebugListUsersResponse")
	proto.RegisterType((*ChatRequest)(nil), "helloworld.ChatRequest")
	proto.RegisterType((*ChatResponse)(nil), "helloworld.ChatResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Greeter service

type GreeterClient interface {
	// Sends a greeting
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/helloworld.Greeter/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	// Sends a greeting
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

// Client API for DatingGame service

type DatingGameClient interface {
	DebugListUsers(ctx context.Context, in *DebugListUsersRequest, opts ...grpc.CallOption) (*DebugListUsersResponse, error)
	SendChat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error)
}

type datingGameClient struct {
	cc *grpc.ClientConn
}

func NewDatingGameClient(cc *grpc.ClientConn) DatingGameClient {
	return &datingGameClient{cc}
}

func (c *datingGameClient) DebugListUsers(ctx context.Context, in *DebugListUsersRequest, opts ...grpc.CallOption) (*DebugListUsersResponse, error) {
	out := new(DebugListUsersResponse)
	err := grpc.Invoke(ctx, "/helloworld.DatingGame/DebugListUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *datingGameClient) SendChat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error) {
	out := new(ChatResponse)
	err := grpc.Invoke(ctx, "/helloworld.DatingGame/SendChat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DatingGame service

type DatingGameServer interface {
	DebugListUsers(context.Context, *DebugListUsersRequest) (*DebugListUsersResponse, error)
	SendChat(context.Context, *ChatRequest) (*ChatResponse, error)
}

func RegisterDatingGameServer(s *grpc.Server, srv DatingGameServer) {
	s.RegisterService(&_DatingGame_serviceDesc, srv)
}

func _DatingGame_DebugListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DebugListUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatingGameServer).DebugListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.DatingGame/DebugListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatingGameServer).DebugListUsers(ctx, req.(*DebugListUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatingGame_SendChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatingGameServer).SendChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.DatingGame/SendChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatingGameServer).SendChat(ctx, req.(*ChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DatingGame_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.DatingGame",
	HandlerType: (*DatingGameServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DebugListUsers",
			Handler:    _DatingGame_DebugListUsers_Handler,
		},
		{
			MethodName: "SendChat",
			Handler:    _DatingGame_SendChat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

// Client API for ClientChat service

type ClientChatClient interface {
	Chat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error)
}

type clientChatClient struct {
	cc *grpc.ClientConn
}

func NewClientChatClient(cc *grpc.ClientConn) ClientChatClient {
	return &clientChatClient{cc}
}

func (c *clientChatClient) Chat(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (*ChatResponse, error) {
	out := new(ChatResponse)
	err := grpc.Invoke(ctx, "/helloworld.ClientChat/Chat", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ClientChat service

type ClientChatServer interface {
	Chat(context.Context, *ChatRequest) (*ChatResponse, error)
}

func RegisterClientChatServer(s *grpc.Server, srv ClientChatServer) {
	s.RegisterService(&_ClientChat_serviceDesc, srv)
}

func _ClientChat_Chat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientChatServer).Chat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.ClientChat/Chat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientChatServer).Chat(ctx, req.(*ChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ClientChat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "helloworld.ClientChat",
	HandlerType: (*ClientChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Chat",
			Handler:    _ClientChat_Chat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x52, 0x41, 0x4f, 0xf2, 0x40,
	0x10, 0xa5, 0x94, 0xef, 0x03, 0x07, 0xd2, 0x98, 0x49, 0x84, 0x86, 0x53, 0xdd, 0x03, 0x72, 0xe2,
	0x50, 0x8f, 0x26, 0x6a, 0x42, 0x13, 0x24, 0x31, 0x1e, 0x4a, 0x3c, 0x78, 0x5c, 0x64, 0x02, 0x4d,
	0x96, 0xb6, 0x76, 0x97, 0x98, 0xfe, 0x23, 0x7f, 0xa6, 0xd9, 0x6d, 0x1b, 0x57, 0x45, 0x0f, 0xde,
	0x76, 0xe6, 0xbd, 0x79, 0x33, 0xef, 0x65, 0xc1, 0xdb, 0x93, 0x94, 0x7c, 0x4b, 0x72, 0x96, 0x17,
	0x99, 0xca, 0x10, 0x76, 0x24, 0x44, 0xf6, 0x9a, 0x15, 0x62, 0xc3, 0x18, 0x0c, 0xee, 0x74, 0x15,
	0xd3, 0xcb, 0x81, 0xa4, 0x42, 0x84, 0x4e, 0xca, 0xf7, 0xe4, 0x3b, 0x81, 0x33, 0x3d, 0x89, 0xcd,
	0x9b, 0x4d, 0x00, 0x6a, 0x4e, 0x2e, 0x4a, 0xf4, 0xa1, 0x5b, 0xeb, 0xd5, 0xa4, 0xa6, 0x64, 0x31,
	0x74, 0x1e, 0x25, 0x15, 0xe8, 0x41, 0x3b, 0x89, 0x0c, 0xe8, 0xc6, 0xed, 0x24, 0xc2, 0x00, 0xfa,
	0x9b, 0x44, 0xe6, 0x82, 0x97, 0x0f, 0x5a, 0xba, 0x6d, 0xa6, 0xec, 0x96, 0xd1, 0xe4, 0xea, 0x79,
	0x47, 0xd2, 0x77, 0x03, 0x77, 0xea, 0xc6, 0x4d, 0xc9, 0x46, 0x70, 0x16, 0xd1, 0xfa, 0xb0, 0xbd,
	0x4f, 0xa4, 0xd2, 0xe2, 0xb2, 0x3e, 0x94, 0xdd, 0xc2, 0xf0, 0x2b, 0x20, 0xf3, 0x2c, 0x95, 0x84,
	0x13, 0xf8, 0x77, 0xd0, 0x0d, 0xdf, 0x09, 0xdc, 0x69, 0x3f, 0x3c, 0x9d, 0x7d, 0xd8, 0x9d, 0x69,
	0x66, 0x5c, 0xc1, 0xec, 0x02, 0xfa, 0xf3, 0x1d, 0x57, 0x8d, 0xf3, 0x9f, 0x7d, 0x79, 0x30, 0xa8,
	0x88, 0xd5, 0x82, 0x70, 0x09, 0xdd, 0x45, 0x41, 0xa4, 0xa8, 0xc0, 0x6b, 0xe8, 0xad, 0x78, 0x69,
	0xd2, 0x41, 0xdf, 0x5e, 0x64, 0x87, 0x3a, 0x1e, 0x1e, 0x41, 0x72, 0x51, 0xb2, 0x56, 0xf8, 0xe6,
	0x00, 0x44, 0x5c, 0x25, 0xe9, 0x76, 0xa1, 0x73, 0x78, 0x02, 0xef, 0xb3, 0x29, 0x3c, 0xb7, 0x47,
	0x8f, 0x26, 0x31, 0x66, 0xbf, 0x51, 0xaa, 0x93, 0x59, 0x0b, 0x6f, 0xa0, 0xb7, 0xa2, 0x74, 0xa3,
	0x8d, 0xe0, 0xc8, 0x9e, 0xb0, 0x32, 0x18, 0xfb, 0xdf, 0x81, 0x46, 0x20, 0x5c, 0x02, 0xcc, 0x45,
	0x42, 0xa9, 0x32, 0x12, 0x57, 0xd0, 0xf9, 0xb3, 0xd4, 0xfa, 0xbf, 0xf9, 0x87, 0x97, 0xef, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x19, 0x0d, 0x9a, 0xbf, 0x99, 0x02, 0x00, 0x00,
}
