package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"project/backend/api/router"
	"project/backend/config"
	"project/backend/internal/database"
	"project/backend/services/auth"
	"project/backend/services/email"
	"project/backend/services/i18n"
	"project/backend/services/jwt"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/config.yaml", "path to config file")
	flag.Parse()

	// 初始化配置
	if err := config.InitConfig(*configPath); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化日志
	if err := config.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化数据库连接
	if err := database.InitMongoDB(ctx); err != nil {
		log.Printf("警告: MongoDB连接失败: %v", err)
		log.Println("继续运行，但数据库功能将不可用...")
	} else {
		log.Println("MongoDB连接成功")
	}

	redisClient := database.InitRedis()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("警告: Redis连接失败: %v", err)
		log.Println("继续运行，但Redis缓存和会话功能将不可用...")
	} else {
		log.Println("Redis连接成功")
	}

	// 初始化服务
	jwtService := jwt.NewService(config.GetConfig().JWT)
	emailService := email.NewService(config.GetConfig().Email)
	oauthProvider := auth.NewMockOAuthProvider()
	
	// 获取用户集合
	var userCollection *mongo.Collection
	if database.MongoClient != nil {
		userCollection = database.MongoClient.Database("app").Collection("users")
	}
	
	// 初始化身份验证服务
	tokenGenerator := auth.NewSimpleTokenGenerator(jwtService)
	authService := auth.NewService(
		userCollection,
		tokenGenerator,
		emailService,
		oauthProvider,
	)
	
	i18nService := i18n.NewService()

	// 初始化路由
	r := gin.Default()
	
	// 健康检查路由
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
			"mongo":  database.MongoClient != nil,
			"redis":  redisClient != nil,
		})
	})

	// 设置路由
	router.InitRouter(r, authService, jwtService, i18nService)

	// 启动服务器
	cfg := config.GetConfig()
	port := cfg.Server.Port
	
	// 清理端口格式
	port = strings.TrimPrefix(port, ":")
	if port == "" {
		port = "8080" // 设置默认端口
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server started on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}