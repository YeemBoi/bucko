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
	// DefaultUser is the Default user for parsing JWT.
	DefaultUser struct {
		bun.BaseModel `bun:"users,alias:user"`
		Username      string `json:"username" bun:"username"`
		Id            uint64 `json:"-" bun:"id,pk"`
	}

	DefaultJwtClaims struct {
		User *DefaultUser `json:"user"`
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
		customClaims, ok := claims.(DefaultJwtClaims)
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

// Defaults ...

func (*DefaultUser) GetParam() string {
	return "username"
}

func (*DefaultUser) GetColumn() bun.Ident {
	return "username"
}
func (m *DefaultUser) GetField() interface{} {
	return m.Username
}

func (*DefaultUser) GetSelectQuery(cq *CtxQuery) *bun.SelectQuery {
	return cq.Q.Column("id", "username")
}

func (b *DefaultUser) Insert(rc *ReqCtx) (m BaseFieldModel, err error) {
	return nil, fmt.Errorf("cannot insert skeletal model %T", b)
}
