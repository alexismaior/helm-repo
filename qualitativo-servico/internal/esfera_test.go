package internal

import (
	"fmt"
	"reflect"
	"testing"

	"git.sof.intra/pkg/dbtesting"
	"github.com/aleksvp2017/wraptypes"
	"github.com/mcesar/dbrx"
)

const sqlInsertEsfera = `insert into esfera 
	(esferaid, esfera, descricao, descricaoabreviada, snativo) 
	values (%v, %q, %q, %q, %v);`

func TestEsferas(t *testing.T) {
	type args struct {
		filtro Esfera
		dml    dbrx.DML
	}

	filtro := Esfera{
		Codigo: wraptypes.WrapString("10"),
	}

	tests := []struct {
		name    string
		script  string
		args    args
		want    []Esfera
		wantErr bool
	}{
		{
			name:   "Busca esfera por código",
			script: fmt.Sprintf(sqlInsertEsfera, 1, "10", "Esfera1", "Esfera1", true),
			args: args{
				filtro: filtro,
			},
			want: []Esfera{{
				ID:            wraptypes.WrapInt32(1),
				Codigo:        wraptypes.WrapString("10"),
				Nome:          wraptypes.WrapString("Esfera1"),
				NomeAbreviado: wraptypes.WrapString("Esfera1"),
				Ativo:         wraptypes.WrapBool(true),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := dbtesting.Setup(tt.script)
			wrap := dbrx.Wrap(session)
			got, err := Esferas(tt.args.filtro, wrap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Esferas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Esferas() = %v, want %v", got, tt.want)
			}

			if tt.want != nil {
				itens, err := consultaItens(wrap)
				if err != nil {
					t.Errorf("Consulta() = %v", err)
				}
				if !reflect.DeepEqual(itens, tt.want) {
					t.Errorf("Consulta() = %v, want %v", itens, tt.want)
				}
			}
		})
	}
}

func consultaItens(dml dbrx.DML) ([]Esfera, error) {
	var itens []Esfera

	_, err := dml.Select("*").
		From("esfera").
		Load(&itens)

	if itens == nil {
		itens = []Esfera{}
	}
	return itens, err
}
