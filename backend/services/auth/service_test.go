package auth

import (
	"project/backend/internal/errors"
	"project/backend/tests/mocks"
	"project/backend/tests/testutil"
	"project/backend/types/auth"
	"project/backend/types/claims"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestService_Login(t *testing.T) {
	// 设置 mock 控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 初始化测试数据库
	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	// 创建 mock
	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	// 创建服务实例
	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功登录已验证邮箱的用户", func(t *testing.T) {
		// 准备测试数据
		user := testutil.GetTestUser(t, "verified")
		deviceID := "test_device_123"

		// 设置token生成的预期行为
		mockTokenGen.EXPECT().
			GenerateTokenPair(user.ID.Hex(), user.Role, deviceID).
			Return("access_token", "refresh_token", nil)

		// 执行登录
		resp, err := svc.Login(context.Background(), &auth.LoginRequest{
			LoginType: auth.EmailLogin,
			Email:     user.Email,
			Password:  "password123",
			DeviceID:  deviceID,
		})

		// 验证结果
		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "access_token", resp.AccessToken)
		assert.Equal(t, "refresh_token", resp.RefreshToken)
		assert.Equal(t, user.ID.Hex(), resp.UserID)
		assert.Equal(t, user.Email, resp.Email)
	})

	t.Run("未验证邮箱的用户登录失败", func(t *testing.T) {
		user := testutil.GetTestUser(t, "unverified")

		resp, err := svc.Login(context.Background(), &auth.LoginRequest{
			LoginType: auth.EmailLogin,
			Email:     user.Email,
			Password:  "password123",
			DeviceID:  "test_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.Unauthorized, errors.GetErrorCode(err))
	})

	t.Run("锁定用户登录失败", func(t *testing.T) {
		user := testutil.GetTestUser(t, "locked")

		resp, err := svc.Login(context.Background(), &auth.LoginRequest{
			LoginType: auth.EmailLogin,
			Email:     user.Email,
			Password:  "password123",
			DeviceID:  "test_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.Forbidden, errors.GetErrorCode(err))
	})

	t.Run("密码错误登录失败", func(t *testing.T) {
		user := testutil.GetTestUser(t, "verified")

		resp, err := svc.Login(context.Background(), &auth.LoginRequest{
			LoginType: auth.EmailLogin,
			Email:     user.Email,
			Password:  "wrong_password",
			DeviceID:  "test_device",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errors.Unauthorized, errors.GetErrorCode(err))
	})
}

func TestService_HandleOAuthLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("OAuth用户首次登录成功", func(t *testing.T) {
		// 准备测试数据
		userInfo := &auth.OAuthUserInfo{
			ID:    "new_google_id",
			Email: "new_oauth@example.com",
			Name:  "New OAuth User",
		}

		deviceID := "oauth_device_123"
		expectedUserID := "" // 将在创建用户后获取

		// 设置token生成的预期行为
		mockTokenGen.EXPECT().
			GenerateTokenPair(gomock.Any(), "user", deviceID).
			DoAndReturn(func(userID, role, deviceID string) (string, string, error) {
				expectedUserID = userID
				return "access_token", "refresh_token", nil
			})

		// 执行OAuth登录
		resp, err := svc.HandleOAuthLogin(context.Background(), userInfo)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "access_token", resp.AccessToken)
		assert.Equal(t, "refresh_token", resp.RefreshToken)
		assert.Equal(t, expectedUserID, resp.UserID)
		assert.Equal(t, userInfo.Email, resp.Email)
		assert.True(t, resp.OAuthConnected)
		assert.Equal(t, "google", resp.OAuthProvider)
	})

	t.Run("已存在OAuth用户登录成功", func(t *testing.T) {
		user := testutil.GetTestUser(t, "oauth")
		userInfo := &auth.OAuthUserInfo{
			ID:    "google_123456",
			Email: user.Email,
			Name:  user.Username,
		}

		deviceID := "oauth_device_123"

		mockTokenGen.EXPECT().
			GenerateTokenPair(user.ID.Hex(), user.Role, deviceID).
			Return("access_token", "refresh_token", nil)

		resp, err := svc.HandleOAuthLogin(context.Background(), userInfo)

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, user.ID.Hex(), resp.UserID)
		assert.Equal(t, user.Email, resp.Email)
		assert.True(t, resp.OAuthConnected)
	})
}

func TestService_SendVerificationEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功发送验证邮件", func(t *testing.T) {
		user := testutil.GetTestUser(t, "unverified")

		mockEmailSender.EXPECT().
			SendVerificationEmail(user.Email, user.Username, gomock.Any()).
			Return(nil)

		err := svc.SendVerificationEmail(context.Background(), user.ID.Hex())
		require.NoError(t, err)

		// 验证数据库中的验证token已更新
		updatedUser, err := svc.GetUserByID(context.Background(), user.ID.Hex())
		require.NoError(t, err)
		assert.NotEmpty(t, updatedUser.Status.VerifyToken)
		assert.True(t, updatedUser.Status.TokenExpires.After(time.Now()))
	})

	t.Run("已验证邮箱的用户请求验证失败", func(t *testing.T) {
		user := testutil.GetTestUser(t, "verified")

		err := svc.SendVerificationEmail(context.Background(), user.ID.Hex())
		assert.Error(t, err)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))
	})
}

