

package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	SaveRefreshToken(ctx context.Context, userID primitive.ObjectID, tokenKey string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, tokenKey string) error
	ValidateRefreshToken(ctx context.Context, tokenKey string, userID primitive.ObjectID) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID primitive.ObjectID) error // Optional: logout all devices
}

type mongoRepository struct {
	collection *mongo.Collection 
}

func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

func (r *mongoRepository) SaveRefreshToken(ctx context.Context, userID primitive.ObjectID, tokenKey string, expiresAt time.Time) error {
	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":        tokenKey,
		"user_id":    userID,
		"expires_at": expiresAt,
		"created_at": time.Now(),
	})
	return err
}

func (r *mongoRepository) DeleteRefreshToken(ctx context.Context, tokenKey string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": tokenKey})
	return err
}

func (r *mongoRepository) ValidateRefreshToken(ctx context.Context, tokenKey string, userID primitive.ObjectID) error {
	var result struct {
		UserID primitive.ObjectID `bson:"user_id"`
	}
	err := r.collection.FindOne(ctx, bson.M{
		"_id":        tokenKey,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return errors.New("invalid or expired refresh token")
	}
	if result.UserID != userID {
		return errors.New("refresh token does not belong to user")
	}
	return err
}


func (r *mongoRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}