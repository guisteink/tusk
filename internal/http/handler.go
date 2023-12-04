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
	log.Println("\nConfiguring")

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
	log.Printf("Request = %+v\n", ctx.Request)

	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		handleError(ctx, http.StatusBadRequest, "Failed to parse JSON", err)
		return
	}

	if err := service.Create(post); err != nil {
		handleError(ctx, http.StatusInternalServerError, "Failed to create post", err)
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
