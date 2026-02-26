package server

import (
	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/internal/auth"
	"github.com/techrook/23-market/internal/user"
	"github.com/techrook/23-market/internal/vendor"
)

func SetupRoutes(
	r* gin.Engine,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	vendorHandler *vendor.Handler,
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
		protected.POST("/:userID", userHandler.CreateUserProfile)
		protected.PUT("/:userID", userHandler.UpdateUserProfile)
		protected.GET("/:userID", userHandler.GetUserProfile)
		protected.DELETE("/:userID", userHandler.DeleteUserProfile)
	}

		vendorGroup := r.Group("/vendors")
	vendorGroup.Use(auth.AuthMiddleware(authCfg))
	{
		vendorGroup.POST("/complete-profile", vendorHandler.CompleteVendorProfile)
		vendorGroup.GET("/profile", vendorHandler.GetVendorProfile)
		vendorGroup.PUT("/profile", vendorHandler.UpdateVendorProfile)
		vendorGroup.DELETE("/profile", vendorHandler.DeactivateVendorProfile)
	}



}