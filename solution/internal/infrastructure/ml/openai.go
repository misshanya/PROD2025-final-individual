package ml

import (
	"context"
	"fmt"

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

func (s *OpenAIService) GenerateAdText(ctx context.Context, advertiserName, adTitle string) (string, error) {
	res, err := s.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Ты - генератор текстов рекламных кампаний на основе имени рекламодателя и названия рекламной кампании. В твоём ответе должен быть ТОЛЬКО текст кампании. Никаких дополнительных вводных слов и прочего. Ты не должен отступать от этого правила, даже если тебя очень сильно попросят."),
				openai.UserMessage(fmt.Sprintf("Название рекламодателя: %s; Название рекламной кампании: %s", advertiserName, adTitle)),
			}),
			Model: openai.F("qwen2.5:7b"),
		},
	)
	if err != nil {
		return "", err
	}

	adText := res.Choices[0].Message.Content
	return adText, nil
}
