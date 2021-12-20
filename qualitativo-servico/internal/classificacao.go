package internal

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	orgaos "git.sof.intra/siop/orgaos-servico"
	qualitativo "git.sof.intra/siop/qualitativo-servico"
	siop_proto "git.sof.intra/siop/siop-proto"
	"github.com/gocraft/dbr"
	"github.com/mcesar/copier"
	"github.com/mcesar/dbrx"
)

var regexpClassificacao = regexp.MustCompile(`(\d\d)\.(\d\d\d\d\d)\.(\d\d)\.(\d\d\d)\.(\d\d\d\d)\.(....)\.(\d\d\d\d)`)

func Classificacoes(
	ctx context.Context,
	exercicio int32,
	classificacoesCodificadas []string,
	dml dbrx.DML,
	orgaosServico orgaos.Service,
) (classificacoes []*qualitativo.Classificacao, err error) {

	if exercicio == 0 {
		return nil, fmt.Errorf("o exercício é obrigatório")
	}

	if len(classificacoesCodificadas) == 0 {
		return nil, nil
	}

	var (
		codigosEsferas,
		codigosOrgaos,
		codigosFuncoes,
		codigosSubfuncoes,
		codigosProgramas,
		codigosAcoes,
		codigosLocalizadores []string
	)
	for _, cc := range classificacoesCodificadas {
		c, err := decodificarClassificacao(cc)
		if err != nil {
			return nil, err
		}
		classificacoes = append(classificacoes, &c)
		codigosEsferas = append(codigosEsferas, c.Esfera.Codigo)
		codigosOrgaos = append(codigosOrgaos, c.Orgao.Codigo)
		codigosFuncoes = append(codigosFuncoes, c.Funcao.Codigo)
		codigosSubfuncoes = append(codigosSubfuncoes, c.SubFuncao.Codigo)
		codigosProgramas = append(codigosProgramas, c.Programa.Codigo)
		codigosAcoes = append(codigosAcoes, c.Acao.Codigo)
		codigosLocalizadores = append(codigosLocalizadores, c.Localizador.Codigo)
	}

	var (
		esferaPorCodigo      = make(map[string]Esfera)
		orgaoPorCodigo       = make(map[string]orgaos.Orgao)
		funcaoPorCodigo      = make(map[string]*siop_proto.Classificador)
		subfuncaoPorCodigo   = make(map[string]*siop_proto.Classificador)
		programaPorCodigo    = make(map[string]*siop_proto.Classificador)
		acaoPorCodigo        = make(map[string]*siop_proto.Classificador)
		localizadorPorCodigo = make(map[string]*siop_proto.Classificador)
	)
	es, err := Esferas(Esfera{}, dml)
	if err != nil {
		return nil, err
	}
	for _, e := range es {
		esferaPorCodigo[*e.Codigo] = e
	}
	os, err := orgaosServico.Buscar(ctx, orgaos.Filtro{Codigos: codigosOrgaos})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar órgãos: %v", err)
	}
	for _, o := range os {
		orgaoPorCodigo[o.Codigo] = o
	}
	_, err = dml.Select("funcao", "funcaoid id", "funcao codigo", "descricao nome").
		From("funcao").
		Where("exercicio = ?", exercicio).
		Where("funcao in ?", codigosFuncoes).
		Load(&funcaoPorCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar funções: %v", err)
	}
	_, err = dml.Select("subfuncao", "subfuncaoid id", "subfuncao codigo", "descricao nome").
		From("subfuncao").
		Where("exercicio = ?", exercicio).
		Where("subfuncao in ?", codigosSubfuncoes).
		Load(&subfuncaoPorCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar subfunções: %v", err)
	}
	_, err = dml.Select("programa", "programaid id", "programa codigo", "titulo nome").
		From("programa").
		Where("exercicio = ?", exercicio).
		Where("programa in ?", codigosProgramas).
		Load(&programaPorCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar programas: %v", err)
	}
	_, err = dml.Select("acao", "acaoid id", "acao codigo", "titulo nome").
		From("acao").
		Where("exercicio = ?", exercicio).
		Where("acao in ?", codigosAcoes).
		Load(&acaoPorCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ações: %v", err)
	}
	_, err = dml.Select("localizador", "localizadorid id", "localizador codigo", "l.descricao nome").
		From("localizador l").Join("acao", "l.acaoid = acao.acaoid").
		Where("acao.exercicio = ?", exercicio).
		Where("localizador in ?", codigosLocalizadores).
		Load(&localizadorPorCodigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar localizadores: %v", err)
	}

	for _, c := range classificacoes {
		toClassificador := func(c interface{}) *siop_proto.Classificador {
			return copier.Copy(&siop_proto.Classificador{}, c).(*siop_proto.Classificador)
		}
		c.Esfera = toClassificador(esferaPorCodigo[c.Esfera.Codigo])
		c.Orgao = toClassificador(orgaoPorCodigo[c.Orgao.Codigo])
		c.Funcao = funcaoPorCodigo[c.Funcao.Codigo]
		c.SubFuncao = subfuncaoPorCodigo[c.SubFuncao.Codigo]
		c.Programa = programaPorCodigo[c.Programa.Codigo]
		c.Acao = acaoPorCodigo[c.Acao.Codigo]
		c.Localizador = localizadorPorCodigo[c.Localizador.Codigo]
		c.ClassificacaoCodificada = c.Codificada()
	}

	return classificacoes, nil
}

