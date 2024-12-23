package echox

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/whitekid/goxp/validate"
)

type Echo struct {
	*echo.Echo
}

// New create new default echo handlers
func New(middlewares ...echo.MiddlewareFunc) *Echo {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &Validator{}
	e.Use(middleware.Recover())
	e.Use(LogErrors())
	e.Use(middlewares...)

	return &Echo{Echo: e}
}

func (e *Echo) Route(g *echo.Group, r Router) {
	if g == nil {
		g = e.Group(r.Path())
	} else {
		g = g.Group(r.Path())
	}
	r.Route(g)
}

// Middlewares

// LogErrors log error when http status error occurred
func LogErrors() echo.MiddlewareFunc { return LogErrorsWithCode(http.StatusInternalServerError) }
func LogErrorsWithCode(logCode int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// log http errors
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				code := http.StatusInternalServerError

				if ee, ok := err.(validator.ValidationErrors); ok {
					err = echo.NewHTTPError(http.StatusBadRequest, ee.Error())
				}

				if he, ok := err.(*echo.HTTPError); ok {
					code = he.Code
				}

				if code >= logCode {
					c.Logger().Errorf("%+v", err)
				}
			}

			return err
		}
	}
}

func ExtractHeader(header string, fn func(c echo.Context, cval string)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fn(c, c.Request().Header.Get(header))

			return next(c)
		}
	}
}

// ExtractParam extract path parameter and callback to use custom context
// Usage:
// 	e.Use(ExtractParam("project_id", func(c echo.Context, val string) { c.(*Context).projectID = val }))

func ExtractParam(param string, callback func(c echo.Context, val string)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			callback(c, c.Param(param))

			return next(c)
		}
	}
}

// Bind bind & validate
func Bind[T any](c echo.Context, val *T) error {
	if err := c.Bind(val); err != nil {
		return echo.ErrBadRequest
	}

	if err := c.Validate(val); err != nil {
		return err
	}

	return nil
}

type Validator struct {
}

func (v *Validator) Validate(i any) error {
	if err := validate.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
