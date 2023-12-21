package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type OpenAIResponse struct {
	Revision string   `json:"revision"`
	Tips     string   `json:"tips"`
	Tags     []string `json:"tags"`
}

type OpenAIClient struct {
	Client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		Client: openai.NewClient(apiKey),
	}
}

func buildChatCompletionRequest(PostBody string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("Faça a correção ortográfica do texto a seguir: { %v }, descreva em inglês o que foi escrito errado pelo aluno e o que pode ser melhorado. Coloque 3 tags avaliando o texto", PostBody),
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Coloque a resposta numa estrutura JSON valida, com revision, tips e tags[]",
			},
		},
	}
}

func (o *OpenAIClient) CreateCompletion(ctx context.Context, PostBody string) (response OpenAIResponse, err error) {
	logrus.Infof("Building revision powered by Openai")

	req := buildChatCompletionRequest(PostBody)

	resp, err := o.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Openai revision error: %v", err)
	}

	content := resp.Choices[0].Message.Content
	logrus.Infof("Content from Openai: %s", content)

	err = json.Unmarshal([]byte(content), &response)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Parsing content error: %v", err)
	}

	return response, nil
}
