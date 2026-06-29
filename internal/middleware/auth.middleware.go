package middleware

import (
	"net/http"
	"strings"

	"blomiva-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func normalizeRole(role string) string {
	role = strings.TrimSpace(strings.ToLower(role))
	role = strings.ReplaceAll(role, "_", "-")

	switch role {
	case "superadmin":
		return "super-admin"
	case "schooladmin":
		return "school-admin"
	default:
		return role
	}
}

func isAdminRole(role string) bool {
	switch normalizeRole(role) {
	case "owner", "admin", "super-admin", "school-admin":
		return true
	default:
		return false
	}
}

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "empty token"})
			c.Abort()
			return
		}

		claims, err := jwtService.Validate(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		role := normalizeRole(claims.Role)

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", role)
		c.Set("schoolId", claims.SchoolID)
		c.Set("isAdmin", isAdminRole(role))

		c.Next()
	}
}
