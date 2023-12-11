package config

import (
	"os"
)

var (
	MONGODB_HOST            = os.Getenv("DATABASE_URI")
	MONGODB_DATABASE        = os.Getenv("DATABASE_NAME")
	MONGODB_CONNECTION_POOL = 5
	API_PORT                = getAPIPort()
)

func getAPIPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
