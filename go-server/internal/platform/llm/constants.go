package llm

type Provider string

const (
	ProviderOpenAI   Provider = "openai"
	ProviderOllama   Provider = "ollama"
	ProviderGoogleAI Provider = "googleai"
	ProviderGroq     Provider = "groq"
)
