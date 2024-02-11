// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: room_message.proto

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
	RoomGRPC_RoomSubscribe_FullMethodName = "/room_message.RoomGRPC/RoomSubscribe"
)

// RoomGRPCClient is the client API for RoomGRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomGRPCClient interface {
	// rpc Electricity(Source) returns (Empty) {}
	// rpc Air(Source) returns (Empty) {}
	// rpc RoomSubscribe(Sub) returns (stream Supply){}
	RoomSubscribe(ctx context.Context, opts ...grpc.CallOption) (RoomGRPC_RoomSubscribeClient, error)
}

type roomGRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomGRPCClient(cc grpc.ClientConnInterface) RoomGRPCClient {
	return &roomGRPCClient{cc}
}

func (c *roomGRPCClient) RoomSubscribe(ctx context.Context, opts ...grpc.CallOption) (RoomGRPC_RoomSubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &RoomGRPC_ServiceDesc.Streams[0], RoomGRPC_RoomSubscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &roomGRPCRoomSubscribeClient{stream}
	return x, nil
}

type RoomGRPC_RoomSubscribeClient interface {
	Send(*Supply) error
	Recv() (*Supply, error)
	grpc.ClientStream
}

type roomGRPCRoomSubscribeClient struct {
	grpc.ClientStream
}

func (x *roomGRPCRoomSubscribeClient) Send(m *Supply) error {
	return x.ClientStream.SendMsg(m)
}

func (x *roomGRPCRoomSubscribeClient) Recv() (*Supply, error) {
	m := new(Supply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RoomGRPCServer is the server API for RoomGRPC service.
// All implementations must embed UnimplementedRoomGRPCServer
// for forward compatibility
type RoomGRPCServer interface {
	// rpc Electricity(Source) returns (Empty) {}
	// rpc Air(Source) returns (Empty) {}
	// rpc RoomSubscribe(Sub) returns (stream Supply){}
	RoomSubscribe(RoomGRPC_RoomSubscribeServer) error
	mustEmbedUnimplementedRoomGRPCServer()
}

// UnimplementedRoomGRPCServer must be embedded to have forward compatible implementations.
type UnimplementedRoomGRPCServer struct {
}

func (UnimplementedRoomGRPCServer) RoomSubscribe(RoomGRPC_RoomSubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method RoomSubscribe not implemented")
}
func (UnimplementedRoomGRPCServer) mustEmbedUnimplementedRoomGRPCServer() {}

// UnsafeRoomGRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomGRPCServer will
// result in compilation errors.
type UnsafeRoomGRPCServer interface {
	mustEmbedUnimplementedRoomGRPCServer()
}

func RegisterRoomGRPCServer(s grpc.ServiceRegistrar, srv RoomGRPCServer) {
	s.RegisterService(&RoomGRPC_ServiceDesc, srv)
}

func _RoomGRPC_RoomSubscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RoomGRPCServer).RoomSubscribe(&roomGRPCRoomSubscribeServer{stream})
}

type RoomGRPC_RoomSubscribeServer interface {
	Send(*Supply) error
	Recv() (*Supply, error)
	grpc.ServerStream
}

type roomGRPCRoomSubscribeServer struct {
	grpc.ServerStream
}

func (x *roomGRPCRoomSubscribeServer) Send(m *Supply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *roomGRPCRoomSubscribeServer) Recv() (*Supply, error) {
	m := new(Supply)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RoomGRPC_ServiceDesc is the grpc.ServiceDesc for RoomGRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomGRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "room_message.RoomGRPC",
	HandlerType: (*RoomGRPCServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RoomSubscribe",
			Handler:       _RoomGRPC_RoomSubscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "room_message.proto",
}
