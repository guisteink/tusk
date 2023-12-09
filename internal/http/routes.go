package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
}

func SetRoutes(g *gin.Engine) {
	v1 := g.Group("/v1")

	v1.GET("/", handleHealthCheck)
	v1.POST("/posts", HandleNewPost)
	v1.GET("/posts/:id", handleListPostById)
	v1.GET("/posts", handleListPosts)
	v1.DELETE("/posts/:id", handleDeletePost)
	v1.PATCH("/posts/:id", handleUpdatePost)

	g.NoRoute(NotFoundHandler)
}
