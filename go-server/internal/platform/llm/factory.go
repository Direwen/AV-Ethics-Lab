package llm

import (
	"fmt"

	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/direwen/go-server/internal/util"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

// TaskConfig holds model configuration for a specific task
type TaskConfig struct {
	Model    string
	Provider Provider
}

// NewClient creates a client for the specified task
func NewClient(task domain.LLMTask, key string) (domain.Client, error) {
	config := getTaskConfig(task)
	model, err := initModel(config, key)
	if err != nil {
		return nil, fmt.Errorf("failed to init model for task %s: %w", task, err)
	}

	switch task {
	case domain.TaskScenario:
		return newScenarioClient(model), nil
	case domain.TaskFeedback:
		return newFeedbackClient(model), nil
	default:
		return nil, fmt.Errorf("unsupported task: %s", task)
	}
}

// getTaskConfig reads config from env based on task
func getTaskConfig(task domain.LLMTask) TaskConfig {
	switch task {
	case domain.TaskScenario:
		return TaskConfig{
			Model:    util.GetEnvOrDefault("SCENARIO_MODEL", "qwen/qwen3-32b"),
			Provider: Provider(util.GetEnvOrDefault("SCENARIO_PROVIDER", "groq")),
		}
	case domain.TaskFeedback:
		return TaskConfig{
			Model:    util.GetEnvOrDefault("FEEDBACK_MODEL", "nvidia/nemotron-nano-9b-v2:free"),
			Provider: Provider(util.GetEnvOrDefault("FEEDBACK_PROVIDER", "openrouter")),
		}
	default:
		return TaskConfig{}
	}
}

// initModel creates the LLM model based on provider
func initModel(config TaskConfig, keys ...string) (llms.Model, error) {

	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	switch config.Provider {
	case ProviderOpenAI:
		return openai.New(
			openai.WithModel(config.Model),
			openai.WithResponseFormat(openai.ResponseFormatJSON),
		)
	case ProviderOllama:
		return ollama.New(
			ollama.WithModel(config.Model),
			ollama.WithFormat("json"),
		)
	case ProviderGroq:
		return openai.New(
			openai.WithModel(config.Model),
			openai.WithBaseURL("https://api.groq.com/openai/v1"),
			openai.WithToken(key),
			openai.WithResponseFormat(openai.ResponseFormatJSON),
		)
	case ProviderOpenRouter:
		return openai.New(
			openai.WithModel(config.Model),
			openai.WithBaseURL("https://openrouter.ai/api/v1"),
			openai.WithToken(key),
			openai.WithResponseFormat(openai.ResponseFormatJSON),
		)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}
