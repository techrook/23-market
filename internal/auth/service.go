package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/techrook/23-market/internal/user"
	"github.com/techrook/23-market/internal/vendor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrTokenGeneration     = errors.New("failed to generate token")
)

type Service interface {
	Signup(ctx context.Context, req SignupRequest) (*TokenPair, error)
	Login(ctx context.Context, req LoginRequest) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string, userID primitive.ObjectID) error
	GetMe(ctx context.Context, userID primitive.ObjectID) (*MeResponse, error)
}

type service struct {
	cfg      *Config
	userRepo user.Repository
	authRepo Repository
	vendorRepo         vendor.Repository
}

func NewService(cfg *Config, userRepo user.Repository, authRepo Repository, vendorRepo vendor.Repository) Service {
	return &service{
		cfg:      cfg,
		userRepo: userRepo,
		authRepo: authRepo,
		vendorRepo: vendorRepo,
	}
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func (s *service) Signup(ctx context.Context, req SignupRequest) (*TokenPair, error) {

	exists, err := s.userRepo.Exists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}


	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}



	newUser := user.NewUser(req.Email, hash, req.Role)
	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	if newUser.Role == user.RoleUser {
		if err := s.userRepo.RegisterProfile(ctx, newUser.ID); err != nil {
						 log.Printf("⚠️ Profile creation failed for user %s: %v", newUser.ID.Hex(), err)
			return nil, fmt.Errorf("failed to initialize profile: %w", err)
		}
	}
	if newUser.Role == user.RoleVendor {
		if err := s.vendorRepo.CreateVendorProfile(ctx, newUser.ID); err != nil {
						 log.Printf("⚠️ Vendor profile creation failed for user %s: %v", newUser.ID.Hex(), err)
			return nil, fmt.Errorf("failed to initialize vendor profile: %w", err)
		}
	}

	return s.generateTokenPair(ctx, newUser)
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*TokenPair, error) {

	u, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}


	if err := CheckPassword(req.Password, u.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Todo: Check email verification
	// if !u.IsVerified {
	// 	return nil, ErrEmailNotVerified
	// }


	return s.generateTokenPair(ctx, u)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	tokenKey := s.cfg.RefreshTokenKey(refreshToken)


	var tokenDoc struct {
		UserID primitive.ObjectID `bson:"user_id"`
	}
	err := s.authRepo.(*mongoRepository).collection.FindOne(
		ctx,
		bson.M{"_id": tokenKey, "expires_at": bson.M{"$gt": time.Now()}},
	).Decode(&tokenDoc)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}


	u, err := s.userRepo.FindByID(ctx, tokenDoc.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	accessToken, err := GenerateAccessToken(s.cfg, u)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, 
		ExpiresIn:    int64(s.cfg.JWTExpiry.Seconds()),
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string, userID primitive.ObjectID) error {
	if refreshToken != "" {
		tokenKey := s.cfg.RefreshTokenKey(refreshToken)
		_ = s.authRepo.DeleteRefreshToken(ctx, tokenKey)
	}
	return nil
}

func (s *service) GetMe(ctx context.Context, userID primitive.ObjectID) (*MeResponse, error) {
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &MeResponse{
		ID:    u.ID.Hex(),
		Email: u.Email,
		Role:  u.Role,
	}, nil
}

func (s *service) generateTokenPair(ctx context.Context, u *user.User) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(s.cfg, u)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, ErrTokenGeneration
	}

	tokenKey := s.cfg.RefreshTokenKey(refreshToken)
	expiresAt := time.Now().Add(s.cfg.RefreshTokenExpiry)
	if err := s.authRepo.SaveRefreshToken(ctx, u.ID, tokenKey, expiresAt); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.cfg.JWTExpiry.Seconds()),
	}, nil
}