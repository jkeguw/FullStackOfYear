package jwt

import (
	"FullStackOfYear/backend/config"
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/types/auth"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// GenerateToken creates a new JWT token
func GenerateToken(claims auth.Claims, expireTime time.Duration) (string, time.Time, error) {
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

// ParseToken validates and parses a JWT token
func ParseToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})

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

	return &auth.Claims{
		UserID:   claims["uid"].(string),
		Role:     claims["role"].(string),
		DeviceID: claims["deviceId"].(string),
		Type:     claims["type"].(string),
	}, nil
}
