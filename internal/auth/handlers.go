// auth/handlers.go

package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/techrook/23-market/pkg/response"
)

type Handler struct {
	service Service
	cfg     *Config 
}


func NewHandler(service Service, cfg *Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}


func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format", gin.H{"errors": err.Error()}, response.IsProduction(c))
		return
	}

	tokens, err := h.service.Signup(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserAlreadyExists):
			response.Conflict(c, "Email already registered", nil, response.IsProduction(c))
		case errors.Is(err, ErrTokenGeneration):
			response.InternalError(c, "Failed to create account", err, response.IsProduction(c))
		default:
			response.InternalError(c, "Signup failed", err, response.IsProduction(c))
		}
		return
	}


	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(h.cfg.RefreshTokenExpiry.Seconds()),
		"/",
		"",
		true,  
		true,  
	)

	response.Created(c, AuthResponse{
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
	}, "Account created successfully")
}


func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid credentials format", nil, response.IsProduction(c))
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {

			response.Unauthorized(c, "Invalid email or password", response.IsProduction(c))
			return
		}
		response.InternalError(c, "Login failed", err, response.IsProduction(c))
		return
	}


	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(h.cfg.RefreshTokenExpiry.Seconds()),
		"/",
		"",
		true,
		true,
	)

	response.OK(c, AuthResponse{
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
	}, "Login successful")
}


func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.Unauthorized(c, "Missing refresh token", response.IsProduction(c))
		return
	}

	tokens, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			response.Unauthorized(c, "Session expired, please login again", response.IsProduction(c))
			return
		}
		response.InternalError(c, "Token refresh failed", err, response.IsProduction(c))
		return
	}

	// Optional: Rotate refresh token (issue new one for security)
	// For now, we keep the same cookie. To rotate, uncomment below:
	// c.SetCookie(
	// 	"refresh_token",
	// 	tokens.RefreshToken,
	// 	int(h.cfg.RefreshTokenExpiry.Seconds()),
	// 	"/",
	// 	"",
	// 	true,
	// 	true,
	// )

	response.OK(c, AuthResponse{
		AccessToken: tokens.AccessToken,
		ExpiresIn:   tokens.ExpiresIn,
	}, "Token refreshed successfully")
}


func (h *Handler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")


	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required", response.IsProduction(c))
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		response.InternalError(c, "Invalid user context", nil, response.IsProduction(c))
		return
	}

	if err := h.service.Logout(c.Request.Context(), refreshToken, userID); err != nil {
		response.InternalError(c, "Logout failed", err, response.IsProduction(c))
		return
	}


	c.SetCookie(
		"refresh_token",
		"",
		-1, 
		"/",
		"",
		true,
		true,
	)

	response.OK(c, nil, "Logged out successfully")
}

// i would take this to the user service
func (h *Handler) Me(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required", response.IsProduction(c))
		return
	}

	userID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		response.InternalError(c, "Invalid user context", nil, response.IsProduction(c))
		return
	}

	me, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.NotFound(c, "User", response.IsProduction(c))
			return
		}
		response.InternalError(c, "Failed to fetch profile", err, response.IsProduction(c))
		return
	}

	response.OK(c, me, "Profile retrieved successfully")
}