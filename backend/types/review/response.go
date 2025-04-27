package review

import "time"

// ReviewResponse 评测响应
type ReviewResponse struct {
	ID             string     `json:"id"`
	ExternalItemID string     `json:"externalItemId"`
	ItemType       string     `json:"itemType"`
	UserID         string     `json:"userId"`
	Content        string     `json:"content"`
	Pros           []string   `json:"pros"`
	Cons           []string   `json:"cons"`
	Score          float64    `json:"score"`
	Usage          string     `json:"usage"`
	Status         string     `json:"status"`
	Type           string     `json:"type"`
	ContentType    string     `json:"contentType"`
	ReviewerID     string     `json:"reviewerId,omitempty"`
	ReviewerNotes  string     `json:"reviewerNotes,omitempty"`
	ReviewedAt     *time.Time `json:"reviewedAt,omitempty"`
	PublishedAt    *time.Time `json:"publishedAt,omitempty"`
	FeaturedRank   *int       `json:"featuredRank,omitempty"`
	ViewCount      int        `json:"viewCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// ReviewListResponse 评测列表响应
type ReviewListResponse struct {
	Total   int              `json:"total"`
	Page    int              `json:"page"`
	PageSize int             `json:"pageSize"`
	Reviews []ReviewResponse `json:"reviews"`
}

// UserReviewStats 用户评测统计
type UserReviewStats struct {
	TotalReviews  int     `json:"totalReviews"`
	ApprovedCount int     `json:"approvedCount"`
	PendingCount  int     `json:"pendingCount"`
	RejectedCount int     `json:"rejectedCount"`
	FeaturedCount int     `json:"featuredCount"`
	AverageScore  float64 `json:"averageScore"`
	TotalViewCount int    `json:"totalViewCount"`
}