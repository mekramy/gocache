package gocache

import (
	"context"
	"errors"
	"time"

	"github.com/mekramy/gocast"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	prefix string
	client *redis.Client
}

func (driver *redisCache) prefixer(key string) string {
	return cacheKey(driver.prefix, key)
}

func (driver *redisCache) Put(key string, value any, ttl *time.Duration) error {
	return driver.client.Set(
		context.Background(),
		driver.prefixer(key),
		value,
		safeValue(ttl, 0),
	).Err()
}

func (driver *redisCache) Set(key string, value any) (bool, error) {
	exists, err := driver.Exists(key)
	if err != nil || !exists {
		return false, err
	}

	err = driver.client.Set(
		context.Background(),
		driver.prefixer(key),
		value,
		redis.KeepTTL,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (driver *redisCache) Override(key string, value any, ttl *time.Duration) error {
	ok, err := driver.Set(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return driver.Put(key, value, ttl)
	}

	return nil
}

func (driver *redisCache) Get(key string) (any, error) {
	val, err := driver.client.Get(
		context.TODO(),
		driver.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (driver *redisCache) Pull(key string) (any, error) {
	val, err := driver.Get(key)
	if err != nil {
		return nil, err
	}

	err = driver.Forget(key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (driver *redisCache) Cast(key string) (gocast.Caster, error) {
	val, err := driver.Get(key)
	if err != nil {
		return nil, err
	}

	return gocast.NewCaster(val), nil
}

func (driver *redisCache) Exists(key string) (bool, error) {
	exists, err := driver.client.Exists(
		context.TODO(),
		driver.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (driver *redisCache) Forget(key string) error {
	err := driver.client.Del(
		context.TODO(),
		driver.prefixer(key),
	).Err()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (driver *redisCache) TTL(key string) (time.Duration, error) {
	ttl, err := driver.client.TTL(
		context.Background(),
		driver.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (driver *redisCache) Increment(key string, value int64) (bool, error) {
	exists, err := driver.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = driver.client.IncrBy(
		context.Background(),
		driver.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (driver *redisCache) Decrement(key string, value int64) (bool, error) {
	exists, err := driver.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = driver.client.DecrBy(
		context.Background(),
		driver.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (driver *redisCache) IncrementFloat(key string, value float64) (bool, error) {
	exists, err := driver.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = driver.client.IncrByFloat(
		context.Background(),
		driver.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (driver *redisCache) DecrementFloat(key string, value float64) (bool, error) {
	exists, err := driver.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = driver.client.IncrByFloat(
		context.Background(),
		driver.prefixer(key),
		-value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
