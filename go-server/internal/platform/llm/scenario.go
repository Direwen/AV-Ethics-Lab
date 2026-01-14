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

//go:embed prompts/system.md
var scenarioSystemPrompt string

//go:embed prompts/template.md
var scenarioPromptTemplate string

type ScenarioClient interface {
	domain.Client
	GenerateScenario(ctx context.Context, req domain.ScenarioLLMRequest) (*domain.ScenarioLLMResponse, error)
}

type scenarioClient struct {
	model llms.Model
}

// Implement Client marker interface
func (c *scenarioClient) IsLLMClient() {}

func newScenarioClient(model llms.Model) ScenarioClient {
	return &scenarioClient{model: model}
}

func (c *scenarioClient) GenerateScenario(ctx context.Context, req domain.ScenarioLLMRequest) (*domain.ScenarioLLMResponse, error) {
	// Prepare template
	template := prompts.PromptTemplate{
		Template:       scenarioPromptTemplate,
		InputVariables: []string{"TemplateName", "Dimensions", "Factors", "EgoPosition", "EgoOrientation", "ZoneA", "ZoneB", "ZoneC"},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}

	// Inject data
	promptStr, err := template.Format(map[string]any{
		"TemplateName":   req.TemplateName,
		"Dimensions":     req.GridDimensions,
		"Factors":        req.Factors,
		"EgoPosition":    formatCoordForLLM(req.EgoPosition),
		"EgoOrientation": req.EgoOrientation,
		"ZoneA":          formatZoneForLLM(req.TridentZones.ZoneA),
		"ZoneB":          formatZoneForLLM(req.TridentZones.ZoneB),
		"ZoneC":          formatZoneForLLM(req.TridentZones.ZoneC),
	})
	if err != nil {
		return nil, err
	}

	// Debug output
	// fmt.Println("========== SYSTEM PROMPT ==========")
	// fmt.Println(scenarioSystemPrompt)
	// fmt.Println("========== USER PROMPT ==========")
	// fmt.Println(promptStr)
	// fmt.Println("===================================")

	// Call LLM
	res, err := c.model.GenerateContent(
		ctx,
		[]llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, scenarioSystemPrompt),
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

	var response domain.ScenarioLLMResponse
	if err := json.Unmarshal([]byte(res.Choices[0].Content), &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}

func formatCoordForLLM(coord domain.Coordinate) string {
	return fmt.Sprintf("[%d, %d]", coord.Row, coord.Col)
}

func formatZoneForLLM(zone domain.TridentZone) string {
	result, _ := json.Marshal(zone.Coordinates)
	return string(result)
}
