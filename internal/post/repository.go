package post

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guisteink/tusk/internal"
)

type Repository struct {
	Conn *mongo.Client
}

func (r *Repository) Insert(post internal.Post) (string, error) {
	collection := r.Conn.Database("tusk").Collection("posts")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		return "", fmt.Errorf("failed to insert post: %v", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to get inserted post ID")
	}

	return insertedID.Hex(), nil
}
