package config

import (
	"os"
)

var (
	MONGODB_HOST                       = os.Getenv("DATABASE_URI")
	MONGODB_DATABASE                   = os.Getenv("DATABASE_NAME")
	MONGODB_CONNECTION_POOL            = 5
	API_PORT                           = getAPIPort()
	PROCESS_WORKER_INTERVAL_IN_SECONDS = os.Getenv("PROCESS_WORKER_INTERVAL_IN_SECONDS")
	OPENAI_APIKEY                      = os.Getenv("OPENAI_APIKEY")
)

func getAPIPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return port
}
