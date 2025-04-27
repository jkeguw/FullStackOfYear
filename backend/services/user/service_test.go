package user

import (
	"project/backend/internal/errors"
	"project/backend/types/user"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

// 创建一个模拟服务
type mockUserService struct {
	mock.Mock
}

// 实现Service接口的方法
func (m *mockUserService) GetUserProfile(ctx context.Context, userID string) (*user.ProfileResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.ProfileResponse), args.Error(1)
}

func (m *mockUserService) UpdateUserProfile(ctx context.Context, userID string, profileData user.UpdateProfileRequest) (*user.ProfileResponse, error) {
	args := m.Called(ctx, userID, profileData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.ProfileResponse), args.Error(1)
}

func (m *mockUserService) GetUserSettings(ctx context.Context, userID string) (*user.SettingsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.SettingsResponse), args.Error(1)
}

func (m *mockUserService) UpdateUserSettings(ctx context.Context, userID string, settingsData user.UpdateSettingsRequest) (*user.SettingsResponse, error) {
	args := m.Called(ctx, userID, settingsData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.SettingsResponse), args.Error(1)
}

func (m *mockUserService) GetGameSettings(ctx context.Context, userID string) (*user.GameSettingsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.GameSettingsResponse), args.Error(1)
}

func (m *mockUserService) UpdateGameSettings(ctx context.Context, userID string, settings user.UpdateGameSettingsRequest) (*user.GameSettingsResponse, error) {
	args := m.Called(ctx, userID, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.GameSettingsResponse), args.Error(1)
}

func (m *mockUserService) AddSensitivityConfig(ctx context.Context, userID string, config user.AddSensitivityConfigRequest) (*user.GameSettingsResponse, error) {
	args := m.Called(ctx, userID, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.GameSettingsResponse), args.Error(1)
}

func (m *mockUserService) UpdateSensitivityConfig(ctx context.Context, userID string, game string, config user.UpdateSensitivityConfigRequest) (*user.GameSettingsResponse, error) {
	args := m.Called(ctx, userID, game, config)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.GameSettingsResponse), args.Error(1)
}

func (m *mockUserService) DeleteSensitivityConfig(ctx context.Context, userID string, game string) error {
	args := m.Called(ctx, userID, game)
	return args.Error(0)
}

func (m *mockUserService) GetPrivacySettings(ctx context.Context, userID string) (*user.PrivacySettingsResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PrivacySettingsResponse), args.Error(1)
}

func (m *mockUserService) UpdatePrivacySettings(ctx context.Context, userID string, settings user.UpdatePrivacySettingsRequest) (*user.PrivacySettingsResponse, error) {
	args := m.Called(ctx, userID, settings)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PrivacySettingsResponse), args.Error(1)
}

func (m *mockUserService) GetDevicePreferences(ctx context.Context, userID string) ([]user.DevicePreferenceResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.DevicePreferenceResponse), args.Error(1)
}

func (m *mockUserService) AddDevicePreference(ctx context.Context, userID string, pref user.AddDevicePreferenceRequest) (*user.DevicePreferenceResponse, error) {
	args := m.Called(ctx, userID, pref)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.DevicePreferenceResponse), args.Error(1)
}

func (m *mockUserService) UpdateDevicePreference(ctx context.Context, userID string, devicePrefID string, pref user.UpdateDevicePreferenceRequest) (*user.DevicePreferenceResponse, error) {
	args := m.Called(ctx, userID, devicePrefID, pref)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.DevicePreferenceResponse), args.Error(1)
}

func (m *mockUserService) DeleteDevicePreference(ctx context.Context, userID string, devicePrefID string) error {
	args := m.Called(ctx, userID, devicePrefID)
	return args.Error(0)
}

func (m *mockUserService) GetUserPreferences(ctx context.Context, userID string) (*user.PreferencesResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PreferencesResponse), args.Error(1)
}

func (m *mockUserService) UpdateUserPreferences(ctx context.Context, userID string, preferencesData user.UpdatePreferencesRequest) (*user.PreferencesResponse, error) {
	args := m.Called(ctx, userID, preferencesData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.PreferencesResponse), args.Error(1)
}

