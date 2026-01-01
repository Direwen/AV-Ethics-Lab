package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/direwen/go-server/internal/platform/llm"
	"github.com/direwen/go-server/internal/scenario"
	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/lpernett/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system env vars")
	}

	modelName := os.Getenv("LLM_MODEL")
	if modelName == "" {
		modelName = "llama-3.3-70b-versatile" // Groq default
	}

	fmt.Printf("Creating LLM client with model: %s, provider: %s\n", modelName, llm.ProviderGroq)

	// Create client
	client, err := llm.NewClient(modelName, llm.ProviderGroq)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}

	fmt.Println("LLM client created successfully!")

	// Test GenerateScenario with the full prompt
	ctx := context.Background()
	req := domain.ScenarioLLMRequest{
		TemplateName:   "4-Way Urban Intersection",
		GridDimensions: "20:11",
		GridData: [][]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 19, 15, 15, 15, 20, 1, 1, 1, 1, 1, 1, 1},
			{3, 3, 3, 3, 3, 3, 3, 3, 5, 11, 11, 11, 7, 3, 3, 3, 3, 3, 3, 3},
			{9, 9, 9, 9, 9, 9, 9, 16, 11, 11, 11, 11, 11, 16, 9, 9, 9, 9, 9, 9},
			{13, 13, 13, 13, 13, 13, 13, 16, 11, 11, 11, 11, 11, 16, 13, 13, 13, 13, 13, 13},
			{9, 9, 9, 9, 9, 9, 9, 16, 11, 11, 11, 11, 11, 16, 9, 9, 9, 9, 9, 9},
			{4, 4, 4, 4, 4, 4, 4, 4, 6, 11, 11, 11, 8, 4, 4, 4, 4, 4, 4, 4},
			{2, 2, 2, 2, 2, 2, 2, 2, 19, 15, 15, 15, 20, 2, 2, 2, 2, 2, 2, 2},
			{0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 19, 10, 14, 10, 20, 0, 0, 0, 0, 0, 0, 0},
		},
		Factors: scenario.ScenarioFactors{
			Visibility:         "Rain",
			RoadCondition:      "Wet",
			Location:           "4-Way Urban Intersection",
			BrakeStatus:        "Fade",
			Speed:              "Medium",
			HasTailgater:       true,
			PrimaryEntity:      "ped_elderly",
			PrimaryBehavior:    "Compliant",
			BackgroundEntities: []string{"obstacle_barrier", "vehicle_bus", "vehicle_car"},
		},
	}

	fmt.Println("Calling GenerateScenario...")
	resp, err := client.GenerateScenario(ctx, req)
	if err != nil {
		log.Fatalf("GenerateScenario failed: %v", err)
	}

	fmt.Printf("\n========== RESPONSE ==========\n")
	rawJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(rawJSON))
}
