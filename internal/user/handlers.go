package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/techrook/23-market/pkg/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	userService Service
}

func NewHandler(userService Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) CreateUserProfile(c *gin.Context ) {
	var req CreateUserProfileRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.BadRequest(c, "Invalide request format", gin.H{"errors": err.Error()}, response.IsProduction(c))
		return
	}
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

	profile, err := h.userService.CreateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserProfileExists):
			response.Conflict(c, "Profile already exists", nil, response.IsProduction(c))
		case errors.Is(err, ErrUserNotFound):
			response.NotFound(c, "User", response.IsProduction(c))
		default:
			response.InternalError(c, "Failed to create profile", err, response.IsProduction(c))
		}
		return
	}
	response.Created(c, profile, "Profile created successfully")
}

func (h *Handler) UpdateUserProfile(c *gin.Context,)  {
	var req UpdateProfileRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.BadRequest(c, "Invalide request format", gin.H{"errors": err.Error()}, response.IsProduction(c))
		return
	}
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
	profile, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.NotFound(c, "User", response.IsProduction(c))
		default:
			response.InternalError(c, "Failed to update profile", err, response.IsProduction(c))
		}
		return
	}
	response.OK(c, profile, "Profile updated successfully")
}


func (h *Handler) GetUserProfile(c *gin.Context, )  {
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
	
	profile, err := h.userService.FindUserProfileByUserId(c.Request.Context(), userID)	
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.NotFound(c, "User", response.IsProduction(c))
		default:
			response.InternalError(c, "Failed to update profile", err, response.IsProduction(c))
		}
		return
	}
	response.OK(c, profile, "Profile retrieved successfully")
}


func (h *Handler) DeleteUserProfile(c *gin.Context, )  {
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

	err := h.userService.DeleteUserProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			response.NotFound(c, "User", response.IsProduction(c))
		default:
			response.InternalError(c, "Failed to delete profile", err, response.IsProduction(c))
		}
		return
	}
	response.OK(c, nil, "Profile deleted successfully")
}