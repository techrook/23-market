package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/techrook/23-market/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database


func Connect(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoTimeout)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(cfg.MongoConnectionString()).
		SetMaxPoolSize(cfg.MongoMaxPoolSize).
		SetMinPoolSize(cfg.MongoMinPoolSize).
		SetConnectTimeout(cfg.MongoTimeout).
		SetServerSelectionTimeout(5 * time.Second)

	var err error
	Client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}


	if err := Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	DB = Client.Database(cfg.MongoDatabase)
	log.Printf("✅ Connected to MongoDB: %s", cfg.MongoDatabase)
	return nil
}


func Close() error {
	if Client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return Client.Disconnect(ctx)
}


func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}

func EnsureIndexes(db *mongo.Database) error {
	ctx := context.Background()
	users := db.Collection("users")

	// Unique index on email
	_, err := users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    primitive.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Optional: index on role for faster queries
	_, err = users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: primitive.M{"role": 1},
	})
	return err
}