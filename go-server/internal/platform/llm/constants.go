package llm

type Provider string

const (
	ProviderOpenAI     Provider = "openai"
	ProviderOllama     Provider = "ollama"
	ProviderGroq       Provider = "groq"
	ProviderOpenRouter Provider = "openrouter"
)

type Task string

const (
	TaskScenario Task = "scenario"
	TaskFeedback Task = "feedback"
)
