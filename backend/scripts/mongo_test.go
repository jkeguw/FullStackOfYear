package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 设置MongoDB连接
	uri := "mongodb://root:example@mongodb:27017/?authSource=admin"
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("无法连接到MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Ping失败:", err)
	}
	fmt.Println("成功连接到MongoDB!")

	// 获取数据库和集合
	db := client.Database("cpc")
	collection := db.Collection("devices")

	// 插入测试文档
	testDoc := bson.M{
		"name":        "Test Mouse Go",
		"brand":       "TestBrand",
		"type":        "mouse",
		"description": "Test document inserted by Go script",
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	}

	_, err = collection.InsertOne(ctx, testDoc)
	if err != nil {
		log.Fatal("插入文档失败:", err)
	}
	fmt.Println("成功插入测试文档")

	// 查询设备
	filter := bson.M{"type": "mouse"}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer cur.Close(ctx)

	var results []bson.M
	if err = cur.All(ctx, &results); err != nil {
		log.Fatal("解析结果失败:", err)
	}

	fmt.Printf("找到 %d 个鼠标设备:\n", len(results))
	for i, result := range results {
		fmt.Printf("%d: %s (%s)\n", i+1, result["name"], result["brand"])
	}
}
