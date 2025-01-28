package cache

import "time"

const (
	UserPrefix   = "user:"
	DevicePrefix = "device:"
	ReviewPrefix = "review:"
	CacheTimeout = 24 * time.Hour
)
