package config

import (
	"fmt"
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

	//Verify OAuth Configuration
	if config.OAuth.Google.ClientID != "" {
		if config.OAuth.Google.ClientSecret == "" {
			return nil, fmt.Errorf("google OAuth client secret is required when client ID is set")
		}
		if config.OAuth.Google.RedirectURL == "" {
			return nil, fmt.Errorf("google OAuth redirect URL is required when client ID is set")
		}
	}

	// Email configuration environment variable override
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		config.Email.SMTP.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		config.Email.SMTP.Password = password
	}

	// Verify Email Configuration
	if config.Email.SMTP.Host == "" {
		return nil, fmt.Errorf("SMTP host is required")
	}
	if config.Email.SMTP.Username == "" {
		return nil, fmt.Errorf("SMTP username is required")
	}
	if config.Email.SMTP.Password == "" {
		return nil, fmt.Errorf("SMTP password is required")
	}

	return config, nil
}
