package post

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guisteink/tusk/internal"
)

type Repository struct {
	Conn *mongo.Client
}

func (r *Repository) Insert(post internal.Post, ctx *gin.Context) (*internal.Post, error) {
	collection := r.Conn.Database("tusk").Collection("posts")

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

func (r *Repository) Find(filter interface{}, ctx *gin.Context) ([]internal.Post, error) {
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

func (r *Repository) Delete(id primitive.ObjectID, ctx *gin.Context) (internal.Post, error) {
	collection := r.Conn.Database("tusk").Collection("posts")

	var deletedPost internal.Post
	err := collection.FindOneAndDelete(ctx, primitive.M{"_id": id}).Decode(&deletedPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return internal.Post{}, ErrPostNotFound
		}
		return internal.Post{}, fmt.Errorf("failed to delete post: %v", err)
	}

	return deletedPost, nil
}

func (r *Repository) Update(id primitive.ObjectID, updatedPost internal.Post, ctx *gin.Context) (internal.Post, error) {
	collection := r.Conn.Database("tusk").Collection("posts")

	update := bson.M{
		"$set": bson.M{
			"username":  updatedPost.Username,
			"title":     updatedPost.Title,
			"body":      updatedPost.Body,
			"createdAt": updatedPost.CreatedAt,
			"tags":      updatedPost.Tags,
		},
	}

	_, err := collection.UpdateOne(ctx, primitive.M{"_id": id}, update)
	if err != nil {
		return internal.Post{}, fmt.Errorf("failed to update post: %v", err)
	}

	// Retrieve the updated post
	var updated internal.Post
	err = collection.FindOne(ctx, primitive.M{"_id": id}).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return internal.Post{}, ErrPostNotFound
		}
		return internal.Post{}, fmt.Errorf("failed to retrieve updated post: %v", err)
	}

	return updated, nil
}
