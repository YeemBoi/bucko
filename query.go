package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
)

type CtxQuery struct {
	R          *ReqCtx
	Q          *bun.SelectQuery
	M          BaseFieldModel
	JoinPrefix string
	TableAlias bun.Safe
}

func (cq *CtxQuery) SafeCol(column bun.Ident) bun.Safe {
	return bun.Safe(fmt.Sprintf("%s.`%s`", cq.TableAlias, column))
}

func (cq *CtxQuery) ColAlias(column bun.Ident) bun.Ident {
	return bun.Ident(cq.JoinPrefix) + column
}

func (cq *CtxQuery) CheckFieldExists() (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "? = ?", cq.SafeCol(cq.M.GetColumn()), cq.M.GetField())
}

func (cq *CtxQuery) CheckIdExists() (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "id = ?", cq.M.GetId())
}

func (cq *CtxQuery) CheckExistsFromField(field interface{}) (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "? = ?", cq.SafeCol(cq.M.GetColumn()), field)
}

func (cq *CtxQuery) CheckExistsFromId(id uint64) (exists bool, err error) {
	return BaseCheckExists(cq.M, cq.R.Ctx, "id = ?", id)
}

func (cq *CtxQuery) DeleteFromId() (err error) {
	return cq.R.delete(cq.M)
}

func (cq *CtxQuery) SetLimitOffset() {
	var limit, offset uint64
	limit, err := strconv.ParseUint(cq.R.Context.QueryParam("limit"), 10, 8)
	if err != nil {
		cq.Q = cq.Q.Limit(25)
	} else {
		if limit > 100 {
			cq.Q = cq.Q.Limit(100)
		} else {
			cq.Q = cq.Q.Limit(int(limit))
		}
	}
	if offset, err = strconv.ParseUint(cq.R.Context.QueryParam("offset"), 10, 64); err == nil {
		cq.Q = cq.Q.Offset(int(offset))
	}
}

func (cq *CtxQuery) ApplySearch(searchCols bun.Safe) {
	cq.Q = cq.R.search(cq.Q, searchCols)
}

func (cq *CtxQuery) WhereParamToCol() {
	cq.Q = cq.Q.Where("? = ?", cq.M.GetColumn(), cq.R.Context.Param(cq.M.GetParam())).Limit(1)
}

func (cq *CtxQuery) FromParam() (err error) {
	cq.Select()
	return cq.Q.Where("? = ?", cq.M.GetColumn(), cq.R.Context.Param(cq.M.GetParam())).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) FromField() (err error) {
	cq.Select()
	return cq.Q.Where("? = ?", cq.M.GetColumn(), cq.M.GetField()).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) FromId() (err error) {
	return cq.M.GetSelectQuery(cq).WherePK().Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromField(field interface{}) (err error) {
	return cq.M.GetSelectQuery(cq).Where("? = ?", cq.SafeCol(cq.M.GetColumn()), field).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromId(id uint64) (err error) {
	return cq.M.GetSelectQuery(cq).WherePK().Scan(cq.R.Ctx)
}
func (cq *CtxQuery) Insert() (err error) {
	cq.M, err = cq.M.Insert(cq.R)
	return
}

func (cq *CtxQuery) SetId() (err error) {
	return DB.NewSelect().Column("id").Model(cq.M).
		Where("? = ?", cq.SafeCol(cq.M.GetColumn()), cq.M.GetField()).Limit(1).Scan(cq.R.Ctx, cq.M.GetId())
}

func (cq *CtxQuery) Select() {
	cq.Q = cq.M.GetSelectQuery(cq)
	for _, rel := range cq.M.GetRelationQueries() {
		cq.selectRel(rel, make([]string, 0), make([]string, 0))
	}
}

func (cq *CtxQuery) selectRel(rel *RelationQuery, oldNames []string, oldPrefixes []string) {
	newNames := append(oldNames, rel.Name)
	newPrefixes := append(oldPrefixes, rel.JoinPrefix)
	cq.Q = cq.Q.Relation(strings.Join(newNames, "."),
		func(q *bun.SelectQuery) *bun.SelectQuery {
			return rel.Model.GetSelectQuery(&CtxQuery{
				R:          cq.R,
				Q:          q,
				JoinPrefix: strings.Join(newPrefixes, "__") + "__",
				TableAlias: bun.Safe(strings.Join(newPrefixes, "__")),
			})
		})
	for _, rel := range rel.Model.GetRelationQueries() {
		cq.selectRel(rel, newNames, newPrefixes)
	}
}
