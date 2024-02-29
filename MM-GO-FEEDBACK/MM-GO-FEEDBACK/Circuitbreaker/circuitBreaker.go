package config

import (
	"time"

	"github.com/sony/gobreaker"
)

func CircuitBreakerConfig() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:    "Feedbacks",
		Timeout: 5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.TotalFailures >= 3
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}
