package models

import "time"

// Role 代表用户角色信息
type Role struct {
	Type                string               `bson:"type" json:"type"`
	ReviewerApplication *ReviewerApplication `bson:"reviewerApplication,omitempty" json:"reviewerApplication,omitempty"`
	InviteCode          string               `bson:"inviteCode,omitempty" json:"inviteCode,omitempty"`
}

// ReviewerApplication 代表审阅者申请状态和信息
type ReviewerApplication struct {
	Status      string    `bson:"status" json:"status"`
	AppliedAt   time.Time `bson:"appliedAt" json:"appliedAt"`
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
}

// UserStats 代表用户统计信息
type UserStats struct {
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
	Violations  int       `bson:"violations" json:"violations"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
	LastLoginIP string    `bson:"lastLoginIP,omitempty" json:"lastLoginIP,omitempty"`
}

// Status 代表用户账户的当前状态
type Status struct {
	EmailVerified bool      `bson:"emailVerified" json:"emailVerified"`
	VerifyToken   string    `bson:"verifyToken,omitempty" json:"-"`
	TokenExpires  time.Time `bson:"tokenExpires,omitempty" json:"-"`
	IsLocked      bool      `bson:"isLocked" json:"isLocked"`
	LockReason    string    `bson:"lockReason,omitempty" json:"lockReason,omitempty"`
	LockExpires   time.Time `bson:"lockExpires,omitempty" json:"lockExpires,omitempty"`
	EmailChange   string    `bson:"emailChange,omitempty" json:"-"`
}
