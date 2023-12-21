package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"

	"github.com/guisteink/tusk/config"
)

type OpenAIResponse struct {
	Revision string   `json:"revision"`
	Tips     string   `json:"tips"`
	Tags     []string `json:"tags"`
}

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(config.OPENAI_APIKEY),
	}
}

func (o *OpenAIClient) CreateCompletion(ctx context.Context, PostBody string) (OpenAIResponse, error) {
	logrus.Infof("Building revision powered by Openai")

	/** todo: insert prompt into varenv **/
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("Faça a correção ortográfica do texto a seguir: { %v }, descreva em inglês o que foi escrito errado pelo aluno e o que pode ser melhorado. Coloque 3 tags avaliando o texto", PostBody),
			}, {
				Role: openai.ChatMessageRoleSystem,
				// Content: "Coloque sempre a resposta na seguinte estrutura: {'revision': string, 'tips': string, 'tags': []string}, evitar o uso de contracções ou apóstrofos possessivo na resposta",
				Content: "Coloque a resposta numa estrutura JSON valida, com revision, tips e tags[]",
			},
		},
	}

	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Openai revision error: %v", err)
	}

	// Extracting revision, tips, and tags from the content string
	content := resp.Choices[0].Message.Content
	logrus.Infof("Content from Openai: %s", content)

	var responseStruct OpenAIResponse
	err = json.Unmarshal([]byte(content), &responseStruct)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Parsing content error: %v", err)
	}

	return responseStruct, nil
}
