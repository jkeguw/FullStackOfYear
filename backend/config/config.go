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
	Email   EmailConfig   `yaml:"email"`
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

type OAuthConfig struct {
	Google struct {
		ClientID     string   `yaml:"clientId"`
		ClientSecret string   `yaml:"clientSecret"`
		RedirectURL  string   `yaml:"redirectUrl"`
		Scopes       []string `yaml:"scopes"`
	} `yaml:"google"`
}

// EmailConfig defines email service configuration
type EmailConfig struct {
	SMTP struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"smtp"`
	From      string            `yaml:"from"`
	BaseURL   string            `yaml:"baseUrl"`
	Templates map[string]string `yaml:"templates"`
}

func LoadConfig() (*Config, error) {
	return LoadConfigFromPath("config/config.yaml")
}

func LoadConfigFromPath(path string) (*Config, error) {
	config := &Config{}

	// read yaml config
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	// OAuth environment variable override
	if clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID"); clientID != "" {
		config.OAuth.Google.ClientID = clientID
	}
	if clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"); clientSecret != "" {
		config.OAuth.Google.ClientSecret = clientSecret
	}
	if redirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"); redirectURL != "" {
		config.OAuth.Google.RedirectURL = redirectURL
	}

	// Simplified validation - skip complex validations for now
	
	// Email configuration environment variable override
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		config.Email.SMTP.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		config.Email.SMTP.Password = password
	}

	// Simplified email validation - skip complex validations for now

	return config, nil
}