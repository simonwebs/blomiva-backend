package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PasswordlessHandler struct {
	service *PasswordlessService
}

func NewPasswordlessHandler(service *PasswordlessService) *PasswordlessHandler {
	return &PasswordlessHandler{
		service: service,
	}
}

type RequestPasswordlessLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyPasswordlessLoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

func (h *PasswordlessHandler) RequestLogin(c *gin.Context) {
	var req RequestPasswordlessLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Valid email is required.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := h.service.RequestCode(ctx, req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "passwordless_request_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login code sent. Check your email.",
	})
}

func (h *PasswordlessHandler) VerifyLogin(c *gin.Context) {
	var req VerifyPasswordlessLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Email and code are required.",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	res, err := h.service.VerifyCode(ctx, req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "passwordless_verify_failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
