package bucko

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

type (
	// CustomRelI provides functionality to customize relations to a BaseFieldModel.
	// To disable a relation, return nil.
	CustomRelI interface {
		BaseFieldModel
		GetCustomRel(rc *ReqCtx, rel *schema.Relation) *RelApplier
	}

	RelApplyFunc func(q *bun.SelectQuery) *bun.SelectQuery

	// RelApplier can customize how to handle a relation.
	RelApplier struct {
		UseDefault  bool
		Ignore      bool
		ApplyFunc   RelApplyFunc
		FollowModel BaseFieldModel
	}
)

var (
	DefaultRel = &RelApplier{UseDefault: true}
	IgnoredRel = &RelApplier{Ignore: true}
)

func (cq *CtxQuery) selectRel(rel *schema.Relation, customs CustomRelI, oldNames []string, oldPrefixes []string) {
	if rel.Type != schema.HasOneRelation {
		return
	}

	newNames := append(oldNames, rel.Field.GoName)
	newPrefixes := append(oldPrefixes, rel.Field.Name)

	var m BaseFieldModel
	var applyFunc RelApplyFunc

	var relApplier *RelApplier
	if customs != nil {
		relApplier = customs.GetCustomRel(cq.R, rel)
		if relApplier.Ignore {
			return
		} else if relApplier.ApplyFunc != nil {
			applyFunc = relApplier.ApplyFunc
		} else if relApplier.FollowModel != nil {
			m = relApplier.FollowModel
		}
	}
	if customs == nil || (relApplier.UseDefault && relApplier.FollowModel == nil) {
		var ok bool
		m, ok = rel.JoinTable.ZeroIface.(BaseFieldModel)
		if !ok {
			err := fmt.Errorf("joined table %v does not implement BaseFieldModel", rel.JoinTable.ModelName)
			fmt.Println(err)
			return
		}
	} else {
		return
	}

	if m != nil && applyFunc == nil {
		applyFunc = func(q *bun.SelectQuery) *bun.SelectQuery {
			return m.GetSelectQuery(&CtxQuery{
				R:           cq.R,
				SelectQuery: q,
				JoinPrefix:  fmt.Sprintf("%s__", strings.Join(newPrefixes, "__")),
				TableAlias:  bun.Safe(fmt.Sprintf("`%s`", strings.Join(newPrefixes, "__"))),
			})
		}
	} else if applyFunc == nil {
		return
	}

	newCustoms, _ := rel.JoinTable.ZeroIface.(CustomRelI)
	cq.Q(cq.Relation(strings.Join(newNames, "."), applyFunc))
	for _, rel := range rel.JoinTable.Relations {
		cq.selectRel(rel, newCustoms, newNames, newPrefixes)
	}
}
