package gocache

import (
	"time"

	"github.com/mekramy/gocast"
	"github.com/redis/go-redis/v9"
)

// Cache nil safe caching.
type Cache interface {
	// Put stores a value in the cache with the specified key and TTL (time-to-live).
	// If ttl is nil, the value is stored indefinitely.
	// Returns an error if the operation fails.
	Put(key string, value any, ttl *time.Duration) error

	// Set set value of existing key in the cache.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	Set(key string, value any) (bool, error)

	// Override stores a value in the cache with the specified key and TTL (time-to-live).
	// If key exists ttl keep.
	// If ttl is nil, the value is stored indefinitely.
	// Returns an error if the operation fails.
	Override(key string, value any, ttl *time.Duration) error

	// Get retrieves a value from the cache by its key.
	// Returns the value and an error if the operation fails.
	Get(key string) (any, error)

	// Pull retrieves a value from the cache by its key and removes it from the cache.
	// Returns the value and an error if the operation fails.
	Pull(key string) (any, error)

	// Cast retrieves a value from the cache by its key and casts it to a gocast.Caster.
	// Returns the casted value and an error if the operation fails.
	Cast(key string) (gocast.Caster, error)

	// Exists checks if a key exists in the cache.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	Exists(key string) (bool, error)

	// Forget removes a value from the cache by its key.
	// Returns an error if the operation fails.
	Forget(key string) error

	// TTL retrieves the time-to-live (TTL) of a value in the cache by its key.
	// Returns the TTL and an error if the operation fails.
	TTL(key string) (time.Duration, error)

	// Increment increases the value of a key by the specified integer value.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	Increment(key string, value int64) (bool, error)

	// Decrement decreases the value of a key by the specified integer value.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	Decrement(key string, value int64) (bool, error)

	// IncrementFloat increases the value of a key by the specified float value.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	IncrementFloat(key string, value float64) (bool, error)

	// DecrementFloat decreases the value of a key by the specified float value.
	// Returns a boolean indicating whether the key exists and an error if the operation fails.
	DecrementFloat(key string, value float64) (bool, error)
}

// NewRedisCache creates a new redis cache instance.
func NewRedisCache(prefix string, client *redis.Client) Cache {
	driver := new(redisCache)
	driver.prefix = prefix
	driver.client = client
	return driver
}

// NewRedisCache creates a new in-memory cache instance.
func NewMemoryCache() Cache {
	driver := new(memoryCache)
	driver.data = make(map[string]memoryRecord)
	return driver
}
