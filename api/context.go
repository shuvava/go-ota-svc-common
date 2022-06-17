package api

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/go-ota-svc-common/data"
)

const (
	// DefaultNamespaceValue is default OTA namespace
	DefaultNamespaceValue = "default"

	headerNamespace = "x-ats-namespace"
)

// GetRequestContext return populated request context.Context
func GetRequestContext(ctx echo.Context) context.Context {
	c := ctx.Request().Context()
	rid := GetRequestID(ctx)
	tenantID := GetNamespace(ctx)

	newCtx := context.
		WithValue(c, logger.ContextKeyRequestID, rid)
	return context.
		WithValue(newCtx, logger.ContextKeyTenantID, tenantID)
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

// GetNamespace returns OTA namespace from header
func GetNamespace(ctx echo.Context) data.Namespace {
	ns := ctx.Request().Header.Get(headerNamespace)
	if ns == "" {
		ns = DefaultNamespaceValue
	}
	return data.NewNamespace(ns)
}

// GetRequestID returns RequestID from header
func GetRequestID(ctx echo.Context) string {
	rid := ctx.Request().Header.Get(echo.HeaderXRequestID)
	if rid == "" {
		rid = data.NewCorrelationID().String()
	}
	return rid
}
