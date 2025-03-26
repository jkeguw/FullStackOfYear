package errors

import (
	"net/http"
)

// NewBadRequestError 创建一个400错误
func NewBadRequestError(message string) *AppError {
	return NewAppError(BadRequest, message)
}

// NewUnauthorizedError 创建一个401错误
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(Unauthorized, message)
}

// NewForbiddenError 创建一个403错误
func NewForbiddenError(message string) *AppError {
	return NewAppError(Forbidden, message)
}

// NewNotFoundError 创建一个404错误
func NewNotFoundError(message string) *AppError {
	return NewAppError(NotFound, message)
}

// NewTooManyRequestsError 创建一个429错误
func NewTooManyRequestsError(message string) *AppError {
	return NewAppError(TooManyRequests, message)
}

// NewInternalServerError 创建一个500错误
func NewInternalServerError(message string) *AppError {
	return NewAppError(InternalError, message)
}

// HTTPStatusFromError 从错误中获取HTTP状态码
func HTTPStatusFromError(err error) int {
	code := GetErrorCode(err)

	switch code {
	case BadRequest:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case TooManyRequests:
		return http.StatusTooManyRequests
	case InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
