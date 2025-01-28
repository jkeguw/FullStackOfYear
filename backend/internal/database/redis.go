package database

import (
	"FullStackOfYear/backend/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})
	return RedisClient
}
