package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/direwen/go-server/internal/platform/llm"
	"github.com/lpernett/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: No .env file found, using system env vars")
	}

	modelName := os.Getenv("LLM_MODEL")
	if modelName == "" {
		modelName = "mistral:7b" // fallback default
	}

	fmt.Printf("Creating LLM client with model: %s, provider: %s\n", modelName, llm.ProviderOllama)

	// Create client
	client, err := llm.NewClient(modelName, llm.ProviderOllama)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	fmt.Println("LLM client created successfully!")

	// Test GenerateScenario
	ctx := context.Background()
	req := llm.ScenarioRequest{
		TemplateName:   "test",
		GridDimensions: "3x3",
		GridData:       []string{"A", "B", "C"},
		Factors: map[string]any{
			"difficulty": "easy",
		},
	}

	fmt.Println("Calling GenerateScenario...")
	resp, err := client.GenerateScenario(ctx, req)
	if err != nil {
		log.Fatalf("GenerateScenario failed: %v", err)
	}

	fmt.Printf("Response:\n  Narrative: %s\n  Entities: %+v\n  Factors: %+v\n",
		resp.Narrative, resp.Entities, resp.Factors)
}
