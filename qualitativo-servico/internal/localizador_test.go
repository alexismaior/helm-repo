package internal

import (
	"context"
	"testing"

	"git.sof.intra/pkg/dbtesting"

	orgaos "git.sof.intra/siop/orgaos-servico"
	mock_orgaos "git.sof.intra/siop/orgaos-servico/mock"
	qualitativo "git.sof.intra/siop/qualitativo-servico"
	siop_proto "git.sof.intra/siop/siop-proto"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/mcesar/dbrx"
)

func TestLocalizadores(t *testing.T) {
	const insertIntoAcao = `
		INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
			programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
			identificadorunico,acao,titulo)
		VALUES(1,2,3,4,5,6,7,8,2019,10,'ABCD','t');`
	type args struct {
		filtro *qualitativo.FiltroLocalizador
		fetch  *qualitativo.FetchLocalizador
	}
	tests := []struct {
		name                string
		script              string
		usoEsperadoDeOrgaos func(*mock_orgaos.MockService)
		args                args
		wantLocalizadores   []*qualitativo.Localizador
		wantErr             bool
	}{
		{
			"nenhum localizador cadastrado, não deve retornar nada",
			"",
			nil,
			args{},
			nil,
			false,
		},
		{
			"um localizador cadastrado, deve retorná-lo",
			insertIntoAcao + `
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid)
			 VALUES(1,2,3,'d',4,5);`,
			nil,
			args{},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				EsferaID:     7,
				OrgaoID:      8,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  2,
				ProgramaID:   5,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
		{
			"dois localizadores cadastrados, deve retornar apenas o especificado",
			insertIntoAcao + `
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid)
			 VALUES(1,2,3,'d',4,5);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid)
			 VALUES(1,20,30,'dd',40,50);`,
			nil,
			args{filtro: &qualitativo.FiltroLocalizador{Ids: []int32{1}}},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				EsferaID:     7,
				OrgaoID:      8,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  2,
				ProgramaID:   5,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
		{
			"um localizador cadastrado e fetch de classificação codificada especificado",
			`INSERT INTO esfera(esfera,descricao,descricaoabreviada,snativo)
				VALUES('10','','',1);
			 INSERT INTO funcao(exercicio,funcao,descricao,descricaoabreviada,snativo)
				VALUES(2019,'33','','',1);
			 INSERT INTO subfuncao(exercicio,subfuncao,funcaoid,descricao,descricaoabreviada,snativo)
				VALUES(2019,'444',1,'','',1);
			 INSERT INTO programa(exercicio,programa,titulo,snexclusaologica,snatual)
			 	VALUES(2019,'5555','',0,1);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
					descricao,momentoid,tipoinclusaolocalizadorid)
				 VALUES(1,'7777',3,'d',4,5);
			 INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
					programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
					identificadorunico,acao,titulo)
				VALUES(1,1,1,1,1,1,1,1,2019,1,'6666','');`,
			func(mockOrgaos *mock_orgaos.MockService) {
				mockOrgaos.EXPECT().
					Buscar(gomock.Any(), orgaos.Filtro{Ids: []int32{1}}).
					Return([]orgaos.Orgao{{Id: 1, Codigo: "22222"}}, nil)
			},
			args{fetch: &qualitativo.FetchLocalizador{ClassificacaoCodificada: true}},
			[]*qualitativo.Localizador{{
				ID:                      1,
				Exercicio:               2019,
				Codigo:                  "7777",
				EsferaID:                1,
				OrgaoID:                 1,
				AcaoID:                  1,
				FuncaoID:                1,
				SubfuncaoID:             1,
				ProgramaID:              1,
				MomentoId:               4,
				TipoInclusao:            qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
				ClassificacaoCodificada: "10.22222.33.444.5555.6666.7777",
			}},
			false,
		},
		{
			"um localizador cadastrado e filtro por classificação codificada",
			`INSERT INTO esfera(esfera,descricao,descricaoabreviada,snativo)
				VALUES('10','','',1);
			 INSERT INTO funcao(exercicio,funcao,descricao,descricaoabreviada,snativo)
				VALUES(2019,'33','','',1);
			 INSERT INTO subfuncao(exercicio,subfuncao,funcaoid,descricao,descricaoabreviada,snativo)
				VALUES(2019,'444',1,'','',1);
			 INSERT INTO programa(exercicio,programa,titulo,snexclusaologica,snatual)
			 	VALUES(2019,'5555','',0,1);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
					descricao,momentoid,tipoinclusaolocalizadorid,snatual)
				 VALUES(1,'7777',3,'d',4,5,1);
			 INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
					programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
					identificadorunico,acao,titulo)
				VALUES(1,1,1,1,1,1,1,1,2019,1,'6666','');`,
			func(mockOrgaos *mock_orgaos.MockService) {
				mockOrgaos.EXPECT().
					Buscar(gomock.Any(), orgaos.Filtro{
						Exercicio: 2019,
						Codigos:   []string{"22222"},
					}).
					Return([]orgaos.Orgao{{Id: 1, Codigo: "22222"}}, nil)
			},
			args{filtro: &qualitativo.FiltroLocalizador{
				Exercicio: 2019,
				ClassificacoesCodificadas: []string{
					"10.22222.33.444.5555.6666.7777",
				},
			}},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				EsferaID:     1,
				OrgaoID:      1,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  1,
				ProgramaID:   1,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
		{
			"um localizador cadastrado e fetch de classificação especificado",
			`INSERT INTO esfera(esfera,descricao,descricaoabreviada,snativo)
				VALUES('10','e','',1);
			 INSERT INTO funcao(exercicio,funcao,descricao,descricaoabreviada,snativo)
				VALUES(2019,'33','f','',1);
			 INSERT INTO subfuncao(exercicio,subfuncao,funcaoid,descricao,descricaoabreviada,snativo)
				VALUES(2019,'444',1,'sf','',1);
			 INSERT INTO programa(exercicio,programa,titulo,snexclusaologica,snatual)
			 	VALUES(2019,'5555','p',0,1);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
					descricao,momentoid,tipoinclusaolocalizadorid)
				 VALUES(1,'7777',3,'l',4,5);
			 INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
					programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
					identificadorunico,acao,titulo)
				VALUES(1,1,1,1,1,1,1,1,2019,1,'6666','a');`,
			func(mockOrgaos *mock_orgaos.MockService) {
				mockOrgaos.EXPECT().
					Buscar(gomock.Any(), orgaos.Filtro{Ids: []int32{1}}).
					Return([]orgaos.Orgao{{Id: 1, Codigo: "22222", Nome: "o"}}, nil)
			},
			args{fetch: &qualitativo.FetchLocalizador{Classificacao: true}},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				Codigo:       "7777",
				EsferaID:     1,
				OrgaoID:      1,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  1,
				ProgramaID:   1,
				MomentoId:    4,
				Descricao:    "l",
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
				Classificacao: &qualitativo.Classificacao{
					Esfera: &siop_proto.Classificador{
						Id:     1,
						Codigo: "10",
						Nome:   "e",
					},
					Orgao: &siop_proto.Classificador{
						Id:     1,
						Codigo: "22222",
						Nome:   "o",
					},
					Funcao: &siop_proto.Classificador{
						Id:     1,
						Codigo: "33",
						Nome:   "f",
					},
					SubFuncao: &siop_proto.Classificador{
						Id:     1,
						Codigo: "444",
						Nome:   "sf",
					},
					Programa: &siop_proto.Classificador{
						Id:     1,
						Codigo: "5555",
						Nome:   "p",
					},
					Acao: &siop_proto.Classificador{
						Id:     1,
						Codigo: "6666",
						Nome:   "a",
					},
					Localizador: &siop_proto.Classificador{
						Id:     1,
						Codigo: "7777",
						Nome:   "l",
					},
				},
			}},
			false,
		},
		{
			`dois localizadores de órgãos diferentes cadastrados,
             deve retornar apenas o do órgão especificado`,
			`INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
                programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
                identificadorunico,acao,titulo)
             VALUES(1,2,3,4,5,6,7,8,2019,10,'ABCD','t');
			 INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
                programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
                identificadorunico,acao,titulo)
             VALUES(10,20,30,40,50,60,70,80,2019,100,'ABCD','t');
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid,snatual)
			 VALUES(1,2,3,'d',4,5,true);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid,snatual)
			 VALUES(2,20,30,'dd',40,50,true);`,
			nil,
			args{filtro: &qualitativo.FiltroLocalizador{OrgaosIds: []int32{8}}},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				EsferaID:     7,
				OrgaoID:      8,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  2,
				ProgramaID:   5,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
		{
			`dois localizadores de órgãos diferentes cadastrados,
             deve retornar apenas o do órgão ou órgão pai especificado`,
			`INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
                programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
                identificadorunico,acao,titulo)
             VALUES(1,2,3,4,5,6,7,8,2019,10,'ABCD','t');
			 INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
                programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
                identificadorunico,acao,titulo)
             VALUES(10,20,30,40,50,60,70,80,2019,100,'ABCD','t');
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid,snatual)
			 VALUES(1,2,3,'d',4,5,true);
			 INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid,snatual)
			 VALUES(2,20,30,'dd',40,50,true);`,
			func(mockOrgaos *mock_orgaos.MockService) {
				mockOrgaos.EXPECT().
					Buscar(gomock.Any(), orgaos.Filtro{IdsOuIdsPai: []int32{8}}).
					Return([]orgaos.Orgao{{Id: 8, Codigo: "22222"}}, nil)
			},
			args{filtro: &qualitativo.FiltroLocalizador{OrgaosIdsOuIdsPais: []int32{8}}},
			[]*qualitativo.Localizador{{
				ID:           1,
				Exercicio:    2019,
				EsferaID:     7,
				OrgaoID:      8,
				AcaoID:       1,
				FuncaoID:     1,
				SubfuncaoID:  2,
				ProgramaID:   5,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess := dbtesting.Setup(tt.script)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockOrgaos := mock_orgaos.NewMockService(mockCtrl)
			if tt.usoEsperadoDeOrgaos != nil {
				tt.usoEsperadoDeOrgaos(mockOrgaos)
			}

			gotLocalizadores, err := Localizadores(
				context.Background(),
				tt.args.filtro,
				tt.args.fetch,
				0,
				dbrx.Wrap(sess),
				mockOrgaos,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Localizadores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotLocalizadores, tt.wantLocalizadores); diff != "" {
				t.Errorf("Localizadores() = -got +want %v", diff)
			}
		})
	}
}
