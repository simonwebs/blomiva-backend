package students

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, handler *Handler) {
	routes := api.Group("/students")

	routes.POST("/", handler.CreateStudent)
	routes.GET("/", handler.ListStudents)
	routes.GET("/:id", handler.GetStudent)
	routes.PUT("/:id", handler.UpdateStudent)
	routes.DELETE("/:id", handler.DeleteStudent)
}
