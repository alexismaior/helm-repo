package gerados

import (
	"context"

	"git.sof.intra/siop/qualitativo-servico"
	"git.sof.intra/siop/qualitativo-servico/internal"

	"git.sof.intra/siop/framework"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type (
	gokitService struct {
		endpoints map[string]endpoint.Endpoint
	}
	gokitServiceGrpcServer struct {
		UnimplementedServiceServer
		handlers map[string]kitgrpc.Handler
	}
)

func init() {
	framework.RegisterServiceFactory(
		func(endpoints map[string]endpoint.Endpoint) interface{} {
			return &gokitService{endpoints: endpoints}
		},
		"git.sof.intra/siop/qualitativo-servico/internal/Service",
		"git.sof.intra/siop/qualitativo-servico/Service",
		"git.sof.intra/siop/qualitativo-servico/mock/MockService",
	)
	framework.RegisterGrpcServerFactory(
		func(handlers map[string]kitgrpc.Handler) interface{} {
			return &gokitServiceGrpcServer{handlers: handlers}
		},
		RegisterServiceServer,
		"qualitativo.Service",
		"git.sof.intra/siop/qualitativo-servico/internal/Service",
		"git.sof.intra/siop/qualitativo-servico/Service",
		"git.sof.intra/siop/qualitativo-servico/mock/MockService",
	)
	framework.RegisterParameterNames(
		"git.sof.intra/siop/qualitativo-servico/internal/Service",
		map[string][2][]string{
			"Classificacoes": {
				{"exercicio", "classificacoesCodificadas"},
				{"classificacoes"},
			},
			"Localizadores": {
				{"filtro", "fetch", "n"},
				{"localizadores"},
			},
			"Decodificar": {
				{"classificacoesCodificadas"},
				{"classificacoes"},
			},
			"Esferas": {
				{"filtro"},
				{"esferas"},
			},
		},
	)
}

func (_s *gokitService) Classificacoes(_ctx context.Context, _exercicio int32, _classificacoesCodificadas []string) (_classificacoes []*qualitativo.Classificacao, _err error) {
	var req = []interface{}{_exercicio, _classificacoesCodificadas}
	rep, err := _s.endpoints["Classificacoes"](_ctx, req)
	if rep != nil {
		out := rep.([]interface{})
		_classificacoes = out[0].([]*qualitativo.Classificacao)
	}
	return _classificacoes, framework.ErrorWithoutGrpcStuff(err)
}

func (_s *gokitService) Localizadores(_ctx context.Context, _filtro *qualitativo.FiltroLocalizador, _fetch *qualitativo.FetchLocalizador, _n int32) (_localizadores []*qualitativo.Localizador, _err error) {
	var req = []interface{}{_filtro, _fetch, _n}
	rep, err := _s.endpoints["Localizadores"](_ctx, req)
	if rep != nil {
		out := rep.([]interface{})
		_localizadores = out[0].([]*qualitativo.Localizador)
	}
	return _localizadores, framework.ErrorWithoutGrpcStuff(err)
}

func (_s *gokitService) Decodificar(_ctx context.Context, _classificacoesCodificadas []string) (_classificacoes []qualitativo.Classificacao, _err error) {
	var req = []interface{}{_classificacoesCodificadas}
	rep, err := _s.endpoints["Decodificar"](_ctx, req)
	if rep != nil {
		out := rep.([]interface{})
		_classificacoes = out[0].([]qualitativo.Classificacao)
	}
	return _classificacoes, framework.ErrorWithoutGrpcStuff(err)
}

func (_s *gokitService) Esferas(_ctx context.Context, _filtro internal.Esfera) (_esferas []internal.Esfera, _err error) {
	var req = []interface{}{_filtro}
	rep, err := _s.endpoints["Esferas"](_ctx, req)
	if rep != nil {
		out := rep.([]interface{})
		_esferas = out[0].([]internal.Esfera)
	}
	return _esferas, framework.ErrorWithoutGrpcStuff(err)
}

func (_s *gokitServiceGrpcServer) Classificacoes(ctx context.Context, req *ClassificacoesRequest) (rep *ClassificacoesReply, err error) {
	_, _rep, err := _s.handlers["Classificacoes"].ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return _rep.(*ClassificacoesReply), err
}

func (_s *gokitServiceGrpcServer) Localizadores(ctx context.Context, req *LocalizadoresRequest) (rep *LocalizadoresReply, err error) {
	_, _rep, err := _s.handlers["Localizadores"].ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return _rep.(*LocalizadoresReply), err
}

func (_s *gokitServiceGrpcServer) Decodificar(ctx context.Context, req *DecodificarRequest) (rep *DecodificarReply, err error) {
	_, _rep, err := _s.handlers["Decodificar"].ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return _rep.(*DecodificarReply), err
}
