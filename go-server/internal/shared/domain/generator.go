package domain

import (
	"math/rand"

	"github.com/direwen/go-server/internal/scenario"
)

func GenerateBalancedDesign(sessionID string, count int) []scenario.ScenarioFactors {

	var deck []scenario.ScenarioFactors

	for i := 0; i < count; i++ {

		// CRITICAL FACTORS (requires the balance)
		vis := Visibilities[i%len(Visibilities)]
		brake := BrakeStatuses[i%len(BrakeStatuses)]

		// Balance Legal Status exactly 50/50
		var behavior Behavior
		if i%2 == 0 {
			behavior = BehaviorViolation
		} else {
			behavior = BehaviorCompliant
		}

		// RANDOMIZED FACTORS
		road := RoadConditions[rand.Intn(len(RoadConditions))]
		loc := Locations[rand.Intn(len(Locations))]
		speed := Speeds[rand.Intn(len(Speeds))]

		// LOGIC CONSTRAINTS
		if road == RoadConditionIcy && speed == SpeedLow {
			speed = SpeedMedium
		}

		// CASTING
		// STAR selection
		primaryEntity := CastPrimaryEntity()
		// Background Noise Selection
		backgroundEntities := CastBackgroundEntities(2, 4)

		factors := scenario.ScenarioFactors{
			Visibility:         string(vis),
			BrakeStatus:        string(brake),
			RoadCondition:      string(road),
			Location:           string(loc),
			Speed:              string(speed),
			HasTailgater:       rand.Intn(2) == 1,
			PrimaryEntity:      primaryEntity,
			PrimaryBehavior:    string(behavior),
			BackgroundEntities: backgroundEntities,
		}

		deck = append(deck, factors)
	}

	// SHUFFLE
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}
