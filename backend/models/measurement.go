package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Measurement 表示用户的手部测量记录
type Measurement struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID  `bson:"userId" json:"userId"`
	Measurements MeasurementData     `bson:"measurements" json:"measurements"`
	Metadata     MeasurementMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Stats        MeasurementStats    `bson:"stats,omitempty" json:"stats,omitempty"`
	Quality      MeasurementQuality  `bson:"quality,omitempty" json:"quality,omitempty"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt    *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// MeasurementData 包含实际的测量数据
type MeasurementData struct {
	Palm   float64 `bson:"palm" json:"palm"`     // 手掌宽度(mm)
	Length float64 `bson:"length" json:"length"` // 手指长度(mm)
	Unit   string  `bson:"unit" json:"unit"`     // 单位 (mm/cm/inch)
}

// MeasurementMetadata 包含测量环境信息
type MeasurementMetadata struct {
	Device    string `bson:"device,omitempty" json:"device,omitempty"`       // 测量设备信息
	Browser   string `bson:"browser,omitempty" json:"browser,omitempty"`     // 浏览器信息
	UserAgent string `bson:"userAgent,omitempty" json:"userAgent,omitempty"` // 用户代理
}

// MeasurementStats 包含与历史数据的统计分析
type MeasurementStats struct {
	Variance  float64 `bson:"variance" json:"variance"`   // 与历史数据的方差
	IsAnomaly bool    `bson:"isAnomaly" json:"isAnomaly"` // 是否异常值
}

// MeasurementQuality 表示测量质量评估
type MeasurementQuality struct {
	Score   int                `bson:"score" json:"score"`     // 质量分数 0-100
	Factors MeasurementFactors `bson:"factors" json:"factors"` // 影响因素
}

// MeasurementFactors 包含影响测量质量的因素
type MeasurementFactors struct {
	Calibration bool    `bson:"calibration" json:"calibration"` // 是否经过校准
	Stability   float64 `bson:"stability" json:"stability"`     // 测量稳定性
	Consistency float64 `bson:"consistency" json:"consistency"` // 与历史数据一致性
}

// MeasurementUserStats 表示用户的测量统计信息
type MeasurementUserStats struct {
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`
	Averages         MeasurementData    `bson:"averages" json:"averages"`
	HandSize         string             `bson:"handSize" json:"handSize"` // 手型分类 (small/medium/large)
	LastMeasuredAt   time.Time          `bson:"lastMeasuredAt,omitempty" json:"lastMeasuredAt,omitempty"`
	MeasurementCount int                `bson:"measurementCount" json:"measurementCount"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// 集合名常量
const (
	MeasurementsCollection         = "measurements"
	MeasurementUserStatsCollection = "measurement_user_stats"
)
