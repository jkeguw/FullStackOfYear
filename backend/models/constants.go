package models

import "time"

const (
	// Common constants
	MaxLoginHistory       = 50
	MaxConcurrentSessions = 5

	// Token related
	VerifyTokenExpiration = 24 * time.Hour

	// Role types
	RoleUser     = "user"
	RoleReviewer = "reviewer"
	RoleAdmin    = "admin"
)
