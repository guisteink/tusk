package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
}

func SetRoutes(g *gin.Engine) {
	g.GET("/", handleHealthCheck)
	g.POST("/posts", HandleNewPost)
	g.GET("/posts/:id", handleListPostById)
	g.GET("/posts", handleListPosts)
	g.DELETE("/posts/:id", handleDeletePost)
	g.PATCH("/posts/:id", handleUpdatePost)
	g.NoRoute(NotFoundHandler)
}
