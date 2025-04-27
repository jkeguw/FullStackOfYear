package device

import (
	"project/backend/models"
	"time"
)

// 基础设备请求和响应类型

// DeviceListFilter 设备列表筛选条件
type DeviceListFilter struct {
	Page     int
	PageSize int
	Type     string
	Brand    string
}

// CreateReviewRequest 创建评测的请求，与SubmitReviewRequest基本一致
type CreateReviewRequest SubmitReviewRequest

// UpdateReviewRequest 更新评测的请求，与UpdateDeviceReviewRequest基本一致
type UpdateReviewRequest UpdateDeviceReviewRequest

// DeviceListRequest 获取设备列表的请求
type DeviceListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Type      string `form:"type" binding:"omitempty,oneof=mouse keyboard monitor mousepad accessory"`
	Brand     string `form:"brand" binding:"omitempty"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=name brand createdAt"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// DeviceListResponse 设备列表的响应
type DeviceListResponse struct {
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Devices  []DevicePreview `json:"devices"`
}

// DevicePreview 设备预览信息
type DevicePreview struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Brand       string    `json:"brand"`
	Type        string    `json:"type"`
	ImageURL    string    `json:"imageUrl,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// 鼠标相关类型

// CreateMouseRequest 创建鼠标设备的请求
type CreateMouseRequest struct {
	Name        string                     `json:"name" binding:"required"`
	Brand       string                     `json:"brand" binding:"required"`
	ImageURL    string                     `json:"imageUrl,omitempty"`
	Description string                     `json:"description,omitempty"`
	Dimensions  models.MouseDimensions     `json:"dimensions" binding:"required"`
	Shape       models.MouseShape          `json:"shape" binding:"required"`
	Technical   models.MouseTechnical      `json:"technical" binding:"required"`
	Recommended models.MouseRecommended    `json:"recommended" binding:"required"`
}

