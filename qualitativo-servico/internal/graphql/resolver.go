package graphql

import (
	"context"

	qualitativo "git.sof.intra/siop/qualitativo-servico"
	"git.sof.intra/siop/qualitativo-servico/internal"
	"github.com/mcesar/copier"
)

// Resolver Resolve o schema.graphql
type Resolver struct {
	qualitativo internal.Unimportable
}

type FiltroLocalizador struct {
	Exercicio                 int32
	ClassificacoesCodificadas *[]string
	ClassificacoesParciais    *[]*ClassificacaoInput
}

type ClassificacaoInput struct {
	Esfera      *ClassificadorInput
	Orgao       *ClassificadorInput
	Funcao      *ClassificadorInput
	SubFuncao   *ClassificadorInput
	Programa    *ClassificadorInput
	Acao        *ClassificadorInput
	Localizador *ClassificadorInput
}

type ClassificadorInput struct {
	Codigo, Nome *string
}

// NewResolver cria um resolver graphql
func NewResolver(svc interface{}) interface{} {
	return &Resolver{
		qualitativo: svc.(internal.Unimportable),
	}
}

//Esferas delega para a função Esferas do serviço qualitativo
func (r *Resolver) Esferas(ctx context.Context, args struct {
	Filtro *internal.Esfera
}) ([]internal.Esfera, error) {
	filtro := internal.Esfera{}
	if args.Filtro != nil {
		filtro = *args.Filtro
	}
	return r.qualitativo.Esferas(ctx, filtro)
}

func (r *Resolver) Classificacoes(
	ctx context.Context,
	args struct {
		Exercicio                 int32
		ClassificacoesCodificadas []string
	},
) (classificacoes []*qualitativo.Classificacao, err error) {
	return r.qualitativo.Classificacoes(ctx, args.Exercicio, args.ClassificacoesCodificadas)
}

func (r *Resolver) Localizadores(
	ctx context.Context,
	args struct {
		Filtro FiltroLocalizador
		N      *int32
	},
) (localizadores []*qualitativo.Localizador, err error) {
	if args.N == nil {
		var infinito int32 = 0
		args.N = &infinito
	}
	return r.qualitativo.Localizadores(
		ctx,
		&qualitativo.FiltroLocalizador{
			Exercicio: args.Filtro.Exercicio,
			ClassificacoesCodificadas: func() []string {
				if args.Filtro.ClassificacoesCodificadas != nil {
					return *args.Filtro.ClassificacoesCodificadas
				}
				return nil
			}(),
			ClassificacoesParciais: func() []*qualitativo.Classificacao {
				if args.Filtro.ClassificacoesParciais != nil {
					return copier.CopyAndDereference(
						&[]*qualitativo.Classificacao{}, args.Filtro.ClassificacoesParciais,
					).([]*qualitativo.Classificacao)
				}
				return nil
			}(),
		},
		&qualitativo.FetchLocalizador{ClassificacaoCodificada: true, Classificacao: true},
		*args.N,
	)
}
