package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// DatabaseCredentials represents database secret structure
type DatabaseCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

var (
	cachedCreds *DatabaseCredentials
	cacheMutex  sync.RWMutex
)

// GetDatabaseCredentials retrieves and caches database credentials from AWS Secrets Manager
func GetDatabaseCredentials(ctx context.Context) (*DatabaseCredentials, error) {
	// Check cache first
	cacheMutex.RLock()
	if cachedCreds != nil {
		cacheMutex.RUnlock()
		return cachedCreds, nil
	}
	cacheMutex.RUnlock()

	// Get secret ARN from environment
	secretArn := os.Getenv("DB_SECRET_ARN")
	if secretArn == "" {
		return nil, fmt.Errorf("DB_SECRET_ARN environment variable not set")
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	// Retrieve secret value
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretArn,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret: %w", err)
	}

	// Parse secret JSON
	var creds DatabaseCredentials
	if err := json.Unmarshal([]byte(*result.SecretString), &creds); err != nil {
		return nil, fmt.Errorf("failed to parse secret: %w", err)
	}

	// Cache credentials
	cacheMutex.Lock()
	cachedCreds = &creds
	cacheMutex.Unlock()

	return &creds, nil
}

// BuildDSN builds MySQL database connection string from credentials
func BuildDSN(creds *DatabaseCredentials) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		creds.Username,
		creds.Password,
		creds.Host,
		creds.Port,
		creds.Database,
	)
}

// ClearCache clears the cached credentials (useful for testing)
func ClearCache() {
	cacheMutex.Lock()
	cachedCreds = nil
	cacheMutex.Unlock()
}
