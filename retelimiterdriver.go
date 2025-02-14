package gocache

import (
	"time"
)

type limiter struct {
	name        string
	maxAttempts uint32
	ttl         time.Duration
	cache       Cache
}

func (l limiter) Hit() error {
	exists, err := l.cache.Decrement(l.name, 1)
	if err != nil {
		return err
	}

	if !exists {
		return l.cache.Put(l.name, l.maxAttempts-1, &l.ttl)
	}

	return nil
}

func (l limiter) Lock() error {
	exists, err := l.cache.Set(l.name, 0)
	if err != nil {
		return err
	}

	if !exists {
		return l.cache.Put(l.name, 0, &l.ttl)
	}

	return nil
}

func (l limiter) Reset() error {
	return l.cache.Put(l.name, l.maxAttempts, &l.ttl)
}

func (l limiter) Clear() error {
	return l.cache.Forget(l.name)
}

func (l limiter) MustLock() (bool, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return true, err
	}

	if caster.IsNil() {
		return false, nil
	}

	num, err := caster.Int()
	if err != nil {
		return true, err
	}

	return num <= 0, nil
}

func (l limiter) TotalAttempts() (uint32, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return 0, err
	}

	if caster.IsNil() {
		return 0, nil
	}

	num, err := caster.Int()
	if err != nil {
		return 0, err
	}

	if num < int(l.maxAttempts) {
		num = int(l.maxAttempts)
	}

	return l.maxAttempts - uint32(num), nil
}

func (l limiter) RetriesLeft() (uint32, error) {
	caster, err := l.cache.Cast(l.name)
	if err != nil {
		return 0, err
	}

	if caster.IsNil() {
		return 0, nil
	}

	num, err := caster.Int()
	if err != nil {
		return 0, err
	}

	if num < 0 {
		num = 0
	}

	return uint32(num), nil

}

func (l limiter) AvailableIn() (time.Duration, error) {
	ttl, err := l.cache.TTL(l.name)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
