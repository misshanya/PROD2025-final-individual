package ml

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(baseURL, apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(apiKey),
		),
	}
}

func (s *OpenAIService) ValidateAdText(ctx context.Context, text string) (bool, error) {
	res, err := s.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Ты - модератор. Ты должен проверять тексты рекламных кампаний на что-то неприличное. В итоге твой ответ должен быть только \"+\" (если текст проходит модерацию) или \"-\" (если текст не проходит модерацию)"),
				openai.UserMessage(text),
			}),
			Model: openai.F("qwen2.5:7b"),
		},
	)
	if err != nil {
		return false, err
	}

	var response bool
	if res.Choices[0].Message.Content == "+" {
		response = true
	}

	return response, nil
}
