package utils

import (
	"project/backend/middleware"
	"project/backend/services/i18n"
	"github.com/gin-gonic/gin"
)

// I18nResponse 是国际化响应的助手结构
type I18nResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data,omitempty"`
	Error   interface{}            `json:"error,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// NewI18nResponse 创建一个新的国际化响应
func NewI18nResponse(c *gin.Context, i18nService i18n.Service, messageID string, code int, data interface{}, err interface{}) I18nResponse {
	lang := middleware.GetLanguage(c)
	message := i18nService.T(lang, messageID, nil)

	return I18nResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Error:   err,
	}
}

// WithMeta 添加元数据到响应
func (r I18nResponse) WithMeta(meta map[string]interface{}) I18nResponse {
	r.Meta = meta
	return r
}

// Send 发送响应到客户端
func (r I18nResponse) Send(c *gin.Context) {
	c.JSON(r.Code, r)
}

// Success 创建一个成功的响应
func Success(c *gin.Context, i18nService i18n.Service, messageID string, data interface{}) {
	NewI18nResponse(c, i18nService, messageID, 200, data, nil).Send(c)
}

// Error 创建一个错误响应
func Error(c *gin.Context, i18nService i18n.Service, messageID string, code int, err interface{}) {
	NewI18nResponse(c, i18nService, messageID, code, nil, err).Send(c)
}