package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from various sources
func LoadConfig() (*Config, error) {
	return LoadConfigWithPath("")
}

// LoadConfigWithPath loads configuration from a specific path
func LoadConfigWithPath(configPath string) (*Config, error) {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")

	// Set default config paths
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("/etc/sarc")
	viper.AddConfigPath("$HOME/.sarc")

	// If a specific config path is provided, use it
	if configPath != "" {
		if filepath.IsAbs(configPath) {
			viper.SetConfigFile(configPath)
		} else {
			viper.AddConfigPath(filepath.Dir(configPath))
			viper.SetConfigName(strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath)))
		}
	}

	// Set default values
	setDefaults()

	// Override with standard database environment variables
	mapDatabaseEnvVars()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Load environment-specific config if available
	env := getEnvironment()
	if env != "" {
		log.Printf("Loading environment-specific config for: %s", env)
		loadEnvironmentConfig(env)
	}

	// Unmarshal config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate config
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	log.Printf("Configuration loaded successfully")
	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.shutdown_timeout", "30s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "sarcng")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.conn_max_idle_time", "2m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiry", "24h")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")

	// API defaults
	viper.SetDefault("api.base_url", "http://localhost:8080")
	viper.SetDefault("api.timeout", "30s")
}

// mapDatabaseEnvVars maps standard database environment variables to viper keys
func mapDatabaseEnvVars() {
	// Map standard DB_* environment variables
	dbEnvMap := map[string]string{
		"DB_HOST":     "database.host",
		"DB_PORT":     "database.port",
		"DB_USER":     "database.user",
		"DB_PASSWORD": "database.password",
		"DB_NAME":     "database.name",
		"DB_CHARSET":  "database.charset",
	}

	// Apply environment variable mappings
	for envVar, configKey := range dbEnvMap {
		if value := os.Getenv(envVar); value != "" {
			viper.Set(configKey, value)
		}
	}

	// Also handle PORT for server (common in Docker/Heroku)
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
}

// loadEnvironmentConfig loads environment-specific configuration
func loadEnvironmentConfig(env string) {
	envConfigName := env

	// Try to read environment-specific config
	viper.SetConfigName(envConfigName)
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Error reading environment config %s: %v", envConfigName, err)
		}
	} else {
		log.Printf("Merged environment config: %s", envConfigName)
	}
}

// getEnvironment determines the current environment
func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "development"
	}
	return strings.ToLower(env)
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port <= 0 || config.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", config.Database.Port)
	}

	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}

	if config.JWT.Secret == "" || config.JWT.Secret == "your-secret-key" {
		log.Printf("Warning: Using default JWT secret, please set a secure secret in production")
	}

	return nil
}

// ReloadConfig reloads the configuration
func ReloadConfig() (*Config, error) {
	log.Printf("Reloading configuration...")
	return LoadConfig()
}

// GetConfigPath returns the path to the configuration file being used
func GetConfigPath() string {
	return viper.ConfigFileUsed()
}
