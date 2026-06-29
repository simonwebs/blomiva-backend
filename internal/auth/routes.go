package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	handler *Handler,
	passwordlessHandler *PasswordlessHandler,
	jwtService *JWTService,
) {
	auth := rg.Group("/auth")

	// ===================== CORE AUTH =====================
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)
	auth.GET("/me", AuthMiddleware(jwtService), handler.Me)

	// ===================== PASSWORDLESS =====================
	auth.POST("/passwordless/request", passwordlessHandler.RequestLogin)
	auth.POST("/passwordless/verify", passwordlessHandler.VerifyLogin)

	// ===================== SYSTEM (DEV ONLY) =====================

	auth.POST("/seed/super-admin", handler.SeedSuperAdmin)

	auth.GET("/debug/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"service": "auth",
		})
	})

	auth.GET("/validate", AuthMiddleware(jwtService), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":   true,
			"user": c.GetString("userId"),
			"role": c.GetString("role"),
		})
	})
}
