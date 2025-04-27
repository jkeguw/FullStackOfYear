package review

// CreateReviewRequest 创建评测的请求
type CreateReviewRequest struct {
	ExternalItemID string   `json:"externalItemId" binding:"required"`
	ItemType       string   `json:"itemType" binding:"required,oneof=device game peripheral software"`
	Content        string   `json:"content" binding:"required,min=50"`
	Pros           []string `json:"pros" binding:"required,min=1"`
	Cons           []string `json:"cons" binding:"required,min=1"`
	Score          float64  `json:"score" binding:"required,min=1,max=5"`
	Usage          string   `json:"usage" binding:"required"`
	Type           string   `json:"type" binding:"required,oneof=mouse keyboard monitor mousepad accessory game software"`
	ContentType    string   `json:"contentType" binding:"required,oneof=single comparison experience gaming buying"`
}

// UpdateReviewRequest 更新评测的请求
type UpdateReviewRequest struct {
	Content     *string   `json:"content,omitempty" binding:"omitempty,min=50"`
	Pros        *[]string `json:"pros,omitempty" binding:"omitempty,min=1"`
	Cons        *[]string `json:"cons,omitempty" binding:"omitempty,min=1"`
	Score       *float64  `json:"score,omitempty" binding:"omitempty,min=1,max=5"`
	Usage       *string   `json:"usage,omitempty"`
	Type        *string   `json:"type,omitempty" binding:"omitempty,oneof=mouse keyboard monitor mousepad accessory game software"`
	ContentType *string   `json:"contentType,omitempty" binding:"omitempty,oneof=single comparison experience gaming buying"`
}

// ReviewListRequest 获取评测列表的请求
type ReviewListRequest struct {
	ExternalItemID string `form:"externalItemId,omitempty"`
	ItemType       string `form:"itemType,omitempty" binding:"omitempty,oneof=device game peripheral software"`
	UserID         string `form:"userId,omitempty"`
	Status         string `form:"status,omitempty" binding:"omitempty,oneof=pending approved rejected featured"`
	Type           string `form:"type,omitempty" binding:"omitempty,oneof=mouse keyboard monitor mousepad accessory game software"`
	ContentType    string `form:"contentType,omitempty" binding:"omitempty,oneof=single comparison experience gaming buying"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	PageSize       int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	SortBy         string `form:"sortBy" binding:"omitempty,oneof=createdAt score viewCount"`
	SortOrder      string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// ReviewerActionRequest 评测员操作请求
type ReviewerActionRequest struct {
	Notes string `json:"notes,omitempty"`
	Rank  int    `json:"rank,omitempty"` // 用于featured操作
}