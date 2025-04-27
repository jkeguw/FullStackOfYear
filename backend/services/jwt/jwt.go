package jwt

import (
	"project/backend/config"
	"project/backend/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID   string
	Role     string
	DeviceID string
	Type     string
}

// Service provides JWT functionality
type Service interface {
	GenerateToken(claims Claims, expireTime time.Duration) (string, time.Time, error)
	ParseToken(tokenString string) (*Claims, error)
}

type jwtService struct {
	config config.JWTConfig
}

// NewService creates a new JWT service
func NewService(config config.JWTConfig) Service {
	return &jwtService{
		config: config,
	}
}

// GenerateToken creates a new JWT token
func (s *jwtService) GenerateToken(claims Claims, expireTime time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(expireTime)

	customClaims := &jwt.MapClaims{
		"uid":      claims.UserID,
		"role":     claims.Role,
		"deviceId": claims.DeviceID,
		"type":     claims.Type,
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"iss":      s.config.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	signedToken, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// ParseToken validates and parses a JWT token
func (s *jwtService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token")
	}

	if !token.Valid {
		return nil, errors.NewAppError(errors.Unauthorized, "Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token claims")
	}

	// Check expiration time explicitly
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.NewAppError(errors.Unauthorized, "Missing expiration time")
	}
	
	if time.Now().Unix() > int64(exp) {
		return nil, errors.NewAppError(errors.Unauthorized, "Token has expired")
	}

	return &Claims{
		UserID:   claims["uid"].(string),
		Role:     claims["role"].(string),
		DeviceID: claims["deviceId"].(string),
		Type:     claims["type"].(string),
	}, nil
}

// 兼容旧代码的函数
func GenerateToken(claims Claims, expireTime time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(expireTime)

	customClaims := &jwt.MapClaims{
		"uid":      claims.UserID,
		"role":     claims.Role,
		"deviceId": claims.DeviceID,
		"type":     claims.Type,
		"exp":      expiresAt.Unix(),
		"iat":      now.Unix(),
		"nbf":      now.Unix(),
		"iss":      config.Cfg.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	signedToken, err := token.SignedString([]byte(config.Cfg.JWT.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiresAt, nil
}

// 兼容旧代码的函数
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token")
	}

	if !token.Valid {
		return nil, errors.NewAppError(errors.Unauthorized, "Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token claims")
	}

	// Check expiration time explicitly
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.NewAppError(errors.Unauthorized, "Missing expiration time")
	}
	
	if time.Now().Unix() > int64(exp) {
		return nil, errors.NewAppError(errors.Unauthorized, "Token has expired")
	}

	return &Claims{
		UserID:   claims["uid"].(string),
		Role:     claims["role"].(string),
		DeviceID: claims["deviceId"].(string),
		Type:     claims["type"].(string),
	}, nil
}