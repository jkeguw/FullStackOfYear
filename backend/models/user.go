package models

import (
	"FullStackOfYear/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Role     UserRole           `bson:"role" json:"role"`
	Stats    UserStats          `bson:"stats" json:"stats"`
	OAuth    *OAuthInfo         `bson:"oauth,omitempty" json:"oauth,omitempty"`
}

type UserRole struct {
	Type                string               `bson:"type" json:"type"`
	ReviewerApplication *ReviewerApplication `bson:"reviewerApplication,omitempty" json:"reviewerApplication,omitempty"`
	InviteCode          string               `bson:"inviteCode,omitempty" json:"inviteCode,omitempty"`
}

type ReviewerApplication struct {
	Status      string    `bson:"status" json:"status"`
	AppliedAt   time.Time `bson:"appliedAt" json:"appliedAt"`
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
}

type UserStats struct {
	ReviewCount int       `bson:"reviewCount" json:"reviewCount"`
	TotalWords  int       `bson:"totalWords" json:"totalWords"`
	Violations  int       `bson:"violations" json:"violations"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
}

type OAuthInfo struct {
	Google *GoogleOAuth `bson:"google,omitempty" json:"google,omitempty"`
}

type GoogleOAuth struct {
	ID          string    `bson:"id" json:"id"`
	Email       string    `bson:"email" json:"email"`
	Connected   bool      `bson:"connected" json:"connected"`
	ConnectedAt time.Time `bson:"connectedAt" json:"connectedAt"`
}

const (
	RoleUser     = "user"
	RoleReviewer = "reviewer"
	RoleAdmin    = "admin"
)

// ValidateRole validate role type
func (u *User) ValidateRole() error {
	validRoles := map[string]bool{
		RoleUser:     true,
		RoleReviewer: true,
		RoleAdmin:    true,
	}

	if !validRoles[u.Role.Type] {
		return errors.NewAppError(errors.BadRequest, "Invalid role type")
	}
	return nil
}

// NewUser create new user
func NewUser(username, email, password string) *User {
	now := time.Now()
	return &User{
		Username: username,
		Email:    email,
		Password: password,
		Role: UserRole{
			Type: RoleUser,
		},
		Stats: UserStats{
			ReviewCount: 0,
			TotalWords:  0,
			Violations:  0,
			CreatedAt:   now,
			LastLoginAt: now,
		},
	}
}
