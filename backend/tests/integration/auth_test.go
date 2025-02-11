package integration

import (
	"FullStackOfYear/backend/internal/errors"
	authsvc "FullStackOfYear/backend/services/auth" // 服务
	"FullStackOfYear/backend/tests/mocks"
	"FullStackOfYear/backend/tests/testutil"
	authtypes "FullStackOfYear/backend/types/auth"    // 类型定义
	authclaims "FullStackOfYear/backend/types/claims" // claims包使用别名
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
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

	t.Run("邮箱验证-登录流程", func(t *testing.T) {
		ctx := context.Background()
		email := "newuser@example.com"
		password := "Password123!"
		deviceID := "test_device_123"

		// 登录请求
		loginReq := &authtypes.LoginRequest{
			LoginType: authtypes.EmailLogin,
			Email:     email,
			Password:  password,
			DeviceID:  deviceID,
		}

		// 设置预期行为
		mockTokenGen.EXPECT().
			GenerateTokenPair(gomock.Any(), "user", deviceID).
			Return("test_access_token", "test_refresh_token", nil)

		// 执行登录
		resp, err := authService.Login(ctx, loginReq)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
	})

	t.Run("OAuth登录流程", func(t *testing.T) {
		ctx := context.Background()
		userInfo := &authtypes.OAuthUserInfo{
			ID:    "google_123",
			Email: "google@example.com",
			Name:  "Google User",
		}

		// 设置OAuth相关的预期行为
		mockOAuthProvider.EXPECT().
			ExchangeCode(ctx, "test_code").
			Return(&authtypes.OAuthToken{
				AccessToken: "oauth_access_token",
			}, nil)

		mockOAuthProvider.EXPECT().
			GetUserInfo(ctx, gomock.Any()).
			Return(userInfo, nil)

		mockTokenGen.EXPECT().
			GenerateTokenPair(gomock.Any(), "user", "oauth_device").
			Return("access_token", "refresh_token", nil)

		// 执行OAuth登录
		resp, err := authService.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "test_state",
			DeviceID:  "oauth_device",
		})

		require.NoError(t, err)
		assert.True(t, resp.OAuthConnected)
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
		updatedUser, err := authService.GetUserByID(ctx, user.ID.Hex())
		require.NoError(t, err)
		assert.False(t, updatedUser.Status.EmailVerified)
	})

	t.Run("OAuth状态验证失败", func(t *testing.T) {
		ctx := context.Background()

		mockOAuthProvider.EXPECT().
			ExchangeCode(ctx, "test_code").
			Return(nil, errors.NewAppError(errors.BadRequest, "Invalid state"))

		loginReq := &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "invalid_state",
			DeviceID:  "test_device",
		}

		resp, err := authService.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))
	})

	t.Run("多次登录失败导致账户锁定", func(t *testing.T) {
		ctx := context.Background()
		user := testutil.GetTestUser(t, "verified")

		loginReq := &authtypes.LoginRequest{
			LoginType: authtypes.EmailLogin,
			Email:     user.Email,
			Password:  "wrong_password",
			DeviceID:  "test_device",
		}

		// 多次尝试错误密码
		for i := 0; i < 5; i++ {
			_, err := authService.Login(ctx, loginReq)
			assert.Error(t, err)
		}

		// 验证账户已锁定
		loginReq.Password = "correct_password"
		_, err := authService.Login(ctx, loginReq)
		assert.Error(t, err)
		assert.Equal(t, errors.Forbidden, errors.GetErrorCode(err))
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
		ctx := context.Background()

		// 创建用户信息时直接使用已存在的Google ID
		userInfo := &authtypes.OAuthUserInfo{
			ID:    "google_123456", // 假设这是已被使用的Google ID
			Email: "new_email@example.com",
			Name:  "New User",
		}

		mockOAuthProvider.EXPECT().
			ExchangeCode(ctx, "test_code").
			Return(&authtypes.OAuthToken{AccessToken: "oauth_token"}, nil)

		mockOAuthProvider.EXPECT().
			GetUserInfo(ctx, gomock.Any()).
			Return(userInfo, nil)

		// 尝试使用已绑定的Google ID创建新账号
		resp, err := authService.Login(ctx, &authtypes.LoginRequest{
			LoginType: authtypes.GoogleLogin,
			Code:      "test_code",
			State:     "test_state",
			DeviceID:  "oauth_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))
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

		// 设置预期行为
		mockTokenGen.EXPECT().
			ValidateRefreshToken(gomock.Any()).
			MinTimes(1).
			Return(validClaims, nil)

		mockTokenGen.EXPECT().
			GenerateTokenPair(validClaims.UserID, validClaims.Role, validClaims.DeviceID).
			MinTimes(1).
			Return("new_access_token", "new_refresh_token", nil)

		// 并发执行token刷新
		const goroutines = 3
		results := make(chan error, goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				_, _, err := authService.RefreshToken(ctx, "test_refresh_token")
				results <- err
			}()
		}

		// 收集并验证结果
		var successCount int
		for i := 0; i < goroutines; i++ {
			if err := <-results; err == nil {
				successCount++
			}
		}

		// 确保只有一次刷新成功
		assert.Equal(t, 1, successCount)
	})

	t.Run("并发邮箱验证", func(t *testing.T) {
		ctx := context.Background()
		user := testutil.GetTestUser(t, "unverified")
		verifyToken := "test_verify_token"

		// 并发执行邮箱验证
		const goroutines = 3
		results := make(chan error, goroutines)
		for i := 0; i < goroutines; i++ {
			go func() {
				err := authService.VerifyEmail(ctx, verifyToken)
				results <- err
			}()
		}

		// 收集并验证结果
		var successCount int
		for i := 0; i < goroutines; i++ {
			if err := <-results; err == nil {
				successCount++
			}
		}

		// 确保只有一次验证成功
		assert.Equal(t, 1, successCount)

		// 验证最终状态
		updatedUser, err := authService.GetUserByID(ctx, user.ID.Hex())
		require.NoError(t, err)
		assert.True(t, updatedUser.Status.EmailVerified)
	})
}
