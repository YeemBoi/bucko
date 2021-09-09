package bucko

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
)

type db struct {
	*bun.DB
	mu sync.Mutex
}

var DB db

func UseDB(db *bun.DB) {
	DB.DB = db
	DB.mu.Unlock()
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
