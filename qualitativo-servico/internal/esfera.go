package internal

import (
	"github.com/gocraft/dbr"
	"github.com/mcesar/dbrx"
)

//Esferas retorna Esferas que atendam ao filtro
func Esferas(filtro Esfera, dml dbrx.DML) ([]Esfera, error) {
	var esferas []Esfera
	smt := dml.Select("*").From("esfera")
	if filtro.Codigo != nil {
		smt.Where(dbr.Like("esfera", *filtro.Codigo+"%"))
	}
	if filtro.ID != nil {
		smt.Where(dbr.Eq("esferaid", *filtro.ID))
	}
	if filtro.Nome != nil {
		smt.Where(dbr.Like("descricao", *filtro.Nome+"%"))
	}
	if filtro.NomeAbreviado != nil {
		smt.Where(dbr.Like("descricaoabreviada", *filtro.NomeAbreviado+"%"))
	}
	// TODO: acrescentar no filtro um campo que permita escolher
	// entre 'mostrar todos', 'somente ativos' e 'somente inativos'.
	// Considerar como default o 'somente ativos'.
	if filtro.Ativo != nil {
		smt.Where(dbr.Eq("snativo", *filtro.Ativo))
	}

	smt.OrderBy("esfera")

	_, err := smt.Load(&esferas)
	return esferas, err
}
