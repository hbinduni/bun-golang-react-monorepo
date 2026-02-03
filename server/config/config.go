package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	FrontendURL string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) GetAllowedOrigins() string {
	origins := []string{
		c.FrontendURL,
		"http://localhost:5173",
		"http://localhost:3000",
	}
	return strings.Join(origins, ",")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
