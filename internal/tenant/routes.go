package tenant

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	tenants := api.Group("/tenants")

	// PUBLIC
	tenants.POST("/", h.CreateTenant)
	tenants.POST("/verify-school", h.VerifySchool)

	// PROTECTED
	tenants.GET("/", h.ListTenants)
	tenants.GET("/:id", h.GetTenant)
	tenants.GET("/slug/:slug", h.GetTenantBySlug)
	tenants.PUT("/:id", h.UpdateTenant)
	tenants.DELETE("/:id", h.DeleteTenant)
}
