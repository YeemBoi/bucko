package bucko

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/uptrace/bun"
)

type (
	BaseUser struct {
		bun.BaseModel `bun:"users,alias:user"`
		Username      string `json:"username" bun:"username"`
		Id            uint64 `json:"-" bun:"id,pk"`
	}
	JwtCustomClaims struct {
		User BaseUser `json:"user"`
		jwt.StandardClaims
	}
)

func (m *BaseUser) GetId() *uint64 {
	return &m.Id
}

func (*BaseUser) GetParam() string {
	return "username"
}

func (*BaseUser) GetColumn() bun.Ident {
	return "username"
}
func (m *BaseUser) GetField() interface{} {
	return m.Username
}

func (*BaseUser) GetSelectQuery(cq *CtxQuery) *bun.SelectQuery {
	return cq.Q.Column("id", "username")
}

func (s *BaseUser) GetRelationQueries() []*RelationQuery {
	return make([]*RelationQuery, 0)
}

func (b *BaseUser) Insert(rc *ReqCtx) (m BaseFieldModel, err error) {
	return nil, fmt.Errorf("Cannot insert skeletal model %T", b)
}
