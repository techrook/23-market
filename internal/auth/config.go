package auth

import (
	"os"
	"time"
)

type Config struct {
	JWTSecret          string
	JWTExpiry          time.Duration
	RefreshTokenExpiry time.Duration
	RefreshTokenPrefix string 
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret-change-in-prod"),
		JWTExpiry:          15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		RefreshTokenPrefix: "rt_",
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}