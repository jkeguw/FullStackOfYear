package middleware

import (
	"project/backend/config"
	"project/backend/internal/errors"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 初始化配置和日志
	if config.Cfg == nil {
		config.Cfg = &config.Config{}
	}
	if config.Logger == nil {
		// 使用空操作日志器，避免测试输出过多日志
		config.Logger = zap.NewNop()
	}

	// 创建带恢复中间件的路由器
	router := gin.New()
	router.Use(Recovery())

	// 定义一个会引发 panic 的处理程序
	router.GET("/panic", func(c *gin.Context) {
		// 故意引发 panic
		panic("test panic situation")
	})

	// 定义一个正常的处理程序
	router.GET("/normal", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 测试正常路由
	t.Run("Normal Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/normal", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	// 测试引发 panic 的路由
	t.Run("Panic Route", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/panic", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 验证状态码是否为 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		// 验证响应是否为预期的错误格式
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(errors.InternalError), response["code"])
		assert.Equal(t, "服务器内部错误", response["message"])
	})
}

func TestRecoveryWithNestedPanic(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 初始化日志
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	// 创建带恢复中间件的路由器
	router := gin.New()
	router.Use(Recovery())

	// 定义一个引发复杂panic的处理程序
	router.GET("/nested-panic", func(c *gin.Context) {
		// 使用复杂对象引发panic，而不是递归调用
		panic(map[string]interface{}{
			"error": "complex panic",
			"data": map[string]interface{}{
				"code": 12345,
				"details": []string{
					"detail1", "detail2", "detail3",
				},
				"nested": map[string]bool{
					"isValid": false,
				},
			},
		})
	})

	// 测试复杂的 panic 情况
	req, _ := http.NewRequest("GET", "/nested-panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证状态码是否为 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// 验证响应是否为预期的错误格式
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(errors.InternalError), response["code"])
}

func TestRecoveryWithDifferentPanicTypes(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 初始化日志
	if config.Logger == nil {
		config.Logger = zap.NewNop()
	}

	// 定义不同类型的 panic 测试
	testCases := []struct {
		name       string
		panicValue interface{}
	}{
		{"String Panic", "string panic"},
		{"Error Panic", fmt.Errorf("error panic")},
		{"Integer Panic", 42},
		{"Struct Panic", struct{ Msg string }{"structured panic"}},
		{"Nil Panic", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 为每个测试用例创建新的路由器
			router := gin.New()
			router.Use(Recovery())

			// 添加引发特定类型 panic 的处理程序
			router.GET("/panic", func(c *gin.Context) {
				panic(tc.panicValue)
			})

			// 执行请求
			req, _ := http.NewRequest("GET", "/panic", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证结果
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, float64(errors.InternalError), response["code"])
			assert.Equal(t, "服务器内部错误", response["message"])
		})
	}
}