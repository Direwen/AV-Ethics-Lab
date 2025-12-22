package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/direwen/go-server/internal/platform/llm"
	"github.com/direwen/go-server/internal/scenario"
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
		TemplateName:   "intersection_01",
		GridDimensions: "10x10",
		GridData: [][]int{
			{0, 0, 0, 3, 9, 10, 3, 0, 0, 0},
			{0, 0, 0, 3, 11, 12, 3, 0, 0, 0},
			{3, 3, 3, 3, 13, 14, 3, 3, 3, 3},
			{9, 11, 13, 15, 16, 17, 15, 13, 11, 9},
			{10, 12, 14, 17, 18, 18, 17, 14, 12, 10},
			{3, 3, 3, 3, 13, 14, 3, 3, 3, 3},
			{0, 0, 0, 3, 11, 12, 3, 0, 0, 0},
			{0, 0, 0, 3, 9, 10, 3, 0, 0, 0},
			{0, 0, 0, 3, 9, 10, 3, 0, 0, 0},
			{0, 0, 0, 3, 9, 10, 3, 0, 0, 0},
		},
		Factors: scenario.ScenarioFactors{
			Visibility:         "Foggy",
			RoadCondition:      "Wet",
			Location:           "Urban Intersection",
			BrakeStatus:        "Functional",
			Speed:              "40 km/h",
			HasTailgater:       true,
			PrimaryEntity:      "ped_child",
			PrimaryBehavior:    "Violation",
			BackgroundEntities: []string{"ped_adult", "car_sedan"},
		},
	}

	fmt.Println("Calling GenerateScenario...")
	resp, err := client.GenerateScenario(ctx, req)
	if err != nil {
		log.Fatalf("GenerateScenario failed: %v", err)
	}

	fmt.Printf("Response:\n  Verification: %s\n  Narrative: %s\n  Entities: %+v\n",
		resp.Verification, resp.Narrative, resp.Entities)
}
