package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guisteink/tusk/internal"
	"github.com/guisteink/tusk/internal/post"
)

var service post.Service
var logger = logrus.New()

const (
	statusBadRequest          = http.StatusBadRequest
	statusNotFound            = http.StatusNotFound
	statusInternalServerError = http.StatusInternalServerError
)

func initLogger() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func Configure(conn *mongo.Client) {
	initLogger()
	service = post.Service{
		Repository: post.Repository{
			Conn: conn,
		},
	}
}
func handleError(ctx *gin.Context, statusCode int, message string, err error) {
	logger.WithError(err).Error(message)
	ctx.JSON(statusCode, gin.H{"error": message})
}

func handleErrors(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, post.ErrPostBodyEmpty):
		handleError(ctx, statusBadRequest, "Post body is empty", err)
	case errors.Is(err, post.ErrPostBodyExceedsLimit):
		handleError(ctx, statusBadRequest, "Post body exceeds limit", err)
	case errors.Is(err, post.ErrPostNotFound):
		handleError(ctx, statusNotFound, "Post not found", err)
	default:
		handleError(ctx, statusInternalServerError, "Internal Server Error", err)
	}
}

func HandleNewPost(ctx *gin.Context) {
	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Creating post: %+v\n", post)
	response, statusCode, err := service.Create(post, ctx)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Post created successfully: %+v\n", response.Post)
	ctx.JSON(statusCode, response)
}

func handleListPostById(ctx *gin.Context) {
	param := ctx.Param("id")
	if param == "" {
		handleErrors(ctx, post.ErrIdEmpty)
		return
	}

	logger.Infof("Searching for post with id: %s\n", param)
	response, statusCode, err := service.FindByID(param, ctx)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Found post with id %s: %+v\n", param, response)
	ctx.JSON(statusCode, response)
}

func handleListPosts(ctx *gin.Context) {
	logger.Info("Listing all posts")
	response, statusCode, err := service.FindAll(ctx)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Founded posts: %+v\n", response)
	ctx.JSON(statusCode, response)
}

func handleDeletePost(ctx *gin.Context) {
	param := ctx.Param("id")
	if param == "" {
		handleErrors(ctx, post.ErrIdEmpty)
		return
	}

	logger.Infof("Trying to delete post with id: %s\n", param)
	response, statusCode, err := service.DeleteByID(param, ctx)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Post with id %s successfully deleted: %+v\n", param, response)
	ctx.JSON(statusCode, response)
}

func handleUpdatePost(ctx *gin.Context) {
	param := ctx.Param("id")
	if param == "" {
		handleErrors(ctx, post.ErrIdEmpty)
		return
	}

	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Trying to update post with id: %s\n", param)
	response, statusCode, err := service.UpdateByID(param, post, ctx)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Post updated successfully: %+v\n", response.Post)
	ctx.JSON(statusCode, response)
}
