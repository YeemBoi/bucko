package bucko

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/uptrace/bun"
)

type CtxQuery struct {
	*bun.SelectQuery
	R          *ReqCtx
	M          BaseFieldModel
	JoinPrefix string
	TableAlias bun.Safe
}

// Q updates the query.
func (cq *CtxQuery) Q(q *bun.SelectQuery) {
	cq.SelectQuery = q
}

// SafeCol returns a query's column, safe to use in joins.
func (cq *CtxQuery) SafeCol(column bun.Ident) bun.Safe {
	return bun.Safe(fmt.Sprintf("%s.`%s`", cq.TableAlias, column))
}

// ColAlias returns a query's column's alias for selection to bun, safe to use in joins.
func (cq *CtxQuery) ColAlias(column bun.Ident) bun.Ident {
	return bun.Ident(cq.JoinPrefix) + column
}

func (cq *CtxQuery) CheckFieldExists() (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "? = ?", cq.SafeCol(cq.M.GetColumn()), cq.M.GetField())
}

func (cq *CtxQuery) CheckPKExists() (exists bool, err error) {
	return BaseCheckPKExists(cq.M, cq.R.Ctx)
}

func (cq *CtxQuery) CheckExistsFromField(field interface{}) (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "? = ?", cq.SafeCol(cq.M.GetColumn()), field)
}

func (cq *CtxQuery) CheckExistsFromPK(id uint64) (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "? = ?", cq.SafeCol(GetPKCol(cq.M)), id)
}

func (cq *CtxQuery) DeleteFromPK() (err error) {
	return cq.R.delete(cq.M)
}

// SetLimitOffset sets a queries limit & offset from given URL parameters.
func (cq *CtxQuery) SetLimitOffset() {
	var limit, offset uint64
	limit, err := strconv.ParseUint(cq.R.Context.QueryParam("limit"), 10, 8)
	if err != nil {
		cq.Q(cq.Limit(25))
	} else {
		if limit > 100 {
			cq.Q(cq.Limit(100))
		} else {
			cq.Q(cq.Limit(int(limit)))
		}
	}
	if offset, err = strconv.ParseUint(cq.R.Context.QueryParam("offset"), 10, 64); err == nil {
		cq.Q(cq.Offset(int(offset)))
	}
}

func (cq *CtxQuery) ApplySearch(searchCols bun.Safe) {
	cq.Q(cq.R.Search(cq.SelectQuery, searchCols))
}

func (cq *CtxQuery) WhereParamToCol() {
	cq.Q(cq.Where("? = ?", cq.M.GetColumn(), cq.R.Context.Param(cq.M.GetParam())).Limit(1))
}

func (cq *CtxQuery) FromParam() (err error) {
	cq.Select()
	return cq.Where("? = ?", cq.M.GetColumn(), cq.R.Context.Param(cq.M.GetParam())).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) FromField() (err error) {
	cq.Select()
	return cq.Where("? = ?", cq.M.GetColumn(), cq.M.GetField()).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) FromPK() (err error) {
	return cq.M.GetSelectQuery(cq).WherePK().Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromField(field interface{}) (err error) {
	return cq.M.GetSelectQuery(cq).Where("? = ?", cq.SafeCol(cq.M.GetColumn()), field).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromPK(id uint64) (err error) {
	return cq.M.GetSelectQuery(cq).Where("? = ?", cq.SafeCol(GetPKCol(cq.M)), id).Limit(1).Scan(cq.R.Ctx)
}

// Insert inserts CtxQuery's model based on its Insert() method.
func (cq *CtxQuery) Insert() (err error) {
	cq.M, err = cq.M.Insert(cq.R)
	return
}

func (cq *CtxQuery) SetPK() (err error) {
	return DB.NewSelect().Column(string(GetPKCol(cq.M))).Model(cq.M).
		Where("? = ?", cq.SafeCol(cq.M.GetColumn()), cq.M.GetField()).Limit(1).Scan(cq.R.Ctx)
}

// Select applies its model's GetSelectQuery and relations to itself.
func (cq *CtxQuery) Select() {
	cq.Q(cq.M.GetSelectQuery(cq))
	customs, _ := cq.M.(CustomRelI)
	for _, rel := range DB.Table(reflect.TypeOf(cq.M)).Relations {
		cq.selectRel(rel, customs, make([]string, 0), make([]string, 0))
	}
}
