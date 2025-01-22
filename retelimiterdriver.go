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

func (driver limiter) Hit() error {
	exists, err := driver.cache.Decrement(driver.name, 1)
	if err != nil {
		return err
	}

	if !exists {
		return driver.cache.Put(driver.name, driver.maxAttempts-1, &driver.ttl)
	}

	return nil
}

func (driver limiter) Lock() error {
	exists, err := driver.cache.Set(driver.name, 0)
	if err != nil {
		return err
	}

	if !exists {
		return driver.cache.Put(driver.name, 0, &driver.ttl)
	}

	return nil
}

func (driver limiter) Reset() error {
	return driver.cache.Put(driver.name, driver.maxAttempts, &driver.ttl)
}

func (driver limiter) Clear() error {
	return driver.cache.Forget(driver.name)
}

func (driver limiter) MustLock() (bool, error) {
	caster, err := driver.cache.Cast(driver.name)
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

func (driver limiter) TotalAttempts() (uint32, error) {
	caster, err := driver.cache.Cast(driver.name)
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

	if num < int(driver.maxAttempts) {
		num = int(driver.maxAttempts)
	}

	return driver.maxAttempts - uint32(num), nil
}

func (driver limiter) RetriesLeft() (uint32, error) {
	caster, err := driver.cache.Cast(driver.name)
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

func (driver limiter) AvailableIn() (time.Duration, error) {
	ttl, err := driver.cache.TTL(driver.name)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
