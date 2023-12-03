package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/guisteink/tusk/internal/database"
	"github.com/guisteink/tusk/internal/http"
)

func main() {
	connectionString := os.Getenv("DATABASE_URI")
	log.Println("DATABASE_URI:", connectionString)

	client, err := database.NewConnection(connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("DATABASE_URI:", client)
	}

	g := gin.Default()
	http.Configure()
	http.SetRoutes(g)
	g.Run(":3000")
}
