package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/config"
)

func main() {

	cfg := config.Load()

	r := gin.Default()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Server starting on http://localhost%s [%s]", addr, cfg.Environment)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}