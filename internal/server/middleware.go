package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := fmt.Sprintf("%d", time.Now().UnixNano())
		c.Set("requestId", id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

func RequestLogger(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.Printf(
			"status=%d method=%s path=%s ip=%s latency=%s reqId=%v",
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			time.Since(start),
			c.GetString("requestId"),
		)
	}
}

func CORSMiddleware(cfg interface{ Env string }) gin.HandlerFunc {
	origins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
	}

	if cfg.Env == "production" {
		origins = []string{
			"https://blomiva.com",
			"https://app.blomiva.com",
		}
	}

	if raw := os.Getenv("ALLOWED_ORIGINS"); raw != "" {
		parts := strings.Split(raw, ",")
		var out []string
		for _, p := range parts {
			if v := strings.TrimSpace(p); v != "" {
				out = append(out, v)
			}
		}
		if len(out) > 0 {
			origins = out
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}