package models

import "time"

// Basic type definitions
type Role struct {
	Type                string               `bson:"type" json:"type"`
	ReviewerApplication *ReviewerApplication `bson:"reviewerApplication,omitempty" json:"reviewerApplication,omitempty"`
	InviteCode          string               `bson:"inviteCode,omitempty" json:"inviteCode,omitempty"`
}

type Status struct {
	EmailVerified bool      `bson:"emailVerified" json:"emailVerified"`
	VerifyToken   string    `bson:"verifyToken,omitempty" json:"-"`
	TokenExpires  time.Time `bson:"tokenExpires,omitempty" json:"-"`
	IsLocked      bool      `bson:"isLocked" json:"isLocked"`
	LockReason    string    `bson:"lockReason,omitempty" json:"lockReason,omitempty"`
	LockExpires   time.Time `bson:"lockExpires,omitempty" json:"lockExpires,omitempty"`
	EmailChange   string    `bson:"emailChange,omitempty" json:"-"`
}

type UserStats struct {
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
	Violations  int       `bson:"violations" json:"violations"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
	LastLoginIP string    `bson:"lastLoginIP,omitempty" json:"lastLoginIP,omitempty"`
}

type ReviewerApplication struct {
	Status      string    `bson:"status" json:"status"`
	AppliedAt   time.Time `bson:"appliedAt" json:"appliedAt"`
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
}
