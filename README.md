# GoCache Library Documentation

## Overview

The `gocache` library provides a set of interfaces and implementations for caching, queuing, rate limiting, and verification code generation. It supports both Redis-based and in-memory caching.

## Installation

To install GoCache, use the following command:

```sh
go get github.com/mekramy/gocache
```

## Cache

Key Features:

- Supports Redis and in-memory caching.
- Provides functionality to store, retrieve, update, and manage cache keys.
- Offers type casting for cached values.
- Includes advanced operations like incrementing, decrementing, and TTL management.

Methods:

- `Put(key string, value any, ttl *time.Duration) error`: Store a value in the cache with a time-to-live (TTL). If nil ttl passed, the value is stored indefinitely.
- `Set(key string, value any) (bool, error)`: Update the value of an existing key, returns `false` if key not exists.
- `Override(key string, value any, ttl *time.Duration) error`: Store or update a value in the cache, retaining TTL if the key exists.
- `Get(key string) (any, error)`: Retrieve a value from the cache.
- `Pull(key string) (any, error)`: Retrieve and remove a value from the cache.
- `Cast(key string) (gocast.Caster, error)`: Retrieve and cast a value to a `gocast.Caster`.
- `Exists(key string) (bool, error)`: Check if a key exists.
- `Forget(key string) error`: Remove a key from the cache.
- `TTL(key string) (time.Duration, error)`: Get the time-to-live of a key.
- `Increment(key string, value int64) (bool, error)`: Increment a key's value by a given integer, returns `false` if key not exists.
- `Decrement(key string, value int64) (bool, error)`: Decrement a key's value by a given integer, returns `false` if key not exists.
- `IncrementFloat(key string, value float64) (bool, error)`: Increment a key's value by a given float, returns `false` if key not exists.
- `DecrementFloat(key string, value float64) (bool, error)`: Decrement a key's value by a given float, returns `false` if key not exists.

Create a Redis-based cache instance:

```go
func NewRedisCache(prefix string, client *redis.Client) Cache
```

Create an in-memory cache instance:

```go
func NewMemoryCache() Cache
```

## Queue

Key Features:

- Implements a safe, nil-resistant queue interface.
- Compatible with Redis for distributed queue management.
- Includes operations to push, pull, and retrieve items with type casting support.

Methods:

- `Push(value any) error`: Add a value to the queue.
- `Pull() (any, error)`: Retrieve and remove the first item in the queue.
- `Pop() (any, error)`: Retrieve and remove the last item in the queue.
- `Cast() (gocast.Caster, error)`: Retrieve and cast the first item in the queue to a `gocast.Caster`.
- `Length() (int64, error)`: Get the queue length.

Create a Redis-based queue instance:

```go
func NewRedisQueue(name string, client *redis.Client) Queue
```

## RateLimiter

Key Features:

- Limit Actions: Tracks and restricts the number of attempts allowed for a specific action, ensuring controlled usage.
- Locking Support: Automatically locks when the limit is exceeded and provides information about when the lock will be lifted.
- Easy Reset: Quickly reset the limiter to allow new attempts or clear its state entirely.
- Retry Tracking: Keeps track of how many retries are left and provides real-time updates.

Methods:

- `Hit() error`: Decrement the remaining attempts.
- `Lock() error`: Lock the rate limiter.
- `Reset() error`: Reset the rate limiter.
- `Clear() error`: Remove the rate limiter from the cache.
- `MustLock() (bool, error)`: Check if the rate limiter is locked.
- `TotalAttempts() (uint32, error)`: Get the total number of attempts allowed.
- `RetriesLeft() (uint32, error)`: Get the remaining attempts.
- `AvailableIn() (time.Duration, error)`: Get the time until the rate limiter unlocks.

Create a rate limiter instance:

```go
func NewRateLimiter(name string, maxAttempts uint32, ttl time.Duration, cache Cache) RateLimiter
```

## VerificationCode

Key Features:

- Handles verification code generation and management.
- Supports TTL and retry attempt tracking.
- Allows easy integration with cache for storing verification codes.

Methods:

- `Set(code string) error`: Set a verification code.
- `Generate(count uint) (string, error)`: Generate a random numeric code of a specified length.
- `Clear() error`: Remove a verification code from the cache.
- `Get() (string, error)`: Retrieve a verification code.
- `Exists() (bool, error)`: Check if a verification code exists.
- `TTL() (time.Duration, error)`: Get the time-to-live of a verification code.

Create a verification code instance:

```go
func NewVerification(name string, maxAttempts uint32, ttl time.Duration, cache Cache) VerificationCode
```

## Dependencies

- `github.com/mekramy/gocast`: Casting library for dynamic type handling.
- `github.com/redis/go-redis/v9`: Redis client library.
