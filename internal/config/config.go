// Package config loads and provides application configuration.
// Configuration is loaded from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration values.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port         int
	ReadTimeout  int // seconds
	WriteTimeout int // seconds
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	FilePath string
}

// LoggerConfig holds logging settings.
type LoggerConfig struct {
	Level    string // DEBUG, INFO, WARN, ERROR
	FilePath string
}

// Load reads configuration from environment variables.
// If an environment variable is not set, a default value is used.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 15),
		},
		Database: DatabaseConfig{
			FilePath: getEnvStr("DB_FILE_PATH", "./data/app.db"),
		},
		Logger: LoggerConfig{
			Level:    getEnvStr("LOG_LEVEL", "DEBUG"),
			FilePath: getEnvStr("LOG_FILE_PATH", "./logs/app.log"),
		},
	}
}

// getEnvStr returns the value of an environment variable or a default.
func getEnvStr(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvInt returns the integer value of an environment variable or a default.
func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
