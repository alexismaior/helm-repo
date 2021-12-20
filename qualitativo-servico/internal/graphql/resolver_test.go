package graphql

import (
	"testing"

	"git.sof.intra/siop/framework"
	"git.sof.intra/siop/qualitativo-servico/internal"

	"git.sof.intra/pkg/dbtesting"
	orgaos "git.sof.intra/siop/orgaos-servico"
	mock_orgaos "git.sof.intra/siop/orgaos-servico/mock"
	"github.com/golang/mock/gomock"
)

func TestHandler(t *testing.T) {
	type args struct {
		request, response string
	}
	tests := []struct {
		name                string
		script              string
		usoEsperadoDeOrgaos func(*mock_orgaos.MockService)
		args                []args
	}{
		{
			name: "Buscar uma esfera por ",
			script: `insert into esfera (esferaid, esfera, descricao, descricaoabreviada, snativo) 
			values (1, '10', 'teste','teste', true);`,
			args: []args{
				{
					`{esferas(filtro: {codigo: "10"}){codigo, nome}}`,
					`{"esferas":[{"codigo": "10", "nome": "teste"}]}`,
				},
			},
		},
		{
			"um localizador cadastrado e fetch de classificação especificado",
			`INSERT INTO esfera(esfera,descricao,descricaoabreviada,snativo)
				VALUES('10','teste','',1);
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
			[]args{
				{
					`{localizadores(filtro: {exercicio: 2019, classificacoesCodificadas: ["10.22222.33.444.5555.6666.7777"]})
				{
					classificacao {
						esfera {
							codigo
							nome
						}
						orgao {
							codigo
						}
						funcao {
							codigo
						}
						subFuncao {
							codigo
						}
						programa {
							codigo
						}
						acao {
							codigo
						}
						localizador {
							codigo
						}
					}
					classificacaoCodificada
				}}`,
					`{"data":{"localizadores":[{"classificacao":{"esfera":{"codigo":"10","nome":"teste"},"orgao":{"codigo":"22222"},"funcao":{"codigo":"33"},"subFuncao":{"codigo":"444"},"programa":{"codigo":"5555"},"acao":{"codigo":"6666"},"localizador":{"codigo":"7777"}},"classificacaoCodificada":"10.22222.33.444.5555.6666.7777"}]}}`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockOrgaos := mock_orgaos.NewMockService(mockCtrl)
			if tt.usoEsperadoDeOrgaos != nil {
				tt.usoEsperadoDeOrgaos(mockOrgaos)
			}
			framework.AddMockServiceFactory(
				new(orgaos.Service),
				func() interface{} { return mockOrgaos },
			)
			dbtesting.Setup(tt.script)
			mux := framework.SetupGraphQLTest(NewResolver(&internal.Service{}))
			for _, args := range tt.args {
				framework.TestGraphQL(t, args.request, args.response, mux)
			}
		})
	}
}
