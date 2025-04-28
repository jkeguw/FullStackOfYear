package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// 获取MongoDB连接URI，如果环境变量中没有，则使用默认值
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// 尝试使用docker容器的默认连接
		mongoURI = "mongodb://root:example@localhost:27017"
		log.Printf("Using default MongoDB URI: %s", mongoURI)
	} else {
		log.Printf("Using MongoDB URI from environment: %s", mongoURI)
	}

	// 连接到MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// 确认连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully!")

	// 获取数据库和集合
	db := client.Database("cpc")
	usersCollection := db.Collection("users")

	// 检查管理员账户是否已存在
	adminEmail := "2353727288@qq.com"
	var existingAdmin bson.M
	err = usersCollection.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existingAdmin)

	if err == nil {
		// 管理员账户已存在
		fmt.Println("Admin account already exists, skipping creation.")
	} else if err == mongo.ErrNoDocuments {
		// 管理员账户不存在，创建它
		now := time.Now()
		adminUser := bson.M{
			"username":     "admin",
			"email":        adminEmail,
			"passwordHash": "$2a$10$GjSIOtMDbVN1XmCvpZXfDuH/DfEQD6jKZ1q8.1Y9TQy.ggZEFEX5y", // 密码为 "unnamed03634614"
			"roles":        []string{"admin"},
			"createdAt":    now,
			"updatedAt":    now,
		}

		// 插入管理员账户
		res, err := usersCollection.InsertOne(ctx, adminUser)
		if err != nil {
			log.Fatalf("Failed to create admin account: %v", err)
		}
		fmt.Printf("Admin account created successfully with ID: %v\n", res.InsertedID)
	} else {
		// 其他错误
		log.Fatalf("Error checking for existing admin account: %v", err)
	}

	// 确保索引存在
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = usersCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: Failed to create email index: %v", err)
	} else {
		fmt.Println("Email index created or already exists")
	}

	fmt.Println("Admin account initialization completed!")
}