// UpdateMouseRequest 更新鼠标设备的请求
type UpdateMouseRequest struct {
	Name        *string                    `json:"name,omitempty"`
	Brand       *string                    `json:"brand,omitempty"`
	ImageURL    *string                    `json:"imageUrl,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Dimensions  *models.MouseDimensions    `json:"dimensions,omitempty"`
	Shape       *models.MouseShape         `json:"shape,omitempty"`
	Technical   *models.MouseTechnical     `json:"technical,omitempty"`
	Recommended *models.MouseRecommended   `json:"recommended,omitempty"`
}

// MouseResponse 鼠标设备的响应
type MouseResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Brand       string                   `json:"brand"`
	Type        string                   `json:"type"`
	ImageURL    string                   `json:"imageUrl,omitempty"`
	Description string                   `json:"description,omitempty"`
	Dimensions  models.MouseDimensions   `json:"dimensions"`
	Shape       models.MouseShape        `json:"shape"`
	Technical   models.MouseTechnical    `json:"technical"`
	Recommended models.MouseRecommended  `json:"recommended"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

// 键盘相关类型

// CreateKeyboardRequest 创建键盘设备的请求
type CreateKeyboardRequest struct {
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Description string `json:"description,omitempty"`
	Layout      string `json:"layout" binding:"required"`
	Switches    []string `json:"switches" binding:"required"`
	Size        string `json:"size" binding:"required"`
	Technical   struct {
		Connectivity []string `json:"connectivity" binding:"required"`
		Keycaps      string   `json:"keycaps" binding:"required"`
		HotSwap      bool     `json:"hotSwap"`
		RGBLighting  bool     `json:"rgbLighting"`
		NKeyRollover bool     `json:"nKeyRollover"`
	} `json:"technical" binding:"required"`
	Recommended struct {
		GameTypes []string `json:"gameTypes"`
		DailyUse  bool     `json:"dailyUse"`
		Portable  bool     `json:"portable"`
	} `json:"recommended" binding:"required"`
}

// 显示器相关类型

// CreateMonitorRequest 创建显示器设备的请求
type CreateMonitorRequest struct {
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Description string `json:"description,omitempty"`
	Size        float64 `json:"size" binding:"required,min=10,max=50"`
	Resolution struct {
		Width  int `json:"width" binding:"required,min=640"`
		Height int `json:"height" binding:"required,min=480"`
	} `json:"resolution" binding:"required"`
	Technical struct {
		RefreshRate    int     `json:"refreshRate" binding:"required,min=30,max=500"`
		ResponseTime   float64 `json:"responseTime" binding:"required"`
		PanelType      string  `json:"panelType" binding:"required"`
		AspectRatio    string  `json:"aspectRatio" binding:"required"`
		Curvature      int     `json:"curvature,omitempty"`
		HDRSupport     bool    `json:"hdrSupport"`
		AdaptiveSync   string  `json:"adaptiveSync,omitempty"`
	} `json:"technical" binding:"required"`
	Recommended struct {
		GameTypes []string `json:"gameTypes"`
		Content   []string `json:"content"`
		ProUse    bool     `json:"proUse"`
	} `json:"recommended" binding:"required"`
}

// 鼠标垫相关类型

// CreateMousepadRequest 创建鼠标垫设备的请求
type CreateMousepadRequest struct {
	Name        string `json:"name" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Description string `json:"description,omitempty"`
	Size      struct {
		Length float64 `json:"length" binding:"required,min=10"`
		Width  float64 `json:"width" binding:"required,min=10"`
		Height float64 `json:"height" binding:"required"`
	} `json:"size" binding:"required"`
	Material    string   `json:"material" binding:"required"`
	Surface     string   `json:"surface" binding:"required"`
	Base        string   `json:"base" binding:"required"`
	Recommended []string `json:"recommended"`
}

// 评测相关类型

// CreateDeviceReviewRequest 创建设备评测的请求
type CreateDeviceReviewRequest struct {
	DeviceID    string   `json:"deviceId" binding:"required"`
	Content     string   `json:"content" binding:"required,min=50"`
	Pros        []string `json:"pros" binding:"required,min=1"`
	Cons        []string `json:"cons" binding:"required,min=1"`
	Score       float64  `json:"score" binding:"required,min=1,max=5"`
	Usage       string   `json:"usage" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=mouse keyboard monitor mousepad accessory"`
	ContentType string   `json:"contentType" binding:"required,oneof=single comparison experience gaming buying"`
}

// UpdateDeviceReviewRequest 更新设备评测的请求
type UpdateDeviceReviewRequest struct {
	Content     *string   `json:"content,omitempty" binding:"omitempty,min=50"`
	Pros        *[]string `json:"pros,omitempty" binding:"omitempty,min=1"`
	Cons        *[]string `json:"cons,omitempty" binding:"omitempty,min=1"`
	Score       *float64  `json:"score,omitempty" binding:"omitempty,min=1,max=5"`
	Usage       *string   `json:"usage,omitempty"`
	Type        *string   `json:"type,omitempty" binding:"omitempty,oneof=mouse keyboard monitor mousepad accessory"`
	ContentType *string   `json:"contentType,omitempty" binding:"omitempty,oneof=single comparison experience gaming buying"`
}

// DeviceReviewResponse 设备评测响应
type DeviceReviewResponse struct {
	ID            string     `json:"id"`
	DeviceID      string     `json:"deviceId"`
	UserID        string     `json:"userId"`
	Content       string     `json:"content"`
	Pros          []string   `json:"pros"`
	Cons          []string   `json:"cons"`
	Score         float64    `json:"score"`
	Usage         string     `json:"usage"`
	Status        string     `json:"status"`
	Type          string     `json:"type"`
	ContentType   string     `json:"contentType"`
	ReviewerID    string     `json:"reviewerId,omitempty"`
	ReviewerNotes string     `json:"reviewerNotes,omitempty"`
	ReviewedAt    *time.Time `json:"reviewedAt,omitempty"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
	FeaturedRank  *int       `json:"featuredRank,omitempty"`
	ViewCount     int        `json:"viewCount"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// DeviceReviewListRequest 获取设备评测列表的请求
type DeviceReviewListRequest struct {
	DeviceID  string `form:"deviceId,omitempty"`
	UserID    string `form:"userId,omitempty"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=createdAt score"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// DeviceReviewListResponse 设备评测列表响应
type DeviceReviewListResponse struct {
	Total   int                   `json:"total"`
	Page    int                   `json:"page"`
	PageSize int                  `json:"pageSize"`
	Reviews []DeviceReviewResponse `json:"reviews"`
}

// 用户设备配置相关类型

// CreateUserDeviceRequest 创建用户设备配置的请求
type CreateUserDeviceRequest struct {
	Name        string                      `json:"name" binding:"required"`
	Description string                      `json:"description,omitempty"`
	Devices     []UserDeviceSettingsRequest `json:"devices" binding:"required,min=1"`
	IsPublic    bool                        `json:"isPublic"`
}

// UpdateUserDeviceRequest 更新用户设备配置的请求
type UpdateUserDeviceRequest struct {
	Name        string                      `json:"name,omitempty"`
	Description string                      `json:"description,omitempty"`
	Devices     []UserDeviceSettingsRequest `json:"devices,omitempty"`
	IsPublic    bool                        `json:"isPublic"`
}

// UserDeviceSettingsRequest 用户设备设置请求
type UserDeviceSettingsRequest struct {
	DeviceID  string         `json:"deviceId" binding:"required"`
	DeviceType string         `json:"deviceType" binding:"required,oneof=mouse keyboard monitor mousepad accessory"`
	Settings  map[string]any `json:"settings,omitempty"`
}

// UserDeviceResponse 用户设备配置响应
type UserDeviceResponse struct {
	ID          string                    `json:"id"`
	UserID      string                    `json:"userId"`
	Name        string                    `json:"name"`
	Description string                    `json:"description,omitempty"`
	Devices     []UserDeviceSettingsResponse `json:"devices"`
	IsPublic    bool                      `json:"isPublic"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
}

// UserDeviceSettingsResponse 用户设备设置响应
type UserDeviceSettingsResponse struct {
	DeviceID   string         `json:"deviceId"`
	DeviceType string         `json:"deviceType"`
	DeviceName string         `json:"deviceName"`
	DeviceBrand string        `json:"deviceBrand"`
	Settings   map[string]any `json:"settings,omitempty"`
}

// UserDeviceListRequest 获取用户设备配置列表的请求
type UserDeviceListRequest struct {
	UserID    string `form:"userId,omitempty"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	IsPublic  *bool  `form:"isPublic,omitempty"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=name createdAt updatedAt"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
}

// UserDeviceListResponse 用户设备配置列表响应
type UserDeviceListResponse struct {
	Total       int                 `json:"total"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"pageSize"`
	UserDevices []UserDeviceResponse `json:"userDevices"`
}

// 评测提交相关类型
type SubmitReviewRequest struct {
	DeviceID    string   `json:"deviceId" binding:"required"`
	Content     string   `json:"content" binding:"required,min=50"`
	Pros        []string `json:"pros" binding:"required,min=1"`
	Cons        []string `json:"cons" binding:"required,min=1"`
	Score       float64  `json:"score" binding:"required,min=1,max=5"`
	Usage       string   `json:"usage" binding:"required"`
	Type        string   `json:"type" binding:"required,oneof=mouse keyboard monitor mousepad accessory"`
	ContentType string   `json:"contentType" binding:"required,oneof=single comparison experience gaming buying"`
}