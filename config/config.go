package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	Environment     string
}

func Load() *Config {

	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		Environment: getEnv("ENVIRONMENT", "development"),

	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf(" Invalid integer for %s: %s, using default %d", key, value, defaultValue)
	}
	return defaultValue
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.GinMode == "release"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.GinMode == "debug"
}