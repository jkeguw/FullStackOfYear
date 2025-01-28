package cache

import (
	"FullStackOfYear/backend/internal/database"
	"context"
	"encoding/json"
	"time"
)

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return database.RedisClient.Set(ctx, key, data, expiration).Err()
}

func Get(ctx context.Context, key string, value interface{}) error {
	data, err := database.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
