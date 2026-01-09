package llm

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

//go:embed prompts/feedback_system.md
var feedbackSystemPrompt string

//go:embed prompts/feedback_template.md
var feedbackPromptTemplate string

type FeedbackClient interface {
	GenerateFeedback(ctx context.Context, req domain.FeedbackLLMRequest) (*domain.FeedbackLLMResponse, error)
}

type feedbackclient struct {
	model llms.Model
}

func newFeedbackClient(model llms.Model) FeedbackClient {
	return &feedbackclient{model: model}
}

func (c *feedbackclient) GenerateFeedback(ctx context.Context, req domain.FeedbackLLMRequest) (*domain.FeedbackLLMResponse, error) {
	template := prompts.PromptTemplate{
		Template:       feedbackPromptTemplate,
		InputVariables: []string{"Demographic", "Responses"},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}

	promptStr, err := template.Format(map[string]any{
		"Demographic": req.Demographic,
		"Responses":   req.Responses,
	})
	if err != nil {
		return nil, err
	}

	// Call LLM
	res, err := c.model.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, feedbackSystemPrompt),
			llms.TextParts(llms.ChatMessageTypeHuman, promptStr),
		},
		llms.WithJSONMode(),
	)
	if err != nil {
		return nil, err
	}

	// Parse response
	if len(res.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	var response domain.FeedbackLLMResponse
	if err := json.Unmarshal([]byte(res.Choices[0].Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}
