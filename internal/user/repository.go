package user

import (
	"context"
	"errors"
	"fmt"
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

	CreateProfile(ctx context.Context, p *UserProfile) error
	GetProfileByUserID(ctx context.Context, userID primitive.ObjectID) (*UserProfile, error)
	UpdateProfile(ctx context.Context, p *UserProfile) error
	DeleteProfile(ctx context.Context, userID primitive.ObjectID) error
	ProfileExists(ctx context.Context, userID primitive.ObjectID) (bool, error)
	RegisterProfile(ctx context.Context, userID primitive.ObjectID) error
}


type UserRepository struct {
	collection *mongo.Collection
	profileCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) Repository {
	return &UserRepository{
		collection: db.Collection("users"),
		profileCollection: db.Collection("user_profiles"),
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

func (r *UserRepository) CreateProfile (ctx context.Context, p *UserProfile)error{
	exists, err := r.ProfileExists(ctx, p.UserID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user profile already exists")
	}

	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()

	_, err = r.profileCollection.InsertOne(ctx, p)
	return err
}

func (r *UserRepository) RegisterProfile (ctx context.Context, userID primitive.ObjectID ) error {
	_,err := r.profileCollection.InsertOne(ctx, userID)
	return err
}
func (r *UserRepository) GetProfileByUserID(ctx context.Context, userID primitive.ObjectID) (*UserProfile, error) {
	var p UserProfile
	fmt.Println("Getting profile for userID:", userID) // Debug log
	err := r.profileCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user profile not found")
	}
	fmt.Printf("Found profile: %+v\n", p) // Debug log
	return &p, err

}

func (r *UserRepository) UpdateProfile(ctx context.Context, p *UserProfile) error {


	p.UpdateTimestamp()

	_, err := r.profileCollection.ReplaceOne(
		ctx,
		bson.M{"user_id": p.UserID},
		p,
		options.Replace().SetUpsert(false),
	)
	return err
}

func (r *UserRepository) DeleteProfile(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.profileCollection.DeleteOne(ctx, bson.M{"user_id": userID})
	return err
}

func (r *UserRepository) ProfileExists(ctx context.Context, userID primitive.ObjectID) (bool, error) {
	count, err := r.profileCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}