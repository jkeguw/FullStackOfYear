package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/services/email"
	"github.com/gin-gonic/gin"
)

type EmailVerificationHandler struct {
	authService  *AuthService
	emailService *email.Service
}

func NewEmailVerificationHandler(authService *AuthService, emailService *email.Service) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		authService:  authService,
		emailService: emailService,
	}
}

// SendVerification sends verification email
func (h *EmailVerificationHandler) SendVerification(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		errors.AbortWithError(c, errors.Unauthorized, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(c, userID)
	if err != nil {
		errors.AbortWithError(c, errors.InternalError, "Failed to get user")
		return
	}

	// Generate verification token
	token, err := h.authService.GenerateEmailVerificationToken(user)
	if err != nil {
		errors.AbortWithError(c, errors.InternalError, "Failed to generate token")
		return
	}

	// Send verification email
	err = h.emailService.SendVerificationEmail(user.Email, user.Username, token)
	if err != nil {
		errors.AbortWithError(c, errors.InternalError, "Failed to send verification email")
		return
	}

	c.JSON(200, gin.H{"message": "Verification email sent"})
}

// VerifyEmail verifies email with token
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		errors.AbortWithError(c, errors.BadRequest, "Missing verification token")
		return
	}

	// Verify token and update user
	err := h.authService.VerifyEmailToken(c, token)
	if err != nil {
		errors.AbortWithError(c, errors.BadRequest, "Invalid or expired token")
		return
	}

	c.JSON(200, gin.H{"message": "Email verified successfully"})
}
