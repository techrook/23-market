package auth

import "github.com/techrook/23-market/internal/user"


type SignupRequest struct {
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=8"`
	Role     user.Role `json:"role" binding:"required,oneof=vendor user"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


type RefreshRequest struct{}


type LogoutRequest struct{}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` 
}


type MeResponse struct {
	ID    string    `json:"id"`
	Email string    `json:"email"`
	Role  user.Role `json:"role"`
}


type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"` 
}