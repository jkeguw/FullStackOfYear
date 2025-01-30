package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// User 代表系统中的用户实体
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`

	Status Status     `bson:"status" json:"status"`
	Role   Role       `bson:"role" json:"role"`
	Stats  UserStats  `bson:"stats" json:"stats"`
	OAuth  *OAuthInfo `bson:"oauth,omitempty" json:"oauth,omitempty"`

	LoginHistory []LoginRecord `bson:"loginHistory" json:"loginHistory"`
	SecurityLogs []SecurityLog `bson:"securityLogs" json:"securityLogs"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// NewUser 创建新用户
func NewUser(username, email, password string) *User {
	now := time.Now()
	return &User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Email:    email,
		Password: password,
		Role: Role{
			Type: RoleUser,
		},
		Status: Status{
			EmailVerified: false,
		},
		Stats: UserStats{
			CreatedAt:   now,
			LastLoginAt: now,
		},
		CreatedAt:    now,
		UpdatedAt:    now,
		LoginHistory: make([]LoginRecord, 0),
		SecurityLogs: make([]SecurityLog, 0),
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
