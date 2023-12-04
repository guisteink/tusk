package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

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

	log.Printf("Post created successfully: %+v\n", response.Post)
	ctx.JSON(statusCode, response)
}

func handleListPostById(ctx *gin.Context) {
	param := ctx.Param("id")

	if param == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": post.ErrIdEmpty,
		})
	}

	log.Printf("Searching for post with id: %s\n", param)
	response, statusCode, err := service.FindByID(param)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	log.Printf("Found post with id %s: %+v\n", param, response)
	ctx.JSON(statusCode, response)
}

func handleListPosts(ctx *gin.Context) {
	log.Printf("Listing all posts")
	response, statusCode, err := service.FindAll()
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	log.Printf("Founded posts: %+v\n", response)
	ctx.JSON(statusCode, response)
}