func Decodificar(
	classificacoesCodificadas ...string,
) (classificacoes []qualitativo.Classificacao, err error) {
	for _, cc := range classificacoesCodificadas {
		c, err := decodificarClassificacao(cc)
		if err != nil {
			return nil, err
		}
		classificacoes = append(classificacoes, c)
	}
	return classificacoes, err
}

func decodificarClassificacao(classificacaoCodificada string) (c qualitativo.Classificacao, err error) {
	submatches := regexpClassificacao.FindStringSubmatch(classificacaoCodificada)
	if len(submatches) < 7 {
		return c, errors.New("formato incorreto")
	}
	c = qualitativo.Classificacao{
		Esfera:      &siop_proto.Classificador{Codigo: submatches[1]},
		Orgao:       &siop_proto.Classificador{Codigo: submatches[2]},
		Funcao:      &siop_proto.Classificador{Codigo: submatches[3]},
		SubFuncao:   &siop_proto.Classificador{Codigo: submatches[4]},
		Programa:    &siop_proto.Classificador{Codigo: submatches[5]},
		Acao:        &siop_proto.Classificador{Codigo: submatches[6]},
		Localizador: &siop_proto.Classificador{Codigo: submatches[7]},
	}
	return c, nil
}

// adicionaCondicoesDasClassificacoesParciais adiciona ao Where do stmt as
// condições correspondentes aos "likes" dos códigos dos classificadores.
// É necessário que o stmt defina os seguintes aliases de tabelas:
// e: esfera
// f: funcao
// sf: subfuncao
// p: programa
// a: acao
// l: localizador
func adicionaCondicoesDasClassificacoesParciais(
	ctx context.Context,
	exercicio int32,
	classificacoesParciais []*qualitativo.Classificacao,
	stmt *dbrx.SelectStmt,
	orgaosServico orgaos.Service,
) error {
	var (
		codigosDosOrgaos []string
		condicoesEsferas,
		condicoesFuncoes,
		condicoesSubfuncoes,
		condicoesProgramas,
		condicoesAcoes,
		condicoesLocalizadores []dbr.Builder
	)
	for _, classificacaoParcial := range classificacoesParciais {
		if classificacaoParcial.Orgao != nil && classificacaoParcial.Orgao.Codigo != "" {
			codigosDosOrgaos = append(codigosDosOrgaos, classificacaoParcial.Orgao.Codigo)
		}
		if classificacaoParcial.Esfera != nil && classificacaoParcial.Esfera.Codigo != "" {
			condicoesEsferas = append(condicoesEsferas,
				dbr.Like("e.esfera", "%"+classificacaoParcial.Esfera.Codigo+"%"),
			)
		}
		if classificacaoParcial.Funcao != nil && classificacaoParcial.Funcao.Codigo != "" {
			condicoesFuncoes = append(condicoesFuncoes,
				dbr.Like("f.funcao", "%"+classificacaoParcial.Funcao.Codigo+"%"),
			)
		}
		if classificacaoParcial.SubFuncao != nil && classificacaoParcial.SubFuncao.Codigo != "" {
			condicoesSubfuncoes = append(condicoesSubfuncoes,
				dbr.Like("sf.subfuncao", "%"+classificacaoParcial.SubFuncao.Codigo+"%"),
			)
		}
		if classificacaoParcial.Programa != nil && classificacaoParcial.Programa.Codigo != "" {
			condicoesProgramas = append(condicoesProgramas,
				dbr.Like("p.programa", "%"+classificacaoParcial.Programa.Codigo+"%"),
			)
		}
		if classificacaoParcial.Acao != nil && classificacaoParcial.Acao.Codigo != "" {
			condicoesAcoes = append(condicoesAcoes,
				dbr.Like("a.acao", "%"+classificacaoParcial.Acao.Codigo+"%"),
			)
		}
		if classificacaoParcial.Localizador != nil && classificacaoParcial.Localizador.Codigo != "" {
			condicoesLocalizadores = append(condicoesLocalizadores,
				dbr.Like("l.localizador", "%"+classificacaoParcial.Localizador.Codigo+"%"),
			)
		}
	}
	if len(codigosDosOrgaos) > 0 {
		filtroOrgao := orgaos.Filtro{
			Exercicio: exercicio,
			Codigos:   codigosDosOrgaos,
		}

		orgaosFiltrados, err := orgaosServico.Buscar(ctx, filtroOrgao)
		if err != nil {
			return err
		}

		var idOrgaos []int32
		for _, orgao := range orgaosFiltrados {
			idOrgaos = append(idOrgaos, orgao.Id)
		}

		if len(idOrgaos) > 0 {
			stmt.Where("a.orgaoid in ?", idOrgaos)
		}
	}

	if len(condicoesEsferas) > 0 {
		stmt.Where(dbr.Or(condicoesEsferas...))
	}
	if len(condicoesFuncoes) > 0 {
		stmt.Where(dbr.Or(condicoesFuncoes...))
	}
	if len(condicoesSubfuncoes) > 0 {
		stmt.Where(dbr.Or(condicoesSubfuncoes...))
	}
	if len(condicoesProgramas) > 0 {
		stmt.Where(dbr.Or(condicoesProgramas...))
	}
	if len(condicoesAcoes) > 0 {
		stmt.Where(dbr.Or(condicoesAcoes...))
	}
	if len(condicoesLocalizadores) > 0 {
		stmt.Where(dbr.Or(condicoesLocalizadores...))
	}
	return nil
}
