package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf(
			"[HTTP] method=%s path=%s status=%d duration=%s",
			c.Request.Method,
			c.Request.URL.RequestURI(),
			statusCode,
			duration,
		)
	}
}
