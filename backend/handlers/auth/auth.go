package auth

import (
	"github.com/gin-gonic/gin"
	"project/backend/internal/errors"
)

// Register handles user registration
func Register(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "User registration is not implemented yet"))
}

// Login handles user login
func Login(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "User login is not implemented yet"))
}

// RefreshToken handles token refresh requests
func RefreshToken(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "Token refresh is not implemented yet"))
}

// Logout handles user logout
func Logout(c *gin.Context) {
	c.JSON(500, errors.NewAppError(errors.NotImplemented, "User logout is not implemented yet"))
}