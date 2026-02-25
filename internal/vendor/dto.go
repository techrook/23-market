package vendor

type CompleteVendorRegistrationRequest struct {
	BusinessName string `json:"business_name" binding:"required,min=2,max=100"`
	Slug         string `json:"slug" binding:"required,min=2,max=100,alphanum"`
}

type UpdateVendorProfileRequest struct {
	BusinessName *string `json:"business_name,omitempty" binding:"omitempty,min=2,max=100"`
	Slug         *string `json:"slug,omitempty" binding:"omitempty,min=2,max=100,alphanum"`
}

type VendorProfileResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	BusinessName string  `json:"business_name"`
	Slug         string  `json:"slug"`
	Status       int     `json:"status"`
	Rating       float64 `json:"rating"`
	TotalReviews int     `json:"total_reviews"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}	
