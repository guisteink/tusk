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
	connectionPort := os.Getenv("PORT")

	conn, err := database.NewConnection(connectionString)
	if conn == nil {
		log.Println("Database connection is nil.")
		os.Exit(1)
	} else if err != nil {
		panic(err)
	}

	// Pass the MongoDB client to the http.Configure function
	http.Configure(conn)

	g := gin.Default()
	http.SetRoutes(g)

	g.Run(":" + connectionPort)

}
