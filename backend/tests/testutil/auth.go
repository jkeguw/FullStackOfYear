package testutil

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

// SetupAuthTest 设置认证测试环境
func SetupAuthTest(t *testing.T) (*mongo.Database, func()) {
	// 1. 创建上下文（添加超时控制）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. 获取数据库连接
	db, err := NewTestDB()
	require.NoError(t, err, "Failed to create test database")

	// 3. 初始化测试数据
	err = InitTestData(ctx, db)
	require.NoError(t, err, "Failed to initialize test data")

	// 4. 验证数据初始化
	var count int64
	count, err = db.Collection("users").CountDocuments(ctx, bson.M{})
	require.NoError(t, err, "Failed to count users")
	require.Equal(t, int64(5), count, "Expected 5 test users to be created")

	// 5. 返回带超时控制的清理函数
	cleanup := func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		// 删除测试数据库
		err := db.Drop(cleanupCtx)
		require.NoError(t, err, "Failed to drop test database")

		// 关闭数据库连接
		err = db.Client().Disconnect(cleanupCtx)
		require.NoError(t, err, "Failed to disconnect from database")
	}

	return db, cleanup
}

// GetTestUser 获取测试用户数据
func GetTestUser(t *testing.T, userType string) TestUser {
	user, exists := DefaultUsers[userType]
	require.True(t, exists, "Test user %s not found", userType)
	return user
}

// 添加辅助函数，用于验证用户数据
func ValidateTestUser(t *testing.T, db *mongo.Database, userType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expectedUser := DefaultUsers[userType]
	var actualUser bson.M

	err := db.Collection("users").FindOne(ctx, bson.M{
		"email": expectedUser.Email,
	}).Decode(&actualUser)

	require.NoError(t, err)
	require.Equal(t, expectedUser.Email, actualUser["email"])
	require.Equal(t, expectedUser.Username, actualUser["username"])

	if userType == "verified" {
		require.Equal(t, true, actualUser["status"].(bson.M)["emailVerified"])
	}
}
