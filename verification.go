package gocache

import "time"

// VerificationCode verification code.
type VerificationCode interface {
	// Set set code.
	// Returns a boolean indicating verification exists in cache and an error if the operation fails.
	Set(code string) error

	// Generate generate a random numeric code with n character length.
	// Returns generated code, a boolean indicating verification exists in cache and an error if the operation fails.
	Generate(count uint) (string, error)

	// Clear remove verification from cache.
	// Returns an error if the operation fails.
	Clear() error

	// Get get code.
	// Returns verification code and an error if the operation fails.
	Get() (string, error)

	// Exists check if code exists.
	// Returns a boolean indicating whether the verification exists and an error if the operation fails.
	Exists() (bool, error)

	// TTL retrieves the time-to-live (TTL) of verification in the cache.
	// Returns the TTL and an error if the operation fails.
	TTL() (time.Duration, error)
}

// NewVerification creates a new verification instance.
func NewVerification(name string, maxAttempts uint32, ttl time.Duration, cache Cache) VerificationCode {
	driver := new(verification)
	driver.name = "verify " + name
	driver.ttl = ttl
	driver.cache = cache
	return driver
}
