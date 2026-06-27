package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// CSRFProtection returns a real CSRF protection middleware.
// It expects the frontend to read the X-XSRF-TOKEN response header/cookie and send it back
// in the X-XSRF-TOKEN request header for state-changing requests.
func CSRFProtection() gin.HandlerFunc {
	authKey := []byte(os.Getenv("CSRF_AUTH_KEY"))
	if len(authKey) == 0 {
		// 默认32字节key，仅用于开发；生产环境必须通过 CSRF_AUTH_KEY 环境变量设置
		authKey = []byte("default-csrf-auth-key-32bytes!!")
	}

	secure := os.Getenv("SECURE_COOKIE") == "true" || os.Getenv("ENV") == "production"

	protect := csrf.Protect(
		authKey,
		csrf.Secure(secure),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Path("/"),
		csrf.RequestHeader("X-XSRF-TOKEN"),
		csrf.CookieName("XSRF-TOKEN"),
	)

	return func(c *gin.Context) {
		served := false
		protect(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			token := csrf.Token(r)
			c.Set("csrf_token", token)
			c.Header("X-XSRF-TOKEN", token)
			if !served {
				served = true
				c.Next()
			}
		})).ServeHTTP(c.Writer, c.Request)
	}
}
