package vendor

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VendorStatus string
const (
	ActivatedVendorStatus VendorStatus = "Activated"
	DeactivatedVendorStatus VendorStatus = "Dectivated"
)
type Vendor struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	BusinessName string `json:"business_name" bson:"business_name"`
	Slug string `json:"slug" bson:"slug"`
	Status VendorStatus  `json:"status" bson:"status"`
	RatingAverage float64 `json:"rating_average" bson:"rating_average"`
	RatingCount int32 `json:"rating_count" bson:"rating_count"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewVendor (userID primitive.ObjectID,businessname,slug string, status VendorStatus,ratingAverage float64, ratingCount float32 ) *Vendor{
	now:=time.Now()
	return &Vendor{
		ID: primitive.NewObjectID(),
		UserID: userID,
		BusinessName: businessname,
		Slug: slug,
		Status: status,
		RatingAverage: ratingAverage,
		RatingCount: int32(ratingCount),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (v *Vendor) UpdateTimestamp() {
	v.UpdatedAt = time.Now()
}

func (v *Vendor) ApplyUpdate(req UpdateVendorProfileRequest) {
	if req.BusinessName != nil {
		v.BusinessName = *req.BusinessName
	}
	if req.Slug != nil {
		v.Slug = *req.Slug
	}
	v.UpdateTimestamp()
}

func (v *Vendor) ToResponse() *VendorProfileResponse {
	return &VendorProfileResponse{
		ID: v.ID.Hex(),
		UserID: v.UserID.Hex(),
		BusinessName: v.BusinessName,
		Slug: v.Slug,
		Status: string(v.Status),
		RatingAverage: v.RatingAverage,
		RatingCount: v.RatingCount,
		CreatedAt: v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}