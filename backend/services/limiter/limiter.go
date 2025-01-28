package limiter

import (
	"FullStackOfYear/backend/internal/database"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	// Login failure limit
	MaxLoginAttempts  = 5
	BaseBlockDuration = time.Minute
	MaxBlockDuration  = 24 * time.Hour

	// Request rate limit
	LoginRateLimit  = 5
	RateLimitWindow = time.Minute
)

// calculateBlockDuration Calculating ban time
func calculateBlockDuration(attempts int64) time.Duration {
	blockDuration := BaseBlockDuration * time.Duration(1<<uint(attempts-MaxLoginAttempts))
	if blockDuration > MaxBlockDuration {
		blockDuration = MaxBlockDuration
	}
	return blockDuration
}

// LoginFailureKey Generate a Redis key for failed login records
func LoginFailureKey(email string) string {
	return fmt.Sprintf("login:failure:%s", email)
}

// LoginBlockKey Generate a Redis key for login ban
func LoginBlockKey(email string) string {
	return fmt.Sprintf("login:block:%s", email)
}

// RateLimitKey Generate a Redis key for request frequency limit
func RateLimitKey(ip string) string {
	return fmt.Sprintf("rate:limit:login:%s", ip)
}

// CheckLoginAttempts Check the number of failed logins and update
func CheckLoginAttempts(ctx context.Context, email string) (bool, time.Duration, error) {
	rdb := database.RedisClient
	blockKey := LoginBlockKey(email)
	failureKey := LoginFailureKey(email)

	pipe := rdb.Pipeline()
	blockTTL := pipe.TTL(ctx, blockKey)
	failureCount := pipe.Get(ctx, failureKey)

	_, err := pipe.Exec(ctx)
	// redis.Nil 表示 key 不存在
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, 0, err
	}

	// 检查是否被封禁
	ttl, err := blockTTL.Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, 0, err
	}
	if ttl > 0 {
		return false, ttl, nil
	}

	// 获取失败次数
	count, err := failureCount.Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, 0, err
	}
	// 如果key不存在，设置count为0
	if errors.Is(err, redis.Nil) {
		count = 0
	}

	if count >= MaxLoginAttempts {
		blockDuration := calculateBlockDuration(count)

		// 设置封禁，使用pipeline保证原子性
		pipe := rdb.Pipeline()
		pipe.Set(ctx, blockKey, 1, blockDuration)
		pipe.Del(ctx, failureKey)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return false, 0, err
		}

		return false, blockDuration, nil
	}

	return true, 0, nil
}

// RecordLoginFailure record login failed
func RecordLoginFailure(ctx context.Context, email string) error {
	redisClient := database.RedisClient
	key := LoginFailureKey(email)

	pipe := redisClient.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	count, err := incr.Result()
	if err != nil {
		return err
	}

	if count >= MaxLoginAttempts {
		blockDuration := calculateBlockDuration(count)
		pipe := redisClient.Pipeline()
		pipe.Set(ctx, LoginBlockKey(email), 1, blockDuration)
		pipe.Del(ctx, key) // 清除失败计数
		_, err = pipe.Exec(ctx)
		return err
	}

	return nil
}

// ClearLoginFailure Clear failed login records
func ClearLoginFailure(ctx context.Context, email string) error {
	redisClient := database.RedisClient
	return redisClient.Del(ctx, LoginFailureKey(email)).Err()
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
