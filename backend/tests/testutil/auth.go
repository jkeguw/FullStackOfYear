package testutil

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

// SetupAuthTest 设置认证测试环境
func SetupAuthTest(t *testing.T) (*mongo.Database, func()) {
	db, err := NewTestDB()
	require.NoError(t, err)

	err = InitTestData(context.Background(), db)
	require.NoError(t, err)

	cleanup := func() {
		err := db.Drop(context.Background())
		require.NoError(t, err)
	}

	return db, cleanup
}

// GetTestUser 获取测试用户数据
func GetTestUser(t *testing.T, userType string) TestUser {
	user, exists := DefaultUsers[userType]
	require.True(t, exists, "Test user %s not found", userType)
	return user
}
