package utils

import (
	"project/backend/internal/errors"
	"project/backend/middleware"
	"project/backend/services/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
)

// I18nError 处理国际化的错误响应
func I18nError(c *gin.Context, i18nService *i18n.Service, err error) {
	// 获取用户语言
	lang := middleware.GetLanguage(c)
	
	// 获取错误代码
	code := errors.GetErrorCode(err)
	
	// 获取HTTP状态码
	httpStatus := errors.HTTPStatusFromError(err)
	
	// 根据错误代码获取i18n消息ID
	var messageID string
	switch code {
	case errors.BadRequest:
		messageID = "errors.bad_request"
	case errors.Unauthorized:
		messageID = "errors.unauthorized"
	case errors.Forbidden:
		messageID = "errors.forbidden"
	case errors.NotFound:
		messageID = "errors.not_found"
	case errors.InternalError:
		messageID = "errors.internal_server"
	case errors.TooManyRequests:
		messageID = "errors.too_many_requests"
	default:
		messageID = "errors.internal_server"
	}
	
	// 获取翻译后的错误消息
	message := i18nService.T(lang, messageID, nil)
	
	// 返回JSON响应
	c.JSON(httpStatus, gin.H{
		"code":    code,
		"message": message,
		"error":   true,
	})
}

// I18nSuccess 处理国际化的成功响应
func I18nSuccess(c *gin.Context, i18nService *i18n.Service, messageID string, data interface{}) {
	// 获取用户语言
	lang := middleware.GetLanguage(c)
	
	// 获取翻译后的消息
	message := i18nService.T(lang, messageID, nil)
	
	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": message,
		"data":    data,
	})
}