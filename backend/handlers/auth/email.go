package auth

import (
	"project/backend/internal/errors"
	"project/backend/services/auth"
	"project/backend/services/email"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EmailVerificationHandler struct {
	authService  auth.Service
	emailService *email.Service
}

func NewEmailVerificationHandler(authService auth.Service, emailService *email.Service) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		authService:  authService,
		emailService: emailService,
	}
}

// SendVerification sends verification email
func (h *EmailVerificationHandler) SendVerification(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppError(errors.Unauthorized, "User not authenticated"),
		})
		return
	}

	// Generate token
	token, err := h.authService.GenerateEmailVerificationToken(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to generate token"),
		})
		return
	}

	// Get user info
	user, err := h.authService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to get user"),
		})
		return
	}

	// Send verification email
	err = h.emailService.SendVerificationEmail(user.Email, user.Username, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to send verification email"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

// VerifyEmail verifies email with token
func (h *EmailVerificationHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppError(errors.BadRequest, "Missing verification token"),
		})
		return
	}

	// Verify token and update user
	err := h.authService.VerifyEmail(c, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppError(errors.BadRequest, "Invalid or expired token"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// UpdateEmail initiates email change process
func (h *EmailVerificationHandler) UpdateEmail(c *gin.Context) {
	var req struct {
		NewEmail string `json:"newEmail" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppError(errors.BadRequest, "Invalid email format"),
		})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.NewAppError(errors.Unauthorized, "User not authenticated"),
		})
		return
	}

	// Generate token for new email
	user, err := h.authService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to get user"),
		})
		return
	}

	// Check if new email is different
	if user.Email == req.NewEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.NewAppError(errors.BadRequest, "New email is same as current"),
		})
		return
	}

	// Generate and send verification
	token, err := h.authService.GenerateEmailChangeToken(user, req.NewEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to generate token"),
		})
		return
	}

	err = h.emailService.SendVerificationEmail(req.NewEmail, user.Username, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.NewAppError(errors.InternalError, "Failed to send verification email"),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent to new address"})
}
