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

	return config, nil
}
