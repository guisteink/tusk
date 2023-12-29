package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"

	"github.com/guisteink/tusk/config"
)

type OpenAIResponse struct {
	Revision     string   `json:"revision"`
	Tips         string   `json:"tips"`
	Tags         []string `json:"tags"`
	WritingScore float64  `json:"writingScore"`
}

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(config.OPENAI_APIKEY),
	}
}

func buildChatCompletionRequest(PostBody string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("Faça a correção ortográfica do texto a seguir: { %v }, descreva em inglês o que foi escrito errado pelo aluno e o que pode ser melhorado. Coloque 3 tags avaliando o texto", PostBody),
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Coloque a resposta numa estrutura JSON valida, com revision, tips, writingScore(0.00~10.00) e tags[]",
			},
		},
	}
}

func (o *OpenAIClient) CreateCompletion(ctx context.Context, PostBody string) (OpenAIResponse, error) {
	logger.Info("Building revision powered by Openai")

	// todo: insert prompt into varenv
	req := buildChatCompletionRequest(PostBody)

	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Openai revision error: %v", err)
	}

	content := resp.Choices[0].Message.Content
	logger.Infof("Content from Openai: %s", content)

	var responseStruct OpenAIResponse
	err = json.Unmarshal([]byte(content), &responseStruct)
	if err != nil {
		return OpenAIResponse{}, fmt.Errorf("Parsing content error: %v", err)
	}

	logger.Infof("Successfully by Openai")

	return responseStruct, nil
}
