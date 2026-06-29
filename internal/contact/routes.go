package contact

import (
	"blomiva-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtService *auth.JWTService) {
	contacts := rg.Group("/contacts")

	// Public: app users/visitors send messages.
	contacts.POST("", handler.Create)

	// Protected: super admin reads messages.
	admin := contacts.Group("")
	admin.Use(auth.AuthMiddleware(jwtService))
	{
		admin.GET("", handler.List)
		admin.GET("/:id", handler.Get)
		admin.PATCH("/:id/status", handler.UpdateStatus)
	}
}
