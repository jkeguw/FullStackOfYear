package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
	JWT     JWTConfig     `yaml:"jwt"`
	OAuth   OAuthConfig   `yaml:"oauth"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret        string        `yaml:"secret"`
	AccessExpire  time.Duration `yaml:"accessExpireTime"`
	RefreshExpire time.Duration `yaml:"refreshExpireTime"`
	Issuer        string        `yaml:"issuer"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	// read yaml config
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	// 环境变量覆盖
	if mongoURI := os.Getenv("MONGO_URI"); mongoURI != "" {
		config.MongoDB.URI = mongoURI
	}
	if mongoDB := os.Getenv("MONGO_DATABASE"); mongoDB != "" {
		config.MongoDB.Database = mongoDB
	}
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		config.Redis.Addr = redisAddr
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}

	return config, nil
}

type OAuthConfig struct {
	Google struct {
		ClientID     string   `yaml:"clientId"`
		ClientSecret string   `yaml:"clientSecret"`
		RedirectURL  string   `yaml:"redirectUrl"`
		Scopes       []string `yaml:"scopes"`
	} `yaml:"google"`
}
