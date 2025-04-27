package auth

import (
	"project/backend/internal/errors"
	"project/backend/models"
	"project/backend/types/auth"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// SecurityService 处理安全相关功能
type SecurityService struct {
	authService Service
}

// GetTwoFactorStatus 获取两因素认证状态
func (s *SecurityService) GetTwoFactorStatus(userID string) (*auth.TwoFactorStatusResponse, error) {
	// 实现获取两因素认证状态逻辑
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	
	enabled := false
	if user.TwoFactor != nil {
		enabled = user.TwoFactor.Enabled
	}
	
	return &auth.TwoFactorStatusResponse{
		Enabled: enabled,
	}, nil
}

// ListDevices 列出用户设备
func (s *SecurityService) ListDevices(userID, currentDeviceID string) (*auth.DeviceListResponse, error) {
	// 实现获取设备列表逻辑
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	
	devices := make([]auth.DeviceInfo, 0)
	for _, device := range user.ActiveDevices {
		deviceInfo := auth.DeviceInfo{
			ID:           device.ID,
			Name:         device.Name,
			Type:         device.Type,
			OS:           device.OS,
			Browser:      device.Browser,
			IsCurrent:    device.ID == currentDeviceID,
			LastUsedAt:   device.LastUsedAt,
			LastUsedIP:   device.IP,
			Trusted:      device.Trusted,
		}
		devices = append(devices, deviceInfo)
	}
	
	return &auth.DeviceListResponse{
		Devices: devices,
	}, nil
}

// RemoveDevice 移除设备
func (s *SecurityService) RemoveDevice(userID, deviceID, currentDeviceID string) error {
	// 实现移除设备逻辑
	if deviceID == currentDeviceID {
		return errors.NewBadRequestError("无法移除当前登录的设备")
	}
	
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}
	
	// 查找并移除设备
	removed := user.RemoveDevice(deviceID)
	if !removed {
		return errors.NewNotFoundError("未找到指定设备")
	}
	
	// 更新用户信息
	err = s.authService.UpdateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("移除设备失败: %w", err)
	}
	
	return nil
}

// UpdateDevice 更新设备信息
func (s *SecurityService) UpdateDevice(userID, deviceID, name string, trusted bool) error {
	// 实现更新设备逻辑
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}
	
	// 查找并更新设备
	found := false
	for i, device := range user.ActiveDevices {
		if device.ID == deviceID {
			user.ActiveDevices[i].Name = name
			user.ActiveDevices[i].Trusted = trusted
			found = true
			break
		}
	}
	
	if !found {
		return errors.NewNotFoundError("未找到指定设备")
	}
	
	// 更新用户信息
	err = s.authService.UpdateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("更新设备失败: %w", err)
	}
	
	return nil
}

// VerifyTwoFactorCode 验证两因素认证码
func (s *SecurityService) VerifyTwoFactorCode(userID, code string) (bool, error) {
	if code == "" {
		return false, errors.NewBadRequestError("验证码不能为空")
	}
	
	err := s.VerifyTwoFactor(userID, code)
	return err == nil, err
}

// VerifyAndEnableTwoFactor 验证并启用两因素认证
func (s *SecurityService) VerifyAndEnableTwoFactor(userID, code string) error {
	if code == "" {
		return errors.NewBadRequestError("验证码不能为空")
	}
	
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}
	
	// 检查用户是否处于TFA待激活状态
	if user.TwoFactor == nil || user.TwoFactor.Secret == "" || user.TwoFactor.Enabled {
		return errors.NewBadRequestError("用户未设置两因素认证或已激活")
	}
	
	// 验证TOTP码
	valid := totp.Validate(code, user.TwoFactor.Secret)
	if !valid {
		return errors.NewBadRequestError("无效的验证码")
	}
	
	// 激活TFA
	err = s.authService.ActivateUserTwoFactor(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("激活两因素认证失败: %w", err)
	}
	
	return nil
}

// NewSecurityService 创建新的安全服务
func NewSecurityService(authService Service) *SecurityService {
	return &SecurityService{
		authService: authService,
	}
}

// SetupTwoFactor 设置两因素认证
func (s *SecurityService) SetupTwoFactor(userID string) (*auth.TwoFactorSetupResponse, error) {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ExternalDeviceReviewPlatform",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("生成TOTP密钥失败: %w", err)
	}

	// 生成QR码
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("生成QR码失败: %w", err)
	}

	// 更新用户TFA状态
	err = s.authService.UpdateUserTwoFactorPending(context.Background(), userID, key.Secret())
	if err != nil {
		return nil, fmt.Errorf("更新用户TFA状态失败: %w", err)
	}

	// 生成备份码
	backupCodes := s.generateRecoveryCodes()
	
	// 返回设置信息
	return &auth.TwoFactorSetupResponse{
		Secret:        key.Secret(),
		QRCode:        qrCode,
		RecoveryCodes: backupCodes,
	}, nil
}

// VerifyTwoFactor 验证两因素认证
func (s *SecurityService) VerifyTwoFactor(userID, code string) error {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// 检查用户是否设置了TFA
	if user.TwoFactor == nil || user.TwoFactor.Secret == "" {
		return errors.NewBadRequestError("用户未设置两因素认证")
	}

	// 验证TOTP码
	valid := totp.Validate(code, user.TwoFactor.Secret)
	if !valid {
		// 检查是否是恢复码
		isRecoveryCode := false
		for _, backupCode := range user.TwoFactor.BackupCodes {
			if backupCode == code {
				isRecoveryCode = true
				// 使用恢复码
				user.UseBackupCode(code)
				// 更新用户信息
				err := s.authService.UpdateUser(context.Background(), user)
				if err != nil {
					return fmt.Errorf("更新恢复码状态失败: %w", err)
				}
				break
			}
		}

		if !isRecoveryCode {
			return errors.NewBadRequestError("无效的验证码或恢复码")
		}
	}

	return nil
}

