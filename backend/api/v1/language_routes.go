package v1

import (
	"project/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// registerLanguageRoutes registers language related routes
func (r *Router) registerLanguageRoutes(router *gin.RouterGroup) {
	languageRoutes := router.Group("/languages")
	{
		languageRoutes.GET("", getLanguages)
		languageRoutes.GET("/current", getCurrentLanguage)
		languageRoutes.POST("/set", setLanguage)
	}
}

// getLanguages returns the list of supported languages
func getLanguages(c *gin.Context) {
	// 简化实现，返回固定的语言列表
	c.JSON(http.StatusOK, gin.H{
		"languages": []string{"en-US", "zh-CN"},
		"default":   "en-US",
		"current":   "en-US",
	})
}

// getCurrentLanguage returns the current language
func getCurrentLanguage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"language": "en-US",
	})
}

// setLanguage sets the language
func setLanguage(c *gin.Context) {
	var request struct {
		Language string `json:"language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Invalid request"))
		return
	}

	// 验证语言是否受支持
	validLanguages := []string{"en-US", "zh-CN"}
	isValid := false
	for _, lang := range validLanguages {
		if request.Language == lang {
			isValid = true
			break
		}
	}
	
	if !isValid {
		c.JSON(http.StatusBadRequest, errors.NewAppError(errors.BadRequest, "Unsupported language"))
		return
	}

	// 设置cookie，用于后续请求
	c.SetCookie(
		"locale", 
		request.Language,
		60*60*24*30, // 30天过期
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"language": request.Language,
	})
}