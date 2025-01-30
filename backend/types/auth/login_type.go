package auth

// LoginType defines the type of login
type LoginType string

const (
	// EmailLogin represents login with email and password
	EmailLogin LoginType = "email"
	// GoogleLogin represents login with Google OAuth
	GoogleLogin LoginType = "google"
)

// IsValid checks if the login type is valid
func (lt LoginType) IsValid() bool {
	switch lt {
	case EmailLogin, GoogleLogin:
		return true
	default:
		return false
	}
}
