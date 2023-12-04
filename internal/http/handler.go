package http

import (
	"net/http"
	"log"

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

func HandleNewPost(ctx *gin.Context) {
	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		handleError(ctx, http.StatusBadRequest, "Failed to parse JSON", err)
		return
	}

	log.Printf("Creating post: %+v\n", post)
	err := service.Create(post)
	if err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	log.Printf("Post created successfully: %+v\n", post)
	ctx.JSON(http.StatusCreated, post)
}
