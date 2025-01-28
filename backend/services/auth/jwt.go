package auth

import (
	"FullStackOfYear/backend/config"
	"FullStackOfYear/backend/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Device string `json:"device"`
	jwt.RegisteredClaims
}

// GenerateToken Generate JWT token
func GenerateToken(userID, role, deviceID string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		Device: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)), // access token 2小时
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWT.Secret))
}

// ParseToken Analysis JWT token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, errors.NewAppError(errors.Unauthorized, "Invalid token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewAppError(errors.Unauthorized, "Invalid token claims")
}
