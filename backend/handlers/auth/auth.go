package auth

import (
	"FullStackOfYear/backend/internal/database"
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/geoip"
	"FullStackOfYear/backend/services/limiter"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/mail"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	User         *models.User `json:"user"`
}

// Register user register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Verify password strength
	if !auth.ValidatePassword(req.Password) {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Password too weak"))
		return
	}

	// Verify email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid email format"))
		return
	}

	// Check if the user already exists
	collection := database.MongoClient.Database("cpc").Collection("users")
	var existingUser models.User
	err := collection.FindOne(c, bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": req.Username},
		},
	}).Decode(&existingUser)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Database error"))
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "User already exists"))
		return
	}

	// Create New User
	hashedPassword, _ := auth.HashPassword(req.Password)
	user := models.NewUser(req.Username, req.Email, hashedPassword)

	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to create user"))
		return
	}

	c.JSON(200, user)
}

// Login user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Get the client IP
	clientIP := c.ClientIP()

	// Check request frequency
	allowed, err := limiter.CheckRateLimit(c, clientIP)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Rate limit check failed"))
		return
	}
	if !allowed {
		c.JSON(429, errors.NewAppError(errors.TooManyRequests, "Too many login attempts"))
		return
	}

	// Check the number of failed login attempts
	allowed, blockDuration, err := limiter.CheckLoginAttempts(c, req.Email)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Login attempt check failed"))
		return
	}
	if !allowed {
		c.JSON(429, errors.NewAppError(errors.TooManyRequests,
			fmt.Sprintf("Account temporary blocked for %v", blockDuration)))
		return
	}

	// Query User
	collection := database.MongoClient.Database("cpc").Collection("users")
	var user models.User
	err = collection.FindOne(c, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Recording failure
		limiter.RecordLoginFailure(c, req.Email)
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid credentials"))
		return
	}

	// validate password
	if !auth.CheckPassword(req.Password, user.Password) {
		// Recording failure
		limiter.RecordLoginFailure(c, req.Email)
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid credentials"))
		return
	}

	// Check IP Geolocation
	locationSafe, err := geoip.CheckLocation(c, user.ID.Hex(), clientIP)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Location check failed"))
		return
	}
	if !locationSafe {
		// TODO: 发送邮件通知或要求额外验证
		c.JSON(403, errors.NewAppError(errors.Forbidden, "Suspicious login location detected"))
		return
	}

	// Login successful, clear failed records
	limiter.ClearLoginFailure(c, req.Email)

	// Generate Token
	token, err := auth.GenerateToken(user.ID.Hex(), user.Role.Type, req.DeviceID)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to generate token"))
		return
	}

	// Update last login time and IP
	_, err = collection.UpdateOne(c, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"stats.lastLoginAt": time.Now(),
			"stats.lastLoginIP": clientIP,
		},
	})

	response := AuthResponse{
		Token: token,
		User:  &user,
	}

	c.JSON(200, response)
}