func (m *mockUserService) ApplyForReviewer(ctx context.Context, userID string, application user.ReviewerApplicationRequest) error {
	args := m.Called(ctx, userID, application)
	return args.Error(0)
}

func (m *mockUserService) GetReviewerApplications(ctx context.Context, status string, page, pageSize int) (*user.ReviewerApplicationListResponse, int64, error) {
	args := m.Called(ctx, status, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).(*user.ReviewerApplicationListResponse), args.Get(1).(int64), args.Error(2)
}

func (m *mockUserService) ApproveReviewerApplication(ctx context.Context, userID string, adminID string) error {
	args := m.Called(ctx, userID, adminID)
	return args.Error(0)
}

func (m *mockUserService) RejectReviewerApplication(ctx context.Context, userID string, adminID string, reason string) error {
	args := m.Called(ctx, userID, adminID, reason)
	return args.Error(0)
}

// TestGetUserProfile_Success 测试正常获取用户档案
func TestGetUserProfile_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建期望的返回值
	expectedProfile := &user.ProfileResponse{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Bio:         "This is a test user",
		Verified:    true,
		JoinedAt:    time.Now(),
	}
	
	// 设置Mock期望
	mockSvc.On("GetUserProfile", ctx, userID).Return(expectedProfile, nil)
	
	// 调用函数
	result, err := mockSvc.GetUserProfile(ctx, userID)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedProfile, result)
	mockSvc.AssertExpectations(t)
}

// TestGetUserProfile_InvalidID 测试无效ID时的获取用户档案
func TestGetUserProfile_InvalidID(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	invalidID := "invalid-id"
	
	// 设置Mock期望
	expectedErr := errors.NewBadRequestError("无效的用户ID")
	mockSvc.On("GetUserProfile", ctx, invalidID).Return(nil, expectedErr)
	
	// 调用函数
	result, err := mockSvc.GetUserProfile(ctx, invalidID)
	
	// 验证结果
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockSvc.AssertExpectations(t)
}

// TestGetUserProfile_UserNotFound 测试用户不存在时的获取用户档案
func TestGetUserProfile_UserNotFound(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	nonExistentID := primitive.NewObjectID().Hex()
	
	// 设置Mock期望
	expectedErr := errors.NewNotFoundError("用户不存在")
	mockSvc.On("GetUserProfile", ctx, nonExistentID).Return(nil, expectedErr)
	
	// 调用函数
	result, err := mockSvc.GetUserProfile(ctx, nonExistentID)
	
	// 验证结果
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	mockSvc.AssertExpectations(t)
}

// TestUpdateUserProfile_Success 测试成功更新用户档案
func TestUpdateUserProfile_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	displayName := "Updated Name"
	bio := "Updated bio information"
	updateRequest := user.UpdateProfileRequest{
		DisplayName: &displayName,
		Bio:         &bio,
	}
	
	// 创建期望的返回值
	expectedProfile := &user.ProfileResponse{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: displayName,
		Bio:         bio,
		Verified:    true,
		JoinedAt:    time.Now(),
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateUserProfile", ctx, userID, updateRequest).Return(expectedProfile, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateUserProfile(ctx, userID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedProfile, result)
	assert.Equal(t, displayName, result.DisplayName)
	assert.Equal(t, bio, result.Bio)
	mockSvc.AssertExpectations(t)
}

