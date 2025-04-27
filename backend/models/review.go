package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// ReviewType 评测类型枚举
type ReviewType string

const (
	ReviewTypeMouse     ReviewType = "mouse"     // 鼠标评测
	ReviewTypeKeyboard  ReviewType = "keyboard"  // 键盘评测
	ReviewTypeMonitor   ReviewType = "monitor"   // 显示器评测
	ReviewTypeMousepad  ReviewType = "mousepad"  // 鼠标垫评测
	ReviewTypeAccessory ReviewType = "accessory" // 配件评测
)

// ReviewContentType 评测内容类型枚举
type ReviewContentType string

const (
	ReviewContentSingle     ReviewContentType = "single"     // 单品评测
	ReviewContentComparison ReviewContentType = "comparison" // 对比评测
	ReviewContentExperience ReviewContentType = "experience" // 使用心得
	ReviewContentGaming     ReviewContentType = "gaming"     // 游戏体验
	ReviewContentBuying     ReviewContentType = "buying"     // 选购建议
)

// ReviewStatus 评测状态枚举
type ReviewStatus string

const (
	ReviewStatusPending  ReviewStatus = "pending"  // 待审核
	ReviewStatusApproved ReviewStatus = "approved" // 已批准
	ReviewStatusRejected ReviewStatus = "rejected" // 已拒绝
	ReviewStatusFeatured ReviewStatus = "featured" // 已推荐
)

// Review 评测数据模型
type Review struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ExternalItemID primitive.ObjectID  `bson:"externalItemId" json:"externalItemId"` // 外部项目ID
	ItemType       string              `bson:"itemType" json:"itemType"`           // 项目类型(device, game等)
	UserID         primitive.ObjectID  `bson:"userId" json:"userId"`
	Content        string              `bson:"content" json:"content"`
	Pros           []string            `bson:"pros" json:"pros"`
	Cons           []string            `bson:"cons" json:"cons"`
	Score          float64             `bson:"score" json:"score"`                // 评分(1-5)
	Usage          string              `bson:"usage" json:"usage"`                // 使用场景
	Status         ReviewStatus        `bson:"status" json:"status"`              // 状态
	Type           ReviewType          `bson:"type" json:"type"`                  // 评测类型
	ContentType    ReviewContentType   `bson:"contentType" json:"contentType"`    // 内容类型
	ReviewerID     *primitive.ObjectID `bson:"reviewerId,omitempty" json:"reviewerId,omitempty"`
	ReviewerNotes  string              `bson:"reviewerNotes,omitempty" json:"reviewerNotes,omitempty"`
	ReviewedAt     *time.Time          `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"`
	PublishedAt    *time.Time          `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	FeaturedRank   *int                `bson:"featuredRank,omitempty" json:"featuredRank,omitempty"`
	ViewCount      int                 `bson:"viewCount" json:"viewCount"`
	CreatedAt      time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// 集合名常量已在 constants.go 中定义