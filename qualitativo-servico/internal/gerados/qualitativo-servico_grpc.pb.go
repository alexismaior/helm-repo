// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package gerados

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	Classificacoes(ctx context.Context, in *ClassificacoesRequest, opts ...grpc.CallOption) (*ClassificacoesReply, error)
	Localizadores(ctx context.Context, in *LocalizadoresRequest, opts ...grpc.CallOption) (*LocalizadoresReply, error)
	Decodificar(ctx context.Context, in *DecodificarRequest, opts ...grpc.CallOption) (*DecodificarReply, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Classificacoes(ctx context.Context, in *ClassificacoesRequest, opts ...grpc.CallOption) (*ClassificacoesReply, error) {
	out := new(ClassificacoesReply)
	err := c.cc.Invoke(ctx, "/qualitativo.Service/Classificacoes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Localizadores(ctx context.Context, in *LocalizadoresRequest, opts ...grpc.CallOption) (*LocalizadoresReply, error) {
	out := new(LocalizadoresReply)
	err := c.cc.Invoke(ctx, "/qualitativo.Service/Localizadores", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) Decodificar(ctx context.Context, in *DecodificarRequest, opts ...grpc.CallOption) (*DecodificarReply, error) {
	out := new(DecodificarReply)
	err := c.cc.Invoke(ctx, "/qualitativo.Service/Decodificar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	Classificacoes(context.Context, *ClassificacoesRequest) (*ClassificacoesReply, error)
	Localizadores(context.Context, *LocalizadoresRequest) (*LocalizadoresReply, error)
	Decodificar(context.Context, *DecodificarRequest) (*DecodificarReply, error)
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) Classificacoes(context.Context, *ClassificacoesRequest) (*ClassificacoesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Classificacoes not implemented")
}
func (*UnimplementedServiceServer) Localizadores(context.Context, *LocalizadoresRequest) (*LocalizadoresReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Localizadores not implemented")
}
func (*UnimplementedServiceServer) Decodificar(context.Context, *DecodificarRequest) (*DecodificarReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decodificar not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Classificacoes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClassificacoesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Classificacoes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/qualitativo.Service/Classificacoes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Classificacoes(ctx, req.(*ClassificacoesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Localizadores_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocalizadoresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Localizadores(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/qualitativo.Service/Localizadores",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Localizadores(ctx, req.(*LocalizadoresRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_Decodificar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecodificarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Decodificar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/qualitativo.Service/Decodificar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Decodificar(ctx, req.(*DecodificarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "qualitativo.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Classificacoes",
			Handler:    _Service_Classificacoes_Handler,
		},
		{
			MethodName: "Localizadores",
			Handler:    _Service_Localizadores_Handler,
		},
		{
			MethodName: "Decodificar",
			Handler:    _Service_Decodificar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/gerados/qualitativo-servico.proto",
}
