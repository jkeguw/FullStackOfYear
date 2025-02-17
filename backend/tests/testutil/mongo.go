package testutil

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// TestUser 代表测试用户数据
type TestUser struct {
	ID       primitive.ObjectID
	Email    string
	Password string
	Username string
	Role     string
}

var (
	// DefaultUsers 包含所有测试用户
	DefaultUsers = map[string]TestUser{
		"verified": {
			ID:       primitive.NewObjectID(),
			Email:    "verified@example.com",
			Password: "$2a$10$7juGndSkFAhhd/J8gmRBBOidUg69fbusoY8lQnLc4hqUXlvvkfz5G", // "password123"
			Username: "verified_user",
			Role:     "user",
		},
		"unverified": {
			ID:       primitive.NewObjectID(),
			Email:    "unverified@example.com",
			Password: "$2a$10$7juGndSkFAhhd/J8gmRBBOidUg69fbusoY8lQnLc4hqUXlvvkfz5G",
			Username: "unverified_user",
			Role:     "user",
		},
		"locked": {
			ID:       primitive.NewObjectID(),
			Email:    "locked@example.com",
			Password: "$2a$10$7juGndSkFAhhd/J8gmRBBOidUg69fbusoY8lQnLc4hqUXlvvkfz5G",
			Username: "locked_user",
			Role:     "user",
		},
		"oauth": {
			ID:       primitive.NewObjectID(),
			Email:    "oauth@example.com",
			Username: "oauth_user",
			Role:     "user",
		},
		"admin": {
			ID:       primitive.NewObjectID(),
			Email:    "admin@example.com",
			Password: "$2a$10$7juGndSkFAhhd/J8gmRBBOidUg69fbusoY8lQnLc4hqUXlvvkfz5G",
			Username: "admin_user",
			Role:     "admin",
		},
	}
)

// CreateTestUsers 创建测试用户数据
func CreateTestUsers(ctx context.Context, collection *mongo.Collection) error {
	now := time.Now()

	users := []interface{}{
		createVerifiedUser(DefaultUsers["verified"], now),
		createUnverifiedUser(DefaultUsers["unverified"], now),
		createLockedUser(DefaultUsers["locked"], now),
		createOAuthUser(DefaultUsers["oauth"], now),
		createAdminUser(DefaultUsers["admin"], now),
	}

	_, err := collection.InsertMany(ctx, users)
	return err
}

func createVerifiedUser(u TestUser, now time.Time) bson.M {
	return bson.M{
		"_id":      u.ID,
		"email":    u.Email,
		"username": u.Username,
		"password": u.Password,
		"status": bson.M{
			"emailVerified": true,
			"isLocked":      false,
		},
		"role": bson.M{
			"type": u.Role,
		},
		"stats": bson.M{
			"reviewCount":         0,
			"totalWords":          0,
			"violations":          0,
			"failedLoginAttempts": 0,
			"createdAt":           now,
			"lastLoginAt":         now,
		},
		"createdAt": now,
		"updatedAt": now,
	}
}

func createUnverifiedUser(u TestUser, now time.Time) bson.M {
	user := createVerifiedUser(u, now)
	user["status"] = bson.M{
		"emailVerified": false,
		"isLocked":      false,
		"verifyToken":   "test_verify_token",
		"tokenExpires":  now.Add(24 * time.Hour),
	}
	return user
}

func createLockedUser(u TestUser, now time.Time) bson.M {
	user := createVerifiedUser(u, now)
	user["status"] = bson.M{
		"emailVerified": true,
		"isLocked":      true,
		"lockReason":    "Too many login attempts",
		"lockExpires":   now.Add(24 * time.Hour),
	}
	user["stats"].(bson.M)["failedLoginAttempts"] = 5
	return user
}

func createOAuthUser(u TestUser, now time.Time) bson.M {
	user := createVerifiedUser(u, now)
	user["oauth"] = bson.M{
		"google": bson.M{
			"id":          "google_123456",
			"email":       u.Email,
			"connected":   true,
			"connectedAt": now,
		},
	}
	delete(user, "password") // OAuth用户没有密码
	return user
}

func createAdminUser(u TestUser, now time.Time) bson.M {
	user := createVerifiedUser(u, now)
	user["role"] = bson.M{
		"type": "admin",
	}
	return user
}

