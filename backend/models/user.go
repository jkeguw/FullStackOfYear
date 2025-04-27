package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"time"
)

// User 代表系统中的用户实体
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`

	Status Status     `bson:"status" json:"status"`
	Role   UserRole   `bson:"role" json:"role"`
	Stats  UserStats  `bson:"stats" json:"stats"`
	OAuth  *OAuthInfo `bson:"oauth,omitempty" json:"oauth,omitempty"`

	TwoFactor    *TwoFactorAuth     `bson:"twoFactor,omitempty" json:"twoFactor,omitempty"`
	ActiveDevices []Device          `bson:"activeDevices,omitempty" json:"activeDevices,omitempty"`

	LoginHistory []LoginRecord `bson:"loginHistory" json:"loginHistory"`
	SecurityLogs []SecurityLog `bson:"securityLogs" json:"securityLogs"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	IsVerified bool     `bson:"isVerified" json:"isVerified"`
}

// NewUser 创建新用户
func NewUser(username, email, password string) *User {
	now := time.Now()
	return &User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    email,
		Password: password,
		Role: UserRole{
			Type: string(RoleUser),
		},
		Status: Status{
			EmailVerified: false,
		},
		Stats: UserStats{
			CreatedAt:   now,
			LastLoginAt: now,
		},
		TwoFactor:     &TwoFactorAuth{Enabled: false},
		ActiveDevices: make([]Device, 0),
		CreatedAt:     now,
		UpdatedAt:     now,
		LoginHistory:  make([]LoginRecord, 0),
		SecurityLogs:  make([]SecurityLog, 0),
		IsVerified:    false,
	}
}

// AddLoginRecord 添加登录记录
func (u *User) AddLoginRecord(record LoginRecord) {
	if len(u.LoginHistory) >= MaxLoginHistory {
		u.LoginHistory = u.LoginHistory[1:]
	}
	u.LoginHistory = append(u.LoginHistory, record)
	u.Stats.LastLoginAt = record.Timestamp
	u.Stats.LastLoginIP = record.IP
	u.UpdatedAt = record.Timestamp
}

// AddSecurityLog 添加安全日志
func (u *User) AddSecurityLog(log SecurityLog) {
	u.SecurityLogs = append(u.SecurityLogs, log)
	u.UpdatedAt = log.Timestamp
}

// UpdateOAuthInfo 更新 OAuth 信息
func (u *User) UpdateOAuthInfo(provider string, info interface{}) {
	if u.OAuth == nil {
		u.OAuth = &OAuthInfo{}
	}

	switch provider {
	case "google":
		if googleInfo, ok := info.(*GoogleOAuth); ok {
			u.OAuth.Google = googleInfo
		}
	}
	u.UpdatedAt = time.Now()
}

// AddDevice 添加设备到活跃设备列表
func (u *User) AddDevice(device Device) {
	// 检查设备是否已存在
	for i, d := range u.ActiveDevices {
		if d.ID == device.ID {
			// 更新现有设备信息
			u.ActiveDevices[i] = device
			u.UpdatedAt = time.Now()
			return
		}
	}
	
	// 如果达到设备限制，移除最旧的设备
	if len(u.ActiveDevices) >= MaxConcurrentSessions {
		// 按最后使用时间排序
		sort.Slice(u.ActiveDevices, func(i, j int) bool {
			return u.ActiveDevices[i].LastUsedAt.Before(u.ActiveDevices[j].LastUsedAt)
		})
		// 移除最旧的设备
		u.ActiveDevices = u.ActiveDevices[1:]
	}
	
	// 添加新设备
	u.ActiveDevices = append(u.ActiveDevices, device)
	u.UpdatedAt = time.Now()
}

// RemoveDevice 从活跃设备列表中移除设备
func (u *User) RemoveDevice(deviceID string) bool {
	for i, device := range u.ActiveDevices {
		if device.ID == deviceID {
			// 移除设备
			u.ActiveDevices = append(u.ActiveDevices[:i], u.ActiveDevices[i+1:]...)
			u.UpdatedAt = time.Now()
			
			// 添加安全日志
			u.AddSecurityLog(SecurityLog{
				Action:      "device_removed",
				Timestamp:   time.Now(),
				Description: fmt.Sprintf("设备已从账户移除: %s", device.Name),
			})
			return true
		}
	}
	return false
}

