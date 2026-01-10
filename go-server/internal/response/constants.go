package response

import (
	"os"
	"strconv"
)

var (
	// Timer duration in milliseconds (default: 12 seconds)
	TimerDurationMs int64 = 12000

	// Network buffer time in milliseconds (default: 10 seconds)
	NetworkBufferMs int64 = 10000
)

func init() {
	if val := os.Getenv("TIMER_DURATION_MS"); val != "" {
		if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
			TimerDurationMs = parsed
		}
	}

	if val := os.Getenv("NETWORK_BUFFER_MS"); val != "" {
		if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
			NetworkBufferMs = parsed
		}
	}
}
