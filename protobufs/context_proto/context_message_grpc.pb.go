// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: context_message.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	SpaceContext_Hit_FullMethodName          = "/context_message.SpaceContext/Hit"
	SpaceContext_RoomRegister_FullMethodName = "/context_message.SpaceContext/RoomRegister"
)

// SpaceContextClient is the client API for SpaceContext service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SpaceContextClient interface {
	// Other player's incoming damage
	Hit(ctx context.Context, in *Damage, opts ...grpc.CallOption) (*Empty, error)
	// Room register itself in Context
	RoomRegister(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SpaceContext_RoomRegisterClient, error)
}

type spaceContextClient struct {
	cc grpc.ClientConnInterface
}

func NewSpaceContextClient(cc grpc.ClientConnInterface) SpaceContextClient {
	return &spaceContextClient{cc}
}

func (c *spaceContextClient) Hit(ctx context.Context, in *Damage, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, SpaceContext_Hit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceContextClient) RoomRegister(ctx context.Context, in *Empty, opts ...grpc.CallOption) (SpaceContext_RoomRegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &SpaceContext_ServiceDesc.Streams[0], SpaceContext_RoomRegister_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &spaceContextRoomRegisterClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SpaceContext_RoomRegisterClient interface {
	Recv() (*Damage, error)
	grpc.ClientStream
}

type spaceContextRoomRegisterClient struct {
	grpc.ClientStream
}

func (x *spaceContextRoomRegisterClient) Recv() (*Damage, error) {
	m := new(Damage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SpaceContextServer is the server API for SpaceContext service.
// All implementations must embed UnimplementedSpaceContextServer
// for forward compatibility
type SpaceContextServer interface {
	// Other player's incoming damage
	Hit(context.Context, *Damage) (*Empty, error)
	// Room register itself in Context
	RoomRegister(*Empty, SpaceContext_RoomRegisterServer) error
	mustEmbedUnimplementedSpaceContextServer()
}

// UnimplementedSpaceContextServer must be embedded to have forward compatible implementations.
type UnimplementedSpaceContextServer struct {
}

func (UnimplementedSpaceContextServer) Hit(context.Context, *Damage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hit not implemented")
}
func (UnimplementedSpaceContextServer) RoomRegister(*Empty, SpaceContext_RoomRegisterServer) error {
	return status.Errorf(codes.Unimplemented, "method RoomRegister not implemented")
}
func (UnimplementedSpaceContextServer) mustEmbedUnimplementedSpaceContextServer() {}

// UnsafeSpaceContextServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SpaceContextServer will
// result in compilation errors.
type UnsafeSpaceContextServer interface {
	mustEmbedUnimplementedSpaceContextServer()
}

func RegisterSpaceContextServer(s grpc.ServiceRegistrar, srv SpaceContextServer) {
	s.RegisterService(&SpaceContext_ServiceDesc, srv)
}

func _SpaceContext_Hit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Damage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceContextServer).Hit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceContext_Hit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceContextServer).Hit(ctx, req.(*Damage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceContext_RoomRegister_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SpaceContextServer).RoomRegister(m, &spaceContextRoomRegisterServer{stream})
}

type SpaceContext_RoomRegisterServer interface {
	Send(*Damage) error
	grpc.ServerStream
}

type spaceContextRoomRegisterServer struct {
	grpc.ServerStream
}

func (x *spaceContextRoomRegisterServer) Send(m *Damage) error {
	return x.ServerStream.SendMsg(m)
}

// SpaceContext_ServiceDesc is the grpc.ServiceDesc for SpaceContext service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SpaceContext_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "context_message.SpaceContext",
	HandlerType: (*SpaceContextServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hit",
			Handler:    _SpaceContext_Hit_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RoomRegister",
			Handler:       _SpaceContext_RoomRegister_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "context_message.proto",
}