package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/config"
	"github.com/techrook/23-market/database"
	"github.com/techrook/23-market/internal/auth"
	"github.com/techrook/23-market/internal/server"
	"github.com/techrook/23-market/internal/user"
	"github.com/techrook/23-market/internal/vendor"
)

func main() {

	cfg := config.Load()
	authCfg := auth.LoadConfig()

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

	userRepo := user.NewUserRepository(database.DB)
	authRepo:= auth.NewAuthRepository(database.DB)
	vendorRepo := vendor.NewVendorRepository(database.DB)

	authService := auth.NewService(authCfg, userRepo, authRepo, vendorRepo)

	userHandler := user.NewHandler(user.NewService(userRepo))
	authHandler := auth.NewHandler(authService, authCfg)

	vendorService := vendor.NewService(vendorRepo)
	vendorHandler := vendor.NewHandler(vendorService)

	r := gin.Default()

	server.SetupRoutes(r,authHandler,userHandler,vendorHandler, userRepo)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Server starting on http://localhost%s [%s]", addr, cfg.Environment)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}