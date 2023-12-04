package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
}

func SetRoutes(g *gin.Engine) {
	g.POST("/post", HandleNewPost)
	g.NoRoute(NotFoundHandler)
}
