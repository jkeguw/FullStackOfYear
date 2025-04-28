package models

// 此文件用于修复类型定义冲突问题
// 它会声明一些全局类型，以便统一使用

// Role 角色类型
type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleReviewer  Role = "reviewer"
	RoleModerator Role = "moderator"
)
