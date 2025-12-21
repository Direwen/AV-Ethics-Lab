package llm

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

//go:embed prompts/system.md
var systemPromptSrc string

//go:embed prompts/template.md
var promptTemplateSrc string

type Client interface {
	GenerateScenario(ctx context.Context, req ScenarioRequest) (*ScenarioResponse, error)
}

type client struct {
	model          llms.Model
	sysMsg         string
	promptTemplate string
}

func (c *client) GenerateScenario(ctx context.Context, req ScenarioRequest) (*ScenarioResponse, error) {
	// Prepare Template
	template := prompts.PromptTemplate{
		Template:       c.promptTemplate,
		InputVariables: []string{},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}
	// Inject Data into Prompt Template
	templateStr, err := template.Format(map[string]any{})
	if err != nil {
		return nil, err
	}

	// Call LLM
	res, err := c.model.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, c.sysMsg),
			llms.TextParts(llms.ChatMessageTypeHuman, templateStr),
		},
		llms.WithJSONMode(),
	)
	if err != nil {
		return nil, err
	}

	// Parse Response
	if len(res.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	var response ScenarioResponse
	if err := json.Unmarshal([]byte(res.Choices[0].Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}

func NewClient(modelName string, provider Provider) (Client, error) {
	var llm llms.Model
	var err error

	switch provider {
	case ProviderOpenAI:
		llm, err = openai.New(
			openai.WithModel(modelName),
			openai.WithResponseFormat(openai.ResponseFormatJSON),
		)
	case ProviderOllama:
		llm, err = ollama.New(
			ollama.WithModel(modelName),
			ollama.WithFormat("json"),
		)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	return &client{
		model:          llm,
		sysMsg:         systemPromptSrc,
		promptTemplate: promptTemplateSrc,
	}, nil
}
