package http

import (
	"log"

	"github.com/gin-gonic/gin"
)


func SetRoutes(g *gin.Engine) {
	log.Println("Setting routes")
	g.POST("/post", HandleNewPost)
}