func TestService_VerifyEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功验证邮箱", func(t *testing.T) {
		user := testutil.GetTestUser(t, "unverified")

		err := svc.VerifyEmail(context.Background(), "test_verify_token")
		require.NoError(t, err)

		// 验证用户状态已更新
		updatedUser, err := svc.GetUserByID(context.Background(), user.ID.Hex())
		require.NoError(t, err)
		assert.True(t, updatedUser.Status.EmailVerified)
		assert.Empty(t, updatedUser.Status.VerifyToken)
		assert.True(t, updatedUser.Status.TokenExpires.IsZero()) // 检查是否重置了过期时间
	})

	t.Run("无效的验证token", func(t *testing.T) {
		err := svc.VerifyEmail(context.Background(), "invalid_token")
		assert.Error(t, err)
		assert.Equal(t, errors.BadRequest, errors.GetErrorCode(err))
	})
}

// services/auth/service_test.go

func TestService_GenerateTokenPair(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功生成Token对", func(t *testing.T) {
		userID := "test_user_id"
		role := "user"
		deviceID := "test_device_id"

		// 设置预期行为
		mockTokenGen.EXPECT().
			GenerateTokenPair(userID, role, deviceID).
			Return("access_token", "refresh_token", nil)

		// 执行测试
		accessToken, refreshToken, err := svc.GenerateTokenPair(context.Background(), userID, role, deviceID)

		// 验证结果
		require.NoError(t, err)
		assert.Equal(t, "access_token", accessToken)
		assert.Equal(t, "refresh_token", refreshToken)
	})
}

func TestService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功刷新Token", func(t *testing.T) {
		refreshToken := "valid_refresh_token"
		claims := &claims.Claims{
			UserID:   "test_user_id",
			Role:     "user",
			DeviceID: "test_device_id",
		}

		// 设置预期行为
		mockTokenGen.EXPECT().
			ValidateRefreshToken(refreshToken).
			Return(claims, nil)

		mockTokenGen.EXPECT().
			GenerateTokenPair(claims.UserID, claims.Role, claims.DeviceID).
			Return("new_access_token", "new_refresh_token", nil)

		// 执行测试
		newAccess, newRefresh, err := svc.RefreshToken(context.Background(), refreshToken)

		// 验证结果
		require.NoError(t, err)
		assert.Equal(t, "new_access_token", newAccess)
		assert.Equal(t, "new_refresh_token", newRefresh)
	})

	t.Run("无效的刷新Token", func(t *testing.T) {
		refreshToken := "invalid_refresh_token"

		// 设置预期行为
		mockTokenGen.EXPECT().
			ValidateRefreshToken(refreshToken).
			Return(nil, errors.NewAppError(errors.Unauthorized, "Invalid refresh token"))

		// 执行测试
		access, refresh, err := svc.RefreshToken(context.Background(), refreshToken)

		// 验证结果
		assert.Error(t, err)
		assert.Empty(t, access)
		assert.Empty(t, refresh)
		assert.Equal(t, errors.Unauthorized, errors.GetErrorCode(err))
	})
}

func TestService_RevokeTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, cleanup := testutil.SetupAuthTest(t)
	defer cleanup()

	mockTokenGen := mocks.NewMockTokenGenerator(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockOAuthProvider := mocks.NewMockOAuthProvider(ctrl)

	svc := NewService(
		db.Collection("users"),
		mockTokenGen,
		mockEmailSender,
		mockOAuthProvider,
	)

	t.Run("成功撤销Token", func(t *testing.T) {
		userID := "test_user_id"
		deviceID := "test_device_id"

		// 设置预期行为
		mockTokenGen.EXPECT().
			RevokeTokens(userID, deviceID).
			Return(nil)

		// 执行测试
		err := svc.RevokeTokens(context.Background(), userID, deviceID)

		// 验证结果
		require.NoError(t, err)
	})

	t.Run("撤销Token失败", func(t *testing.T) {
		userID := "test_user_id"
		deviceID := "test_device_id"

		// 设置预期行为
		mockTokenGen.EXPECT().
			RevokeTokens(userID, deviceID).
			Return(errors.NewAppError(errors.InternalError, "Failed to revoke tokens"))

		// 执行测试
		err := svc.RevokeTokens(context.Background(), userID, deviceID)

		// 验证结果
		assert.Error(t, err)
		assert.Equal(t, errors.InternalError, errors.GetErrorCode(err))
	})
}
