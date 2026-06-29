package students

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *StudentService
}

func NewHandler(service *StudentService) *Handler {
	return &Handler{
		service: service,
	}
}

// POST /students
func (h *Handler) CreateStudent(c *gin.Context) {
	var req CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	student, err := h.service.CreateStudent(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	successResponse(c, http.StatusCreated, "student created", student)
}

// GET /students
func (h *Handler) ListStudents(c *gin.Context) {
	// Temporary until repository/service list implemented
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "students fetched",
		"data":    []Student{},
	})
}

// GET /students/:id
func (h *Handler) GetStudent(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "student fetched",
		"id":      id,
	})
}

// PUT /students/:id
func (h *Handler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "student updated",
		"id":      id,
		"payload": payload,
	})
}

// DELETE /students/:id
func (h *Handler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "student deleted",
		"id":      id,
	})
}

func successResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}
