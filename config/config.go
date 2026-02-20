package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	GinMode         string
	Environment     string
	MongoURI        string
	MongoDatabase   string
	MongoMaxPoolSize uint64
	MongoMinPoolSize uint64
	MongoTimeout    time.Duration
}

func Load() *Config {

	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		Environment: getEnv("ENVIRONMENT", "development"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGO_DB", "23market"),
		MongoMaxPoolSize: uint64(getEnvInt("MONGO_MAX_POOL_SIZE", 100)),
		MongoMinPoolSize: uint64(getEnvInt("MONGO_MIN_POOL_SIZE", 10)),
		MongoTimeout:    time.Duration(getEnvInt("MONGO_TIMEOUT_SECONDS", 10)) * time.Second,

	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("⚠️ Invalid integer for %s: %s, using default %d", key, value, defaultValue)
	}
	return defaultValue
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.GinMode == "release"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.GinMode == "debug"
}

func (c *Config) MongoConnectionString() string {
	uri := strings.TrimSpace(c.MongoURI)
	
	if uri != "" && uri != "mongodb://localhost:27017" {
		return uri
	}
	

	user := os.Getenv("MONGO_USER")
	pass := os.Getenv("MONGO_PASS")
	host := getEnv("MONGO_HOST", "localhost")
	port := getEnv("MONGO_PORT", "27017")
	db := c.MongoDatabase

	if user != "" && pass != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
			user, pass, host, port, db)
	}
	return fmt.Sprintf("mongodb://%s:%s/%s", host, port, db)
}