package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt limitation
	bcryptCost        = 12
)

// ValidatePassword checks password strength based on rules
func ValidatePassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= minPasswordLength && len(password) <= maxPasswordLength {
		hasMinLen = true
	}

	// Check each character for meeting complexity requirements
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// All requirements must be met
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

// HashPassword encrypts password using bcrypt
func HashPassword(password string) (string, error) {
	if !ValidatePassword(password) {
		return "", errors.NewAppError(errors.BadRequest, "Password does not meet requirements")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", errors.NewAppError(errors.InternalError, "Password hashing failed")
	}

	return string(bytes), nil
}

// CheckPassword verifies if the provided password matches the hash
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength returns detailed password validation errors
func ValidatePasswordStrength(password string) []string {
	var validationErrors []string

	if len(password) < minPasswordLength {
		validationErrors = append(validationErrors, "Password must be at least 8 characters long")
	}

	if len(password) > maxPasswordLength {
		validationErrors = append(validationErrors, "Password exceeds maximum length")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		validationErrors = append(validationErrors, "Password must contain at least one uppercase letter")
	}
	if !hasLower {
		validationErrors = append(validationErrors, "Password must contain at least one lowercase letter")
	}
	if !hasNumber {
		validationErrors = append(validationErrors, "Password must contain at least one number")
	}
	if !hasSpecial {
		validationErrors = append(validationErrors, "Password must contain at least one special character")
	}

	return validationErrors
}
