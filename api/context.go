package api

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/shuvava/go-logging/logger"
)

// GetRequestContext return populated request context.Context
func GetRequestContext(ctx echo.Context) context.Context {
	c := ctx.Request().Context()
	rid := ctx.Response().Header().Get(echo.HeaderXRequestID)
	return context.
		WithValue(c, logger.ContextKeyRequestID, rid)
}

// GetContentType returns value of ContentType header
func GetContentType(ctx echo.Context) string {
	return ctx.Request().Header.Get(echo.HeaderContentType)
}

// GetContentSize returns Content length from request header
func GetContentSize(ctx echo.Context) int64 {
	s := ctx.Request().Header.Get(echo.HeaderContentLength)
	if s == "" {
		return 0
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return size
}
