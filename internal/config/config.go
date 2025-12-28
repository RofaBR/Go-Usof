package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	Port        string       `validate:"required"`
	LogLevel    string       `validate:"required,oneof=debug, info warn error"`
	Mode        string       `validate:"required,oneof=debug release test"`
	DatabaseURL string       `validate:"required"`
	Redis       RedisConfig  `validate:"required"`
	JWT         JWTConfig    `validate:"required"`
	Sender      SenderConfig `validate:"required"`
}

type RedisConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required,numeric"`
	Password string
	DB       int `validate:"required"`
}

type JWTConfig struct {
	AccessSecret  string `validate:"required"`
	RefreshSecret string `validate:"required"`
	AccessTTL     int    `validate:"required,gt=0"`
	RefreshTTL    int    `validate:"required,gt=0"`
	RefreshMaxTTL int    `validate:"required,gt=0,gtefield=RefreshTTL"`
}

type SenderConfig struct {
	FromEmail string `validate:"required,email"`
	Password  string `validate:"required"`
	SMTPHost  string `validate:"required"`
	SMTPPort  int    `validate:"required,min=1,max=65535"`
}

var validate = validator.New()

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
		Sender: SenderConfig{
			FromEmail: getEnv("SENDER_EMAIL", ""),
			Password:  getEnv("SENDER_PASSWORD", ""),
			SMTPHost:  getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:  getEnvAsInt("SMTP_PORT", 587),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if err := validate.Struct(c); err != nil {
		return formatValidationErrors(err)
	}
	return nil
}

func formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			return fmt.Errorf("%s: validation failed on '%s' tag", err.Field(), err.Tag())
		}
	}
	return err
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
