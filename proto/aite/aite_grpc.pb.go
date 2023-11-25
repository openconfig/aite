// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: aite.proto

// Package aite implements a disruptive service which can manipulate resources
// within a KNE pod.
//
// It is named after Aite (Até) the green goddess of mischief, delusion and ruin.

package aite

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
	Aite_InterfaceState_FullMethodName  = "/net_sdn_decentralized.aite.Aite/InterfaceState"
	Aite_ImpairInterface_FullMethodName = "/net_sdn_decentralized.aite.Aite/ImpairInterface"
)

// AiteClient is the client API for Aite service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AiteClient interface {
	// InterfaceState changes the state of an interface within the target pod.
	InterfaceState(ctx context.Context, in *InterfaceStateRequest, opts ...grpc.CallOption) (*InterfaceStateResponse, error)
	// ImpairInterface adds a traffic impairment (loss, latency, reordering etc.)
	// to the specified interface.
	ImpairInterface(ctx context.Context, in *ImpairInterfaceRequest, opts ...grpc.CallOption) (*ImpairInterfaceResponse, error)
}

type aiteClient struct {
	cc grpc.ClientConnInterface
}

func NewAiteClient(cc grpc.ClientConnInterface) AiteClient {
	return &aiteClient{cc}
}

func (c *aiteClient) InterfaceState(ctx context.Context, in *InterfaceStateRequest, opts ...grpc.CallOption) (*InterfaceStateResponse, error) {
	out := new(InterfaceStateResponse)
	err := c.cc.Invoke(ctx, Aite_InterfaceState_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aiteClient) ImpairInterface(ctx context.Context, in *ImpairInterfaceRequest, opts ...grpc.CallOption) (*ImpairInterfaceResponse, error) {
	out := new(ImpairInterfaceResponse)
	err := c.cc.Invoke(ctx, Aite_ImpairInterface_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiteServer is the server API for Aite service.
// All implementations must embed UnimplementedAiteServer
// for forward compatibility
type AiteServer interface {
	// InterfaceState changes the state of an interface within the target pod.
	InterfaceState(context.Context, *InterfaceStateRequest) (*InterfaceStateResponse, error)
	// ImpairInterface adds a traffic impairment (loss, latency, reordering etc.)
	// to the specified interface.
	ImpairInterface(context.Context, *ImpairInterfaceRequest) (*ImpairInterfaceResponse, error)
	mustEmbedUnimplementedAiteServer()
}

// UnimplementedAiteServer must be embedded to have forward compatible implementations.
type UnimplementedAiteServer struct {
}

func (UnimplementedAiteServer) InterfaceState(context.Context, *InterfaceStateRequest) (*InterfaceStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InterfaceState not implemented")
}
func (UnimplementedAiteServer) ImpairInterface(context.Context, *ImpairInterfaceRequest) (*ImpairInterfaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImpairInterface not implemented")
}
func (UnimplementedAiteServer) mustEmbedUnimplementedAiteServer() {}

// UnsafeAiteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AiteServer will
// result in compilation errors.
type UnsafeAiteServer interface {
	mustEmbedUnimplementedAiteServer()
}

func RegisterAiteServer(s grpc.ServiceRegistrar, srv AiteServer) {
	s.RegisterService(&Aite_ServiceDesc, srv)
}

func _Aite_InterfaceState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InterfaceStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiteServer).InterfaceState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Aite_InterfaceState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiteServer).InterfaceState(ctx, req.(*InterfaceStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Aite_ImpairInterface_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImpairInterfaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiteServer).ImpairInterface(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Aite_ImpairInterface_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiteServer).ImpairInterface(ctx, req.(*ImpairInterfaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Aite_ServiceDesc is the grpc.ServiceDesc for Aite service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Aite_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "net_sdn_decentralized.aite.Aite",
	HandlerType: (*AiteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InterfaceState",
			Handler:    _Aite_InterfaceState_Handler,
		},
		{
			MethodName: "ImpairInterface",
			Handler:    _Aite_ImpairInterface_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aite.proto",
}
