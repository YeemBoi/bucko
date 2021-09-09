package bucko

import (
	"errors"
	"fmt"
	"sync"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

type (

	// defaultUser is the default user for parsing JWT.
	defaultUser struct {
		bun.BaseModel `bun:"users,alias:user"`
		Username      string `json:"username" bun:"username"`
		Id            uint64 `json:"-" bun:"id,pk"`
	}

	defaultJwtClaims struct {
		User *defaultUser `json:"user"`
		jwt.StandardClaims
	}

	syncedJwtConfig struct {
		middleware.JWTConfig
		userGetter func(claims jwt.Claims) (u BaseFieldModel, err error)
		mu         sync.Mutex
	}
)

var jwtConfig = syncedJwtConfig{
	JWTConfig: middleware.DefaultJWTConfig,
	userGetter: func(claims jwt.Claims) (u BaseFieldModel, err error) {
		customClaims, ok := claims.(defaultJwtClaims)
		if !ok {
			err = errors.New("could not parse user")
			return
		}
		return customClaims.User, nil
	},
}

// SetJWTConfig sets the primary JWT config from an echo.middleware.JWTConfig instance.
func SetJWTConfig(config middleware.JWTConfig) {
	jwtConfig.mu.Lock()
	jwtConfig.JWTConfig = config
	jwtConfig.mu.Unlock()
}

// SetUserGetter takes a function to transform claims into a user as a `BaseFieldModel`.
func SetUserGetter(userGetter func(claims jwt.Claims) (u BaseFieldModel, err error)) {
	jwtConfig.mu.Lock()
	jwtConfig.userGetter = userGetter
	jwtConfig.mu.Unlock()
}

// defaults ...

func (m *defaultUser) GetId() *uint64 {
	return &m.Id
}

func (*defaultUser) GetParam() string {
	return "username"
}

func (*defaultUser) GetColumn() bun.Ident {
	return "username"
}
func (m *defaultUser) GetField() interface{} {
	return m.Username
}

func (*defaultUser) GetSelectQuery(cq *CtxQuery) *bun.SelectQuery {
	return cq.Q.Column("id", "username")
}

func (s *defaultUser) GetRelationQueries() []*RelationQuery {
	return make([]*RelationQuery, 0)
}

func (b *defaultUser) Insert(rc *ReqCtx) (m BaseFieldModel, err error) {
	return nil, fmt.Errorf("cannot insert skeletal model %T", b)
}
