package bucko

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/gobuffalo/nulls"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type ReqCtx struct {
	echo.Context
	Ctx    context.Context
	UserId nulls.UInt32
}

// NewCtx creates a new ReqCtx instance from an echo and context Context object.
// Uses context.Background if ctx is nil.
func NewCtx(c echo.Context, ctx context.Context) *ReqCtx {
	rc := ReqCtx{
		Context: c,
		Ctx:     ctx,
	}
	if ctx == nil {
		rc.Ctx = context.Background()
	}

	if err := rc.SetUser(); err != nil {
		rc.UserId.Valid = false
		fmt.Println(err)
	}

	return &rc
}

func (rc *ReqCtx) SetUser() (err error) {
	token, ok := rc.Context.Get(jwtConfig.ContextKey).(*jwt.Token)
	if !ok || !token.Valid {
		return errors.New("could not get token")
	}
	user, err := jwtConfig.userGetter(token.Claims)
	if err != nil {
		return err
	}
	if err = rc.Query(user).SetPK(); err != nil {
		return err
	}
	rc.UserId = nulls.NewUInt32(uint32(GetPK(user)))
	return
}

func (rc *ReqCtx) Query(m BaseFieldModel) *CtxQuery {
	fmt.Println(m)
	cq := CtxQuery{
		R:          rc,
		M:          m,
		JoinPrefix: "",
		TableAlias: DB.Table(reflect.TypeOf(m)).SQLAlias,
		Q:          DB.NewSelect().Model(m),
	}
	return &cq
}

// AuthorInsert sets the user_id column to the request's user id on insert statements.
func (rc *ReqCtx) AuthorInsert(q *bun.InsertQuery) *bun.InsertQuery {
	return q.Value("user_id", "?", rc.UserId)
}

func (rc *ReqCtx) delete(m BaseFieldModel) (err error) {
	_, err = DB.NewDelete().Model(m).WherePK().Exec(rc.Ctx)
	return
}

func (rc *ReqCtx) search(q *bun.SelectQuery, searchCols bun.Safe) *bun.SelectQuery {
	searchString := rc.Context.QueryParam("search")
	if len(searchString) > 2 {
		return q.Where("MATCH(?) AGAINST (?)", searchCols, searchString)
	}
	return q
}
