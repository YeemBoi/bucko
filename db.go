package bucko

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

var DB *bun.DB

func UseDB(db *bun.DB) {
	DB = db
}

func BaseCheckExists(model interface{}, ctx context.Context, query string, args ...interface{}) (exists bool, err error) {
	err = DB.NewSelect().Model(model).ColumnExpr("1").Where(query, args...).Limit(1).Scan(ctx, &exists)
	if err == sql.ErrNoRows {
		exists = false
		err = nil
	}
	return
}

func BaseCheckPKExists(model interface{}, ctx context.Context) (exists bool, err error) {
	err = DB.NewSelect().Model(model).ColumnExpr("1").WherePK().Limit(1).Scan(ctx, &exists)
	if err == sql.ErrNoRows {
		exists = false
		err = nil
	}
	return
}
