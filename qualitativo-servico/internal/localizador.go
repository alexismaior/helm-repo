package internal

import (
	"context"
	"fmt"
	"sort"

	orgaos "git.sof.intra/siop/orgaos-servico"
	"git.sof.intra/siop/qualitativo-servico"
	siop_proto "git.sof.intra/siop/siop-proto"

	"github.com/gocraft/dbr"

	"github.com/mcesar/copier"
	"github.com/mcesar/dbrx"
)

// Localizadores retorna os localizadores de acordo com o filtro especificado
func Localizadores(
	ctx context.Context,
	filtro *qualitativo.FiltroLocalizador,
	fetch *qualitativo.FetchLocalizador,
	n int32,
	dml dbrx.DML,
	orgaosServico orgaos.Service,
) (localizadores []*qualitativo.Localizador, err error) {
	type localizador struct {
		qualitativo.Localizador
		Esfera, Funcao, Subfuncao, Programa, Acao                     string
		NomeEsfera, NomeFuncao, NomeSubfuncao, NomePrograma, NomeAcao string
	}
	cols := []string{
		"localizadorid id",
		"localizador codigo",
		"a.exercicio",
		"a.esferaid esfera_id",
		"a.orgaoid orgao_id",
		"a.funcaoid funcao_id",
		"a.subfuncaoid subfuncao_id",
		"a.programaid programa_id",
		"a.acaoid acao_id",
		"l.momentoid momento_id",
		"l.tipoinclusaolocalizadorid tipo_inclusao",
		"coalesce(l.municipioid, 0) municipio_id",
		"coalesce(l.ufid, 0) uf_id",
		"coalesce(l.regiaoid, 0) regiao_id",
		"(localizador = '0001') nacional",
	}
	if fetch != nil && (fetch.ClassificacaoCodificada || fetch.Classificacao) {
		cols = append(
			cols,
			"e.esfera",
			"f.funcao",
			"sf.subfuncao",
			"p.programa",
			"a.acao",
		)
		if fetch.Classificacao {
			cols = append(
				cols,
				"l.descricao",
				"e.descricao nome_esfera",
				"f.descricao nome_funcao",
				"sf.descricao nome_subfuncao",
				"p.titulo nome_programa",
				"a.titulo nome_acao",
			)
		}
	} else {
		cols = append(
			cols,
			"'' descricao",
			"'' esfera",
			"'' funcao",
			"'' subfuncao",
			"'' programa",
			"'' acao",
			"'' nome_esfera",
			"'' nome_funcao",
			"'' nome_subfuncao",
			"'' nome_programa",
			"'' nome_acao",
			"'' codigo",
		)
	}
	stmt := dml.Select(cols...).
		From("localizador l").
		Join(dbr.I("acao").As("a"), "l.acaoid = a.acaoid")

	if (filtro != nil &&
		(len(filtro.Classificacoes) > 0 ||
			len(filtro.ClassificacoesCodificadas) > 0 ||
			len(filtro.ClassificacoesParciais) > 0)) ||
		(fetch != nil && (fetch.ClassificacaoCodificada || fetch.Classificacao)) {
		stmt.Join(dbr.I("esfera").As("e"), "a.esferaid=e.esferaid")
		stmt.Join(dbr.I("funcao").As("f"), "a.funcaoid=f.funcaoid")
		stmt.Join(dbr.I("subfuncao").As("sf"), "a.subfuncaoid=sf.subfuncaoid")
		stmt.Join(dbr.I("programa").As("p"), "a.programaid=p.programaid")
	}

	if filtro != nil {
		if len(filtro.Ids) > 0 {
			stmt.Where("localizadorid in ?", filtro.Ids)
		} else {
			if filtro.Exercicio > 0 {
				stmt.Where("a.exercicio = ?", filtro.Exercicio)
			}
			stmt.Where("not l.snexclusaologica")
			stmt.Where("l.snatual")
		}
		if len(filtro.OrgaosIds) > 0 {
			stmt.Where("orgaoid in ?", filtro.OrgaosIds)
		}
	}

	if filtro != nil && len(filtro.OrgaosIdsOuIdsPais) > 0 {
		orgaos, err := orgaosServico.Buscar(ctx, orgaos.Filtro{
			IdsOuIdsPai: filtro.OrgaosIdsOuIdsPais,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar orgaos: %v", err)
		}
		var idsOrgaos []int32
		for _, o := range orgaos {
			idsOrgaos = append(idsOrgaos, o.Id)
		}
		stmt.Where("a.orgaoid in ?", idsOrgaos)
	}

	var (
		classificacoes   []qualitativo.Classificacao
		orgaoIDPorCodigo map[string]int32
		orgaoPorID       map[int32]orgaos.Orgao
	)
	if filtro != nil && len(filtro.ClassificacoesCodificadas) > 0 {
		classificacoes, err = Decodificar(
			deduplicateStrings(filtro.ClassificacoesCodificadas)...,
		)
		if err != nil {
			return nil, err
		}
	}
	if filtro != nil && len(filtro.Classificacoes) > 0 {
		copier.Copy(&classificacoes, filtro.Classificacoes)
	}
	if len(classificacoes) > 0 {
		var codigos []string
		for _, c := range classificacoes {
			codigos = append(codigos, c.Orgao.Codigo)
		}
		os, err := orgaosServico.Buscar(ctx, orgaos.Filtro{
			Exercicio: filtro.Exercicio,
			Codigos:   codigos,
		})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar orgaos: %v", err)
		}
		orgaoIDPorCodigo = make(map[string]int32)
		orgaoPorID = make(map[int32]orgaos.Orgao)
		for _, o := range os {
			orgaoIDPorCodigo[o.Codigo] = o.Id
			orgaoPorID[o.Id] = o
		}
		var class []dbr.Builder
		for _, c := range classificacoes {
			class = append(class, dbr.And(
				dbr.Eq("e.esfera", c.Esfera.Codigo),
				dbr.Eq("a.orgaoid", orgaoIDPorCodigo[c.Orgao.Codigo]),
				dbr.Eq("f.funcao", c.Funcao.Codigo),
				dbr.Eq("sf.subfuncao", c.SubFuncao.Codigo),
				dbr.Eq("p.programa", c.Programa.Codigo),
				dbr.Eq("a.acao", c.Acao.Codigo),
				dbr.Eq("l.localizador", c.Localizador.Codigo),
			))
		}
		stmt.Where(dbr.Or(class...))
	}

	if filtro != nil && len(filtro.ClassificacoesParciais) > 0 {
		err = adicionaCondicoesDasClassificacoesParciais(
			ctx,
			filtro.Exercicio,
			filtro.ClassificacoesParciais,
			stmt,
			orgaosServico,
		)
		if err != nil {
			return nil, err
		}
	}

	if n > 0 {
		stmt.Limit(uint64(n))
	}

	var ls []localizador
	_, err = stmt.Load(&ls)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar localizadores: %v", err)
	}

	if fetch != nil &&
		(fetch.ClassificacaoCodificada || fetch.Classificacao) &&
		orgaoPorID == nil {
		var ids []int32
		for _, l := range ls {
			ids = append(ids, l.OrgaoID)
		}
		os, err := orgaosServico.Buscar(ctx, orgaos.Filtro{Ids: ids})
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar orgaos: %v", err)
		}
		orgaoPorID = make(map[int32]orgaos.Orgao)
		for _, o := range os {
			orgaoPorID[o.Id] = o
		}
	}

	for _, l := range ls {
		l := l
		if fetch != nil && (fetch.ClassificacaoCodificada || fetch.Classificacao) {
			c := &qualitativo.Classificacao{
				Esfera: &siop_proto.Classificador{
					Id:     l.EsferaID,
					Codigo: l.Esfera,
					Nome:   l.NomeEsfera,
				},
				Orgao: &siop_proto.Classificador{
					Id:     l.OrgaoID,
					Codigo: orgaoPorID[l.OrgaoID].Codigo,
					Nome:   orgaoPorID[l.OrgaoID].Nome,
				},
				Funcao: &siop_proto.Classificador{
					Id:     l.FuncaoID,
					Codigo: l.Funcao,
					Nome:   l.NomeFuncao,
				},
				SubFuncao: &siop_proto.Classificador{
					Id:     l.SubfuncaoID,
					Codigo: l.Subfuncao,
					Nome:   l.NomeSubfuncao,
				},
				Programa: &siop_proto.Classificador{
					Id:     l.ProgramaID,
					Codigo: l.Programa,
					Nome:   l.NomePrograma,
				},
				Acao: &siop_proto.Classificador{
					Id:     l.AcaoID,
					Codigo: l.Acao,
					Nome:   l.NomeAcao,
				},
				Localizador: &siop_proto.Classificador{
					Id:     l.ID,
					Codigo: l.Codigo,
					Nome:   l.Descricao,
				},
			}
			if fetch.Classificacao {
				l.Localizador.Classificacao = c
			}
			if fetch.ClassificacaoCodificada {
				l.Localizador.ClassificacaoCodificada =
					((*qualitativo.Classificacao)(c)).Codificada()
			}
		}
		localizadores = append(localizadores, &l.Localizador)
	}

	return localizadores, nil
}

func deduplicateStrings(in []string) []string {
	if in == nil {
		return nil
	}
	sort.Slice(in, func(i, j int) bool { return in[i] < in[j] })
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	return in[:j+1]
}
