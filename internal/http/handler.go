package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guisteink/tusk/internal"
	"github.com/guisteink/tusk/internal/database"
	"github.com/guisteink/tusk/internal/post"
)

var service post.Service

func Configure() {
	service = post.Service{
		Repository: post.Repository{
			Conn: database.Conn,
		},
	}
}

func HandleNewPost(ctx *gin.Context) {
	var post internal.Post
	if err := ctx.BindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := service.Create(post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
