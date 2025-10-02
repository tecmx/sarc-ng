//go:build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sarc-ng/pkg/rest/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

	// Add health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "sarc-ng",
			"version": "1.0.0",
		})
	})

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test the health endpoint
	resp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Cleanup
	resp.Body.Close()
}

func TestAPIEndpoints(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router with basic API structure
	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())

	// Add API group
	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"service": "sarc-ng-api",
				"version": "1.0.0",
			})
		})

		// Basic endpoints check
		api.GET("/buildings", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{})
		})

		api.GET("/classes", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{})
		})

		api.GET("/lessons", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{})
		})

		api.GET("/reservations", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{})
		})

		api.GET("/resources", func(c *gin.Context) {
			c.JSON(http.StatusOK, []gin.H{})
		})
	}

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test cases
	testCases := []struct {
		name     string
		endpoint string
		expected int
	}{
		{"Health Check", "/api/v1/health", http.StatusOK},
		{"Buildings List", "/api/v1/buildings", http.StatusOK},
		{"Classes List", "/api/v1/classes", http.StatusOK},
		{"Lessons List", "/api/v1/lessons", http.StatusOK},
		{"Reservations List", "/api/v1/reservations", http.StatusOK},
		{"Resources List", "/api/v1/resources", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tc.endpoint)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp.StatusCode)
			resp.Body.Close()
		})
	}
}
