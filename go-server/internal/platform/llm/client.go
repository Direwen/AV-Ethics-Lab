package llm

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/direwen/go-server/internal/shared/domain"
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
	GenerateScenario(ctx context.Context, req domain.ScenarioLLMRequest) (*domain.ScenarioLLMResponse, error)
}

type client struct {
	model          llms.Model
	sysMsg         string
	promptTemplate string
}

func (c *client) GenerateScenario(ctx context.Context, req domain.ScenarioLLMRequest) (*domain.ScenarioLLMResponse, error) {
	// Prepare Template
	template := prompts.PromptTemplate{
		Template:       c.promptTemplate,
		InputVariables: []string{"TemplateID", "Dimensions", "Factors", "WalkableCells", "DrivableCells", "BuildingCells"},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}
	// Inject Data into Prompt Template
	templateStr, err := template.Format(map[string]any{
		"TemplateName":  req.TemplateName,
		"Dimensions":    req.GridDimensions,
		"Factors":       req.Factors,
		"WalkableCells": formatCellsForLLM(req.WalkableCells),
		"DrivableCells": formatCellsForLLM(req.DrivableCells),
		"BuildingCells": formatCellsForLLM(req.BuildingCells),
	})
	if err != nil {
		return nil, err
	}

	// TO DEBUG
	fmt.Println("========== SYSTEM PROMPT ==========")
	fmt.Println(c.sysMsg)
	fmt.Println("========== USER PROMPT ==========")
	fmt.Println(templateStr)
	fmt.Println("===================================")

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

	var response domain.ScenarioLLMResponse
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
	case ProviderGroq:
		llm, err = openai.New(
			openai.WithModel(modelName),
			openai.WithBaseURL("https://api.groq.com/openai/v1"),
			openai.WithToken(os.Getenv("GROQ_API_KEY")),
			openai.WithResponseFormat(openai.ResponseFormatJSON),
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

func formatGridForLLM(grid [][]int) string {
	var sb strings.Builder
	sb.WriteString("[\n")
	for _, row := range grid {
		// Format each row as "[0, 0, 9, ...],"
		rowJSON, _ := json.Marshal(row)
		sb.WriteString("  ")
		sb.Write(rowJSON)
		sb.WriteString(",\n")
	}
	sb.WriteString("]")
	return sb.String()
}

func formatCellsForLLM(cells [][2]int) string {
	if len(cells) == 0 {
		return "[]"
	}
	result, _ := json.Marshal(cells)
	return string(result)
}
