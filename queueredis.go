package gocache

import (
	"context"
	"errors"

	"github.com/mekramy/gocast"
	"github.com/redis/go-redis/v9"
)

type redisQueue struct {
	name   string
	client *redis.Client
}

func (r redisQueue) Push(value any) error {
	return r.client.LPush(
		context.Background(),
		r.name,
		value,
	).Err()
}

func (r redisQueue) Pull() (any, error) {
	val, err := r.client.LPop(
		context.Background(),
		r.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r redisQueue) Pop() (any, error) {
	val, err := r.client.RPop(
		context.Background(),
		r.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (r redisQueue) Length() (int64, error) {
	val, err := r.client.LLen(
		context.Background(),
		r.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return val, nil
}

func (r redisQueue) Cast() (gocast.Caster, error) {
	val, err := r.Pull()
	return gocast.NewCaster(val), err
}
