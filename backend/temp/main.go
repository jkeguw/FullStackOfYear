package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API路由
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// 健康检查
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Server is running",
		})
	})

	// 模拟登录API
	v1.POST("/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"access_token": "mock-access-token",
			"refresh_token": "mock-refresh-token",
			"expires_in": 3600,
			"token_type": "Bearer",
			"user_id": "123456",
			"email": "user@example.com",
			"username": "mockuser",
		})
	})

	// 设备列表API
	v1.GET("/devices", func(c *gin.Context) {
		devices := []gin.H{
			{
				"id": "1",
				"name": "Logitech G Pro X Superlight",
				"brand": "Logitech",
				"type": "mouse",
				"imageUrl": "https://example.com/gpro.jpg",
				"description": "Ultra-lightweight wireless gaming mouse",
				"createdAt": "2023-01-01T00:00:00Z",
			},
			{
				"id": "2",
				"name": "Razer Viper Ultimate",
				"brand": "Razer",
				"type": "mouse",
				"imageUrl": "https://example.com/viper.jpg",
				"description": "High-end wireless gaming mouse",
				"createdAt": "2023-01-02T00:00:00Z",
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"data": devices,
			"total": 2,
			"page": 1,
			"pageSize": 10,
		})
	})

	// 用户配置文件API
	v1.GET("/user/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"id": "123456",
			"username": "mockuser",
			"email": "user@example.com",
			"status": "active",
			"createdAt": "2023-01-01T00:00:00Z",
		})
	})

	// 启动服务器
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}