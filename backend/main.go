package main

import (
	"FullStackOfYear/backend/api/v1"
	"FullStackOfYear/backend/config"
	"FullStackOfYear/backend/internal/database"
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/email"
	"FullStackOfYear/backend/services/oauth"
	"FullStackOfYear/backend/services/token"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func initServices(cfg *config.Config, logger *zap.Logger) (auth.Service, *email.Service, error) {
	if err := database.InitMongoDB(context.Background()); err != nil {
		return nil, nil, fmt.Errorf("failed to init mongodb: %v", err)
	}

	if err := database.InitRedis(); err != nil {
		return nil, nil, fmt.Errorf("failed to init redis: %v", err)
	}

	emailCfg := &email.Config{
		SMTP: struct {
			Host     string
			Port     int
			Username string
			Password string
		}{
			Host:     cfg.Email.SMTP.Host,
			Port:     cfg.Email.SMTP.Port,
			Username: cfg.Email.SMTP.Username,
			Password: cfg.Email.SMTP.Password,
		},
		From:      cfg.Email.From,
		BaseURL:   cfg.Email.BaseURL,
		Templates: cfg.Email.Templates,
	}

	emailService, err := email.NewEmailService(emailCfg, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init email service: %v", err)
	}

	tokenManager := token.NewManager(database.RedisClient)

	googleProvider := oauth.NewGoogleProvider(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
	)

	authService := auth.NewService(
		database.MongoClient.Database(cfg.MongoDB.Database).Collection("users"),
		tokenManager,
		emailService,
		googleProvider,
	)

	return authService, emailService, nil
}

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("Failed to init config:", err)
	}

	authService, emailService, err := initServices(config.Cfg, config.Logger)
	if err != nil {
		log.Fatal("Failed to init services:", err)
	}

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	v1.NewRouter(authService, emailService).RegisterRoutes(apiV1)

	if err := router.Run(config.Cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
