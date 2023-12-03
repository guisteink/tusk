package http

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(g *gin.Engine) {
	g.POST("/post", new)
	// g.DELETE("/post/:id", DeletePosts)
	// g.GET("/post/:id", GetPosts)
	// g.UPDATE("/post/:id", UpdatePosts)
}
