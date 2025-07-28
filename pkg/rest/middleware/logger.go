package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter is a custom ResponseWriter that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// Logger returns a request logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log in debug mode
		if gin.Mode() == gin.ReleaseMode {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Log request details
		log.Printf("\n[REQUEST] %s %s\n", method, path)
		log.Printf("Headers: %v\n", c.Request.Header)

		// Read and log request body for methods that might have a body
		if method == "POST" || method == "PUT" {
			var requestBody []byte
			if c.Request.Body != nil {
				requestBody, _ = io.ReadAll(c.Request.Body)
				// Restore the request body for the handlers
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
				log.Printf("Body: %s\n", string(requestBody))
			}
		}

		// Create a custom response writer to capture the response body
		responseBody := &bytes.Buffer{}
		writer := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Log response details
		status := c.Writer.Status()
		latency := time.Since(start)
		log.Printf("\n[RESPONSE] %s %s - Status: %d - Time: %s\n", method, path, status, latency)

		// Only log response body for JSON responses to avoid logging binary data
		contentType := c.Writer.Header().Get("Content-Type")
		if contentType == "application/json" {
			log.Printf("Response Body: %s\n", responseBody.String())
		} else {
			log.Printf("Response Content-Type: %s (body not logged)\n", contentType)
		}

		log.Println("-----------------------------------------------------")
	}
}
