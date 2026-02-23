package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfile struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"userid" bson:"userid"`
	Fullname string 	`json:"fullname" bson:"fullname"`
	Phone string `json:"phone" bson:"phone"`
	Street string `json:"street" bson:"street"`
	City string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
	IsDefault   bool `json:"is_default" bson:"is_default"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUserProfile (userID primitive.ObjectID,fullname,phone,street,city,country string, IsDefault bool) *UserProfile {
	now := time.Now()
	return &UserProfile{
		ID: primitive.NewObjectID(),
		UserId: userID,
		Fullname: fullname,
		Phone: phone,
		Street: street,
		City: city,
		Country: country,
		IsDefault: true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p *UserProfile) UpdateTimestamp() {
	p.UpdatedAt = time.Now()
}

type UpdateProfileRequest struct {
	Fullname  *string `json:"fullname,omitempty" binding:"omitempty"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty"`
	Street    *string `json:"street,omitempty" binding:"omitempty"`
	City      *string `json:"city,omitempty" binding:"omitempty"`
	Country   *string `json:"country,omitempty" binding:"omitempty"`
	IsDefault *bool   `json:"is_default,omitempty" binding:"omitempty"`
}


func (p *UserProfile) Apply(req UpdateProfileRequest) {
	if req.Fullname != nil {
		p.Fullname = *req.Fullname
	}
	if req.Phone != nil {
		p.Phone = *req.Phone
	}
	if req.Street != nil {
		p.Street = *req.Street
	}
	if req.City != nil {
		p.City = *req.City
	}
	if req.Country != nil {
		p.Country = *req.Country
	}
	if req.IsDefault != nil {
		p.IsDefault = *req.IsDefault
	}
	p.UpdateTimestamp()
}

func (p *UserProfile) ToResponse() UserProfileResponse {
    return UserProfileResponse{
        ID: p.ID.Hex(),
        UserID: p.UserId.Hex(),
        FullName: p.Fullname,
        Phone: p.Phone,
        Street: p.Street,
        City: p.City,
        Country: p.Country,
        IsDefault: p.IsDefault,
        CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}
