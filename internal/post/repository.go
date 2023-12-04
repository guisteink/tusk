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

func (r *Repository) Insert(post internal.Post, timeout int) (*internal.Post, error) {
	collection := r.Conn.Database("tusk").Collection("posts")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("failed to insert post: %v", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to get inserted post ID")
	}

	var insertedPost internal.Post
	err = collection.FindOne(ctx, primitive.M{"_id": insertedID}).Decode(&insertedPost)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve inserted post: %v", err)
	}

	return &insertedPost, nil
}
