package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/auth"
	"context"
	"crypto/rand"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// AuthService handles authentication related operations
type AuthService struct {
	userCollection *mongo.Collection
	tokenService   TokenService
	emailService   EmailService
}

// TokenService interface for token operations
type TokenService interface {
	GenerateTokenPair(userID string, role string, deviceID string) (accessToken string, refreshToken string, err error)
}

type EmailService interface {
	SendVerificationEmail(to, username, token string) error
}

// NewAuthService creates a new authentication service
func NewAuthService(userCollection *mongo.Collection, tokenService TokenService, emailService EmailService) *AuthService {
	return &AuthService{
		userCollection: userCollection,
		tokenService:   tokenService,
		emailService:   emailService,
	}
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

// HandleOAuthLogin processes OAuth login and returns auth response
func (s *AuthService) HandleOAuthLogin(ctx context.Context, userInfo *auth.OAuthUserInfo) (*AuthResponse, error) {
	// Try to find user by OAuth Google ID
	filter := bson.M{"oauth.google.id": userInfo.ID}
	var user models.User
	err := s.userCollection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, errors.NewAppError(errors.InternalError, "Database error")
		}

		// User not found, create new user
		user = models.User{
			ID:       primitive.NewObjectID(),
			Username: userInfo.Email, // Use email as initial username
			Email:    userInfo.Email,
			Role: models.UserRole{
				Type: models.RoleUser,
			},
			Stats: models.UserStats{
				ReviewCount: 0,
				TotalWords:  0,
				Violations:  0,
				CreatedAt:   time.Now(),
				LastLoginAt: time.Now(),
			},
			OAuth: &models.OAuthInfo{
				Google: &models.GoogleOAuth{
					ID:          userInfo.ID,
					Email:       userInfo.Email,
					Connected:   true,
					ConnectedAt: time.Now(),
				},
			},
		}

		_, err = s.userCollection.InsertOne(ctx, user)
		if err != nil {
			return nil, errors.NewAppError(errors.InternalError, "Failed to create user")
		}
	} else {
		// Update existing user's OAuth info
		update := bson.M{
			"$set": bson.M{
				"oauth.google.email":       userInfo.Email,
				"oauth.google.connected":   true,
				"oauth.google.connectedAt": time.Now(),
				"stats.lastLoginAt":        time.Now(),
			},
		}

		_, err = s.userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, errors.NewAppError(errors.InternalError, "Failed to update user")
		}
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(
		user.ID.Hex(),
		user.Role.Type,
		"oauth_"+user.ID.Hex(),
	)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to generate tokens")
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &user,
	}, nil
}

// GenerateEmailVerificationToken generates token for email verification
func (s *AuthService) GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error) {
	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", errors.NewAppError(errors.InternalError, "Failed to generate token")
	}
	token := base64.URLEncoding.EncodeToString(b)

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", errors.NewAppError(errors.BadRequest, "Invalid user ID")
	}

	// Update user with verification token
	update := bson.M{
		"$set": bson.M{
			"status.verifyToken":  token,
			"status.tokenExpires": time.Now().Add(24 * time.Hour),
		},
	}

	result, err := s.userCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return "", errors.NewAppError(errors.InternalError, "Failed to save token")
	}

	if result.MatchedCount == 0 {
		return "", errors.NewAppError(errors.NotFound, "User not found")
	}

	return token, nil
}

// VerifyEmailToken validates the email verification token
func (s *AuthService) VerifyEmailToken(ctx context.Context, token string) error {
	// Find and update user with atomic operation
	update := bson.M{
		"$set": bson.M{
			"status.emailVerified": true,
			"status.verifyToken":   "",
			"status.tokenExpires":  time.Time{},
		},
	}

	result, err := s.userCollection.UpdateOne(ctx,
		bson.M{
			"status.verifyToken":  token,
			"status.tokenExpires": bson.M{"$gt": time.Now()},
		},
		update,
	)

	if err != nil {
		return errors.NewAppError(errors.InternalError, "Database error")
	}

	if result.MatchedCount == 0 {
		return errors.NewAppError(errors.BadRequest, "Invalid or expired token")
	}

	return nil
}

// SendVerificationEmail sends verification email to user
func (s *AuthService) SendVerificationEmail(ctx context.Context, userID string) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Don't send verification if already verified
	if user.Status.EmailVerified {
		return errors.NewAppError(errors.BadRequest, "Email already verified")
	}

	// Generate new token
	token, err := s.GenerateEmailVerificationToken(ctx, userID)
	if err != nil {
		return err
	}

	// Send email
	return s.emailService.SendVerificationEmail(user.Email, user.Username, token)
}

// GetUserByID retrieves user by ID (保留现有的HandleOAuthLogin方法)
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(errors.BadRequest, "Invalid user ID")
	}

	var user models.User
	err = s.userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewAppError(errors.NotFound, "User not found")
		}
		return nil, errors.NewAppError(errors.InternalError, "Database error")
	}

	return &user, nil
}
