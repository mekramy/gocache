package gocache

import (
	"github.com/mekramy/gocast"
	"github.com/redis/go-redis/v9"
)

// Queue indefinitely nil safe queue.
type Queue interface {
	// Push stores value in queue.
	// Returns an error if the operation fails.
	Push(value any) error

	// Pull retrieves first queue item and removes it from the queue.
	// Returns the value and an error if the operation fails.
	Pull() (any, error)

	// Pop retrieves last queue item and removes it from the queue.
	// Returns the value and an error if the operation fails.
	Pop() (any, error)

	// Cast retrieves first queue item and casts it to a gocast.Caste and removes it from the queue.
	// Returns the Caster instance and an error if the operation fails.
	Cast() (gocast.Caster, error)

	// Length get the length of queue.
	// Returns the queue length and an error if the operation fails.
	Length() (int64, error)
}

// NewRedisQueue creates a new redis queue instance.
func NewRedisQueue(name string, client *redis.Client) Queue {
	driver := new(redisQueue)
	driver.name = name
	driver.client = client
	return driver
}
