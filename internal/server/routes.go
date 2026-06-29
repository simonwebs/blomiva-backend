package server

import (
	"net/http"

	"blomiva-backend/internal/auth"
	"blomiva-backend/internal/config"
	"blomiva-backend/internal/contact"
	"blomiva-backend/internal/geo"
	"blomiva-backend/internal/profiles"
	"blomiva-backend/internal/settings"
	"blomiva-backend/internal/students"
	"blomiva-backend/internal/tenant"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(router *gin.Engine, db *mongo.Database, cfg config.Config) {

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"ok":    false,
			"error": "route not found",
			"path":  c.Request.URL.Path,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true, "service": "blomiva"})
	})

	api := router.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok":       true,
			"env":      cfg.Env,
			"database": cfg.DBName,
		})
	})

	// ================= AUTH =================
	authRepo := auth.NewRepository(db)
	jwt := auth.NewJWTService()
	authService := auth.NewService(authRepo, jwt)
	authHandler := auth.NewHandler(authService)

	registerAuth(api, authHandler)

	// ================= MODULES =================
	registerGeo(api)
	registerContact(api, db, jwt)
	registerTenant(api, db, cfg, jwt)

	protected := api.Group("")
	protected.Use(auth.AuthMiddleware(jwt))

	protected.GET("/auth/me", authHandler.Me)

	registerStudents(protected, db)
	registerSettings(protected, db)
	registerProfile(router, db, jwt)
}