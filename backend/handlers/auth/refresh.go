package auth

import (
	"FullStackOfYear/backend/internal/database"
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/services/jwt"
	"FullStackOfYear/backend/services/token"
	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	DeviceID     string `json:"deviceId" binding:"required"`
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Parse and validate refresh token
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid refresh token"))
		return
	}

	// Verify token type
	if claims.Type != "refresh" {
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid token type"))
		return
	}

	// Verify device ID matches
	if claims.DeviceID != req.DeviceID {
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Device ID mismatch"))
		return
	}

	tokenManager := token.NewManager(database.RedisClient)

	// Check if refresh token exists
	exists, err := tokenManager.CheckTokenExists(c, claims.UserID, claims.DeviceID, "refresh")
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Token verification failed"))
		return
	}
	if !exists {
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Token has been revoked"))
		return
	}

	// Generate new token pair
	accessToken, refreshToken, err := tokenManager.GenerateTokenPair(claims.UserID, claims.Role, claims.DeviceID)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to generate tokens"))
		return
	}

	c.JSON(200, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
