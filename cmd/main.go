package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/guisteink/tusk/config"
	"github.com/guisteink/tusk/internal/database"
	"github.com/guisteink/tusk/internal/http"
)

func main() {
	mongoHost := config.MONGODB_HOST
	mongoDatabase := config.MONGODB_DATABASE
	port := config.API_PORT

	connectionString := fmt.Sprintf("%s/%s", mongoHost, mongoDatabase)

	conn, err := database.NewConnection(connectionString)
	if conn == nil {
		log.Println("Database connection is nil.")
		os.Exit(1)
	} else if err != nil {
		log.Fatal("Failed to establish database connection:", err)
		panic(err)
	}

	http.Configure(conn)

	g := gin.Default()
	http.SetRoutes(g)

	g.Run(":" + port)
}
