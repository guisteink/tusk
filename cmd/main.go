package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/guisteink/tusk/config"
	"github.com/guisteink/tusk/internal/database"
	"github.com/guisteink/tusk/internal/http"
)

var logger = logrus.New()

func main() {
	mongoHost := config.MONGODB_HOST
	mongoDatabase := config.MONGODB_DATABASE
	port := config.API_PORT

	connectionString := fmt.Sprintf("%s/%s", mongoHost, mongoDatabase)

	conn, err := database.NewConnection(connectionString)
	if conn == nil {
		logger.Errorf("Database connection is nil.")
		os.Exit(1)
	} else if err != nil {
		logger.Errorf("Failed to establish database connection:", err)
		panic(err)
	}

	http.Configure(conn, config.OPENAI_APIKEY)

	g := gin.Default()

	g.Use(cors.Default())

	http.SetRoutes(g)

	g.Run(":" + port)
}
