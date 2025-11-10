package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	LogLevel string // Log level: "debug", "info", "warn", "error"
	Mode     string // Gin mode: "debug", "release", "test"
}

func New() (*Config, error) {
	config := &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Mode:     getEnv("GIN_MODE", "debug"),
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	validModes := map[string]bool{"debug": true, "release": true, "test": true}
	if !validModes[c.Mode] {
		return fmt.Errorf("invalid mode: %s (must be debug, release, or test)", c.Mode)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
