package server

import (
	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/internal/auth"
	"github.com/techrook/23-market/internal/user"
)

func SetupRoutes(
	r* gin.Engine,
	authHandler *auth.Handler,
	userRepo user.Repository,
) {
	authCfg := auth.LoadConfig()
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandler.Signup)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
	}

		protected := r.Group("/users")
	protected.Use(auth.AuthMiddleware(authCfg))
	{
		protected.GET("/me", authHandler.Me)
	}
}