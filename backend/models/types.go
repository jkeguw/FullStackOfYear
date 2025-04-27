package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Basic type definitions
// 注意：UserRole, Status, UserStats 类型已在 auth.go 中定义
// 这里定义额外需要的字段

type ReviewerApplication struct {
	Status        string            `bson:"status" json:"status"`
	AppliedAt     time.Time         `bson:"appliedAt" json:"appliedAt"`
	UpdatedAt     time.Time         `bson:"updatedAt" json:"updatedAt"`
	ReviewCount   int               `bson:"reviewCount" json:"reviewCount"`
	TotalWords    int               `bson:"totalWords" json:"totalWords"`
	Experience    string            `bson:"experience" json:"experience"`
	ExpertiseAreas []string          `bson:"expertiseAreas" json:"expertiseAreas"`
	Samples       []string          `bson:"samples" json:"samples"`
	Motivation    string            `bson:"motivation" json:"motivation"`
	ReviewNotes   string            `bson:"reviewNotes,omitempty" json:"reviewNotes,omitempty"`
	ReviewedBy    primitive.ObjectID `bson:"reviewedBy,omitempty" json:"reviewedBy,omitempty"`
}

// 申请状态常量
const (
	ApplicationStatusPending  = "pending"
	ApplicationStatusApproved = "approved"
	ApplicationStatusRejected = "rejected"
)