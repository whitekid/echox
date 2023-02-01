package echox

import "github.com/labstack/echo/v4"

type Router interface {
	Path() string
	Name() string
	Route(*echo.Group)
}

func NewRouter(name, path string, route func(*echo.Group)) Router {
	return &routeWith{
		name:  name,
		path:  path,
		route: route,
	}
}

type routeWith struct {
	name  string
	path  string
	route func(*echo.Group)
}

func (r *routeWith) Path() string        { return r.path }
func (r *routeWith) Name() string        { return r.name }
func (r *routeWith) Route(g *echo.Group) { r.route(g) }
