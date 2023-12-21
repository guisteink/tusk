package infra

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/guisteink/tusk/config"
	"github.com/guisteink/tusk/internal"
	"github.com/guisteink/tusk/internal/post"
)

var (
	service              post.Service
	globalQueue          *Queue
	logger               = logrus.New()
	openAIClientInstance *OpenAIClient
)

func initLogger() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

func initOpenAIClient(openAIClient *OpenAIClient) {
	openAIClientInstance = openAIClient
}

func Configure(queueInstance *Queue, openAIClient *OpenAIClient, svc post.Service) {
	globalQueue = queueInstance
	service = svc
	initLogger()
	initOpenAIClient(openAIClient)
	initProcessWorker()
}

func processQueueWorker(processWorkerIntervalInSeconds int) {
	const startingProcessWorkerMsg = "Starting process worker"
	logger.Infof(startingProcessWorkerMsg)

	for {
		data, err := globalQueue.Dequeue()
		if err != nil {
			time.Sleep(time.Duration(processWorkerIntervalInSeconds) * time.Second)
			continue
		}

		var deserializedPost internal.Post
		err = json.Unmarshal(data, &deserializedPost)
		if err != nil {
			logger.Infof("Deserialization error: %v", err)
			continue
		}

		response, err := openAIClientInstance.CreateCompletion(context.Background(), deserializedPost.Body)
		if err != nil {
			logger.Infof("OpenAI error: %v", err)
			continue
		}

		updatedPost := internal.Post{
			Body:      deserializedPost.Body,
			Title:     deserializedPost.Title,
			Username:  deserializedPost.Username,
			CreatedAt: deserializedPost.CreatedAt,
			Revision:  response.Revision,
			Tips:      response.Tips,
			Tags:      response.Tags,
		}

		_, _, err = service.UpdateByID(deserializedPost.ID.Hex(), updatedPost)
		if err != nil {
			logger.Infof("Error updating post: %v", err)
			continue
		}

		time.Sleep(time.Duration(processWorkerIntervalInSeconds) * time.Second)
	}
}

func initProcessWorker() {
	processWorkerIntervalInSecondsStr := config.PROCESS_WORKER_INTERVAL_IN_SECONDS
	processWorkerIntervalInSeconds, err := strconv.Atoi(processWorkerIntervalInSecondsStr)
	if err != nil {
		logger.Fatalf("Failed to convert process worker interval to int: %v", err)
		return
	}

	go processQueueWorker(processWorkerIntervalInSeconds)
}
