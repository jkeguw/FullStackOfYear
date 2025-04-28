package errors

import "errors"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// HTTPStatus returns the appropriate HTTP status code based on error code
func (e *AppError) HTTPStatus() int {
	if e.Code >= 400 && e.Code < 600 {
		return e.Code // If the code is already an HTTP status code, return it directly
	}

	// Otherwise, map the custom error codes to HTTP status codes
	switch e.Code {
	case BadRequest:
		return 400
	case Unauthorized:
		return 401
	case Forbidden:
		return 403
	case NotFound:
		return 404
	case Conflict:
		return 409
	case Validation:
		return 422
	case RateLimit:
		return 429
	case InternalError:
		return 500
	case NotImplemented:
		return 501
	default:
		return 500
	}
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// GetErrorCode 从错误中获取错误代码
func GetErrorCode(err error) int {
	if err == nil {
		return Success
	}

	// 尝试将错误转换为 AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}

	// 如果不是 AppError，返回内部错误
	return InternalError
}