// TestUpdateUserSettings_Success 测试成功更新用户设置
func TestUpdateUserSettings_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	language := "zh-CN"
	theme := "dark"
	unit := "mm"
	updateRequest := user.UpdateSettingsRequest{
		Language:        &language,
		Theme:           &theme,
		MeasurementUnit: &unit,
		NotificationPref: map[string]bool{
			"email_notifications": true,
			"push_notifications":  false,
		},
	}
	
	// 创建期望的返回值
	expectedSettings := &user.SettingsResponse{
		Language:        language,
		Theme:           theme,
		MeasurementUnit: unit,
		NotificationPref: map[string]bool{
			"email_notifications": true,
			"push_notifications":  false,
		},
		PrivacySettings:  map[string]string{},
		SecuritySettings: map[string]interface{}{},
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateUserSettings", ctx, userID, updateRequest).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateUserSettings(ctx, userID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, language, result.Language)
	assert.Equal(t, theme, result.Theme)
	assert.Equal(t, unit, result.MeasurementUnit)
	assert.Equal(t, true, result.NotificationPref["email_notifications"])
	assert.Equal(t, false, result.NotificationPref["push_notifications"])
	mockSvc.AssertExpectations(t)
}

// TestApplyForReviewer_Success 测试成功申请评测员
func TestApplyForReviewer_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建申请请求
	application := user.ReviewerApplicationRequest{
		Experience:     "5 years of experience with gaming peripherals and extensive testing background.",
		ExpertiseAreas: []string{"gaming mice", "gaming keyboards", "mouse sensors"},
		Samples:        []string{"https://example.com/review1", "https://example.com/review2"},
		Motivation:     "I want to contribute to the community by providing detailed and objective reviews of gaming peripherals.",
	}
	
	// 设置Mock期望
	mockSvc.On("ApplyForReviewer", ctx, userID, application).Return(nil)
	
	// 调用函数
	err := mockSvc.ApplyForReviewer(ctx, userID, application)
	
	// 验证结果
	require.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestGetGameSettings_Success 测试成功获取游戏设置
