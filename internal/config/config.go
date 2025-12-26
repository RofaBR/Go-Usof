package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	LogLevel    string // Log level: "debug", "info", "warn", "error"
	Mode        string // Gin mode: "debug", "release", "test"
	DatabaseURL string
	Redis       RedisConfig
	JWT         JWTConfig
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     int
	RefreshTTL    int
	RefreshMaxTTL int
}

func New() (*Config, error) {
	config := &Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Mode:        getEnv("GIN_MODE", "debug"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessTTL:     getEnvAsInt("JWT_ACCESS_TTL", 15),
			RefreshTTL:    getEnvAsInt("JWT_REFRESH_TTL", 7),
			RefreshMaxTTL: getEnvAsInt("JWT_REFRESH_MAX_TTL", 30),
		},
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

	if c.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}

	if c.Redis.Host == "" {
		return fmt.Errorf("redis host is required")
	}
	if c.Redis.Port == "" {
		return fmt.Errorf("redis port is required")
	}
	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		return fmt.Errorf("redis DB must be between 0 and 15")
	}

	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}
	if c.JWT.AccessTTL <= 0 {
		return fmt.Errorf("JWT access TTL must be positive")
	}
	if c.JWT.RefreshTTL <= 0 {
		return fmt.Errorf("JWT refresh TTL must be positive")
	}
	if c.JWT.RefreshMaxTTL <= 0 {
		return fmt.Errorf("JWT refresh max TTL must be positive")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
