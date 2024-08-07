// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.2
// source: proto/mon.proto

package mon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	MonitorService_GetMetric_FullMethodName          = "/mon.MonitorService/GetMetric"
	MonitorService_UpdateMetric_FullMethodName       = "/mon.MonitorService/UpdateMetric"
	MonitorService_UpdateMetricsBatch_FullMethodName = "/mon.MonitorService/UpdateMetricsBatch"
	MonitorService_PingDB_FullMethodName             = "/mon.MonitorService/PingDB"
)

// MonitorServiceClient is the client API for MonitorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MonitorServiceClient interface {
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error)
	UpdateMetricsBatch(ctx context.Context, in *UpdateMetricsBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type monitorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMonitorServiceClient(cc grpc.ClientConnInterface) MonitorServiceClient {
	return &monitorServiceClient{cc}
}

func (c *monitorServiceClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, MonitorService_GetMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitorServiceClient) UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMetricResponse)
	err := c.cc.Invoke(ctx, MonitorService_UpdateMetric_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitorServiceClient) UpdateMetricsBatch(ctx context.Context, in *UpdateMetricsBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MonitorService_UpdateMetricsBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *monitorServiceClient) PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MonitorService_PingDB_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MonitorServiceServer is the server API for MonitorService service.
// All implementations must embed UnimplementedMonitorServiceServer
// for forward compatibility
type MonitorServiceServer interface {
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error)
	UpdateMetricsBatch(context.Context, *UpdateMetricsBatchRequest) (*emptypb.Empty, error)
	PingDB(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedMonitorServiceServer()
}

// UnimplementedMonitorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMonitorServiceServer struct {
}

func (UnimplementedMonitorServiceServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMonitorServiceServer) UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMonitorServiceServer) UpdateMetricsBatch(context.Context, *UpdateMetricsBatchRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetricsBatch not implemented")
}
func (UnimplementedMonitorServiceServer) PingDB(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingDB not implemented")
}
func (UnimplementedMonitorServiceServer) mustEmbedUnimplementedMonitorServiceServer() {}

// UnsafeMonitorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MonitorServiceServer will
// result in compilation errors.
type UnsafeMonitorServiceServer interface {
	mustEmbedUnimplementedMonitorServiceServer()
}

func RegisterMonitorServiceServer(s grpc.ServiceRegistrar, srv MonitorServiceServer) {
	s.RegisterService(&MonitorService_ServiceDesc, srv)
}

func _MonitorService_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServiceServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitorService_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServiceServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MonitorService_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServiceServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitorService_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServiceServer).UpdateMetric(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MonitorService_UpdateMetricsBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricsBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServiceServer).UpdateMetricsBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitorService_UpdateMetricsBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServiceServer).UpdateMetricsBatch(ctx, req.(*UpdateMetricsBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MonitorService_PingDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MonitorServiceServer).PingDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MonitorService_PingDB_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MonitorServiceServer).PingDB(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// MonitorService_ServiceDesc is the grpc.ServiceDesc for MonitorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MonitorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mon.MonitorService",
	HandlerType: (*MonitorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMetric",
			Handler:    _MonitorService_GetMetric_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _MonitorService_UpdateMetric_Handler,
		},
		{
			MethodName: "UpdateMetricsBatch",
			Handler:    _MonitorService_UpdateMetricsBatch_Handler,
		},
		{
			MethodName: "PingDB",
			Handler:    _MonitorService_PingDB_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/mon.proto",
}
