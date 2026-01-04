package domain

import "math/rand"

// Select One Unique Primary Entity (Star)
func CastPrimaryEntity() string {
	return StarPool[rand.Intn(len(StarPool))]
}

func CastTridentKit(minNoise, maxNoise int) []string {
	var kit []string

	kit = append(kit, VehiclePool[rand.Intn(len(VehiclePool))])
	kit = append(kit, PedestrianPool[rand.Intn(len(PedestrianPool))])
	// Make min/max inclusive
	rangeSize := maxNoise - minNoise
	count := minNoise
	if rangeSize > 0 {
		count += rand.Intn(rangeSize + 1)
	}
	noiseOptions := []string{}
	noiseOptions = append(noiseOptions, VehiclePool...)
	noiseOptions = append(noiseOptions, PedestrianPool...)
	noiseOptions = append(noiseOptions, ObstaclePool...)

	for i := 0; i < count; i++ {
		kit = append(kit, noiseOptions[rand.Intn(len(noiseOptions))])
	}
	// Shuffle so the LLM doesn't always see [Car, Ped, ...] in that order
	rand.Shuffle(len(kit), func(i, j int) {
		kit[i], kit[j] = kit[j], kit[i]
	})

	return kit
}
