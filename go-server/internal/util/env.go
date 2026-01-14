package util

import (
	"os"
	"strings"
)

func GetEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// Flexibly Find all environment variables with a given prefix (GROQ_API_KEY_{index})
func CollectEnvKeys(prefix string) []string {
	var keys []string

	for _, env := range os.Environ() {
		// env is in format "KEY=value"
		if strings.HasPrefix(env, prefix) {
			// Split by = and take the value part
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 && parts[1] != "" {
				keys = append(keys, parts[1])
			}
		}
	}

	return keys
}
