package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := fmt.Sprintf("%d", time.Now().UnixNano())

		c.Set("requestId", id)
		c.Writer.Header().Set("X-Request-ID", id)

		c.Next()
	}
}
