// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: metrics.proto

package generated

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
	Metrics_StreamCPUClock_FullMethodName       = "/system.Metrics/StreamCPUClock"
	Metrics_StreamCPUUtilization_FullMethodName = "/system.Metrics/StreamCPUUtilization"
	Metrics_StreamHwmon_FullMethodName          = "/system.Metrics/StreamHwmon"
)

// MetricsClient is the client API for Metrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The Metrics service definition
type MetricsClient interface {
	StreamCPUClock(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CPUClock], error)
	StreamCPUUtilization(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CPUUtilization], error)
	StreamHwmon(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Hwmon], error)
}

type metricsClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsClient(cc grpc.ClientConnInterface) MetricsClient {
	return &metricsClient{cc}
}

func (c *metricsClient) StreamCPUClock(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CPUClock], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Metrics_ServiceDesc.Streams[0], Metrics_StreamCPUClock_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Empty, CPUClock]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamCPUClockClient = grpc.ServerStreamingClient[CPUClock]

func (c *metricsClient) StreamCPUUtilization(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[CPUUtilization], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Metrics_ServiceDesc.Streams[1], Metrics_StreamCPUUtilization_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Empty, CPUUtilization]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamCPUUtilizationClient = grpc.ServerStreamingClient[CPUUtilization]

func (c *metricsClient) StreamHwmon(ctx context.Context, in *Empty, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Hwmon], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Metrics_ServiceDesc.Streams[2], Metrics_StreamHwmon_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Empty, Hwmon]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamHwmonClient = grpc.ServerStreamingClient[Hwmon]

// MetricsServer is the server API for Metrics service.
// All implementations must embed UnimplementedMetricsServer
// for forward compatibility.
//
// The Metrics service definition
type MetricsServer interface {
	StreamCPUClock(*Empty, grpc.ServerStreamingServer[CPUClock]) error
	StreamCPUUtilization(*Empty, grpc.ServerStreamingServer[CPUUtilization]) error
	StreamHwmon(*Empty, grpc.ServerStreamingServer[Hwmon]) error
	mustEmbedUnimplementedMetricsServer()
}

// UnimplementedMetricsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMetricsServer struct{}

func (UnimplementedMetricsServer) StreamCPUClock(*Empty, grpc.ServerStreamingServer[CPUClock]) error {
	return status.Errorf(codes.Unimplemented, "method StreamCPUClock not implemented")
}
func (UnimplementedMetricsServer) StreamCPUUtilization(*Empty, grpc.ServerStreamingServer[CPUUtilization]) error {
	return status.Errorf(codes.Unimplemented, "method StreamCPUUtilization not implemented")
}
func (UnimplementedMetricsServer) StreamHwmon(*Empty, grpc.ServerStreamingServer[Hwmon]) error {
	return status.Errorf(codes.Unimplemented, "method StreamHwmon not implemented")
}
func (UnimplementedMetricsServer) mustEmbedUnimplementedMetricsServer() {}
func (UnimplementedMetricsServer) testEmbeddedByValue()                 {}

// UnsafeMetricsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServer will
// result in compilation errors.
type UnsafeMetricsServer interface {
	mustEmbedUnimplementedMetricsServer()
}

func RegisterMetricsServer(s grpc.ServiceRegistrar, srv MetricsServer) {
	// If the following call pancis, it indicates UnimplementedMetricsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Metrics_ServiceDesc, srv)
}

func _Metrics_StreamCPUClock_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetricsServer).StreamCPUClock(m, &grpc.GenericServerStream[Empty, CPUClock]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamCPUClockServer = grpc.ServerStreamingServer[CPUClock]

func _Metrics_StreamCPUUtilization_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetricsServer).StreamCPUUtilization(m, &grpc.GenericServerStream[Empty, CPUUtilization]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamCPUUtilizationServer = grpc.ServerStreamingServer[CPUUtilization]

func _Metrics_StreamHwmon_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MetricsServer).StreamHwmon(m, &grpc.GenericServerStream[Empty, Hwmon]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Metrics_StreamHwmonServer = grpc.ServerStreamingServer[Hwmon]

// Metrics_ServiceDesc is the grpc.ServiceDesc for Metrics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Metrics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system.Metrics",
	HandlerType: (*MetricsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamCPUClock",
			Handler:       _Metrics_StreamCPUClock_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamCPUUtilization",
			Handler:       _Metrics_StreamCPUUtilization_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamHwmon",
			Handler:       _Metrics_StreamHwmon_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "metrics.proto",
}
