// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: aite.proto

// Package aite implements a disruptive service which can manipulate resources
// within a KNE pod.
//
// It is named after Aite (Até) the Greek goddess of mischief, delusion and ruin.

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
	Aite_SetInterface_FullMethodName = "/openconfig.aite.Aite/SetInterface"
)

// AiteClient is the client API for Aite service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AiteClient interface {
	// SetInterface changes the state of an interface within the target pod.
	SetInterface(ctx context.Context, in *SetInterfaceRequest, opts ...grpc.CallOption) (*SetInterfaceResponse, error)
}

type aiteClient struct {
	cc grpc.ClientConnInterface
}

func NewAiteClient(cc grpc.ClientConnInterface) AiteClient {
	return &aiteClient{cc}
}

func (c *aiteClient) SetInterface(ctx context.Context, in *SetInterfaceRequest, opts ...grpc.CallOption) (*SetInterfaceResponse, error) {
	out := new(SetInterfaceResponse)
	err := c.cc.Invoke(ctx, Aite_SetInterface_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AiteServer is the server API for Aite service.
// All implementations must embed UnimplementedAiteServer
// for forward compatibility
type AiteServer interface {
	// SetInterface changes the state of an interface within the target pod.
	SetInterface(context.Context, *SetInterfaceRequest) (*SetInterfaceResponse, error)
	mustEmbedUnimplementedAiteServer()
}

// UnimplementedAiteServer must be embedded to have forward compatible implementations.
type UnimplementedAiteServer struct {
}

func (UnimplementedAiteServer) SetInterface(context.Context, *SetInterfaceRequest) (*SetInterfaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetInterface not implemented")
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

func _Aite_SetInterface_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetInterfaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AiteServer).SetInterface(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Aite_SetInterface_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AiteServer).SetInterface(ctx, req.(*SetInterfaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Aite_ServiceDesc is the grpc.ServiceDesc for Aite service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Aite_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "openconfig.aite.Aite",
	HandlerType: (*AiteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetInterface",
			Handler:    _Aite_SetInterface_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "aite.proto",
}
