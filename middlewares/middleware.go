package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
)

// CustomLogger is a middleware that logs request details similar to the Fiber logger
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Before request
		c.Next()

		// After request
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		requestID := c.GetString("requestid")

		log.Info().
			Str("time", time.Now().Format(time.RFC3339)).
			Str("requestid", requestID).
			Int("status", status).
			Str("method", method).
			Str("path", path).
			Dur("latency", duration).
			Msg("request")
	}
}

func AddGinMiddlewares(r *gin.Engine) {
	// CORS
	r.Use(cors.Default())

	// Recover middleware (built-in)
	r.Use(gin.Recovery())

	// Request ID middleware (uses "X-Request-ID" or generates one)
	r.Use(requestid.New())

	// Compression: Gin doesn’t have built-in, but use a third-party like "github.com/gin-contrib/gzip"
	// Uncomment if you add this dependency:
	// import "github.com/gin-contrib/gzip"
	// r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Logger
	r.Use(CustomLogger())

	// Optional: set binding timeouts
	binding.EnableDecoderUseNumber = true
}
