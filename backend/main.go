package main

import (
	"FullStackOfYear/backend/api/v1"
	"FullStackOfYear/backend/config"
	authHandler "FullStackOfYear/backend/handlers/auth"
	"FullStackOfYear/backend/internal/database"
	authService "FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/oauth"
	"FullStackOfYear/backend/services/token"
	"github.com/gin-gonic/gin"
	"log"
)

func initOAuthComponents(cfg *config.Config) error {
	stateManager := oauth.NewStateManager(database.RedisClient)
	provider := oauth.NewGoogleProvider(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
	)

	tokenManager := token.NewManager(database.RedisClient)
	authSvc := authService.NewAuthService(
		database.MongoClient.Database(cfg.MongoDB.Database).Collection("users"),
		tokenManager,
	)

	authHandler.InitOAuthHandler(stateManager, provider, authSvc)
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("init config failed:", err)
	}
	defer config.Logger.Sync()

	// 初始化OAuth组件
	if err := initOAuthComponents(cfg); err != nil {
		log.Fatal("Failed to initialize OAuth components:", err)
	}

	// 设置路由
	router := gin.Default()
	api := router.Group("/api/v1")
	v1.RegisterRoutes(api)

	// 启动服务器
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
