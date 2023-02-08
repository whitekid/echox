package echox

import "github.com/labstack/echo/v4"

// BindHeader bind header & validate it
// NOTE Echo.Bind()는 header를 처리하지 않는다.
// TODO move to echox
func BindHeader(c echo.Context, val any) error {
	if b, ok := c.Echo().Binder.(*echo.DefaultBinder); ok {
		if err := b.BindHeaders(c, val); err != nil {
			return err
		}
	}

	if err := c.Validate(val); err != nil {
		return err
	}

	return nil
}
