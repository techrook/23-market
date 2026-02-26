package vendor

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	CompleteVendorProfile(ctx context.Context, userID primitive.ObjectID, req CompleteVendorRegistrationRequest) (*VendorProfileResponse, error)
	GetVendorProfile(ctx context.Context, userId primitive.ObjectID) (*VendorProfileResponse, error)
	UpdateVendorProfile(ctx context.Context, userId primitive.ObjectID, req UpdateVendorProfileRequest) (*VendorProfileResponse, error)
	DeactivateVendorProfile(ctx context.Context, userID primitive.ObjectID) error
}

type service struct{
	vendorRepo Repository
}

func NewService(vendorRepo Repository) Service {
	return &service{
		vendorRepo: vendorRepo,
	}
}

func (s *service) CompleteVendorProfile(ctx context.Context, userID primitive.ObjectID, req CompleteVendorRegistrationRequest) (*VendorProfileResponse, error) {
	err := s.vendorRepo.CompleteVendorRegistration(ctx, userID, req.BusinessName, req.Slug)
	if err != nil {
		return nil, err
	}
	
	vendor, err := s.vendorRepo.GetVendorByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	

	return vendor.ToResponse(), nil
}

func (s *service) GetVendorProfile(ctx context.Context, userId primitive.ObjectID) (*VendorProfileResponse, error) {
	vendor, err := s.vendorRepo.GetVendorByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return vendor.ToResponse(), nil
}

func (s *service) UpdateVendorProfile(ctx context.Context, userId primitive.ObjectID, req UpdateVendorProfileRequest) (*VendorProfileResponse, error) {
	err := s.vendorRepo.UpdateVendor(ctx, userId, req.BusinessName, req.Slug)
	if err != nil {
		return nil, err
	}
	
	vendor, err := s.vendorRepo.GetVendorByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return vendor.ToResponse(), nil
}

func (s *service) DeactivateVendorProfile(ctx context.Context, userID primitive.ObjectID) error {
	return s.vendorRepo.DeactivateVendor(ctx, userID)
}