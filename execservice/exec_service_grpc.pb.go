// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: exec_service.proto

package execservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ExecService_Execute_FullMethodName = "/execservice.ExecService/Execute"
)

// ExecServiceClient is the client API for ExecService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExecServiceClient interface {
	Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ExecuteResponse], error)
}

type execServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExecServiceClient(cc grpc.ClientConnInterface) ExecServiceClient {
	return &execServiceClient{cc}
}

func (c *execServiceClient) Execute(ctx context.Context, in *ExecuteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ExecuteResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ExecService_ServiceDesc.Streams[0], ExecService_Execute_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ExecuteRequest, ExecuteResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ExecService_ExecuteClient = grpc.ServerStreamingClient[ExecuteResponse]

// ExecServiceServer is the server API for ExecService service.
// All implementations must embed UnimplementedExecServiceServer
// for forward compatibility.
type ExecServiceServer interface {
	Execute(*ExecuteRequest, grpc.ServerStreamingServer[ExecuteResponse]) error
	mustEmbedUnimplementedExecServiceServer()
}

// UnimplementedExecServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedExecServiceServer struct{}

func (UnimplementedExecServiceServer) Execute(*ExecuteRequest, grpc.ServerStreamingServer[ExecuteResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedExecServiceServer) mustEmbedUnimplementedExecServiceServer() {}
func (UnimplementedExecServiceServer) testEmbeddedByValue()                     {}

// UnsafeExecServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExecServiceServer will
// result in compilation errors.
type UnsafeExecServiceServer interface {
	mustEmbedUnimplementedExecServiceServer()
}

func RegisterExecServiceServer(s grpc.ServiceRegistrar, srv ExecServiceServer) {
	// If the following call pancis, it indicates UnimplementedExecServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ExecService_ServiceDesc, srv)
}

func _ExecService_Execute_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecuteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExecServiceServer).Execute(m, &grpc.GenericServerStream[ExecuteRequest, ExecuteResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ExecService_ExecuteServer = grpc.ServerStreamingServer[ExecuteResponse]

// ExecService_ServiceDesc is the grpc.ServiceDesc for ExecService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExecService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "execservice.ExecService",
	HandlerType: (*ExecServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Execute",
			Handler:       _ExecService_Execute_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "exec_service.proto",
}
