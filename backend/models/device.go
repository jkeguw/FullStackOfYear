package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// DeviceTypeEnum 设备类型枚举
type DeviceTypeEnum string

const (
	DeviceTypeMouse     DeviceTypeEnum = "mouse"     // 鼠标
	DeviceTypeKeyboard  DeviceTypeEnum = "keyboard"  // 键盘
	DeviceTypeMonitor   DeviceTypeEnum = "monitor"   // 显示器
	DeviceTypeMousepad  DeviceTypeEnum = "mousepad"  // 鼠标垫
	DeviceTypeAccessory DeviceTypeEnum = "accessory" // 配件
)

// HardwareDevice 外设基础信息
type HardwareDevice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`               // 设备名称
	Brand       string             `bson:"brand" json:"brand"`             // 品牌
	Type        DeviceTypeEnum     `bson:"type" json:"type"`               // 设备类型
	ImageURL    string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"` // 图片URL
	Description string             `bson:"description,omitempty" json:"description,omitempty"` // 设备描述
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// MouseDevice 鼠标设备
type MouseDevice struct {
	HardwareDevice      `bson:",inline"`
	Dimensions  MouseDimensions  `bson:"dimensions" json:"dimensions"`   // 尺寸信息
	Shape       MouseShape       `bson:"shape" json:"shape"`             // 形状信息
	Technical   MouseTechnical   `bson:"technical" json:"technical"`     // 技术参数
	Recommended MouseRecommended `bson:"recommended" json:"recommended"` // 推荐信息
}

// MouseDimensions 鼠标尺寸信息
type MouseDimensions struct {
	Length float64 `bson:"length" json:"length"` // 长度(mm)
	Width  float64 `bson:"width" json:"width"`   // 宽度(mm)
	Height float64 `bson:"height" json:"height"` // 高度(mm)
	Weight float64 `bson:"weight" json:"weight"` // 重量(g)
}

// MouseShape 鼠标形状信息
type MouseShape struct {
	Type              string `bson:"type" json:"type"`                             // 形状类型 (ergonomic, ambidextrous)
	HumpPlacement     string `bson:"humpPlacement" json:"humpPlacement"`           // 坑位位置 (front, center, back)
	FrontFlare        string `bson:"frontFlare" json:"frontFlare"`                 // 前端开叉 (narrow, medium, wide)
	SideCurvature     string `bson:"sideCurvature" json:"sideCurvature"`           // 侧面曲线 (straight, curved)
	HandCompatibility string `bson:"handCompatibility" json:"handCompatibility"`   // 手型适配 (small, medium, large)
}

// MouseTechnical 鼠标技术参数
type MouseTechnical struct {
	Connectivity []string `bson:"connectivity" json:"connectivity"` // 连接方式 (wired, wireless, bluetooth)
	Sensor       string   `bson:"sensor" json:"sensor"`             // 传感器型号
	MaxDPI       int      `bson:"maxDPI" json:"maxDPI"`             // 最大DPI
	PollingRate  int      `bson:"pollingRate" json:"pollingRate"`   // 轮询率(Hz)
	SideButtons  int      `bson:"sideButtons" json:"sideButtons"`   // 侧键数量
	Weight       float64  `bson:"weight,omitempty" json:"weight,omitempty"`           // 重量(g)
	Battery      *Battery `bson:"battery,omitempty" json:"battery,omitempty"`         // 电池信息
}

// Battery 电池信息
type Battery struct {
	Type     string `bson:"type" json:"type"`         // 电池类型
	Capacity int    `bson:"capacity" json:"capacity"` // 电池容量
	Life     int    `bson:"life" json:"life"`         // 电池寿命(小时)
}

// MouseRecommended 鼠标推荐信息
type MouseRecommended struct {
	GameTypes    []string `bson:"gameTypes" json:"gameTypes"`       // 适合游戏类型
	GripStyles   []string `bson:"gripStyles" json:"gripStyles"`     // 适合握持方式
	HandSizes    []string `bson:"handSizes" json:"handSizes"`       // 适合手型大小
	DailyUse     bool     `bson:"dailyUse" json:"dailyUse"`         // 适合日常使用
	Professional bool     `bson:"professional" json:"professional"` // 专业级
}

// KeyboardDevice 键盘设备
type KeyboardDevice struct {
	HardwareDevice    `bson:",inline"`
	Layout    string   `bson:"layout" json:"layout"`       // 配列
	Switches  []string `bson:"switches" json:"switches"`   // 轴体
	Size      string   `bson:"size" json:"size"`           // 尺寸规格
	Technical struct {
		Connectivity []string `bson:"connectivity" json:"connectivity"` // 连接方式
		Keycaps      string   `bson:"keycaps" json:"keycaps"`           // 键帽材质
		HotSwap      bool     `bson:"hotSwap" json:"hotSwap"`           // 热插拔支持
		RGBLighting  bool     `bson:"rgbLighting" json:"rgbLighting"`   // RGB灯光
		NKeyRollover bool     `bson:"nKeyRollover" json:"nKeyRollover"` // N键无冲
	} `bson:"technical" json:"technical"`
	Recommended struct {
		GameTypes []string `bson:"gameTypes" json:"gameTypes"` // 适合游戏类型
		DailyUse  bool     `bson:"dailyUse" json:"dailyUse"`   // 适合日常使用
		Portable  bool     `bson:"portable" json:"portable"`   // 便携性
	} `bson:"recommended" json:"recommended"`
}

