package bucko

import (
	"github.com/labstack/echo/v4"
)

// Middleware provides echo middleware for `echo.Echo.Use()`
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(NewCtx(c, nil))
	}
}
