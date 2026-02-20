

package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	Update(ctx context.Context, u *User) error
	Verify(ctx context.Context, id primitive.ObjectID) error
	Exists(ctx context.Context, email string) (bool, error) 
}


type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) Repository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *User) error {

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var u User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	u.UpdateTimestamp()
	
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": u.ID},
		u,
		options.Replace().SetUpsert(false),
	)
	return err
}

func (r *UserRepository) Verify(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"is_verified": true,
			"updated_at":  time.Now(),
		}},
	)
	return err
}


func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}