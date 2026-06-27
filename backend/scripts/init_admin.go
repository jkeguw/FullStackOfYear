package scripts

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func InitAdmin() {
	// 加载环境变量（从多个可能的位置尝试加载）
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	// 获取MongoDB连接URI，优先使用环境变量
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://root:example@localhost:27017"
		log.Printf("Using default MongoDB URI: %s", mongoURI)
	} else {
		log.Printf("Using MongoDB URI from environment")
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
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "cpc"
	}
	db := client.Database(dbName)
	usersCollection := db.Collection("users")

	// 管理员账户信息优先从环境变量读取
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "root@example.com"
	}

	adminPassword := os.Getenv("ADMIN_INITIAL_PASSWORD")
	if adminPassword == "" {
		// 如果没有设置初始密码，生成一个强随机密码并打印到日志
		adminPassword = generateSecurePassword()
		log.Printf("============================================================")
		log.Printf("WARNING: No ADMIN_INITIAL_PASSWORD set.")
		log.Printf("Generated random admin password for %s: %s", adminEmail, adminPassword)
		log.Printf("Please change this password immediately after first login.")
		log.Printf("============================================================")
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), 12)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	// 检查管理员账户是否已存在
	var existingAdmin bson.M
	err = usersCollection.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existingAdmin)

	if err == nil {
		// 管理员账户已存在
		fmt.Println("Admin account already exists, skipping creation.")
	} else if err == mongo.ErrNoDocuments {
		// 管理员账户不存在，创建它
		now := time.Now()
		adminUser := bson.M{
			"username": "admin",
			"email":    adminEmail,
			"password": string(hashedPassword),
			"role": bson.M{
				"type": "admin",
			},
			"status": bson.M{
				"emailVerified": true,
			},
			"stats": bson.M{
				"createdAt":   now,
				"lastLoginAt": now,
			},
			"createdAt":  now,
			"updatedAt":  now,
			"isVerified": true,
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

func generateSecurePassword() string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		// 如果随机生成失败，使用一个明确的占位符（不应该在生产环境发生）
		return "CHANGE_ME_IMMEDIATELY_" + fmt.Sprintf("%d", time.Now().Unix())
	}
	return base64.URLEncoding.EncodeToString(b)
}
