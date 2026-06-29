package settings

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(api *gin.RouterGroup, db *mongo.Database, authMiddleware gin.HandlerFunc) {
	service := NewService(db)
	handler := NewHandler(service)

	group := api.Group("/settings")
	group.Use(authMiddleware)

	group.GET("/me", handler.GetMySettings)
	group.PATCH("/me", handler.UpdateMySettings)
	group.DELETE("/me", handler.DeleteMyAccount)
	group.POST("/me/restore", handler.RestoreMyAccount)
}
