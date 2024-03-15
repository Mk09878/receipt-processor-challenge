package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(c *gin.Context) {
	// Log the HTTP method and URL of the incoming request
	log.Printf("Logging middleware: %s %s\n", c.Request.Method, c.Request.URL)

	// Call the next middleware or handler
	c.Next()
}
