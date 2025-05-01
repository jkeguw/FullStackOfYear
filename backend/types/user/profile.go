package user

import (
	"project/backend/models"
	"time"
)

// ProfileResponse 个人资料响应
type ProfileResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	IsVerified  bool      `json:"isVerified"`
	LastLoginAt time.Time `json:"lastLoginAt,omitempty"`
	LastLoginIP string    `json:"lastLoginIP,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TwoFactor   bool      `json:"twoFactor,omitempty"`
}

// SecurityLogResponse 安全日志响应
type SecurityLogResponse struct {
	Action      string    `json:"action"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
}

// LoginHistoryResponse 登录历史响应
type LoginHistoryResponse struct {
	IP          string    `json:"ip"`
	Location    string    `json:"location,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	DeviceName  string    `json:"deviceName,omitempty"`
	DeviceType  string    `json:"deviceType,omitempty"`
	LoginMethod string    `json:"loginMethod,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
}

// ActiveDeviceResponse 用户活跃设备响应
type ActiveDeviceResponse struct {
	ID          string    `json:"id"`
	DeviceName  string    `json:"deviceName"`
	DeviceType  string    `json:"deviceType"`
	OS          string    `json:"os"`
	Browser     string    `json:"browser"`
	Location    string    `json:"location,omitempty"`
	LastUsedAt  time.Time `json:"lastUsedAt"`
	FirstUsedAt time.Time `json:"firstUsedAt"`
}

// 转换函数
func MapToProfileResponse(user *models.User) ProfileResponse {
	return ProfileResponse{
		ID:          user.ID.Hex(),
		Username:    user.Username,
		Email:       user.Email,
		Role:        string(user.Role.Type),
		IsVerified:  user.IsVerified,
		LastLoginAt: user.Stats.LastLoginAt,
		LastLoginIP: user.Stats.LastLoginIP,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		TwoFactor:   user.TwoFactor != nil && user.TwoFactor.Enabled,
	}
}

// MapToSecurityLogsResponse 转换安全日志为响应格式
func MapToSecurityLogsResponse(logs []models.SecurityLog) []SecurityLogResponse {
	result := make([]SecurityLogResponse, len(logs))
	for i, log := range logs {
		result[i] = SecurityLogResponse{
			Action:      log.Action,
			Timestamp:   log.Timestamp,
			Description: log.Description,
		}
	}
	return result
}

// MapToLoginHistoryResponse 转换登录历史为响应格式
func MapToLoginHistoryResponse(records []models.LoginRecord) []LoginHistoryResponse {
	result := make([]LoginHistoryResponse, len(records))
	for i, record := range records {
		result[i] = LoginHistoryResponse{
			IP:          record.IP,
			Location:    "", // No Location field in models.LoginRecord
			UserAgent:   record.UserAgent,
			DeviceName:  "", // No DeviceName field in models.LoginRecord
			DeviceType:  "",
			LoginMethod: "",
			Timestamp:   record.Timestamp,
			Success:     record.Success,
		}
	}
	return result
}

// MapToActiveDevicesResponse 转换活跃设备列表为响应格式
func MapToActiveDevicesResponse(devices []models.Device) []ActiveDeviceResponse {
	result := make([]ActiveDeviceResponse, len(devices))
	for i, device := range devices {
		result[i] = ActiveDeviceResponse{
			ID:          device.ID,
			DeviceName:  device.Name,
			DeviceType:  device.Type,
			OS:          device.OS,
			Browser:     device.Browser,
			Location:    device.IP, // Use IP instead of non-existent Location field
			LastUsedAt:  device.LastUsedAt,
			FirstUsedAt: device.CreatedAt, // Use CreatedAt instead of non-existent FirstUsedAt field
		}
	}
	return result
}