// ActivateTwoFactor 激活两因素认证
func (s *SecurityService) ActivateTwoFactor(userID, code string) error {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// 检查用户是否处于TFA待激活状态
	if user.TwoFactor == nil || user.TwoFactor.Secret == "" || user.TwoFactor.Enabled {
		return errors.NewBadRequestError("用户未设置两因素认证或已激活")
	}

	// 验证TOTP码
	valid := totp.Validate(code, user.TwoFactor.Secret)
	if !valid {
		return errors.NewBadRequestError("无效的验证码")
	}

	// 激活TFA
	err = s.authService.ActivateUserTwoFactor(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("激活两因素认证失败: %w", err)
	}

	return nil
}

// DisableTwoFactor 禁用两因素认证
func (s *SecurityService) DisableTwoFactor(userID, password string) error {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// 检查用户是否设置了TFA
	if user.TwoFactor == nil || user.TwoFactor.Secret == "" {
		return errors.NewBadRequestError("用户未设置两因素认证")
	}

	// 验证密码
	isValid, err := s.authService.VerifyPassword(context.Background(), userID, password)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.NewUnauthorizedError("密码错误")
	}

	// 禁用TFA
	err = s.authService.DisableUserTwoFactor(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("禁用两因素认证失败: %w", err)
	}

	return nil
}

// ResetTwoFactorWithRecoveryCode 使用备份码重置两因素认证
func (s *SecurityService) ResetTwoFactorWithRecoveryCode(userID, backupCode string) error {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return err
	}

	// 检查用户是否设置了TFA
	if user.TwoFactor == nil || user.TwoFactor.Secret == "" {
		return errors.NewBadRequestError("用户未设置两因素认证")
	}

	// 验证恢复码
	validCode := false
	for _, code := range user.TwoFactor.BackupCodes {
		if code == backupCode {
			validCode = true
			break
		}
	}

	if !validCode {
		return errors.NewBadRequestError("无效的恢复码")
	}

	// 使用备份码
	used := user.UseBackupCode(backupCode)
	if !used {
		return errors.NewBadRequestError("无法使用备份码")
	}
	
	// 更新用户信息
	err = s.authService.UpdateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("更新恢复码状态失败: %w", err)
	}

	// 禁用TFA
	err = s.authService.DisableUserTwoFactor(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("重置两因素认证失败: %w", err)
	}

	return nil
}

// UpdatePassword 更新密码
func (s *SecurityService) UpdatePassword(userID, currentPassword, newPassword string) error {
	// 验证当前密码
	isValid, err := s.authService.VerifyPassword(context.Background(), userID, currentPassword)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.NewUnauthorizedError("当前密码错误")
	}

	// 更新密码
	err = s.authService.UpdateUserPassword(context.Background(), userID, newPassword)
	if err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// ChangePassword 修改密码
func (s *SecurityService) ChangePassword(userID, currentPassword, newPassword string) error {
	return s.UpdatePassword(userID, currentPassword, newPassword)
}

// GetSecurityLogs 获取用户安全日志
func (s *SecurityService) GetSecurityLogs(userID string) (*auth.SecurityLogResponse, error) {
	user, err := s.authService.GetUserByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	
	logs := make([]auth.SecurityLogEntry, 0)
	for _, log := range user.SecurityLogs {
		entry := auth.SecurityLogEntry{
			Action:      log.Action,
			Timestamp:   log.Timestamp,
			Description: log.Description,
			IP:          log.IP,
			DeviceInfo:  log.DeviceInfo,
		}
		logs = append(logs, entry)
	}
	
	return &auth.SecurityLogResponse{
		Logs: logs,
	}, nil
}

// ResetPassword 重置密码
func (s *SecurityService) ResetPassword(email, token, newPassword string) error {
	// 验证重置令牌
	user, err := s.authService.ValidatePasswordResetToken(context.Background(), email, token)
	if err != nil {
		return err
	}

	// 更新密码
	err = s.authService.UpdateUserPassword(context.Background(), user.ID.Hex(), newPassword)
	if err != nil {
		return fmt.Errorf("重置密码失败: %w", err)
	}

	return nil
}

// 辅助方法

// generateRecoveryCodes 生成备份码
func (s *SecurityService) generateRecoveryCodes() []string {
	codes := make([]string, models.TwoFactorRecoveryCodeCount)
	for i := 0; i < models.TwoFactorRecoveryCodeCount; i++ {
		codes[i] = generateRandomCode(10)
	}
	return codes
}

// generateRandomCode 生成随机字符串
func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	
	// 读取随机字节
	n, err := rand.Read(code)
	if err != nil || n != length {
		// 发生错误时使用备用方案
		for i := 0; i < length; i++ {
			code[i] = charset[0]
		}
		return string(code)
	}
	
	// 将随机字节映射到字符集
	for i := 0; i < length; i++ {
		code[i] = charset[int(code[i])%len(charset)]
	}
	return string(code)
}