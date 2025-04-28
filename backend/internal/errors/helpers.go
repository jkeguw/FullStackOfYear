package errors

import (
	"github.com/gin-gonic/gin"
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

// HandleError 处理错误并发送适当的响应
func HandleError(c interface{}, err error) {
	if err == nil {
		return
	}

	// 转换为gin上下文
	ginCtx, ok := c.(*gin.Context)
	if !ok {
		// 日志记录错误，但无法发送响应
		return
	}

	// 获取HTTP状态码
	status := HTTPStatusFromError(err)

	// 无论错误类型如何，都使用统一的格式
	code := GetErrorCode(err)
	message := err.Error()

	// 构建统一响应格式，包含data: nil字段
	ginCtx.JSON(status, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// IsNotFoundError 判断错误是否为NotFound错误
func IsNotFoundError(err error) bool {
	return GetErrorCode(err) == NotFound
}

// IsForbiddenError 判断错误是否为Forbidden错误
func IsForbiddenError(err error) bool {
	return GetErrorCode(err) == Forbidden
}
