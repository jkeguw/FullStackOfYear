package measurement

import (
	"project/backend/models"
	"time"
)

// CreateMeasurementRequest 创建测量记录的请求
type CreateMeasurementRequest struct {
	Palm       float64 `json:"palm" binding:"required,min=50,max=150"`   // 手掌宽度
	Length     float64 `json:"length" binding:"required,min=40,max=120"` // 手指长度
	Unit       string  `json:"unit" binding:"required,oneof=mm cm inch"` // 单位
	Device     string  `json:"device"`                                   // 测量设备
	Calibrated bool    `json:"calibrated"`                               // 是否校准
}

// UpdateMeasurementRequest 更新测量记录的请求
type UpdateMeasurementRequest struct {
	Palm       *float64 `json:"palm" binding:"omitempty,min=50,max=150"`
	Length     *float64 `json:"length" binding:"omitempty,min=40,max=120"`
	Unit       *string  `json:"unit" binding:"omitempty,oneof=mm cm inch"`
	Calibrated *bool    `json:"calibrated"`
}

// MeasurementResponse 单个测量记录的响应
type MeasurementResponse struct {
	ID        string                     `json:"id"`
	Palm      float64                    `json:"palm"`
	Length    float64                    `json:"length"`
	Unit      string                     `json:"unit"`
	Quality   *models.MeasurementQuality `json:"quality,omitempty"`
	CreatedAt time.Time                  `json:"createdAt"`
	UpdatedAt time.Time                  `json:"updatedAt"`
}

// MeasurementListRequest 获取测量记录列表的请求
type MeasurementListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=createdAt palm length quality"`
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`
	StartDate string `form:"startDate" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"endDate" binding:"omitempty,datetime=2006-01-02"`
}

// MeasurementListResponse 测量记录列表的响应
type MeasurementListResponse struct {
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"pageSize"`
	Measurements []MeasurementResponse `json:"measurements"`
}

// MeasurementStatsResponse 用户测量统计的响应
type MeasurementStatsResponse struct {
	AveragePalm      float64   `json:"averagePalm"`
	AverageLength    float64   `json:"averageLength"`
	HandSize         string    `json:"handSize"`
	MeasurementCount int       `json:"measurementCount"`
	LastMeasuredAt   time.Time `json:"lastMeasuredAt,omitempty"`
}

// MeasurementRecommendationResponse 设备推荐的响应
type MeasurementRecommendationResponse struct {
	HandSize string                 `json:"handSize"`
	GripType string                 `json:"gripType"`
	Devices  []DeviceRecommendation `json:"devices"`
}

// DeviceRecommendation 单个设备推荐
type DeviceRecommendation struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	MatchScore int    `json:"matchScore"` // 匹配得分 0-100
	Reason     string `json:"reason"`     // 推荐理由
}
