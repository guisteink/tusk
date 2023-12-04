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

func (r *Repository) Find(filter interface{}) ([]internal.Post, error) {
	ctx := context.Background()
	collection := r.Conn.Database("tusk").Collection("posts")

	var posts []internal.Post
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %v", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var post internal.Post
		if err := cur.Decode(&post); err != nil {
			return nil, fmt.Errorf("failed to decode post: %v", err)
		}
		posts = append(posts, post)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return posts, nil
}
