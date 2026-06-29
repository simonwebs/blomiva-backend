package profiles

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler, auth gin.HandlerFunc) {
	api := router.Group("/api/v1")

	public := api.Group("/profiles")
	{
		public.GET("/public/:ownerId", handler.PublicByOwnerID)
		public.GET("/slug/:slug", handler.PublicBySlug)
	}

	protected := api.Group("/profiles")
	protected.Use(auth)
	{
		protected.POST("/ensure", handler.EnsureProfile)
		protected.GET("/me", handler.GetMe)
		protected.PUT("/me", handler.UpdateMe)
		protected.POST("/touch", handler.Touch)
		protected.POST("/avatar", handler.UploadAvatar)
		protected.POST("/banner", handler.UploadBanner)
		protected.DELETE("/custom/:key", handler.UnsetCustomKey)
		protected.POST("/schedule-delete", handler.ScheduleDelete)
	}

	admin := api.Group("/admin/profiles")
	admin.Use(auth)
	{
		admin.GET("", handler.AdminList)
		admin.POST("/status", handler.AdminSetUserStatus)
		admin.DELETE("/:ownerId", handler.AdminDeleteUser)
	}
}