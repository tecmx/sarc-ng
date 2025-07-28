//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL           = "http://localhost:8080"
	healthEndpoint    = "/health"
	apiHealthEndpoint = "/api/v1/health"
	buildingsEndpoint = "/api/v1/buildings"
)

func TestMain(m *testing.M) {
	// Wait for service to be ready
	if !waitForService(baseURL+healthEndpoint, 30*time.Second) {
		fmt.Println("Service not ready, skipping integration tests")
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func waitForService(url string, timeout time.Duration) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + healthEndpoint)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var health map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&health)
	require.NoError(t, err)

	assert.Equal(t, "healthy", health["status"])
}

func TestAPIHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + apiHealthEndpoint)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBuildingsEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + buildingsEndpoint)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Should return 200 (empty list is fine for integration test)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify it returns valid JSON
	var buildings []interface{}
	err = json.NewDecoder(resp.Body).Decode(&buildings)
	require.NoError(t, err)
}

func TestAPIEndpoints(t *testing.T) {
	endpoints := []struct {
		name string
		path string
	}{
		{"Buildings", "/api/v1/buildings"},
		{"Classes", "/api/v1/classes"},
		{"Lessons", "/api/v1/lessons"},
		{"Reservations", "/api/v1/reservations"},
		{"Resources", "/api/v1/resources"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			resp, err := http.Get(baseURL + endpoint.path)
			require.NoError(t, err)
			defer resp.Body.Close()

			// All endpoints should return 200 for GET requests
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			// Verify response is valid JSON
			var result interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
		})
	}
}
