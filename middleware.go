package main

import (
	"github.com/labstack/echo/v4"
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(NewCtx(c, nil))
	}
}
