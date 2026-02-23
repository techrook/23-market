package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserProfileExists = errors.New("profile already exist")
	ErrUserNotFound      = errors.New("user not found")
)

type Service interface {
	CreateUserProfile(ctx context.Context, userID primitive.ObjectID, req CreateUserProfileRequest) (UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req UpdateProfileRequest) (UserProfileResponse, error)
	FindUserProfileByUserId(ctx context.Context, userID primitive.ObjectID) (UserProfileResponse, error)
	DeleteUserProfile(ctx context.Context, userID primitive.ObjectID) error
}

type service struct {
	userRepo Repository
}

func NewService(userRepo Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) CreateUserProfile(ctx context.Context, userID primitive.ObjectID, req CreateUserProfileRequest) (UserProfileResponse, error) {
	u, err := s.userRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return UserProfileResponse{}, ErrUserNotFound
	}

	exists, err := s.userRepo.ProfileExists(ctx, u.UserId)
	if err != nil {
		return UserProfileResponse{}, err
	}
	if exists {
		return UserProfileResponse{}, ErrUserProfileExists
	}

	profile := NewUserProfile(
		userID,
		req.FullName,
		req.Phone,
		req.Street,
		req.City,
		req.Country,
		req.IsDefault,
	)

	if err := s.userRepo.CreateProfile(ctx, profile); err != nil {
		return UserProfileResponse{}, err
	}

	return profile.ToResponse(), nil
}

func (s *service) UpdateUserProfile(ctx context.Context, userID primitive.ObjectID, req UpdateProfileRequest) (UserProfileResponse, error) {
	profile, err := s.userRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return UserProfileResponse{}, ErrUserNotFound
	}

	profile.Apply(req)

	if err := s.userRepo.UpdateProfile(ctx, profile); err != nil {
		return UserProfileResponse{}, err
	}

	return profile.ToResponse(), nil
}

func (s *service) FindUserProfileByUserId(ctx context.Context, userID primitive.ObjectID) (UserProfileResponse, error) {
	profile, err := s.userRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		return UserProfileResponse{}, ErrUserNotFound
	}
	return profile.ToResponse(), nil
}

func (s *service) DeleteUserProfile(ctx context.Context, userID primitive.ObjectID) error {
	return s.userRepo.DeleteProfile(ctx, userID)
}