func TestGetGameSettings_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建期望的返回值
	expectedSettings := &user.GameSettingsResponse{
		PreferredGames:     []string{"CS2", "Valorant", "Apex Legends"},
		DefaultDPI:         800,
		PreferredGripStyle: "claw",
		MouseAcceleration:  false,
		PollRate:           1000,
		SensitivityConfigs: []user.SensitivityConfigResponse{
			{
				Game:        "CS2",
				Sensitivity: 1.2,
				DPI:         800,
				IsActive:    true,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	
	// 设置Mock期望
	mockSvc.On("GetGameSettings", ctx, userID).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.GetGameSettings(ctx, userID)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, 800, result.DefaultDPI)
	assert.Equal(t, "claw", result.PreferredGripStyle)
	assert.Equal(t, 1000, result.PollRate)
	assert.Equal(t, 1, len(result.SensitivityConfigs))
	assert.Equal(t, "CS2", result.SensitivityConfigs[0].Game)
	mockSvc.AssertExpectations(t)
}

// TestUpdateGameSettings_Success 测试成功更新游戏设置
func TestUpdateGameSettings_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	dpi := 1600
	gripStyle := "palm"
	acceleration := true
	pollRate := 500
	preferredGames := []string{"CS2", "Valorant", "Overwatch"}
	updateRequest := user.UpdateGameSettingsRequest{
		DefaultDPI:         &dpi,
		PreferredGripStyle: &gripStyle,
		MouseAcceleration:  &acceleration,
		PollRate:           &pollRate,
		PreferredGames:     &preferredGames,
	}
	
	// 创建期望的返回值
	expectedSettings := &user.GameSettingsResponse{
		DefaultDPI:         dpi,
		PreferredGripStyle: gripStyle,
		MouseAcceleration:  acceleration,
		PollRate:           pollRate,
		PreferredGames:     preferredGames,
		SensitivityConfigs: []user.SensitivityConfigResponse{},
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateGameSettings", ctx, userID, updateRequest).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateGameSettings(ctx, userID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, dpi, result.DefaultDPI)
	assert.Equal(t, gripStyle, result.PreferredGripStyle)
	assert.Equal(t, acceleration, result.MouseAcceleration)
	assert.Equal(t, pollRate, result.PollRate)
	assert.Equal(t, preferredGames, result.PreferredGames)
	mockSvc.AssertExpectations(t)
}

// TestAddSensitivityConfig_Success 测试成功添加灵敏度配置
func TestAddSensitivityConfig_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建添加请求
	addRequest := user.AddSensitivityConfigRequest{
		Game:        "Valorant",
		Sensitivity: 0.5,
		DPI:         800,
		IsActive:    true,
	}
	
	// 创建期望的返回值
	now := time.Now()
	expectedSettings := &user.GameSettingsResponse{
		DefaultDPI:         800,
		PreferredGripStyle: "claw",
		MouseAcceleration:  false,
		PollRate:           1000,
		SensitivityConfigs: []user.SensitivityConfigResponse{
			{
				Game:        "Valorant",
				Sensitivity: 0.5,
				DPI:         800,
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	
	// 设置Mock期望
	mockSvc.On("AddSensitivityConfig", ctx, userID, addRequest).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.AddSensitivityConfig(ctx, userID, addRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, 1, len(result.SensitivityConfigs))
	assert.Equal(t, "Valorant", result.SensitivityConfigs[0].Game)
	assert.Equal(t, 0.5, result.SensitivityConfigs[0].Sensitivity)
	assert.Equal(t, 800, result.SensitivityConfigs[0].DPI)
	assert.Equal(t, true, result.SensitivityConfigs[0].IsActive)
	mockSvc.AssertExpectations(t)
}

// TestUpdateSensitivityConfig_Success 测试成功更新灵敏度配置
func TestUpdateSensitivityConfig_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	game := "Valorant"
	
	// 创建更新请求
	sensitivity := 0.7
	dpi := 1600
	isActive := false
	updateRequest := user.UpdateSensitivityConfigRequest{
		Sensitivity: &sensitivity,
		DPI:         &dpi,
		IsActive:    &isActive,
	}
	
	// 创建期望的返回值
	now := time.Now()
	expectedSettings := &user.GameSettingsResponse{
		DefaultDPI:         800,
		PreferredGripStyle: "claw",
		MouseAcceleration:  false,
		PollRate:           1000,
		SensitivityConfigs: []user.SensitivityConfigResponse{
			{
				Game:        game,
				Sensitivity: sensitivity,
				DPI:         dpi,
				IsActive:    isActive,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateSensitivityConfig", ctx, userID, game, updateRequest).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateSensitivityConfig(ctx, userID, game, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, 1, len(result.SensitivityConfigs))
	assert.Equal(t, game, result.SensitivityConfigs[0].Game)
	assert.Equal(t, sensitivity, result.SensitivityConfigs[0].Sensitivity)
	assert.Equal(t, dpi, result.SensitivityConfigs[0].DPI)
	assert.Equal(t, isActive, result.SensitivityConfigs[0].IsActive)
	mockSvc.AssertExpectations(t)
}

// TestDeleteSensitivityConfig_Success 测试成功删除灵敏度配置
func TestDeleteSensitivityConfig_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	game := "Valorant"
	
	// 设置Mock期望
	mockSvc.On("DeleteSensitivityConfig", ctx, userID, game).Return(nil)
	
	// 调用函数
	err := mockSvc.DeleteSensitivityConfig(ctx, userID, game)
	
	// 验证结果
	require.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestGetPrivacySettings_Success 测试成功获取隐私设置
func TestGetPrivacySettings_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建期望的返回值
	expectedSettings := &user.PrivacySettingsResponse{
		ProfileVisibility:       "public",
		DeviceListVisibility:    "friends",
		ReviewHistoryVisibility: "private",
		ShowOnlineStatus:        true,
		ShowActivity:            false,
	}
	
	// 设置Mock期望
	mockSvc.On("GetPrivacySettings", ctx, userID).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.GetPrivacySettings(ctx, userID)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, "public", result.ProfileVisibility)
	assert.Equal(t, "friends", result.DeviceListVisibility)
	assert.Equal(t, "private", result.ReviewHistoryVisibility)
	assert.Equal(t, true, result.ShowOnlineStatus)
	assert.Equal(t, false, result.ShowActivity)
	mockSvc.AssertExpectations(t)
}

// TestUpdatePrivacySettings_Success 测试成功更新隐私设置
func TestUpdatePrivacySettings_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	profileVisibility := "friends"
	deviceListVisibility := "private"
	reviewHistoryVisibility := "public"
	showOnlineStatus := false
	showActivity := true
	updateRequest := user.UpdatePrivacySettingsRequest{
		ProfileVisibility:       &profileVisibility,
		DeviceListVisibility:    &deviceListVisibility,
		ReviewHistoryVisibility: &reviewHistoryVisibility,
		ShowOnlineStatus:        &showOnlineStatus,
		ShowActivity:            &showActivity,
	}
	
	// 创建期望的返回值
	expectedSettings := &user.PrivacySettingsResponse{
		ProfileVisibility:       profileVisibility,
		DeviceListVisibility:    deviceListVisibility,
		ReviewHistoryVisibility: reviewHistoryVisibility,
		ShowOnlineStatus:        showOnlineStatus,
		ShowActivity:            showActivity,
	}
	
	// 设置Mock期望
	mockSvc.On("UpdatePrivacySettings", ctx, userID, updateRequest).Return(expectedSettings, nil)
	
	// 调用函数
	result, err := mockSvc.UpdatePrivacySettings(ctx, userID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedSettings, result)
	assert.Equal(t, profileVisibility, result.ProfileVisibility)
	assert.Equal(t, deviceListVisibility, result.DeviceListVisibility)
	assert.Equal(t, reviewHistoryVisibility, result.ReviewHistoryVisibility)
	assert.Equal(t, showOnlineStatus, result.ShowOnlineStatus)
	assert.Equal(t, showActivity, result.ShowActivity)
	mockSvc.AssertExpectations(t)
}

// TestGetDevicePreferences_Success 测试成功获取设备偏好
func TestGetDevicePreferences_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建期望的返回值
	deviceID := primitive.NewObjectID().Hex()
	now := time.Now()
	expectedPreferences := []user.DevicePreferenceResponse{
		{
			ID:          primitive.NewObjectID().Hex(),
			DeviceID:    deviceID,
			DeviceType:  "mouse",
			DeviceName:  "Logitech G Pro X Superlight",
			DeviceBrand: "Logitech",
			IsFavorite:  true,
			IsWishlist:  false,
			Rating:      5,
			Notes:       "Best mouse I've ever used",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	
	// 设置Mock期望
	mockSvc.On("GetDevicePreferences", ctx, userID).Return(expectedPreferences, nil)
	
	// 调用函数
	result, err := mockSvc.GetDevicePreferences(ctx, userID)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedPreferences, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, deviceID, result[0].DeviceID)
	assert.Equal(t, "mouse", result[0].DeviceType)
	assert.Equal(t, true, result[0].IsFavorite)
	assert.Equal(t, 5, result[0].Rating)
	mockSvc.AssertExpectations(t)
}

// TestAddDevicePreference_Success 测试成功添加设备偏好
func TestAddDevicePreference_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	deviceID := primitive.NewObjectID().Hex()
	
	// 创建添加请求
	rating := 5
	addRequest := user.AddDevicePreferenceRequest{
		DeviceID:   deviceID,
		DeviceType: "mouse",
		IsFavorite: true,
		IsWishlist: false,
		Rating:     &rating,
		Notes:      "Excellent wireless gaming mouse",
	}
	
	// 创建期望的返回值
	now := time.Now()
	expectedResponse := &user.DevicePreferenceResponse{
		ID:          primitive.NewObjectID().Hex(),
		DeviceID:    deviceID,
		DeviceType:  "mouse",
		DeviceName:  "Logitech G Pro X Superlight",
		DeviceBrand: "Logitech",
		IsFavorite:  true,
		IsWishlist:  false,
		Rating:      rating,
		Notes:       "Excellent wireless gaming mouse",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// 设置Mock期望
	mockSvc.On("AddDevicePreference", ctx, userID, addRequest).Return(expectedResponse, nil)
	
	// 调用函数
	result, err := mockSvc.AddDevicePreference(ctx, userID, addRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Equal(t, deviceID, result.DeviceID)
	assert.Equal(t, "mouse", result.DeviceType)
	assert.Equal(t, true, result.IsFavorite)
	assert.Equal(t, rating, result.Rating)
	assert.Equal(t, "Excellent wireless gaming mouse", result.Notes)
	mockSvc.AssertExpectations(t)
}

// TestUpdateDevicePreference_Success 测试成功更新设备偏好
func TestUpdateDevicePreference_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	devicePrefID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	isFavorite := false
	isWishlist := true
	rating := 4
	notes := "Updated notes for this device"
	updateRequest := user.UpdateDevicePreferenceRequest{
		IsFavorite: &isFavorite,
		IsWishlist: &isWishlist,
		Rating:     &rating,
		Notes:      &notes,
	}
	
	// 创建期望的返回值
	now := time.Now()
	deviceID := primitive.NewObjectID().Hex()
	expectedResponse := &user.DevicePreferenceResponse{
		ID:          devicePrefID,
		DeviceID:    deviceID,
		DeviceType:  "mouse",
		DeviceName:  "Logitech G Pro X Superlight",
		DeviceBrand: "Logitech",
		IsFavorite:  isFavorite,
		IsWishlist:  isWishlist,
		Rating:      rating,
		Notes:       notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateDevicePreference", ctx, userID, devicePrefID, updateRequest).Return(expectedResponse, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateDevicePreference(ctx, userID, devicePrefID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Equal(t, isFavorite, result.IsFavorite)
	assert.Equal(t, isWishlist, result.IsWishlist)
	assert.Equal(t, rating, result.Rating)
	assert.Equal(t, notes, result.Notes)
	mockSvc.AssertExpectations(t)
}

// TestDeleteDevicePreference_Success 测试成功删除设备偏好
func TestDeleteDevicePreference_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	devicePrefID := primitive.NewObjectID().Hex()
	
	// 设置Mock期望
	mockSvc.On("DeleteDevicePreference", ctx, userID, devicePrefID).Return(nil)
	
	// 调用函数
	err := mockSvc.DeleteDevicePreference(ctx, userID, devicePrefID)
	
	// 验证结果
	require.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestGetUserPreferences_Success 测试成功获取用户偏好
func TestGetUserPreferences_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建期望的返回值
	expectedPreferences := &user.PreferencesResponse{
		GameSettings: map[string]interface{}{
			"preferredGenres": []string{"FPS", "MOBA"},
			"favoriteGames":   []string{"CS2", "Valorant"},
		},
		PeripheralPrefs: map[string]interface{}{
			"preferredBrands": []string{"Logitech", "Zowie"},
			"mouseWeight":     "light",
		},
		ContentPrefs: map[string]interface{}{
			"showReviews":   true,
			"showTutorials": true,
		},
		RecommendedForYou: map[string]interface{}{
			"enableRecommendations": true,
			"basedOn":               "browser_history",
		},
	}
	
	// 设置Mock期望
	mockSvc.On("GetUserPreferences", ctx, userID).Return(expectedPreferences, nil)
	
	// 调用函数
	result, err := mockSvc.GetUserPreferences(ctx, userID)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedPreferences, result)
	mockSvc.AssertExpectations(t)
}

// TestUpdateUserPreferences_Success 测试成功更新用户偏好
func TestUpdateUserPreferences_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	
	// 创建更新请求
	updateRequest := user.UpdatePreferencesRequest{
		GameSettings: map[string]interface{}{
			"preferredGenres": []string{"FPS", "RPG", "Strategy"},
			"favoriteGames":   []string{"CS2", "Valorant", "Baldur's Gate 3"},
		},
		PeripheralPrefs: map[string]interface{}{
			"preferredBrands": []string{"Logitech", "Razer"},
			"mouseWeight":     "ultralight",
		},
	}
	
	// 创建期望的返回值
	expectedPreferences := &user.PreferencesResponse{
		GameSettings: updateRequest.GameSettings,
		PeripheralPrefs: updateRequest.PeripheralPrefs,
		ContentPrefs: map[string]interface{}{
			"showReviews":   true,
			"showTutorials": true,
		},
		RecommendedForYou: map[string]interface{}{
			"enableRecommendations": true,
		},
	}
	
	// 设置Mock期望
	mockSvc.On("UpdateUserPreferences", ctx, userID, updateRequest).Return(expectedPreferences, nil)
	
	// 调用函数
	result, err := mockSvc.UpdateUserPreferences(ctx, userID, updateRequest)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedPreferences, result)
	assert.Equal(t, updateRequest.GameSettings, result.GameSettings)
	assert.Equal(t, updateRequest.PeripheralPrefs, result.PeripheralPrefs)
	mockSvc.AssertExpectations(t)
}

