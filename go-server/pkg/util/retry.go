package util

import (
	"context"
	"time"
)

// Retry attempts to execute fn up to the specified number of attempts
func Retry(ctx context.Context, attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {

		if err = fn(); err == nil {
			return nil
		}

		// skip waiting after the last attempt
		if i < attempts-1 {
			select {
			case <-ctx.Done():
				// context cancelled, bail out early
				return ctx.Err()
			case <-time.After(delay):
				// wait finished, try again
			}
		}
	}
	return err
}