// NewTestDB 创建测试数据库连接
func NewTestDB() (*mongo.Database, error) {
	ctx := context.Background()
	// 只创建连接，不做其他操作
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	// 只返回数据库引用
	return client.Database("cpc_test"), nil
}

// InitTestData 初始化测试数据
func InitTestData(ctx context.Context, db *mongo.Database) error {
	// 1. 确保数据库为空
	if err := db.Drop(ctx); err != nil {
		return fmt.Errorf("failed to drop database: %v", err)
	}

	// 2. 创建 collection 和验证规则
	if err := createCollections(ctx, db); err != nil {
		return fmt.Errorf("failed to create collections: %v", err)
	}

	// 3. 创建用户数据
	collection := db.Collection("users")
	if err := CreateTestUsers(ctx, collection); err != nil {
		return fmt.Errorf("failed to create test users: %v", err)
	}

	return nil
}

// createCollections 创建测试所需的集合和验证规则
func createCollections(ctx context.Context, db *mongo.Database) error {
	// 创建users集合
	err := db.CreateCollection(ctx, "users", options.CreateCollection().SetValidator(bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{
				"username",
				"email",
				"role",
				"stats",
			},
			"properties": bson.M{
				"username": bson.M{
					"bsonType": "string",
				},
				"email": bson.M{
					"bsonType": "string",
				},
				"password": bson.M{
					"bsonType": "string",
				},
				"role": bson.M{
					"bsonType": "object",
					"required": []string{"type"},
					"properties": bson.M{
						"type": bson.M{
							"bsonType": "string",
							"enum":     []string{"user", "reviewer", "admin"},
						},
						"reviewerApplication": bson.M{
							"bsonType": "object",
							"properties": bson.M{
								"status": bson.M{
									"bsonType": "string",
								},
								"appliedAt": bson.M{
									"bsonType": "date",
								},
								"reviewCount": bson.M{
									"bsonType": "number",
								},
								"totalWords": bson.M{
									"bsonType": "number",
								},
							},
						},
						"inviteCode": bson.M{
							"bsonType": "string",
						},
					},
				},
				"stats": bson.M{
					"bsonType": "object",
					"required": []string{
						"reviewCount",
						"totalWords",
						"violations",
						"failedLoginAttempts",
						"createdAt",
						"lastLoginAt",
					},
					"properties": bson.M{
						"reviewCount": bson.M{
							"bsonType": "number",
						},
						"totalWords": bson.M{
							"bsonType": "number",
						},
						"violations": bson.M{
							"bsonType": "number",
						},
						"createdAt": bson.M{
							"bsonType": "date",
						},
						"lastLoginAt": bson.M{
							"bsonType": "date",
						},
						"failedLoginAttempts": bson.M{
							"bsonType": "int",
							"minimum":  0,
						},
					},
				},
				"oauth": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"google": bson.M{
							"bsonType": "object",
							"properties": bson.M{
								"id": bson.M{
									"bsonType": "string",
								},
								"email": bson.M{
									"bsonType": "string",
								},
								"connected": bson.M{
									"bsonType": "bool",
								},
								"connectedAt": bson.M{
									"bsonType": "date",
								},
							},
						},
					},
				},
				"status": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"emailVerified": bson.M{
							"bsonType": "bool",
						},
						"isLocked": bson.M{
							"bsonType": "bool",
						},
						"lockReason": bson.M{
							"bsonType": "string",
						},
						"lockExpires": bson.M{
							"bsonType": "date",
						},
						"verifyToken": bson.M{
							"bsonType": "string",
						},
						"tokenExpires": bson.M{
							"bsonType": "date",
						},
						"emailChange": bson.M{
							"bsonType": "string",
						},
					},
				},
			},
		},
	}))
	if err != nil {
		return fmt.Errorf("failed to create users collection: %v", err)
	}

	// 创建索引
	usersCollection := db.Collection("users")
	_, err = usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{"email", 1}},
			Options: options.Index().
				SetUnique(true),
		},
		{
			Keys: bson.D{{"oauth.google.id", 1}},
			Options: options.Index().
				SetUnique(true).
				SetSparse(true),
		},
		{
			Keys: bson.D{{"status.verifyToken", 1}},
		},
		{
			Keys: bson.D{{"createdAt", 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create indexes: %v", err)
	}

	return nil
}
