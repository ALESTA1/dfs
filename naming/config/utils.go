package config

import (
	"math/rand"
	"time"
)

func GetRandomKey(m map[string]int) string {
	rand.NewSource(time.Now().UnixNano())
	keys := make([]string, 0, len(m))

	
	for key := range m {
		keys = append(keys, key)
	}

	// Get a random index
	randomIndex := rand.Intn(len(keys))

	// Return the random key
	return keys[randomIndex]
}
