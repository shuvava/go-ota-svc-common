package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// ServerHeader middleware adds a `Server` header to the response.
func ServerHeader(name, version string) echo.MiddlewareFunc {
	appVer := fmt.Sprintf("%s/%s", name, version)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, appVer)
			return next(c)
		}
	}
}
