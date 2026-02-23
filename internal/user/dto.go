package user

type CreateUserProfileRequest struct {
	FullName  string `json:"fullname" binding:"required,min=2,max=100"`
	Phone     string `json:"phone" binding:"required,min=10,max=20"`
	Street    string `json:"street" binding:"required,min=5,max=200"`
	City      string `json:"city" binding:"required,min=2,max=100"`
	Country   string `json:"country" binding:"required,min=2,max=100"`
	IsDefault bool   `json:"is_default"`
}

type UpdateUserProfileRequest struct {
	FullName  *string `json:"fullname,omitempty" binding:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,min=10,max=20"`
	Street    *string `json:"street,omitempty" binding:"omitempty,min=5,max=200"`
	City      *string `json:"city,omitempty" binding:"omitempty,min=2,max=100"`
	Country   *string `json:"country,omitempty" binding:"omitempty,min=2,max=100"`
	IsDefault *bool   `json:"is_default,omitempty"`
}

type UserProfileResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FullName  string `json:"fullname"`
	Phone     string `json:"phone"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Country   string `json:"country"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}