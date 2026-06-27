package middleware

import (
	"html"
	"strings"

	"github.com/gin-gonic/gin"
)

// XSSProtection sets security headers that help mitigate XSS attacks.
// Note: headers alone are not sufficient; user input must still be sanitized before storage.
func XSSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Next()
	}
}

// SanitizeHTML escapes HTML special characters to prevent stored XSS.
// Use this on any user-provided text that will be rendered as HTML or embedded in pages.
func SanitizeHTML(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

// SanitizeText is an alias for SanitizeHTML for plain text fields.
func SanitizeText(input string) string {
	return SanitizeHTML(input)
}
