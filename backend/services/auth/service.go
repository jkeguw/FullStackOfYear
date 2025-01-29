package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/auth"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// AuthService handles authentication related operations
type AuthService struct {
	userCollection *mongo.Collection
	tokenService   TokenService
}

// TokenService interface for token operations
type TokenService interface {
	GenerateTokenPair(userID string, role string, deviceID string) (accessToken string, refreshToken string, err error)
}

// NewAuthService creates a new authentication service
func NewAuthService(userCollection *mongo.Collection, tokenService TokenService) *AuthService {
	return &AuthService{
		userCollection: userCollection,
		tokenService:   tokenService,
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
