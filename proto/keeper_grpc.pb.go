// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: keeper.proto

package proto

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
	GophKeeper_Register_FullMethodName   = "/keeper.GophKeeper/Register"
	GophKeeper_Login_FullMethodName      = "/keeper.GophKeeper/Login"
	GophKeeper_CreateData_FullMethodName = "/keeper.GophKeeper/CreateData"
	GophKeeper_GetAllData_FullMethodName = "/keeper.GophKeeper/GetAllData"
	GophKeeper_DeleteData_FullMethodName = "/keeper.GophKeeper/DeleteData"
	GophKeeper_UpdateData_FullMethodName = "/keeper.GophKeeper/UpdateData"
)

// GophKeeperClient is the client API for GophKeeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	// загрузка данных
	CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error)
	// Получение всех данных пользователя
	GetAllData(ctx context.Context, in *GetAllDataRequest, opts ...grpc.CallOption) (*GetAllDataResponse, error)
	// удаление данных пользователя
	DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error)
	// обновление данных пользователя
	UpdateData(ctx context.Context, in *UpdateDataRequest, opts ...grpc.CallOption) (*UpdateDataResponse, error)
}

type gophKeeperClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperClient(cc grpc.ClientConnInterface) GophKeeperClient {
	return &gophKeeperClient{cc}
}

func (c *gophKeeperClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, GophKeeper_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDataResponse)
	err := c.cc.Invoke(ctx, GophKeeper_CreateData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) GetAllData(ctx context.Context, in *GetAllDataRequest, opts ...grpc.CallOption) (*GetAllDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllDataResponse)
	err := c.cc.Invoke(ctx, GophKeeper_GetAllData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteDataResponse)
	err := c.cc.Invoke(ctx, GophKeeper_DeleteData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperClient) UpdateData(ctx context.Context, in *UpdateDataRequest, opts ...grpc.CallOption) (*UpdateDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateDataResponse)
	err := c.cc.Invoke(ctx, GophKeeper_UpdateData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServer is the server API for GophKeeper service.
// All implementations must embed UnimplementedGophKeeperServer
// for forward compatibility.
type GophKeeperServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	// загрузка данных
	CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error)
	// Получение всех данных пользователя
	GetAllData(context.Context, *GetAllDataRequest) (*GetAllDataResponse, error)
	// удаление данных пользователя
	DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error)
	// обновление данных пользователя
	UpdateData(context.Context, *UpdateDataRequest) (*UpdateDataResponse, error)
	mustEmbedUnimplementedGophKeeperServer()
}

// UnimplementedGophKeeperServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGophKeeperServer struct{}

func (UnimplementedGophKeeperServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGophKeeperServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedGophKeeperServer) CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateData not implemented")
}
func (UnimplementedGophKeeperServer) GetAllData(context.Context, *GetAllDataRequest) (*GetAllDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllData not implemented")
}
func (UnimplementedGophKeeperServer) DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteData not implemented")
}
func (UnimplementedGophKeeperServer) UpdateData(context.Context, *UpdateDataRequest) (*UpdateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateData not implemented")
}
func (UnimplementedGophKeeperServer) mustEmbedUnimplementedGophKeeperServer() {}
func (UnimplementedGophKeeperServer) testEmbeddedByValue()                    {}

// UnsafeGophKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServer will
// result in compilation errors.
type UnsafeGophKeeperServer interface {
	mustEmbedUnimplementedGophKeeperServer()
}

func RegisterGophKeeperServer(s grpc.ServiceRegistrar, srv GophKeeperServer) {
	// If the following call pancis, it indicates UnimplementedGophKeeperServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GophKeeper_ServiceDesc, srv)
}

func _GophKeeper_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_CreateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).CreateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_CreateData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).CreateData(ctx, req.(*CreateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_GetAllData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).GetAllData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_GetAllData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).GetAllData(ctx, req.(*GetAllDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_DeleteData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).DeleteData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_DeleteData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).DeleteData(ctx, req.(*DeleteDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeper_UpdateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServer).UpdateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeper_UpdateData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServer).UpdateData(ctx, req.(*UpdateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeper_ServiceDesc is the grpc.ServiceDesc for GophKeeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keeper.GophKeeper",
	HandlerType: (*GophKeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _GophKeeper_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _GophKeeper_Login_Handler,
		},
		{
			MethodName: "CreateData",
			Handler:    _GophKeeper_CreateData_Handler,
		},
		{
			MethodName: "GetAllData",
			Handler:    _GophKeeper_GetAllData_Handler,
		},
		{
			MethodName: "DeleteData",
			Handler:    _GophKeeper_DeleteData_Handler,
		},
		{
			MethodName: "UpdateData",
			Handler:    _GophKeeper_UpdateData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "keeper.proto",
}
