package domain

import (
	"math/rand"
	"os"
	"strconv"
)

func GenerateBalancedDesign(count int) []ScenarioFactors {

	var deck []ScenarioFactors

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
		minEntities, _ := strconv.Atoi(os.Getenv("BACKGROUND_ENTITIES_MIN"))
		maxEntities, _ := strconv.Atoi(os.Getenv("BACKGROUND_ENTITIES_MAX"))
		backgroundEntities := CastTridentKit(minEntities, maxEntities)

		factors := ScenarioFactors{
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

func CalculateTridentZones(tridentSpawn TridentSpawn) (fRow, fCol, lRow, lCol, rRow, rCol int) {
	switch tridentSpawn.Orientation {
	case DirectionNorth:
		fRow, fCol = -1, 0
		lRow, lCol = 0, -1
		rRow, rCol = 0, 1
	case DirectionSouth:
		fRow, fCol = 1, 0
		lRow, lCol = 0, 1
		rRow, rCol = 0, -1
	case DirectionEast:
		fRow, fCol = 0, 1
		lRow, lCol = -1, 0
		rRow, rCol = 1, 0
	case DirectionWest:
		fRow, fCol = 0, -1
		lRow, lCol = 1, 0
		rRow, rCol = -1, 0
	}
	return
}
