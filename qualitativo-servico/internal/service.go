package internal

import (
	"context"
	"fmt"

	"git.sof.intra/siop/framework"
	orgaos "git.sof.intra/siop/orgaos-servico"
	"git.sof.intra/siop/qualitativo-servico"
)

type (
	Service      struct{}
	Unimportable interface {
		qualitativo.Service
		Esferas(ctx context.Context, filtro Esfera) (esferas []Esfera, err error)
	}
	//Esfera representa esfera da programação qualitativa
	Esfera struct {
		ID            *int32  `db:"esferaid"`
		Codigo        *string `db:"esfera"`
		Nome          *string `db:"descricao"`
		NomeAbreviado *string `db:"descricaoabreviada"`
		Ativo         *bool   `db:"snativo"`
	}
)

var (
	_ qualitativo.Service = &Service{}
	_ Unimportable        = &Service{}
)

func (b *Service) Esferas(
	ctx context.Context,
	filtro Esfera) (esferas []Esfera, err error) {
	return Esferas(filtro, framework.NewDML())
}

func (b *Service) Classificacoes(
	ctx context.Context,
	exercicio int32,
	classificacoesCodificadas []string,
) (classificacoes []*qualitativo.Classificacao, err error) {
	var orgaosServico orgaos.Service
	err = framework.Connect(&orgaosServico)
	if err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao serviço de órgãos: %v", err)
	}
	return Classificacoes(ctx, exercicio, classificacoesCodificadas, framework.NewDML(), orgaosServico)
}

func (b *Service) Localizadores(
	ctx context.Context,
	filtro *qualitativo.FiltroLocalizador,
	fetch *qualitativo.FetchLocalizador,
	n int32,
) (localizadores []*qualitativo.Localizador, err error) {
	var orgaosServico orgaos.Service
	err = framework.Connect(&orgaosServico)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao serviço de órgãos: %v", err)
	}
	return Localizadores(ctx, filtro, fetch, n, framework.NewDML(), orgaosServico)
}

func (b *Service) Decodificar(ctx context.Context, classificacoesCodificadas []string) (classificacoes []qualitativo.Classificacao, err error) {
	return Decodificar(classificacoesCodificadas...)
}
