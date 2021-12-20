package internal

import (
	"context"
	"fmt"
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

const insertIntoAcao = `
		INSERT INTO esfera VALUES (1,'01','e','e',true);
		INSERT INTO funcao VALUES (1,'01',2019,'f','f',true);
		INSERT INTO subfuncao VALUES (1,'001',1,2019,'sf','sf',true);
		INSERT INTO programa VALUES (1,'0001',1,2019,'p',true,false,true);

		INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
			programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
			identificadorunico,acao,titulo,snatual,snexclusaologica)
		VALUES(1,1,1,1,1,1,1,1,2019,1,'0001','a',true,false);
	`

func insertIntoLocalizador(codigo, descricao string) string {
	return fmt.Sprintf(`
		INSERT INTO localizador(acaoid,localizador,identificadorunico,
		   descricao,momentoid,tipoinclusaolocalizadorid,snatual,snexclusaologica)
		VALUES(1,'%v',1,'%v',1,1,true,false);`, codigo, descricao)
}

func mockOrgaos(t *testing.T) orgaos.Service {
	mockCtrl := gomock.NewController(t)
	mockOrgaos := mock_orgaos.NewMockService(mockCtrl)

	mockOrgaos.EXPECT().
		Buscar(
			gomock.Any(),
			gomock.Any(),
		).
		Return(
			[]orgaos.Orgao{{
				Id:     1,
				Codigo: "00001",
				Nome:   "o",
			}},
			nil,
		).
		AnyTimes()
	return mockOrgaos
}

func TestClassificacoes(t *testing.T) {

	mockOrgaos := mockOrgaos(t)

	type args struct {
		ctx                       context.Context
		exercicio                 int32
		classificacoesCodificadas []string
		dml                       dbrx.DML
		orgaosServico             orgaos.Service
	}
	tests := []struct {
		name               string
		script             string
		args               args
		wantClassificacoes []*qualitativo.Classificacao
		wantErr            bool
	}{
		{
			"Todas os classificadores cadastrados",
			insertIntoAcao + insertIntoLocalizador("0001", "l"),
			args{
				orgaosServico:             mockOrgaos,
				exercicio:                 2019,
				classificacoesCodificadas: []string{"01.00001.01.001.0001.0001.0001"},
			},
			[]*qualitativo.Classificacao{{
				Esfera:                  &siop_proto.Classificador{Id: 1, Codigo: "01", Nome: "e"},
				Orgao:                   &siop_proto.Classificador{Id: 1, Codigo: "00001", Nome: "o"},
				Funcao:                  &siop_proto.Classificador{Id: 1, Codigo: "01", Nome: "f"},
				SubFuncao:               &siop_proto.Classificador{Id: 1, Codigo: "001", Nome: "sf"},
				Programa:                &siop_proto.Classificador{Id: 1, Codigo: "0001", Nome: "p"},
				Acao:                    &siop_proto.Classificador{Id: 1, Codigo: "0001", Nome: "a"},
				Localizador:             &siop_proto.Classificador{Id: 1, Codigo: "0001", Nome: "l"},
				ClassificacaoCodificada: "01.00001.01.001.0001.0001.0001",
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess := dbtesting.Setup(tt.script)
			gotClassificacoes, err := Classificacoes(
				tt.args.ctx,
				tt.args.exercicio,
				tt.args.classificacoesCodificadas,
				dbrx.Wrap(sess),
				tt.args.orgaosServico,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Classificacoes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantClassificacoes, gotClassificacoes); diff != "" {
				t.Errorf("Classificacoes() diff: %v", diff)
			}
		})
	}
}
