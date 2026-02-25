package vendor

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface{
	CreateVendorProfile(ctx context.Context,  userID primitive.ObjectID) error
	CompleteVendorRegistration(ctx context.Context, userID primitive.ObjectID, businessName, slug string) error
	GetVendorByUserID(ctx context.Context, userID primitive.ObjectID) (*Vendor, error)
	UpdateVendor(ctx context.Context, v *Vendor) error
	DeactivateVendor(ctx context.Context, id primitive.ObjectID) error
	VendorExist(ctx context.Context, userID primitive.ObjectID) (bool, error)
}

type VendorRepository struct {
	vendorCollection *mongo.Collection
}

func NewVendorRepository(db *mongo.Database) Repository {
	return &VendorRepository{
		vendorCollection: db.Collection("vendors"),
	}
}

func (r *VendorRepository) CreateVendorProfile(ctx context.Context, userID primitive.ObjectID) error {
	exists, err := r.VendorExist(ctx, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("vendor profile already exists")
	}
	vendor := NewVendor(userID, "", "", ActivatedVendorStatus, 0.0, 0)	
	_, err = r.vendorCollection.InsertOne(ctx, vendor)
	return err
}

func (r *VendorRepository) VendorExist(ctx context.Context, userID primitive.ObjectID) (bool, error) {
	count, err := r.vendorCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *VendorRepository) CompleteVendorRegistration(ctx context.Context, userID primitive.ObjectID, businessName, slug string) error {
	_,err := r.vendorCollection.UpdateMany(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{
			"business_name": businessName,
			"slug": slug,
			"status": ActivatedVendorStatus,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		}},
	)
	
	return err
}

func (r *VendorRepository) GetVendorByUserID(ctx context.Context, userID primitive.ObjectID) (*Vendor, error) {
	var v Vendor
	err := r.vendorCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&v)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("vendor not found")
	}
	return &v, err
}

func (r *VendorRepository) UpdateVendor(ctx context.Context, v *Vendor) error {
	v.UpdateTimestamp()
	_, err := r.vendorCollection.ReplaceOne(
		ctx,
		bson.M{"_id": v.ID},
		v,
	)
	return err	
}

func (r *VendorRepository) DeactivateVendor(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.vendorCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"status": DeactivatedVendorStatus,
			"updated_at": primitive.NewDateTimeFromTime(time.Now()),
		}},
	)
	return err
}

