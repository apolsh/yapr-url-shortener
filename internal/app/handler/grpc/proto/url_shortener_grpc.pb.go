// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: url_shortener.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// URLShortenerClient is the client API for URLShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
// URLShortenerClient
type URLShortenerClient interface {
	// PingDB
	PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingDBResponse, error)
	// GetShortenURLByID
	GetShortenURLByID(ctx context.Context, in *GetShortenURLByIDRequest, opts ...grpc.CallOption) (*GetShortenURLByIDResponse, error)
	// GetShortenURLsByUser
	GetShortenURLsByUser(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetShortenURLsByUserResponse, error)
	// SaveShortenURL
	SaveShortenURL(ctx context.Context, in *SaveShortenURLRequest, opts ...grpc.CallOption) (*SaveShortenURLResponse, error)
	// GetAppStats
	GetAppStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAppStatsResponse, error)
	// SaveShortenURLsInBatch
	SaveShortenURLsInBatch(ctx context.Context, in *SaveShortenURLsInBatchRequest, opts ...grpc.CallOption) (*SaveShortenURLsInBatchResponse, error)
	// DeleteShortenURLsInBatch
	DeleteShortenURLsInBatch(ctx context.Context, in *DeleteShortenURLsInBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type uRLShortenerClient struct {
	cc grpc.ClientConnInterface
}

// NewURLShortenerClient
func NewURLShortenerClient(cc grpc.ClientConnInterface) URLShortenerClient {
	return &uRLShortenerClient{cc}
}

// PingDB
func (c *uRLShortenerClient) PingDB(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingDBResponse, error) {
	out := new(PingDBResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/PingDB", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetShortenURLByID
func (c *uRLShortenerClient) GetShortenURLByID(ctx context.Context, in *GetShortenURLByIDRequest, opts ...grpc.CallOption) (*GetShortenURLByIDResponse, error) {
	out := new(GetShortenURLByIDResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/GetShortenURLByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetShortenURLsByUser
func (c *uRLShortenerClient) GetShortenURLsByUser(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetShortenURLsByUserResponse, error) {
	out := new(GetShortenURLsByUserResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/GetShortenURLsByUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SaveShortenURL
func (c *uRLShortenerClient) SaveShortenURL(ctx context.Context, in *SaveShortenURLRequest, opts ...grpc.CallOption) (*SaveShortenURLResponse, error) {
	out := new(SaveShortenURLResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/SaveShortenURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetAppStats
func (c *uRLShortenerClient) GetAppStats(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAppStatsResponse, error) {
	out := new(GetAppStatsResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/GetAppStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SaveShortenURLsInBatch
func (c *uRLShortenerClient) SaveShortenURLsInBatch(ctx context.Context, in *SaveShortenURLsInBatchRequest, opts ...grpc.CallOption) (*SaveShortenURLsInBatchResponse, error) {
	out := new(SaveShortenURLsInBatchResponse)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/SaveShortenURLsInBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteShortenURLsInBatch
func (c *uRLShortenerClient) DeleteShortenURLsInBatch(ctx context.Context, in *DeleteShortenURLsInBatchRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.URLShortener/DeleteShortenURLsInBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServer is the server API for URLShortener service.
// All implementations must embed UnimplementedURLShortenerServer
// for forward compatibility
// URLShortenerServer
type URLShortenerServer interface {
	// PingDB
	PingDB(context.Context, *emptypb.Empty) (*PingDBResponse, error)
	// GetShortenURLByID
	GetShortenURLByID(context.Context, *GetShortenURLByIDRequest) (*GetShortenURLByIDResponse, error)
	// GetShortenURLsByUser
	GetShortenURLsByUser(context.Context, *emptypb.Empty) (*GetShortenURLsByUserResponse, error)
	// SaveShortenURL
	SaveShortenURL(context.Context, *SaveShortenURLRequest) (*SaveShortenURLResponse, error)
	// GetAppStats
	GetAppStats(context.Context, *emptypb.Empty) (*GetAppStatsResponse, error)
	// SaveShortenURLsInBatch
	SaveShortenURLsInBatch(context.Context, *SaveShortenURLsInBatchRequest) (*SaveShortenURLsInBatchResponse, error)
	// DeleteShortenURLsInBatch
	DeleteShortenURLsInBatch(context.Context, *DeleteShortenURLsInBatchRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedURLShortenerServer()
}

// UnimplementedURLShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedURLShortenerServer struct {
}

// PingDB
func (UnimplementedURLShortenerServer) PingDB(context.Context, *emptypb.Empty) (*PingDBResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingDB not implemented")
}

// GetShortenURLByID
func (UnimplementedURLShortenerServer) GetShortenURLByID(context.Context, *GetShortenURLByIDRequest) (*GetShortenURLByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortenURLByID not implemented")
}

// GetShortenURLsByUser
func (UnimplementedURLShortenerServer) GetShortenURLsByUser(context.Context, *emptypb.Empty) (*GetShortenURLsByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetShortenURLsByUser not implemented")
}

// SaveShortenURL
func (UnimplementedURLShortenerServer) SaveShortenURL(context.Context, *SaveShortenURLRequest) (*SaveShortenURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveShortenURL not implemented")
}

// GetAppStats
func (UnimplementedURLShortenerServer) GetAppStats(context.Context, *emptypb.Empty) (*GetAppStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAppStats not implemented")
}

// SaveShortenURLsInBatch
func (UnimplementedURLShortenerServer) SaveShortenURLsInBatch(context.Context, *SaveShortenURLsInBatchRequest) (*SaveShortenURLsInBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveShortenURLsInBatch not implemented")
}

// DeleteShortenURLsInBatch
func (UnimplementedURLShortenerServer) DeleteShortenURLsInBatch(context.Context, *DeleteShortenURLsInBatchRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteShortenURLsInBatch not implemented")
}
func (UnimplementedURLShortenerServer) mustEmbedUnimplementedURLShortenerServer() {}

// UnsafeURLShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServer will
// result in compilation errors.
type UnsafeURLShortenerServer interface {
	mustEmbedUnimplementedURLShortenerServer()
}

// RegisterURLShortenerServer
func RegisterURLShortenerServer(s grpc.ServiceRegistrar, srv URLShortenerServer) {
	s.RegisterService(&URLShortener_ServiceDesc, srv)
}

func _URLShortener_PingDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).PingDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/PingDB",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).PingDB(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetShortenURLByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetShortenURLByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetShortenURLByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/GetShortenURLByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetShortenURLByID(ctx, req.(*GetShortenURLByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetShortenURLsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetShortenURLsByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/GetShortenURLsByUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetShortenURLsByUser(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_SaveShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveShortenURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).SaveShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/SaveShortenURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).SaveShortenURL(ctx, req.(*SaveShortenURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetAppStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetAppStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/GetAppStats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetAppStats(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_SaveShortenURLsInBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveShortenURLsInBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).SaveShortenURLsInBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/SaveShortenURLsInBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).SaveShortenURLsInBatch(ctx, req.(*SaveShortenURLsInBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteShortenURLsInBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteShortenURLsInBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteShortenURLsInBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.URLShortener/DeleteShortenURLsInBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteShortenURLsInBatch(ctx, req.(*DeleteShortenURLsInBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortener_ServiceDesc is the grpc.ServiceDesc for URLShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.URLShortener",
	HandlerType: (*URLShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingDB",
			Handler:    _URLShortener_PingDB_Handler,
		},
		{
			MethodName: "GetShortenURLByID",
			Handler:    _URLShortener_GetShortenURLByID_Handler,
		},
		{
			MethodName: "GetShortenURLsByUser",
			Handler:    _URLShortener_GetShortenURLsByUser_Handler,
		},
		{
			MethodName: "SaveShortenURL",
			Handler:    _URLShortener_SaveShortenURL_Handler,
		},
		{
			MethodName: "GetAppStats",
			Handler:    _URLShortener_GetAppStats_Handler,
		},
		{
			MethodName: "SaveShortenURLsInBatch",
			Handler:    _URLShortener_SaveShortenURLsInBatch_Handler,
		},
		{
			MethodName: "DeleteShortenURLsInBatch",
			Handler:    _URLShortener_DeleteShortenURLsInBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "url_shortener.proto",
}