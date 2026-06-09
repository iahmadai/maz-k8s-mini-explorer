package apperrors

import (
	"errors"
	"net/http"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func ServiceUnavailable(msg string) *AppError {
	return &AppError{Code: http.StatusServiceUnavailable, Message: msg}
}

func FromK8s(err error) *AppError {
	if err == nil {
		return nil
	}
	if k8serrors.IsNotFound(err) {
		return NotFound(err.Error())
	}
	if k8serrors.IsForbidden(err) {
		return &AppError{Code: http.StatusForbidden, Message: err.Error()}
	}
	return Internal(err.Error())
}

func HTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

func Message(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return err.Error()
}
