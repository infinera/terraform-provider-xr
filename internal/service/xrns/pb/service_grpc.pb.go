// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: service.proto

package pb

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

// NamingServiceClient is the client API for NamingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NamingServiceClient interface {
	GetDeviceCount(ctx context.Context, in *GetDeviceCountRequest, opts ...grpc.CallOption) (*GetDeviceCountResponse, error)
	GetDeviceByName(ctx context.Context, in *GetDeviceByNameRequest, opts ...grpc.CallOption) (*Device, error)
	GetAllDevices(ctx context.Context, in *GetAllDevicesRequest, opts ...grpc.CallOption) (NamingService_GetAllDevicesClient, error)
}

type namingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNamingServiceClient(cc grpc.ClientConnInterface) NamingServiceClient {
	return &namingServiceClient{cc}
}

func (c *namingServiceClient) GetDeviceCount(ctx context.Context, in *GetDeviceCountRequest, opts ...grpc.CallOption) (*GetDeviceCountResponse, error) {
	out := new(GetDeviceCountResponse)
	err := c.cc.Invoke(ctx, "/xrns.NamingService/GetDeviceCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namingServiceClient) GetDeviceByName(ctx context.Context, in *GetDeviceByNameRequest, opts ...grpc.CallOption) (*Device, error) {
	out := new(Device)
	err := c.cc.Invoke(ctx, "/xrns.NamingService/GetDeviceByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *namingServiceClient) GetAllDevices(ctx context.Context, in *GetAllDevicesRequest, opts ...grpc.CallOption) (NamingService_GetAllDevicesClient, error) {
	stream, err := c.cc.NewStream(ctx, &NamingService_ServiceDesc.Streams[0], "/xrns.NamingService/GetAllDevices", opts...)
	if err != nil {
		return nil, err
	}
	x := &namingServiceGetAllDevicesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NamingService_GetAllDevicesClient interface {
	Recv() (*Device, error)
	grpc.ClientStream
}

type namingServiceGetAllDevicesClient struct {
	grpc.ClientStream
}

func (x *namingServiceGetAllDevicesClient) Recv() (*Device, error) {
	m := new(Device)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// NamingServiceServer is the server API for NamingService service.
// All implementations must embed UnimplementedNamingServiceServer
// for forward compatibility
type NamingServiceServer interface {
	GetDeviceCount(context.Context, *GetDeviceCountRequest) (*GetDeviceCountResponse, error)
	GetDeviceByName(context.Context, *GetDeviceByNameRequest) (*Device, error)
	GetAllDevices(*GetAllDevicesRequest, NamingService_GetAllDevicesServer) error
	mustEmbedUnimplementedNamingServiceServer()
}

// UnimplementedNamingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNamingServiceServer struct {
}

func (UnimplementedNamingServiceServer) GetDeviceCount(context.Context, *GetDeviceCountRequest) (*GetDeviceCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceCount not implemented")
}
func (UnimplementedNamingServiceServer) GetDeviceByName(context.Context, *GetDeviceByNameRequest) (*Device, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeviceByName not implemented")
}
func (UnimplementedNamingServiceServer) GetAllDevices(*GetAllDevicesRequest, NamingService_GetAllDevicesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllDevices not implemented")
}
func (UnimplementedNamingServiceServer) mustEmbedUnimplementedNamingServiceServer() {}

// UnsafeNamingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NamingServiceServer will
// result in compilation errors.
type UnsafeNamingServiceServer interface {
	mustEmbedUnimplementedNamingServiceServer()
}

func RegisterNamingServiceServer(s grpc.ServiceRegistrar, srv NamingServiceServer) {
	s.RegisterService(&NamingService_ServiceDesc, srv)
}

func _NamingService_GetDeviceCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamingServiceServer).GetDeviceCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xrns.NamingService/GetDeviceCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamingServiceServer).GetDeviceCount(ctx, req.(*GetDeviceCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamingService_GetDeviceByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeviceByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NamingServiceServer).GetDeviceByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/xrns.NamingService/GetDeviceByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NamingServiceServer).GetDeviceByName(ctx, req.(*GetDeviceByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NamingService_GetAllDevices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllDevicesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NamingServiceServer).GetAllDevices(m, &namingServiceGetAllDevicesServer{stream})
}

type NamingService_GetAllDevicesServer interface {
	Send(*Device) error
	grpc.ServerStream
}

type namingServiceGetAllDevicesServer struct {
	grpc.ServerStream
}

func (x *namingServiceGetAllDevicesServer) Send(m *Device) error {
	return x.ServerStream.SendMsg(m)
}

// NamingService_ServiceDesc is the grpc.ServiceDesc for NamingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NamingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "xrns.NamingService",
	HandlerType: (*NamingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDeviceCount",
			Handler:    _NamingService_GetDeviceCount_Handler,
		},
		{
			MethodName: "GetDeviceByName",
			Handler:    _NamingService_GetDeviceByName_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAllDevices",
			Handler:       _NamingService_GetAllDevices_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "service.proto",
}
