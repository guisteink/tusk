package http

import (
	"net/http"
	"log"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"

	"github.com/guisteink/tusk/internal"
	"github.com/guisteink/tusk/internal/post"
)

var service post.Service

func Configure(conn *mongo.Client) {
	service = post.Service{
		Repository: post.Repository{
			Conn: conn,
		},
	}
}
func handleError(ctx *gin.Context, statusCode int, message string, err error) {
	log.Printf("Error: %v", err)
	ctx.JSON(statusCode, gin.H{"error": message})
}

func handleErrors(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, post.ErrPostBodyEmpty):
		handleError(ctx, http.StatusBadRequest, "Post body is empty", err)
	case errors.Is(err, post.ErrPostBodyExceedsLimit):
		handleError(ctx, http.StatusBadRequest, "Post body exceeds limit", err)
	case errors.Is(err, post.ErrPostNotFound):
		handleError(ctx, http.StatusNotFound, "Post not found", err)
	default:
		handleError(ctx, http.StatusInternalServerError, "Internal Server Error", err)
	}
}

func HandleNewPost(ctx *gin.Context) {
	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		handleErrors(ctx, err)
		return
	}

	log.Printf("Creating post: %+v\n", post)
	response, statusCode, err := service.Create(post)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	log.Printf("Post created successfully with ID %s: %+v\n", response.Post.ID, response.Post)
	ctx.JSON(statusCode, response)
}
