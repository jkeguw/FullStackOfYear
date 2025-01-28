package auth

import (
	"FullStackOfYear/backend/internal/database"
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/services/auth"
	"FullStackOfYear/backend/services/limiter"
	"FullStackOfYear/backend/services/token"
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
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *models.User `json:"user"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Validate password strength
	if !auth.ValidatePassword(req.Password) {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Password too weak"))
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid email format"))
		return
	}

	// Check if user exists
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

	// Create new user
	hashedPassword, _ := auth.HashPassword(req.Password)
	user := models.NewUser(req.Username, req.Email, hashedPassword)

	_, err = collection.InsertOne(c, user)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to create user"))
		return
	}

	c.JSON(200, user)
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Get client IP
	clientIP := c.ClientIP()

	// Check rate limit
	allowed, err := limiter.CheckRateLimit(c, clientIP)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Rate limit check failed"))
		return
	}
	if !allowed {
		c.JSON(429, errors.NewAppError(errors.TooManyRequests, "Too many login attempts"))
		return
	}

	// Check login attempts
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

	collection := database.MongoClient.Database("cpc").Collection("users")
	var user models.User
	err = collection.FindOne(c, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Record failure
		limiter.RecordLoginFailure(c, req.Email)
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid credentials"))
		return
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		// Record failure
		limiter.RecordLoginFailure(c, req.Email)
		c.JSON(401, errors.NewAppError(errors.Unauthorized, "Invalid credentials"))
		return
	}

	// Generate tokens
	tokenManager := token.NewManager(database.RedisClient)
	accessToken, refreshToken, err := tokenManager.GenerateTokenPair(user.ID.Hex(), user.Role.Type, req.DeviceID)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to generate token"))
		return
	}

	// Clear login failures
	limiter.ClearLoginFailure(c, req.Email)

	// Update last login time
	_, err = collection.UpdateOne(c, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"stats.lastLoginAt": time.Now(),
			"stats.lastLoginIP": clientIP,
		},
	})

	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to update last login"))
		return
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &user,
	}

	c.JSON(200, response)
}

type LogoutRequest struct {
	DeviceID string `json:"deviceId" binding:"required"`
}

// Logout handles user logout request
func Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, errors.NewAppError(errors.BadRequest, "Invalid request parameters"))
		return
	}

	// Get user info from context (set by Auth middleware)
	userID := c.GetString("userID")
	deviceID := req.DeviceID

	// Create token manager
	tokenManager := token.NewManager(database.RedisClient)

	// Invalidate tokens for this device
	err := tokenManager.InvalidateTokens(c, userID, deviceID)
	if err != nil {
		c.JSON(500, errors.NewAppError(errors.InternalError, "Failed to logout"))
		return
	}

	c.JSON(200, gin.H{"message": "Successfully logged out"})
}
