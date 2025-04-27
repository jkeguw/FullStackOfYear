package models

import "time"

// Device 在登录流程中使用，用于记录用户登录设备信息
// 注意：这个结构体与 device.go 中的硬件设备结构体不同
type Device struct {
	ID         string    `bson:"id" json:"id"`                   // 设备ID
	Name       string    `bson:"name,omitempty" json:"name,omitempty"`           // 设备名称
	Type       string    `bson:"type,omitempty" json:"type,omitempty"`           // 设备类型
	OS         string    `bson:"os,omitempty" json:"os,omitempty"`               // 操作系统
	Browser    string    `bson:"browser,omitempty" json:"browser,omitempty"`     // 浏览器
	IP         string    `bson:"ip,omitempty" json:"ip,omitempty"`               // IP地址
	LastUsedAt time.Time `bson:"lastUsedAt" json:"lastUsedAt"`   // 最后使用时间
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`     // 创建时间
	Trusted    bool      `bson:"trusted" json:"trusted"`         // 是否信任设备
}