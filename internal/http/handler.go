package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guisteink/tusk/infra"
	"github.com/guisteink/tusk/internal"
	"github.com/guisteink/tusk/internal/post"
)

type Queue interface {
	Enqueue(data []byte) error
	Dequeue() ([]byte, error)
	IsEmpty() bool
	GetSize() int
	Peek() ([]byte, error)
	Clear() error
}

type OpenAIClient struct {
	client *openai.Client
}

var globalQueue Queue
var openAIClientInstance OpenAIClient
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
	var queueInstance = &infra.Queue{}
	var openAIClient = infra.NewOpenAIClient("your-openai-api-key")

	initLogger()
	service = post.Service{
		Repository: post.Repository{
			Conn: conn,
		},
	}
	globalQueue = queueInstance
	infra.Configure(queueInstance, openAIClient, service)
}

func handleError(ctx *gin.Context, statusCode int, message string, err error) {
	logger.WithError(err).Error(message)
	ctx.JSON(statusCode, gin.H{"error": message})
}

func handleHealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Health check passed",
	})
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
	response, statusCode, err := service.Create(post)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	serializedPost, err := json.Marshal(response.Post)
	if err != nil {
		logger.Infof("Serialization error: %v", err)
		return
	}

	err = globalQueue.Enqueue(serializedPost)
	if err != nil {
		logger.Infof("Enqueue error: %v", err)
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
	response, statusCode, err := service.FindByID(param)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Found post with id %s: %+v\n", param, response)
	ctx.JSON(statusCode, response)
}

func handleListPosts(ctx *gin.Context) {
	logger.Info("Listing all posts")
	response, statusCode, err := service.FindAll()
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
	response, statusCode, err := service.DeleteByID(param)
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

	objID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	updatedPost := internal.Post{
		ID:        objID,
		Username:  post.Username,
		Title:     post.Title,
		Body:      post.Body,
		CreatedAt: post.CreatedAt,
		Tags:      post.Tags,
	}

	response, statusCode, err := service.UpdateByID(param, updatedPost)
	if err != nil {
		handleErrors(ctx, err)
		return
	}

	logger.Infof("Post updated successfully: %+v\n", response.Post)
	ctx.JSON(statusCode, response)
}
