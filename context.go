package echox

import (
	"github.com/labstack/echo/v4"
)

// TODO move to gormx
type ContextFactory interface {
	ContextIs(echo.Context) bool
	NewContext(echo.Context) echo.Context
}

func CustomContext(factory ContextFactory) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !factory.ContextIs(c) {
				c = factory.NewContext(c)
			}

			return next(c)
		}
	}
}
