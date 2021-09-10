package bucko

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
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
	cq.Q = cq.R.Search(cq.Q, searchCols)
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

func (cq *CtxQuery) FromPK() (err error) {
	return cq.M.GetSelectQuery(cq).WherePK().Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromField(field interface{}) (err error) {
	return cq.M.GetSelectQuery(cq).Where("? = ?", cq.SafeCol(cq.M.GetColumn()), field).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) GetFromPK(id uint64) (err error) {
	return cq.M.GetSelectQuery(cq).Where("? = ?", cq.SafeCol(GetPKCol(cq.M)), id).Limit(1).Scan(cq.R.Ctx)
}
func (cq *CtxQuery) Insert() (err error) {
	cq.M, err = cq.M.Insert(cq.R)
	return
}

func (cq *CtxQuery) SetPK() (err error) {
	return DB.NewSelect().Column(string(GetPKCol(cq.M))).Model(cq.M).
		Where("? = ?", cq.SafeCol(cq.M.GetColumn()), cq.M.GetField()).Limit(1).Scan(cq.R.Ctx)
}

func (cq *CtxQuery) Select() {
	cq.Q = cq.M.GetSelectQuery(cq)
	for _, rel := range DB.Table(reflect.TypeOf(cq.M)).Relations {
		cq.selectRel(rel, make([]string, 0), make([]string, 0))
	}
}

func (cq *CtxQuery) selectRel(rel *schema.Relation, oldNames []string, oldPrefixes []string) {
	if rel.Type != schema.BelongsToRelation {
		return
	}

	newNames := append(oldNames, rel.Field.GoName)
	newPrefixes := append(oldPrefixes, rel.Field.Name)

	var m, ok = rel.JoinTable.ZeroIface.(BaseFieldModel)
	if !ok {
		err := fmt.Errorf("joined table %v does not implement BaseFieldModel", rel.JoinTable.ModelName)
		fmt.Println(err)
	}
	cq.Q = cq.Q.Relation(strings.Join(newNames, "."),
		func(q *bun.SelectQuery) *bun.SelectQuery {
			return m.GetSelectQuery(&CtxQuery{
				R:          cq.R,
				Q:          q,
				JoinPrefix: fmt.Sprintf("%s__", strings.Join(newPrefixes, "__")),
				TableAlias: bun.Safe(fmt.Sprintf("`%s`", strings.Join(newPrefixes, "__"))),
			})
		})
	for _, rel := range rel.JoinTable.Relations {
		cq.selectRel(rel, newNames, newPrefixes)
	}
}
