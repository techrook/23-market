package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/config"
	"github.com/techrook/23-market/database"
)

func main() {

	cfg := config.Load()

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err := database.EnsureIndexes(database.DB); err != nil {
		log.Fatalf("Failed to ensure MongoDB indexes: %v", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("⚠️ Error closing DB connection: %v", err)
		}
	}()


	r := gin.Default()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Server starting on http://localhost%s [%s]", addr, cfg.Environment)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}