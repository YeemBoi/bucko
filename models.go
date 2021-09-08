package bucko

import (
	"github.com/uptrace/bun"
)

type (
	BaseIdModel interface {
		// GetId must return a pointer to the model's numerical PK.
		GetId() *uint64
	}

	// BaseFieldModel provides an interface to models commonly referenced to help provide shortcuts.
	BaseFieldModel interface {

		// GetSelectQuery must return a SelectQuery suited for selecting the model based off of cq.Q.
		GetSelectQuery(cq *CtxQuery) *bun.SelectQuery

		// Insert must insert itself and the result or an error.
		Insert(rc *ReqCtx) (BaseFieldModel, error)

		// GetColumn must return the column name of the main field used to retrieve the model (not PK).
		GetColumn() bun.Ident

		// GetField must return the field used to retrieve the model (not PK).
		GetField() interface{}

		// GetParam must return the URL param name of the main field used to retrieve the model (not PK).
		GetParam() string

		// GetRelationQueries must return all of the relational queries to select.
		GetRelationQueries() []*RelationQuery

		BaseIdModel
	}

	RelationQuery struct {

		// Name must be the name of the bun relation, eg "Author".
		Name string

		// JoinPrefix is the alias of the new join, eg "author".
		JoinPrefix string

		// Model reference the new model being joined.
		Model BaseFieldModel
	}
)
