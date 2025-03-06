package ml

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIService struct {
	client          *openai.Client
	moderationModel string
	generationModel string
}

func NewOpenAIService(
	baseURL, apiKey string,
	moderationModel, generationModel string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(
			option.WithBaseURL(baseURL),
			option.WithAPIKey(apiKey),
		),
		moderationModel: moderationModel,
		generationModel: generationModel,
	}
}

func (s *OpenAIService) ValidateAdText(ctx context.Context, text string) (bool, error) {
	res, err := s.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("Ты - модератор. Ты должен проверять тексты рекламных кампаний на что-то неприличное (маты и оскорбления). В итоге твой ответ должен быть только + (если текст проходит модерацию) или - (если текст не проходит модерацию) и критерий запрета через двоеточие (Например, -:Критерий)"),
				openai.UserMessage(text),
			}),
			Model: openai.F(s.moderationModel),
		},
	)
	if err != nil {
		return false, err
	}

	var response bool
	if strings.Split(res.Choices[0].Message.Content, ":")[0] == "+" {
		response = true
	}

	log.Println(res.Choices[0].Message.Content)

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
			Model: openai.F(s.generationModel),
		},
	)
	if err != nil {
		return "", err
	}

	adText := res.Choices[0].Message.Content
	return adText, nil
}
