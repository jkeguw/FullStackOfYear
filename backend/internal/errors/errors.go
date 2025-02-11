package errors

import "errors"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
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
