package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// MeasurementQuality 测量质量
type MeasurementQuality string

const (
	QualityLow     MeasurementQuality = "low"
	QualityMedium  MeasurementQuality = "medium"
	QualityHigh    MeasurementQuality = "high"
	QualityUnknown MeasurementQuality = "unknown"
)

// Measurement 测量记录模型
type Measurement struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID  `bson:"userId" json:"userId"`
	Palm        float64             `bson:"palm" json:"palm"`           // 手掌宽度
	Length      float64             `bson:"length" json:"length"`       // 中指长度
	Unit        string              `bson:"unit" json:"unit"`           // 单位: mm, cm, inch
	Device      string              `bson:"device" json:"device"`       // 测量设备
	Calibrated  bool                `bson:"calibrated" json:"calibrated"` // 是否校准
	Quality     *MeasurementQuality `bson:"quality" json:"quality"`     // 测量质量: low, medium, high
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt" json:"updatedAt"`
}

// HandSize 手型大小
type HandSize string

const (
	HandSizeSmall  HandSize = "small"
	HandSizeMedium HandSize = "medium"
	HandSizeLarge  HandSize = "large"
)

// GripStyle 握持风格
type GripStyle string

const (
	GripStylePalm       GripStyle = "palm"
	GripStyleClaw       GripStyle = "claw"
	GripStyleFingertip  GripStyle = "fingertip"
	GripStyleUnknown    GripStyle = "unknown"
)

// MeasurementStats 测量统计
type MeasurementStats struct {
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`
	AveragePalm      float64            `bson:"averagePalm" json:"averagePalm"`
	AverageLength    float64            `bson:"averageLength" json:"averageLength"`
	HandSize         HandSize           `bson:"handSize" json:"handSize"`
	MeasurementCount int                `bson:"measurementCount" json:"measurementCount"`
	LastMeasuredAt   time.Time          `bson:"lastMeasuredAt" json:"lastMeasuredAt"`
	UpdatedAt        time.Time          `bson:"updatedAt" json:"updatedAt"`
}