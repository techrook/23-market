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

func NewVendor (userID primitive.ObjectID,businessname,slug, status string,ratingAverage float64, ratingCount float32 ) *Vendor{
	now:=time.Now()
	return &Vendor{
		ID: primitive.NewObjectID(),
		UserID: userID,
		BusinessName: businessname,
		Slug: slug,
		Status: status,
		ratingAverage: ratingAverage,
		ratingCount: ratingCount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (v *Vendor) UpdateTimestamp() {
	v.UpdatedAt = time.Now()
}