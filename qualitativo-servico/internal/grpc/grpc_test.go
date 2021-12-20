package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"git.sof.intra/pkg/dbtesting"
	"git.sof.intra/siop/framework"
	orgaos "git.sof.intra/siop/orgaos-servico"
	qualitativo "git.sof.intra/siop/qualitativo-servico"
	_ "git.sof.intra/siop/qualitativo-servico/client"
	"git.sof.intra/siop/qualitativo-servico/internal"
	"github.com/golang/mock/gomock"

	mock_orgaos "git.sof.intra/siop/orgaos-servico/mock"
	siop_proto "git.sof.intra/siop/siop-proto"
)

func TestClassificacoes(t *testing.T) {
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
	mockOrgaos := mockOrgaos(t)

	type C = siop_proto.Classificador

	casos := []struct {
		nome    string
		script  string
		saida   []*qualitativo.Classificacao
		wantErr bool
	}{
		{
			"um localizador cadastrado, deve retorná-lo",
			insertIntoAcao + `
				INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid,snatual,snexclusaologica)
			 	VALUES(1,'0001',1,'l',1,1,true,false);`,
			[]*qualitativo.Classificacao{{
				Esfera:                  &C{Id: 1, Codigo: "01", Nome: "e"},
				Orgao:                   &C{Id: 1, Codigo: "00001", Nome: "o"},
				Funcao:                  &C{Id: 1, Codigo: "01", Nome: "f"},
				SubFuncao:               &C{Id: 1, Codigo: "001", Nome: "sf"},
				Programa:                &C{Id: 1, Codigo: "0001", Nome: "p"},
				Acao:                    &C{Id: 1, Codigo: "0001", Nome: "a"},
				Localizador:             &C{Id: 1, Codigo: "0001", Nome: "l"},
				ClassificacaoCodificada: "01.00001.01.001.0001.0001.0001",
			}},
			false,
		},
	}
	for _, caso := range casos {
		t.Run(caso.nome, grpcRun(caso.script, func(client qualitativo.Service, t *testing.T) {
			framework.AddMockServiceFactory(
				new(orgaos.Service),
				func() interface{} { return mockOrgaos },
			)
			valores, err := client.Classificacoes(
				context.Background(),
				2019,
				[]string{"01.00001.01.001.0001.0001.0001"},
			)
			if (err != nil) != caso.wantErr {
				t.Errorf("error = %v, wantErr %v", err, caso.wantErr)
				return
			}
			if diff := cmp.Diff(caso.saida, valores); diff != "" {
				t.Errorf("diff %v", diff)
			}
		}))
	}
}

func TestLocalizadores(t *testing.T) {
	casos := []struct {
		nome    string
		script  string
		saida   []*qualitativo.Localizador
		wantErr bool
	}{
		{
			"um localizador cadastrado, deve retorná-lo",
			`INSERT INTO acao(funcaoid,subfuncaoid,momentoid,tipoacaoid,
				programaid,tipoinclusaoacaoid,esferaid,orgaoid,exercicio,
				identificadorunico,acao,titulo)
			VALUES(1,2,3,4,5,6,7,8,2019,10,'ABCD','t');
			INSERT INTO localizador(acaoid,localizador,identificadorunico,
				descricao,momentoid,tipoinclusaolocalizadorid)
			 VALUES(1,2,3,'d',4,5);`,
			[]*qualitativo.Localizador{{
				ID:           1,
				EsferaID:     7,
				OrgaoID:      8,
				FuncaoID:     1,
				SubfuncaoID:  2,
				ProgramaID:   5,
				AcaoID:       1,
				Exercicio:    2019,
				MomentoId:    4,
				TipoInclusao: qualitativo.TipoInclusaoAcao_CREDITO_ESPECIAL,
			}},
			false,
		},
	}
	for _, caso := range casos {
		t.Run(caso.nome, grpcRun(caso.script, func(client qualitativo.Service, t *testing.T) {
			valores, err := client.Localizadores(
				context.Background(),
				&qualitativo.FiltroLocalizador{Ids: []int32{1}},
				nil,
				0,
			)
			if (err != nil) != caso.wantErr {
				t.Errorf("error = %v, wantErr %v", err, caso.wantErr)
				return
			}
			if !reflect.DeepEqual(valores, caso.saida) {
				t.Errorf(
					"O retorno do request deveria ser %v e foi %v\n",
					caso.saida,
					valores,
				)
			}
		}))
	}
}

func grpcRun(script string, tf func(qualitativo.Service, *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		dbtesting.Setup(script)
		svc := &internal.Service{}
		var client qualitativo.Service
		done := framework.SetupGrpcTest(svc, &client)
		tf(client, t)
		done <- true
	}
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

func TestDecodificar(t *testing.T) {
	casos := []struct {
		nome    string
		entrada []string
		saida   []qualitativo.Classificacao
		wantErr bool
	}{
		{
			"uma classificacao codificada, deve decodifica-la",
			[]string{"10.26298.12.368.5011.0509.0029"},
			[]qualitativo.Classificacao{{
				Esfera:      &siop_proto.Classificador{Codigo: "10"},
				Orgao:       &siop_proto.Classificador{Codigo: "26298"},
				Funcao:      &siop_proto.Classificador{Codigo: "12"},
				SubFuncao:   &siop_proto.Classificador{Codigo: "368"},
				Programa:    &siop_proto.Classificador{Codigo: "5011"},
				Acao:        &siop_proto.Classificador{Codigo: "0509"},
				Localizador: &siop_proto.Classificador{Codigo: "0029"},
			}},
			false,
		},
	}
	for _, caso := range casos {
		t.Run(caso.nome, grpcRun("", func(client qualitativo.Service, t *testing.T) {
			classificacoes, err := client.Decodificar(
				context.Background(),
				caso.entrada,
			)
			if (err != nil) != caso.wantErr {
				t.Errorf("error = %v, wantErr %v", err, caso.wantErr)
				return
			}
			if !reflect.DeepEqual(classificacoes, caso.saida) {
				t.Errorf(
					"O retorno do request deveria ser %v e foi %v\n",
					caso.saida,
					classificacoes,
				)
			}
		}))
	}
}
