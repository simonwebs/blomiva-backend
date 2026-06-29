package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
				"path":  c.Request.URL.Path,
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			return
		}

		token := strings.TrimSpace(parts[1])
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "empty bearer token",
			})
			return
		}

		claims, err := jwtService.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		role := normalizeAuthRole(claims.Role)

		c.Set("userId", claims.UserID)
		c.Set("email", normalizeEmail(claims.Email))
		c.Set("role", role)
		c.Set("schoolId", claims.SchoolID)
		c.Set("isAdmin", isAdminRole(role))

		c.Next()
	}
}

func isAdminRole(role string) bool {
	role = normalizeAuthRole(role)

	switch role {
	case "owner", "admin", "super-admin", "school-admin":
		return true
	default:
		return false
	}
}
