package qualitativo

import (
	"context"
	"fmt"

	siop_proto "git.sof.intra/siop/siop-proto"
)

type Service interface {
	Classificacoes(ctx context.Context, exercicio int32, classificacoesCodificadas []string) (classificacoes []*Classificacao, err error)
	Localizadores(ctx context.Context, filtro *FiltroLocalizador, fetch *FetchLocalizador, n int32) (localizadores []*Localizador, err error)
	Decodificar(ctx context.Context, classificacoesCodificadas []string) (classificacoes []Classificacao, err error)
}

type Classificacao struct {
	Esfera                  *siop_proto.Classificador
	Orgao                   *siop_proto.Classificador
	Funcao                  *siop_proto.Classificador
	SubFuncao               *siop_proto.Classificador
	Programa                *siop_proto.Classificador
	Acao                    *siop_proto.Classificador
	Localizador             *siop_proto.Classificador
	ClassificacaoCodificada string
}

type FiltroClassificacao struct {
	Exercicio                 int32
	ClassificacoesCodificadas []string
	Classificacoes            []*Classificacao
	ClassificacoesParciais    []*Classificacao
}

type Localizador struct {
	ID                      int32
	EsferaID                int32
	OrgaoID                 int32
	FuncaoID                int32
	SubfuncaoID             int32
	ProgramaID              int32
	AcaoID                  int32
	ClassificacaoCodificada string
	Exercicio               int32
	MomentoId               int32
	Classificacao           *Classificacao
	TipoInclusao            TipoInclusaoAcao
	MunicipioId             int32
	UfId                    int32
	RegiaoId                int32
	Codigo                  string
	Nacional                bool
	Descricao               string
}

type TipoInclusaoAcao int32

const (
	TipoInclusaoAcao_A_DEFINIR              TipoInclusaoAcao = 0
	TipoInclusaoAcao_PLOA                   TipoInclusaoAcao = 1
	TipoInclusaoAcao_EMENDA                 TipoInclusaoAcao = 2
	TipoInclusaoAcao_CREDITO_ADICIONAL      TipoInclusaoAcao = 3
	TipoInclusaoAcao_PPA                    TipoInclusaoAcao = 4
	TipoInclusaoAcao_CREDITO_ESPECIAL       TipoInclusaoAcao = 5
	TipoInclusaoAcao_CREDITO_EXTRAORDINARIO TipoInclusaoAcao = 6
)

type FiltroLocalizador struct {
	Ids                       []int32
	ClassificacoesCodificadas []string
	Exercicio                 int32
	OrgaosIds                 []int32
	OrgaosIdsOuIdsPais        []int32
	Classificacoes            []*Classificacao
	ClassificacoesParciais    []*Classificacao
}
type FetchLocalizador struct {
	ClassificacaoCodificada bool
	Classificacao           bool
}

// Codificada retorna a classificação codificada
func (c *Classificacao) Codificada() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf(
		"%v.%v.%v.%v.%v.%v.%v",
		codigo(c.Esfera),
		codigo(c.Orgao),
		codigo(c.Funcao),
		codigo(c.SubFuncao),
		codigo(c.Programa),
		codigo(c.Acao),
		codigo(c.Localizador),
	)
}

func (c *Classificacao) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%v", *c)
}

func codigo(c *siop_proto.Classificador) string {
	if c == nil {
		return ""
	}
	return c.Codigo
}
