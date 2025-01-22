package gocache

import (
	"time"
)

// RateLimiter rate limiter.
type RateLimiter interface {
	// Hit decrement user left tries.
	// Returns an error if the operation fails.
	Hit() error

	// Lock lock rate limiter.
	// Returns an error if the operation fails.
	Lock() error

	// Reset reset rate limiter.
	// Returns an error if the operation fails.
	Reset() error

	// Clear remove rate limiter from cache.
	// Returns an error if the operation fails.
	Clear() error

	// MustLock check if rate limiter must lock.
	// Returns a boolean indicating lock state and an error if the operation fails.
	MustLock() (bool, error)

	// TotalAttempts get user attempts count.
	// Returns the total attempts and an error if the operation fails.
	TotalAttempts() (uint32, error)

	// RetriesLeft get user retries left.
	// Returns the remaining attempts and an error if the operation fails.
	RetriesLeft() (uint32, error)

	// AvailableIn get time until unlock
	// Returns the time until unlock and an error if the operation fails.
	AvailableIn() (time.Duration, error)
}

// NewRateLimiter creates a new rate limiter instance.
func NewRateLimiter(name string, maxAttempts uint32, ttl time.Duration, cache Cache) RateLimiter {
	driver := new(limiter)
	driver.name = "limiter " + name
	driver.maxAttempts = maxAttempts
	driver.ttl = ttl
	driver.cache = cache
	return driver
}
