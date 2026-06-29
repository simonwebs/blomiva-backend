package profiles

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Handler struct {
	service *Service
	repo    *Repository
}

func NewHandler(service *Service, repo *Repository) *Handler {
	return &Handler{
		service: service,
		repo:    repo,
	}
}

func userID(c *gin.Context) string {
	id, _ := c.Get("userId")
	value, _ := id.(string)
	return value
}

func isAdmin(c *gin.Context) bool {
	value, _ := c.Get("isAdmin")
	ok, _ := value.(bool)
	return ok
}

func requireAuth(c *gin.Context) (string, bool) {
	id := userID(c)
	if id == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "not authenticated",
		})
		return "", false
	}
	return id, true
}

func requireAdmin(c *gin.Context) (string, bool) {
	id, ok := requireAuth(c)
	if !ok {
		return "", false
	}

	if !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "admin access required",
		})
		return "", false
	}

	return id, true
}

func (h *Handler) EnsureProfile(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	profile, err := h.service.EnsureProfile(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"profile": profile,
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	profile, err := h.service.GetMe(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateProfile(c.Request.Context(), ownerID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"profile": profile,
	})
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	var req UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateAvatar(c.Request.Context(), ownerID, req.Base64Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"profile":   profile,
		"avatarUrl": profile.AvatarURL,
	})
}

func (h *Handler) UploadBanner(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	var req UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateBanner(c.Request.Context(), ownerID, req.Base64Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":        true,
		"profile":   profile,
		"bannerUrl": profile.BannerURL,
	})
}

func (h *Handler) Touch(c *gin.Context) {
	ownerID := userID(c)

	err := h.service.Touch(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) UnsetCustomKey(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	key := c.Param("key")

	err := h.service.UnsetCustomKey(c.Request.Context(), ownerID, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) ScheduleDelete(c *gin.Context) {
	ownerID, ok := requireAuth(c)
	if !ok {
		return
	}

	err := h.service.ScheduleDelete(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) PublicByOwnerID(c *gin.Context) {
	ownerID := c.Param("ownerId")

	profile, err := h.repo.FindProfileByOwnerID(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if profile == nil || profile.Status != "active" || profile.IsBlocked {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ownerId":      profile.OwnerID,
		"username":     profile.Username,
		"slug":         profile.Slug,
		"displayName":  profile.DisplayName,
		"bio":          profile.Bio,
		"profileImage": profile.ProfileImage,
		"avatar":       profile.Avatar,
		"avatarUrl":    profile.AvatarURL,
		"banner":       profile.Banner,
		"bannerUrl":    profile.BannerURL,
	})
}

func (h *Handler) PublicBySlug(c *gin.Context) {
	slug := c.Param("slug")

	profile, err := h.repo.FindProfileBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ownerId":      profile.OwnerID,
		"username":     profile.Username,
		"slug":         profile.Slug,
		"displayName":  profile.DisplayName,
		"bio":          profile.Bio,
		"profileImage": profile.ProfileImage,
		"avatar":       profile.Avatar,
		"avatarUrl":    profile.AvatarURL,
		"banner":       profile.Banner,
		"bannerUrl":    profile.BannerURL,
	})
}

func (h *Handler) AdminList(c *gin.Context) {
	_, ok := requireAdmin(c)
	if !ok {
		return
	}

	limit := 200
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	cursor, err := h.repo.Profiles.Find(
		c.Request.Context(),
		bson.M{},
		nil,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(c.Request.Context())

	var rows []Profile
	for cursor.Next(c.Request.Context()) {
		if len(rows) >= limit {
			break
		}

		var p Profile
		if err := cursor.Decode(&p); err == nil {
			rows = append(rows, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rows,
		"total": len(rows),
	})
}

func (h *Handler) AdminSetUserStatus(c *gin.Context) {
	adminID, ok := requireAdmin(c)
	if !ok {
		return
	}

	var req SetUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SetUserStatus(c.Request.Context(), adminID, req.OwnerID, req.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	adminID, ok := requireAdmin(c)
	if !ok {
		return
	}

	ownerID := c.Param("ownerId")

	err := h.service.DeleteUser(c.Request.Context(), adminID, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}