// EnableTwoFactor 启用两因素认证
func (u *User) EnableTwoFactor(secret string, backupCodes []string) {
	if u.TwoFactor == nil {
		u.TwoFactor = &TwoFactorAuth{}
	}
	u.TwoFactor.Enabled = true
	u.TwoFactor.Secret = secret
	u.TwoFactor.BackupCodes = backupCodes
	u.TwoFactor.VerifiedAt = time.Now()
	u.UpdatedAt = time.Now()
	
	// 添加安全日志
	u.AddSecurityLog(SecurityLog{
		Action:      "2fa_enabled",
		Timestamp:   time.Now(),
		Description: "两因素认证已启用",
	})
}

// DisableTwoFactor 禁用两因素认证
func (u *User) DisableTwoFactor() {
	if u.TwoFactor != nil {
		u.TwoFactor.Enabled = false
		u.TwoFactor.Secret = ""
		u.TwoFactor.BackupCodes = nil
		u.UpdatedAt = time.Now()
		
		// 添加安全日志
		u.AddSecurityLog(SecurityLog{
			Action:      "2fa_disabled",
			Timestamp:   time.Now(),
			Description: "两因素认证已禁用",
		})
	}
}

// UseBackupCode 使用备份恢复码
func (u *User) UseBackupCode(code string) bool {
	if u.TwoFactor == nil || !u.TwoFactor.Enabled {
		return false
	}
	
	for i, backupCode := range u.TwoFactor.BackupCodes {
		if backupCode == code {
			// 删除已使用的备份码
			u.TwoFactor.BackupCodes = append(
				u.TwoFactor.BackupCodes[:i], 
				u.TwoFactor.BackupCodes[i+1:]...
			)
			u.UpdatedAt = time.Now()
			
			// 添加安全日志
			u.AddSecurityLog(SecurityLog{
				Action:      "backup_code_used",
				Timestamp:   time.Now(),
				Description: "两因素认证备份码已使用",
			})
			return true
		}
	}
	return false
}

// UpdateUserTwoFactorPending 更新用户两因素认证待激活状态
func (u *User) UpdateUserTwoFactorPending(secret string) {
	if u.TwoFactor == nil {
		u.TwoFactor = &TwoFactorAuth{Enabled: false}
	}
	u.TwoFactor.Secret = secret
	u.UpdatedAt = time.Now()
}

// ActivateUserTwoFactor 激活用户两因素认证
func (u *User) ActivateUserTwoFactor() {
	if u.TwoFactor == nil {
		u.TwoFactor = &TwoFactorAuth{Enabled: false}
	}
	now := time.Now()
	u.TwoFactor.Enabled = true
	u.TwoFactor.VerifiedAt = now
	u.UpdatedAt = now
	
	// 添加安全日志
	u.AddSecurityLog(SecurityLog{
		Action:      "2fa_activated",
		Timestamp:   now,
		Description: "两因素认证已激活",
	})
}

// DisableUserTwoFactor 禁用用户两因素认证
func (u *User) DisableUserTwoFactor() {
	if u.TwoFactor == nil {
		return
	}
	u.TwoFactor.Enabled = false
	u.TwoFactor.Secret = ""
	u.TwoFactor.BackupCodes = nil
	u.UpdatedAt = time.Now()
	
	// 添加安全日志
	u.AddSecurityLog(SecurityLog{
		Action:      "2fa_disabled",
		Timestamp:   time.Now(),
		Description: "两因素认证已禁用",
	})
}

// UpdateUserRecoveryCodes 更新用户恢复码
func (u *User) UpdateUserRecoveryCodes(recoveryCodes []string) {
	if u.TwoFactor == nil {
		u.TwoFactor = &TwoFactorAuth{Enabled: false}
	}
	u.TwoFactor.BackupCodes = recoveryCodes
	u.UpdatedAt = time.Now()
}