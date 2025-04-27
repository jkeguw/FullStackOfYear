package middleware

import (
	"github.com/gin-gonic/gin"
	"strings"
	"project/backend/services/i18n"
)

// I18nKey is the key used to store the locale in the context
const I18nKey = "locale"

// IsValidLocale checks if a locale is supported
func IsValidLocale(locale string, supportedLocales []string) bool {
	for _, supported := range supportedLocales {
		if locale == supported {
			return true
		}
	}
	return false
}

// I18n middleware extracts the locale from the request
func I18n(i18nService i18n.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var locale string

		// 优先从请求头中获取语言设置
		locale = c.GetHeader("Accept-Language")
		
		// 如果没有或格式不正确，检查cookie
		if locale == "" || !IsValidLocale(locale, i18nService.GetSupportedLocales()) {
			cookie, err := c.Cookie("locale")
			if err == nil && IsValidLocale(cookie, i18nService.GetSupportedLocales()) {
				locale = cookie
			}
		}

		// 如果仍然没有有效的语言设置，使用默认语言
		if locale == "" || !IsValidLocale(locale, i18nService.GetSupportedLocales()) {
			locale = i18nService.GetDefaultLocale()
		}

		// 如果是完整的语言代码（如zh-CN;q=0.9），提取主要部分
		locale = normalizeLocale(locale)

		// 将语言设置存储在上下文中
		c.Set(I18nKey, locale)

		// 设置响应头，告诉客户端服务器使用的语言
		c.Header("Content-Language", locale)

		c.Next()
	}
}

// normalizeLocale normalizes the locale string
// Example: "zh-CN;q=0.9" -> "zh-CN"
func normalizeLocale(locale string) string {
	// 首先处理质量值
	parts := strings.Split(locale, ";")
	base := strings.TrimSpace(parts[0])
	return base
}

// T is a helper function to translate messages in handlers
func T(c *gin.Context, key string, params map[string]interface{}) string {
	i18nService, exists := c.Get("i18n")
	if !exists {
		return key
	}

	locale, _ := c.Get(I18nKey)
	localeStr, ok := locale.(string)
	if !ok {
		localeStr = "en-US" // 默认为英文
	}

	return i18nService.(i18n.Service).T(localeStr, key, params)
}

// ErrorWithTranslation sends a translated error response
func ErrorWithTranslation(c *gin.Context, httpStatus int, errorKey string, params map[string]interface{}) {
	message := T(c, errorKey, params)
	c.JSON(httpStatus, gin.H{
		"error": message,
	})
	c.Abort()
}