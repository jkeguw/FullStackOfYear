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
		mongoURI = "mongodb://root:example@mongodb:27017"
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

	// 获取数据库
	db := client.Database("cpc")

	// 创建集合
	collections := []string{
		"devices",
		"device_reviews",
		"user_devices",
		"users",
		"orders",
		"carts",
	}

	for _, collName := range collections {
		err := db.CreateCollection(ctx, collName)
		if err != nil {
			// 检查是否是因为集合已存在而报错
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 { // 48是集合已存在的错误码
				fmt.Printf("Collection '%s' already exists, skipping creation\n", collName)
			} else {
				log.Printf("Warning: Failed to create collection '%s': %v\n", collName, err)
			}
		} else {
			fmt.Printf("Collection '%s' created successfully\n", collName)
		}
	}

	// 创建索引
	indexes := map[string][]mongo.IndexModel{
		"devices": {
			{Keys: bson.D{{Key: "name", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "brand", Value: 1}}, Options: options.Index()},
			{Keys: bson.D{{Key: "type", Value: 1}}, Options: options.Index()},
		},
		"device_reviews": {
			{Keys: bson.D{{Key: "deviceId", Value: 1}}, Options: options.Index()},
			{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index()},
		},
		"user_devices": {
			{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index()},
		},
		"users": {
			{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		"orders": {
			{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index()},
		},
		"carts": {
			{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index()},
		},
	}

	for collName, collIndexes := range indexes {
		collection := db.Collection(collName)
		for _, indexModel := range collIndexes {
			// 提取索引字段的名称用于日志显示
			var fieldNames string
			for _, key := range indexModel.Keys.(bson.D) {
				fieldNames += key.Key + " "
			}

			_, err := collection.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				log.Printf("Warning: Failed to create index on '%s' (fields: %s): %v\n", collName, fieldNames, err)
			} else {
				fmt.Printf("Index created successfully on '%s' (fields: %s)\n", collName, fieldNames)
			}
		}
	}

	fmt.Println("MongoDB collections and indexes initialization completed successfully!")
}
