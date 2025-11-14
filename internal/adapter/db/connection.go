package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection with the given configuration
// Uses GORM's built-in connection handling and pool management
func Connect(config Config) (*gorm.DB, error) {
	// Build DSN (Data Source Name) with connection timeout
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	// Configure GORM with improved settings
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		// Let GORM handle connection issues during operations
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	log.Printf("Connecting to database at %s:%d", config.Host, config.Port)

	// Open database connection - let GORM handle connection establishment
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool for optimal performance and reliability
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings with shorter lifetimes to avoid stale connections
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection once
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database at %s:%d", config.Host, config.Port)
	return db, nil
}

// ConnectForLambda establishes a Lambda-optimized database connection
func ConnectForLambda() (*gorm.DB, error) {
	config := Config{
		Host:            getEnvWithDefault("DB_HOST", "localhost"),
		Port:            getEnvWithDefaultInt("DB_PORT", 3306),
		User:            getEnvWithDefault("DB_USER", "root"),
		Password:        getEnvWithDefault("DB_PASSWORD", ""),
		Database:        getEnvWithDefault("DB_NAME", "sarc"),
		MaxOpenConns:    1, // Lambda-optimized: single connection
		MaxIdleConns:    1, // Lambda-optimized: single connection
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	return Connect(config)
}

// Helper functions
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvWithDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
