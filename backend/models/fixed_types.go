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

// 额外的集合名常量，以避免重复定义
const (
	// 其他集合名已在各自的文件中定义
	// 这里我们确保 ReviewsCollection 只定义一次
	// ReviewsCollection = "reviews" // 已在 constants.go 中定义
)