package echox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/requests"
)

func TestRoute(t *testing.T) {
	e := New()
	e.Route(nil, NewRouter("app1", "/app1", func(g *echo.Group) {
		g.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "app1") })
	}))
	e.Route(nil, NewRouter("app2", "/app2", func(g *echo.Group) {
		g.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "app2") })
	}))
	e.Route(e.Group("/api"), NewRouter("app3", "/app3", func(g *echo.Group) {
		g.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "app3") })
	}))

	ts := httptest.NewServer(e)
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	{
		resp, err := requests.Get("%s/app1/", ts.URL).Do(ctx)
		require.NoError(t, err)
		require.NoErrorf(t, resp.Success(), "failed with status %v", resp.StatusCode)
		require.Equal(t, "app1", resp.String())
	}

	{
		resp, err := requests.Get("%s/app2/", ts.URL).Do(ctx)
		require.NoError(t, err)
		require.NoErrorf(t, resp.Success(), "failed with status %v", resp.StatusCode)
		require.Equal(t, "app2", resp.String())
	}

	{
		resp, err := requests.Get("%s/api/app3/", ts.URL).Do(ctx)
		require.NoError(t, err)
		require.NoErrorf(t, resp.Success(), "failed with status %v", resp.StatusCode)
		require.Equal(t, "app3", resp.String())
	}
}