// MonitorDevice 显示器设备
type MonitorDevice struct {
	HardwareDevice    `bson:",inline"`
	Size      float64 `bson:"size" json:"size"`             // 尺寸(英寸)
	Resolution struct {
		Width  int `bson:"width" json:"width"`
		Height int `bson:"height" json:"height"`
	} `bson:"resolution" json:"resolution"`
	Technical struct {
		RefreshRate    int     `bson:"refreshRate" json:"refreshRate"`       // 刷新率(Hz)
		ResponseTime   float64 `bson:"responseTime" json:"responseTime"`     // 响应时间(ms)
		PanelType      string  `bson:"panelType" json:"panelType"`           // 面板类型
		AspectRatio    string  `bson:"aspectRatio" json:"aspectRatio"`       // 纵横比
		Curvature      int     `bson:"curvature,omitempty" json:"curvature,omitempty"` // 曲率
		HDRSupport     bool    `bson:"hdrSupport" json:"hdrSupport"`         // HDR支持
		AdaptiveSync   string  `bson:"adaptiveSync,omitempty" json:"adaptiveSync,omitempty"` // 自适应同步技术
	} `bson:"technical" json:"technical"`
	Recommended struct {
		GameTypes []string `bson:"gameTypes" json:"gameTypes"` // 适合游戏类型
		Content   []string `bson:"content" json:"content"`     // 适合内容类型
		ProUse    bool     `bson:"proUse" json:"proUse"`       // 专业用途
	} `bson:"recommended" json:"recommended"`
}

// MousepadDevice 鼠标垫设备
type MousepadDevice struct {
	HardwareDevice    `bson:",inline"`
	Size      struct {
		Length float64 `bson:"length" json:"length"` // 长度(mm)
		Width  float64 `bson:"width" json:"width"`   // 宽度(mm)
		Height float64 `bson:"height" json:"height"` // 高度/厚度(mm)
	} `bson:"size" json:"size"`
	Material    string   `bson:"material" json:"material"`       // 材质
	Surface     string   `bson:"surface" json:"surface"`         // 表面类型
	Base        string   `bson:"base" json:"base"`               // 底座类型
	Recommended []string `bson:"recommended" json:"recommended"` // 推荐场景
}

// DeviceReview 设备评测
type DeviceReview struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID       primitive.ObjectID `bson:"deviceId" json:"deviceId"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Content        string             `bson:"content" json:"content"`
	Pros           []string           `bson:"pros" json:"pros"`
	Cons           []string           `bson:"cons" json:"cons"`
	Score          float64            `bson:"score" json:"score"` // 评分(1-5)
	Usage          string             `bson:"usage" json:"usage"` // 使用场景
	Status         string             `bson:"status" json:"status"` // 状态(pending, approved, rejected)
	Type           ReviewType         `bson:"type" json:"type"`     // 评测类型(mouse, keyboard等)
	ContentType    ReviewContentType  `bson:"contentType" json:"contentType"` // 内容类型(single, comparison等)
	ReviewerID     *primitive.ObjectID `bson:"reviewerId,omitempty" json:"reviewerId,omitempty"` // 评审员ID
	ReviewerNotes  string              `bson:"reviewerNotes,omitempty" json:"reviewerNotes,omitempty"` // 评审员备注
	ReviewedAt     *time.Time          `bson:"reviewedAt,omitempty" json:"reviewedAt,omitempty"` // 评审时间
	PublishedAt    *time.Time          `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"` // 发布时间
	FeaturedRank   *int                `bson:"featuredRank,omitempty" json:"featuredRank,omitempty"` // 推荐排名
	ViewCount      int                 `bson:"viewCount" json:"viewCount"` // 查看次数
	CreatedAt      time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time           `bson:"updatedAt" json:"updatedAt"`
	DeletedAt      *time.Time          `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// UserDevice 用户设备配置
type UserDevice struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID   `bson:"userId" json:"userId"`
	Name        string               `bson:"name" json:"name"` // 配置名称
	Description string               `bson:"description,omitempty" json:"description,omitempty"`
	Devices     []UserDeviceSettings `bson:"devices" json:"devices"`
	IsPublic    bool                 `bson:"isPublic" json:"isPublic"` // 是否公开
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// UserDeviceSettings 用户设备设置
type UserDeviceSettings struct {
	DeviceID  primitive.ObjectID `bson:"deviceId" json:"deviceId"`
	DeviceType DeviceTypeEnum    `bson:"deviceType" json:"deviceType"`
	Settings  map[string]any     `bson:"settings,omitempty" json:"settings,omitempty"` // 自定义设置
}

// 集合名常量
const (
	DevicesCollection      = "devices"
	DeviceReviewsCollection = "device_reviews"
	UserDevicesCollection   = "user_devices"
)