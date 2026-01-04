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
		backgroundEntities := CastBackgroundEntities(minEntities, maxEntities)

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

func CalculateTridentZones(tridentSpawn TridentSpawn) TridentZones {
	distance := 3
	depth := 3

	var fRow, lRow, rRow, fCol, lCol, rCol int

	switch tridentSpawn.Orientation {
	case DirectionNorth: // Forward is Up (-Row)
		fRow, fCol = -1, 0
		lRow, lCol = 0, -1 // Left is West (-Col)
		rRow, rCol = 0, 1  // Right is East (+Col)
	case DirectionSouth: // Forward is Down (+Row)
		fRow, fCol = 1, 0
		lRow, lCol = 0, 1  // Left is East (+Col)
		rRow, rCol = 0, -1 // Right is West (-Col)
	case DirectionEast: // Forward is Right (+Col)
		fRow, fCol = 0, 1
		lRow, lCol = -1, 0 // Left is North (-Row)
		rRow, rCol = 1, 0  // Right is South (+Row)
	case DirectionWest: // Forward is Left (-Col)
		fRow, fCol = 0, -1
		lRow, lCol = 1, 0  // Left is South (+Row)
		rRow, rCol = -1, 0 // Right is North (-Row)
	}

	baseRow := tridentSpawn.Row + (fRow * distance)
	baseCol := tridentSpawn.Col + (fCol * distance)

	// Generate strip of coordinates (surface/orientation filled later)
	generateStrip := func(startRow, startCol int) TridentZone {
		coords := make([]EnrichedCoordinate, depth)
		for i := 0; i < depth; i++ {
			coords[i] = EnrichedCoordinate{
				Coordinate: Coordinate{
					Row: startRow + (fRow * i),
					Col: startCol + (fCol * i),
				},
				// Surface and Orientation will be filled by EnrichTridentZones
			}
		}
		return TridentZone{Coordinates: coords}
	}

	return TridentZones{
		ZoneA: generateStrip(baseRow, baseCol),
		ZoneB: generateStrip(baseRow+lRow, baseCol+lCol),
		ZoneC: generateStrip(baseRow+rRow, baseCol+rCol),
	}
}
