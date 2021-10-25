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
	mValue := reflect.Indirect(reflect.ValueOf(m))

	for _, pk := range table.PKs {
		pkv := pk.Value(mValue)
		switch pkv.Kind() {
		case reflect.Uintptr:
			return reflect.Indirect(pkv).Uint()
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			return pkv.Uint()
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			return uint64(pkv.Int())
		}
	}
	return 0
}

func GetPKCol(m interface{}) bun.Ident {
	table := DB.Table(reflect.TypeOf(m))
	mValue := reflect.Indirect(reflect.ValueOf(m))
	acceptableKinds := []reflect.Kind{
		reflect.Uintptr, reflect.Uint, reflect.Int,
		reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8,
		reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8,
	}
	for _, pk := range table.PKs {
		pkKind := pk.Value(mValue).Kind()
		for _, kind := range acceptableKinds {
			if pkKind == kind {
				return bun.Ident(string(pk.Name))
			}
		}
	}
	return "id"
}
