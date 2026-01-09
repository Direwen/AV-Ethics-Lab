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
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	fmt.Println("Creating feedback LLM client...")

	client, err := llm.NewClient(llm.TaskFeedback)
	if err != nil {
		log.Fatalf("Failed to create feedback client: %v", err)
	}
	feedbackClient := client.(llm.FeedbackClient)

	fmt.Println("Feedback client created!")

	// Dummy demographic data
	demographic := domain.Demographic{
		AgeRange:          3, // 25-34
		Gender:            1, // Male
		Country:           "US",
		Occupation:        "Software Engineer",
		DrivingExperience: 2, // 3-10 years
	}

	// Dummy responses - simulating a user who prioritizes pedestrians
	responses := []domain.EnrichedResponse{
		{
			Narrative: "You're driving through a foggy intersection when an elderly woman steps onto the crosswalk. A truck is tailgating you aggressively.",
			Factors: domain.ScenarioFactors{
				Visibility:      "Fog",
				RoadCondition:   "Wet",
				Speed:           "Medium",
				BrakeStatus:     "Active",
				HasTailgater:    true,
				PrimaryEntity:   "ped_elderly",
				PrimaryBehavior: "Compliant",
			},
			RankedOptions:  []string{"swerve_left", "maintain", "swerve_right"},
			HasInteracted:  true,
			ResponseTimeMs: 4500,
			IsTimeout:      false,
		},
		{
			Narrative: "A child runs into the street chasing a ball. Your brakes are failing and there's a car in the oncoming lane.",
			Factors: domain.ScenarioFactors{
				Visibility:      "Clear",
				RoadCondition:   "Dry",
				Speed:           "High",
				BrakeStatus:     "Failed",
				HasTailgater:    false,
				PrimaryEntity:   "ped_child",
				PrimaryBehavior: "Violation",
			},
			RankedOptions:  []string{"swerve_right", "swerve_left", "maintain"},
			HasInteracted:  true,
			ResponseTimeMs: 2100,
			IsTimeout:      false,
		},
		{
			Narrative: "A jaywalking criminal crosses in front of you while police chase him. Swerving left hits a bus stop with waiting passengers.",
			Factors: domain.ScenarioFactors{
				Visibility:      "Night",
				RoadCondition:   "Icy",
				Speed:           "Low",
				BrakeStatus:     "Fade",
				HasTailgater:    false,
				PrimaryEntity:   "ped_criminal",
				PrimaryBehavior: "Violation",
			},
			RankedOptions:  []string{"maintain", "swerve_right", "swerve_left"},
			HasInteracted:  true,
			ResponseTimeMs: 8900,
			IsTimeout:      false,
		},
		{
			Narrative: "Heavy rain obscures your vision. A pregnant woman is crossing legally while a motorcyclist speeds in the opposite lane.",
			Factors: domain.ScenarioFactors{
				Visibility:      "Rain",
				RoadCondition:   "Wet",
				Speed:           "Medium",
				BrakeStatus:     "Active",
				HasTailgater:    true,
				PrimaryEntity:   "ped_pregnant",
				PrimaryBehavior: "Compliant",
			},
			RankedOptions:  []string{"swerve_left", "swerve_right", "maintain"},
			HasInteracted:  false,
			ResponseTimeMs: 12000,
			IsTimeout:      true,
		},
	}

	req := domain.FeedbackLLMRequest{
		Demographic: demographic,
		Responses:   responses,
	}

	fmt.Println("Calling GenerateFeedback...")
	ctx := context.Background()

	resp, err := feedbackClient.GenerateFeedback(ctx, req)
	if err != nil {
		log.Fatalf("GenerateFeedback failed: %v", err)
	}

	fmt.Printf("\n========== FEEDBACK RESPONSE ==========\n")
	rawJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(rawJSON))
}
