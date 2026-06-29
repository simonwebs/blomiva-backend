package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Printf(
			"http status=%d method=%s path=%s ip=%s latency=%s reqId=%v",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			time.Since(start),
			c.GetString("requestId"),
		)
	}
}
