package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project/backend/config"
	"time"
)

var MongoClient *mongo.Client

func InitMongoDB(ctx context.Context) error {
	// 设置MongoDB客户端选项
	clientOptions := options.Client().ApplyURI(config.Cfg.MongoDB.URI)

	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)

	fmt.Printf("正在连接MongoDB，URI: %s, 数据库: %s\n", config.Cfg.MongoDB.URI, config.Cfg.MongoDB.Database)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("MongoDB连接错误: %w", err)
	}

	MongoClient = client

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB Ping失败: %w", err)
	}

}
