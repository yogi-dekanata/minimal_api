package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logger Middleware untuk logging request dan response
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[Middleware] Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[Middleware] Response: %d", c.Writer.Status())
	}
}

// CORS Middleware untuk menangani request dari frontend (React, Vue, dll.)
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika metode OPTIONS, langsung return 200 OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
