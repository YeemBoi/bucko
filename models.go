package bucko

import (
	"reflect"

	"github.com/uptrace/bun"
)

type (
	// BaseFieldModel provides an interface to models commonly referenced to help provide shortcuts.
	BaseFieldModel interface {

		// GetSelectQuery must return a SelectQuery (excluding relations) suited for selecting the model based off of cq.Q.
		GetSelectQuery(cq *CtxQuery) *bun.SelectQuery

		// Insert must insert itself and the result or an error.
		Insert(rc *ReqCtx) (BaseFieldModel, error)

		// GetColumn must return the column name of the main field used to retrieve the model (not PK).
		GetColumn() bun.Ident

		// GetField must return the field used to retrieve the model (not PK).
		GetField() interface{}

		// GetParam must return the URL param name of the main field used to retrieve the model (not PK).
		GetParam() string
	}
)

func GetPK(m interface{}) uint64 {
	table := DB.Table(reflect.TypeOf(m))

	for _, pk := range table.PKs {
		pkv := pk.Value(reflect.ValueOf(m))
		if pkv.Kind() == reflect.Uint {
			return pkv.Uint()
		} else if pkv.Kind() == reflect.Int {
			return uint64(pkv.Int())
		}
	}
	return 0
}

func GetPKCol(m interface{}) bun.Ident {
	table := DB.Table(reflect.TypeOf(m))

	for _, pk := range table.PKs {
		pkKind := pk.Value(reflect.ValueOf(m)).Kind()
		if pkKind == reflect.Uint || pkKind == reflect.Int {
			return bun.Ident(string(pk.Name))
		}
	}
	return "id"
}
