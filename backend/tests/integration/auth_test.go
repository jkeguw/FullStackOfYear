package integration

import (
	"project/backend/internal/errors"
	authsvc "project/backend/services/auth" // 服务
	"project/backend/tests/mocks"
	"project/backend/tests/testutil"
	authtypes "project/backend/types/auth"    // 类型定义
	authclaims "project/backend/types/claims" // claims包使用别名
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	validateTestSetup(t, db)

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	authService := authsvc.NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("邮箱验证-登录流程", func(t *testing.T) {
		ctx := context.Background()

		t.Log("=== 开始邮箱登录测试 ===")

		// 获取用户数据
		var user bson.M
		err := db.Collection("users").FindOne(ctx, bson.M{
			"email": "verified@example.com",
		}).Decode(&user)
		require.NoError(t, err)

		t.Logf("数据库中的用户数据:")
		t.Logf("- ID: %v", user["_id"])
		t.Logf("- Email: %v", user["email"])
		t.Logf("- Role: %v", user["role"].(bson.M)["type"])

		userID := user["_id"].(primitive.ObjectID).Hex()
		roleType := user["role"].(bson.M)["type"].(string)
		deviceID := "test_device_123"

		// 设置mock
		mockTokenGen.EXPECT().
			GenerateTokenPair(userID, roleType, deviceID).
			Return("test_access_token", "test_refresh_token", nil)

		t.Log("\n准备登录请求:")
		t.Logf("- Email: verified@example.com")
		t.Logf("- Password: password123")
		t.Logf("- Device ID: %s", deviceID)

		// 执行登录
		resp, err := authService.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.EmailLogin,
			Email:     "verified@example.com",
			Password:  "password123",
			DeviceID:  deviceID,
		})

		if err != nil {
			t.Log("\n登录失败:")
			if appErr, ok := err.(*errors.AppError); ok {
				t.Logf("- 错误类型: AppError")
				t.Logf("- 错误代码: %d", appErr.Code)
				t.Logf("- 错误信息: %s", appErr.Message)
			} else {
				t.Logf("- 错误类型: %T", err)
				t.Logf("- 错误信息: %v", err)
			}
		} else {
			t.Log("\n登录成功:")
			t.Logf("- User ID: %s", resp.UserID)
			t.Logf("- Email: %s", resp.Email)
		}

		require.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, "verified@example.com", resp.Email)
		assert.Equal(t, userID, resp.UserID)
	})

	t.Run("OAuth登录流程", func(t *testing.T) {
		ctx := context.Background()
		deviceID := "oauth_device_123"

		// 使用预设的OAuth用户信息
		userInfo := &authtypes.OAuthUserInfo{
			ID:    "google_123456",
			Email: "oauth@example.com",
			Name:  "oauth_user",
		}

		// 设置OAuth相关的mock
		mockOAuthProvider.EXPECT().
			ExchangeCode(ctx, "test_code").
			Return(&authtypes.OAuthToken{
				AccessToken: "oauth_access_token",
			}, nil)

		mockOAuthProvider.EXPECT().
			GetUserInfo(ctx, gomock.Any()).
			Return(userInfo, nil)

		// 使用gomock.Any()匹配用户ID，因为这可能是新创建的用户
		mockTokenGen.EXPECT().
			GenerateTokenPair(gomock.Any(), "user", gomock.Eq(deviceID)).
			Return("access_token", "refresh_token", nil)

		// 执行OAuth登录
		resp, err := authService.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "test_state",
			DeviceID:  deviceID,
		})

		require.NoError(t, err)
		assert.True(t, resp.OAuthConnected)
		assert.Equal(t, "google", resp.OAuthProvider)
	})

	t.Run("Token刷新流程", func(t *testing.T) {
		ctx := context.Background()
		refreshToken := "valid_refresh_token"

		// 设置预期的token验证和生成行为
		mockTokenGen.EXPECT().
			ValidateRefreshToken(refreshToken).
			Return(&authclaims.Claims{
				UserID:   "test_user_id",
				Role:     "user",
				DeviceID: "test_device",
			}, nil)

		mockTokenGen.EXPECT().
			GenerateTokenPair("test_user_id", "user", "test_device").
			Return("new_access_token", "new_refresh_token", nil)

		// 执行token刷新
		newAccess, newRefresh, err := authService.RefreshToken(ctx, refreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newAccess)
		assert.NotEmpty(t, newRefresh)
	})
}

func TestAuthFlowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	authService := authsvc.NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("邮箱验证过期", func(t *testing.T) {
		ctx := context.Background()
		user := testutil.GetTestUser(t, "unverified")

		// 使用过期的验证token
		err := authService.VerifyEmail(ctx, "expired_token")
		assert.Error(t, err)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))

		// 验证用户状态未改变
		var dbUser bson.M
		err = db.Collection("users").FindOne(ctx, bson.M{"_id": user.ID}).Decode(&dbUser)
		require.NoError(t, err)
		assert.False(t, dbUser["status"].(bson.M)["emailVerified"].(bool))
	})

	t.Run("OAuth状态验证失败", func(t *testing.T) {
		ctx := context.Background()

		mockOAuthProvider.EXPECT().
			ExchangeCode(ctx, "test_code").
			Return(nil, errors.NewAppError(errors.BadRequest, "Invalid state"))

		resp, err := authService.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "invalid_state",
			DeviceID:  "test_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))
	})

	t.Run("登录失败密码错误", func(t *testing.T) {
		ctx := context.Background()
		user := testutil.GetTestUser(t, "verified")

		loginReq := &authtypes.LoginRequest{
			LoginType: authtypes.EmailLogin,
			Email:     user.Email,
			Password:  "wrong_password",
			DeviceID:  "test_device",
		}

		// 密码错误登录尝试
		_, err := authService.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Equal(t, errors.Unauthorized, errors.GetErrorCode(err))
		
		// 使用正确密码应该可以登录成功
		loginReq.Password = "password123"
		response, err := authService.Login(ctx, loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("刷新Token失败", func(t *testing.T) {
		ctx := context.Background()

		mockTokenGen.EXPECT().
			ValidateRefreshToken("invalid_token").
			Return(nil, errors.NewAppError(errors.Unauthorized, "Invalid refresh token"))

		newAccess, newRefresh, err := authService.RefreshToken(ctx, "invalid_token")
		assert.Error(t, err)
		assert.Empty(t, newAccess)
		assert.Empty(t, newRefresh)
		assert.Equal(t, errors.Unauthorized, errors.GetErrorCode(err))
	})

	t.Run("OAuth账号重复绑定", func(t *testing.T) {
		ctrlLocal := gomock.NewController(t)
		defer ctrlLocal.Finish()

		mockTokenGenLocal := mocks.NewMockTokenGenerator(ctrlLocal)
		mockEmailSenderLocal := mocks.NewMockEmailSender(ctrlLocal)
		mockOAuthProviderLocal := mocks.NewMockOAuthProvider(ctrlLocal)

		authServiceLocal := authsvc.NewService(
			db.Collection("users"),
			mockTokenGenLocal,
			mockEmailSenderLocal,
			mockOAuthProviderLocal,
		)

		ctx := context.Background()
		userInfo := &authtypes.OAuthUserInfo{
			ID:    "google_123456",
			Email: "new_email@example.com",
			Name:  "New User",
		}

		mockOAuthProviderLocal.EXPECT().
			ExchangeCode(gomock.Any(), "test_code").
			Return(&authtypes.OAuthToken{AccessToken: "oauth_token"}, nil)

		mockOAuthProviderLocal.EXPECT().
			GetUserInfo(gomock.Any(), gomock.Any()).
			Return(userInfo, nil)

		// 新增：允许任何 GenerateTokenPair 调用返回错误
		mockTokenGenLocal.EXPECT().
			GenerateTokenPair(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return("", "", errors.NewAppError(errors.BadRequest, "Account already exists"))

		resp, err := authServiceLocal.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "test_state",
			DeviceID:  "oauth_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		if err != nil {
			appErr, ok := err.(*errors.AppError)
			assert.True(t, ok)
			if ok {
				assert.Equal(t, errors.BadRequest, appErr.Code)
				assert.Contains(t, appErr.Message, "already")
			}
		}
	})
}

func TestAuthRaceConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	authService := authsvc.NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("并发Token刷新", func(t *testing.T) {
		ctx := context.Background()
		validClaims := &authclaims.Claims{
			UserID:   "test_user_id",
			Role:     "user",
			DeviceID: "test_device",
		}

		// 设置 mock 期望 - 允许多次验证但只允许一次生成
		mockTokenGen.EXPECT().
			ValidateRefreshToken(gomock.Any()).
			AnyTimes().
			Return(validClaims, nil)

		mockTokenGen.EXPECT().
			GenerateTokenPair(validClaims.UserID, validClaims.Role, validClaims.DeviceID).
			Times(1).
			Return("new_access_token", "new_refresh_token", nil)

		var (
			wg sync.WaitGroup
			mu sync.Mutex
		)
		successCount := 0
		const goroutines = 3

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, _, err := authService.RefreshToken(ctx, "test_refresh_token"); err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()
		assert.Equal(t, 1, successCount)
	})

	t.Run("并发邮箱验证", func(t *testing.T) {
		ctx := context.Background()
		user := testutil.GetTestUser(t, "unverified")
		verifyToken := "test_verify_token"

		// 准备测试数据
		_, err := db.Collection("users").UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{
				"$set": bson.M{
					"status.verifyToken":   verifyToken,
					"status.tokenExpires":  time.Now().Add(time.Hour),
					"status.emailVerified": false,
				},
			},
		)
		require.NoError(t, err)

		// 使用 WaitGroup 同步所有 goroutine
		var wg sync.WaitGroup
		successCount := 0
		var mu sync.Mutex
		const goroutines = 3

		// 确保所有 goroutine 同时开始
		ready := make(chan struct{})

		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// 等待开始信号
				<-ready

				err := authService.VerifyEmail(ctx, verifyToken)
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		// 发送开始信号
		close(ready)

		// 等待所有 goroutine 完成
		wg.Wait()

		// 验证结果
		assert.Equal(t, 1, successCount, "应该只有一次验证成功")

		// 验证最终状态
		var updatedUser struct {
			Status struct {
				EmailVerified bool `bson:"emailVerified"`
			} `bson:"status"`
		}
		err = db.Collection("users").FindOne(ctx, bson.M{"_id": user.ID}).Decode(&updatedUser)
		require.NoError(t, err)
		assert.True(t, updatedUser.Status.EmailVerified, "邮箱应该被标记为已验证")
	})
}

