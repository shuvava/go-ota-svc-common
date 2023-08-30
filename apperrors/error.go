package apperrors

import (
	"errors"
	"fmt"
)

// AppErrorCode is a type of AppError
// each code should be in format "namespace:error_code"
// namespace value should be unique across application
// error_code should be unique across package
type AppErrorCode string

// AppError application error with additional details
type AppError struct {
	// ErrorCode application error code
	ErrorCode AppErrorCode `json:"error_code"`
	// Description description of error
	Description string `json:"description"`
}

func (err AppError) Error() string {
	return fmt.Sprintf("(%s) : %s", err.ErrorCode, err.Description)
}

// NewAppError creates new AppError
func NewAppError(code AppErrorCode, descr string) error {
	return AppError{
		ErrorCode:   code,
		Description: descr,
	}
}

// CreateError create new AppError
func CreateError(code AppErrorCode, descr string, err error) error {
	return NewAppError(code, fmt.Sprintf("%s (%v)", descr, err))
}

// ToAppErrorWithCode unwrap generic error to AppError or create AppErrorCode with provided code
func ToAppErrorWithCode(err error, code AppErrorCode) *AppError {
	if err == nil {
		return nil
	}
	var appErr AppError
	if errors.As(err, &appErr) {
		return &appErr
	}
	appErr.ErrorCode = code
	appErr.Description = err.Error()
	return &appErr
}

// ToAppError unwrap generic error to AppError or create AppErrorCode with ErrorGeneric code
func ToAppError(err error) *AppError {
	return ToAppErrorWithCode(err, ErrorGeneric)
}