// TestGetReviewerApplications_Success 测试成功获取评测员申请列表
func TestGetReviewerApplications_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	status := "pending"
	page := 1
	pageSize := 10
	
	// 创建期望的返回值
	userID1 := primitive.NewObjectID().Hex()
	userID2 := primitive.NewObjectID().Hex()
	now := time.Now()
	
	expectedResponse := &user.ReviewerApplicationListResponse{
		Applications: []user.ReviewerApplicationDetails{
			{
				UserID:         userID1,
				Username:       "user1",
				Email:          "user1@example.com",
				DisplayName:    "User One",
				Experience:     "5 years of experience with gaming peripherals",
				ExpertiseAreas: []string{"mice", "keyboards"},
				Samples:        []string{"https://example.com/sample1"},
				Motivation:     "I want to help the community",
				Status:         "pending",
				AppliedAt:      now,
				UpdatedAt:      now,
			},
			{
				UserID:         userID2,
				Username:       "user2",
				Email:          "user2@example.com",
				DisplayName:    "User Two",
				Experience:     "3 years of experience with gaming peripherals",
				ExpertiseAreas: []string{"mice", "mousepads"},
				Samples:        []string{"https://example.com/sample2"},
				Motivation:     "I want to share my knowledge",
				Status:         "pending",
				AppliedAt:      now,
				UpdatedAt:      now,
			},
		},
		Page:     page,
		PageSize: pageSize,
	}
	
	expectedTotal := int64(2)
	
	// 设置Mock期望
	mockSvc.On("GetReviewerApplications", ctx, status, page, pageSize).Return(expectedResponse, expectedTotal, nil)
	
	// 调用函数
	result, total, err := mockSvc.GetReviewerApplications(ctx, status, page, pageSize)
	
	// 验证结果
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, result)
	assert.Equal(t, expectedTotal, total)
	assert.Equal(t, 2, len(result.Applications))
	assert.Equal(t, userID1, result.Applications[0].UserID)
	assert.Equal(t, userID2, result.Applications[1].UserID)
	mockSvc.AssertExpectations(t)
}

// TestApproveReviewerApplication_Success 测试成功批准评测员申请
func TestApproveReviewerApplication_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	adminID := primitive.NewObjectID().Hex()
	
	// 设置Mock期望
	mockSvc.On("ApproveReviewerApplication", ctx, userID, adminID).Return(nil)
	
	// 调用函数
	err := mockSvc.ApproveReviewerApplication(ctx, userID, adminID)
	
	// 验证结果
	require.NoError(t, err)
	mockSvc.AssertExpectations(t)
}

// TestRejectReviewerApplication_Success 测试成功拒绝评测员申请
func TestRejectReviewerApplication_Success(t *testing.T) {
	mockSvc := new(mockUserService)
	ctx := context.Background()
	userID := primitive.NewObjectID().Hex()
	adminID := primitive.NewObjectID().Hex()
	reason := "Insufficient experience in the required areas"
	
	// 设置Mock期望
	mockSvc.On("RejectReviewerApplication", ctx, userID, adminID, reason).Return(nil)
	
	// 调用函数
	err := mockSvc.RejectReviewerApplication(ctx, userID, adminID, reason)
	
	// 验证结果
	require.NoError(t, err)
	mockSvc.AssertExpectations(t)
}