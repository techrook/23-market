package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleVendor Role = "vendor"
	RoleUser  Role = "user"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"` // Fixed typo: Passwordhash → PasswordHash
	Role         Role               `json:"role" bson:"role"`
	IsVerified   bool               `json:"is_verified" bson:"is_verified"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewUser(email, passwordHash string, role Role) *User {
	now := time.Now()
	return &User{
		ID:           primitive.NewObjectID(), 
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) UpdateTimestamp() {
	u.UpdatedAt = time.Now()
}


func (u *User) TableName() string {
	return "users"
}