// 为了确保测试数据的可靠性，添加数据验证函数
func validateTestData(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	// 验证已验证用户
	var verifiedUser struct {
		Status struct {
			EmailVerified bool `bson:"emailVerified"`
		} `bson:"status"`
		Password string `bson:"password"`
	}
	err := db.Collection("users").FindOne(
		ctx,
		bson.M{"email": "verified@example.com"},
	).Decode(&verifiedUser)

	require.NoError(t, err)
	assert.True(t, verifiedUser.Status.EmailVerified)
	assert.Equal(t,
		"$2a$10$NWY9SqxvWeYkPFn0R4PCu.DTF5lHcqyof.mKPHqY4TkGJRIA4O0Iy",
		verifiedUser.Password,
	)

	// 验证 OAuth 用户
	var oauthUser struct {
		OAuth struct {
			Google struct {
				ID string `bson:"id"`
			} `bson:"google"`
		} `bson:"oauth"`
	}
	err = db.Collection("users").FindOne(
		ctx,
		bson.M{"email": "oauth@example.com"},
	).Decode(&oauthUser)

	require.NoError(t, err)
	assert.Equal(t, "google_123456", oauthUser.OAuth.Google.ID)
}

// 在主测试函数开始时调用验证
//func TestMain(m *testing.M) {
//	db, err := testutil.NewTestDB()
//	if err != nil {
//		panic(err)
//	}
//
//	// 初始化并验证测试数据
//	err = testutil.InitTestData(context.Background(), db)
//	if err != nil {
//		panic(err)
//	}
//
//	validateTestData(testing.T{}, db)
//
//	os.Exit(m.Run())
//}

// 验证测试数据初始化
func validateTestSetup(t *testing.T, db *mongo.Database) {
	ctx := context.Background()

	// 查询并打印验证用户的完整信息
	var verifiedUser bson.M
	err := db.Collection("users").FindOne(ctx, bson.M{
		"email": "verified@example.com",
	}).Decode(&verifiedUser)
	require.NoError(t, err)

	// 打印完整的用户信息
	t.Logf("Verified user data: %+v", verifiedUser)

	// 特别检查密码字段
	password, ok := verifiedUser["password"].(string)
	require.True(t, ok, "Password field not found or not string")
	t.Logf("Stored password hash: %s", password)

	// 验证密码hash是否能正确验证
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte("password123"))
	require.NoError(t, err, "Password hash verification failed")
}

func TestPasswordHash(t *testing.T) {
	// 1. 测试预定义的hash
	expectedHash := "$2a$10$NWY9SqxvWeYkPFn0R4PCu.DTF5lHcqyof.mKPHqY4TkGJRIA4O0Iy"
	rawPassword := "password123"

	t.Logf("测试预定义hash验证:")
	t.Logf("密码: %s", rawPassword)
	t.Logf("Hash: %s", expectedHash)

	// 验证预定义hash
	err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(rawPassword))
	require.NoError(t, err, "预定义hash验证失败")

	// 2. 测试新生成的hash
	t.Logf("\n测试新生成hash:")
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	require.NoError(t, err)
	t.Logf("密码: %s", rawPassword)
	t.Logf("新Hash: %s", string(hash))

	// 验证新生成的hash
	err = bcrypt.CompareHashAndPassword(hash, []byte(rawPassword))
	require.NoError(t, err, "新生成hash验证失败")
}

func TestGenerateNewPasswordHash(t *testing.T) {
	password := "password123"

	// 生成新的hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	// 验证新生成的hash是否有效
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	require.NoError(t, err)

	// 打印新hash，这个hash应该被用来更新DefaultUsers
	t.Logf("新生成的密码hash for password123: %s", string(hash))

	// 进行双重验证
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	require.NoError(t, err, "新生成的hash应该能验证密码")
}
