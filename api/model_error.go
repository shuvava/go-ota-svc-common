package api

import (
	"context"
	"errors"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// ErrorResponse is http error response model
type ErrorResponse struct {
	// ErrorCode application error code
	ErrorCode string `json:"error_code"`
	// StatusCode HTTP response status code
	StatusCode int `json:"status_code"`
	// Description description of error
	Description string `json:"description"`
	// RequestID HTTP requestID go from header of request
	RequestID string `json:"request_id"`
}

// NewErrorResponse creates new error response from error
func NewErrorResponse(ctx context.Context, statusCode int, err error) ErrorResponse {
	requestID := logger.GetRequestID(ctx)
	resp := ErrorResponse{
		StatusCode: statusCode,
		RequestID:  requestID,
	}

	var typedErr apperrors.AppError
	if errors.As(err, &typedErr) {
		resp.ErrorCode = string(typedErr.ErrorCode)
		resp.Description = typedErr.Description
	} else {
		resp.ErrorCode = apperrors.ErrorGeneric
		resp.Description = err.Error()
	}

	return resp
}
