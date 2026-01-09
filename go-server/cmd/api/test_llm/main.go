package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/direwen/go-server/internal/platform/llm"
	"github.com/direwen/go-server/internal/shared/domain"
	"github.com/lpernett/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system env vars")
	}

	fmt.Println("Creating scenario LLM client...")

	// Create client using strategy pattern
	client, err := llm.NewClient(llm.TaskScenario)
	if err != nil {
		log.Fatalf("Failed to create LLM client: %v", err)
	}
	scenarioClient := client.(llm.ScenarioClient)

	fmt.Println("LLM client created successfully!")

	// Test GenerateScenario with the Trident framework
	// Simulating a spawn point on the 4-Way Urban Intersection
	ctx := context.Background()

	// Mock Ego position (as if from GetRandomTridentSpawn)
	egoPosition := domain.Coordinate{Row: 6, Col: 5}
	egoOrientation := domain.DirectionWest

	// Mock enriched trident zones
	tridentZones := domain.TridentZones{
		ZoneA: domain.TridentZone{
			Coordinates: []domain.EnrichedCoordinate{
				{Coordinate: domain.Coordinate{Row: 6, Col: 2}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
				{Coordinate: domain.Coordinate{Row: 6, Col: 1}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
				{Coordinate: domain.Coordinate{Row: 6, Col: 0}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
			},
		},
		ZoneB: domain.TridentZone{
			Coordinates: []domain.EnrichedCoordinate{
				{Coordinate: domain.Coordinate{Row: 7, Col: 2}, Surface: domain.SurfaceWalkable, Orientation: ""},
				{Coordinate: domain.Coordinate{Row: 7, Col: 1}, Surface: domain.SurfaceWalkable, Orientation: ""},
				{Coordinate: domain.Coordinate{Row: 7, Col: 0}, Surface: domain.SurfaceWalkable, Orientation: ""},
			},
		},
		ZoneC: domain.TridentZone{
			Coordinates: []domain.EnrichedCoordinate{
				{Coordinate: domain.Coordinate{Row: 5, Col: 2}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
				{Coordinate: domain.Coordinate{Row: 5, Col: 1}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
				{Coordinate: domain.Coordinate{Row: 5, Col: 0}, Surface: domain.SurfaceDrivable, Orientation: domain.DirectionWest},
			},
		},
	}

	req := domain.ScenarioLLMRequest{
		TemplateName:   "4-Way Urban Intersection",
		GridDimensions: "20:11",
		Factors: domain.ScenarioFactors{
			Visibility:         "Rain",
			RoadCondition:      "Wet",
			Location:           "US",
			BrakeStatus:        "Fade",
			Speed:              "Medium",
			HasTailgater:       true,
			PrimaryEntity:      "ped_elderly",
			PrimaryBehavior:    "Violation",
			BackgroundEntities: []string{"ped_child", "vehicle_car"},
		},
		EgoPosition:    egoPosition,
		EgoOrientation: egoOrientation,
		TridentZones:   tridentZones,
	}

	fmt.Println("Calling GenerateScenario...")
	resp, err := scenarioClient.GenerateScenario(ctx, req)
	if err != nil {
		log.Fatalf("GenerateScenario failed: %v", err)
	}

	fmt.Printf("\n========== RESPONSE ==========\n")
	rawJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(rawJSON))
}
