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

func (driver redisQueue) Push(value any) error {
	return driver.client.LPush(
		context.Background(),
		driver.name,
		value,
	).Err()
}

func (driver redisQueue) Pull() (any, error) {
	val, err := driver.client.LPop(
		context.Background(),
		driver.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (driver redisQueue) Pop() (any, error) {
	val, err := driver.client.RPop(
		context.Background(),
		driver.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (driver redisQueue) Length() (int64, error) {
	val, err := driver.client.LLen(
		context.Background(),
		driver.name,
	).Result()

	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return val, nil
}

func (driver redisQueue) Cast() (gocast.Caster, error) {
	val, err := driver.Pull()
	return gocast.NewCaster(val), err
}
