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

func (r *redisCache) prefixer(key string) string {
	return cacheKey(r.prefix, key)
}

func (r *redisCache) Put(key string, value any, ttl *time.Duration) error {
	return r.client.Set(
		context.Background(),
		r.prefixer(key),
		value,
		safeValue(ttl, 0),
	).Err()
}

func (r *redisCache) Set(key string, value any) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return false, err
	}

	err = r.client.Set(
		context.Background(),
		r.prefixer(key),
		value,
		redis.KeepTTL,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisCache) Override(key string, value any, ttl *time.Duration) error {
	ok, err := r.Set(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return r.Put(key, value, ttl)
	}

	return nil
}

func (r *redisCache) Get(key string) (any, error) {
	val, err := r.client.Get(
		context.TODO(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisCache) Pull(key string) (any, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	err = r.Forget(key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r *redisCache) Cast(key string) (gocast.Caster, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	return gocast.NewCaster(val), nil
}

func (r *redisCache) Exists(key string) (bool, error) {
	exists, err := r.client.Exists(
		context.TODO(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (r *redisCache) Forget(key string) error {
	err := r.client.Del(
		context.TODO(),
		r.prefixer(key),
	).Err()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	return err
}

func (r *redisCache) TTL(key string) (time.Duration, error) {
	ttl, err := r.client.TTL(
		context.Background(),
		r.prefixer(key),
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (r *redisCache) Increment(key string, value int64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrBy(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisCache) Decrement(key string, value int64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.DecrBy(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisCache) IncrementFloat(key string, value float64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrByFloat(
		context.Background(),
		r.prefixer(key),
		value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisCache) DecrementFloat(key string, value float64) (bool, error) {
	exists, err := r.Exists(key)
	if err != nil || !exists {
		return exists, err
	}

	err = r.client.IncrByFloat(
		context.Background(),
		r.prefixer(key),
		-value,
	).Err()
	if err != nil {
		return false, err
	}

	return true, nil
}
