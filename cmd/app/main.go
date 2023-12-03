package main

import (
	"log"
	// "net/http"
	// "io"
	"os"

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

	defer conn.Close()

	g := gin.Default()
	http.Configure()
	http.SetRoutes(g)
	g.Run(":3000")
}
