package domain

import "math/rand"

// Select One Unique Primary Entity (Star)
func CastPrimaryEntity() string {
	return StarPool[rand.Intn(len(StarPool))]
}

// Select entities randomly for the background noise
func CastBackgroundEntities(min, max int) []string {
	// Determine the total count
	count := rand.Intn(max-min) + min

	// Pick random entities from Background Pool
	var entities []string
	for i := 0; i < count; i++ {
		entities = append(entities, BackgroundPool[rand.Intn(len(BackgroundPool))])
	}

	return entities
}
