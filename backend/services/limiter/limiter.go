package limiter

import (
	"project/backend/internal/database"
	"context"
	"fmt"
	"time"
)

const (
	// Request rate limit
	LoginRateLimit  = 5
	RateLimitWindow = time.Minute
)

// RateLimitKey Generate a Redis key for request frequency limit
func RateLimitKey(ip string) string {
	return fmt.Sprintf("rate:limit:login:%s", ip)
}

// CheckRateLimit Check request rate limit
func CheckRateLimit(ctx context.Context, ip string) (bool, error) {
	redisClient := database.RedisClient
	key := RateLimitKey(ip)

	pipe := redisClient.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, RateLimitWindow)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count, err := incr.Result()
	if err != nil {
		return false, err
	}

	return count <= LoginRateLimit, nil
}