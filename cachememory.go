package gocache

import (
	"errors"
	"math"
	"time"

	"github.com/mekramy/gocast"
)

type memoryRecord struct {
	data   any
	expiry *time.Time
}

type memoryCache struct {
	data map[string]memoryRecord
}

func (driver *memoryCache) read(key string) (*memoryRecord, bool) {
	val, ok := driver.data[key]
	if !ok {
		return nil, false
	}

	if val.expiry != nil && val.expiry.Before(time.Now()) {
		delete(driver.data, key)
		return nil, false
	}

	return &val, true
}

func (driver *memoryCache) Put(key string, value any, ttl *time.Duration) error {
	var expiry *time.Time = nil
	if ttl != nil {
		exp := time.Now().Add(*ttl)
		expiry = &exp
	}
	driver.data[key] = memoryRecord{
		data:   value,
		expiry: expiry,
	}
	return nil
}

func (driver *memoryCache) Set(key string, value any) (bool, error) {
	record, exists := driver.read(key)
	if !exists {
		return false, nil
	}

	record.data = value
	driver.data[key] = *record
	return true, nil
}

func (driver *memoryCache) Override(key string, value any, ttl *time.Duration) error {
	ok, err := driver.Set(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return driver.Put(key, value, ttl)
	}

	return nil
}

func (driver *memoryCache) Get(key string) (any, error) {
	record, exists := driver.read(key)
	if !exists {
		return nil, nil
	}

	return record.data, nil
}

func (driver *memoryCache) Pull(key string) (any, error) {
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

func (driver *memoryCache) Cast(key string) (gocast.Caster, error) {
	val, err := driver.Get(key)
	if err != nil {
		return gocast.NewCaster(nil), err
	}

	return gocast.NewCaster(val), nil
}

func (driver *memoryCache) Exists(key string) (bool, error) {
	_, exists := driver.read(key)
	return exists, nil
}

func (driver *memoryCache) Forget(key string) error {
	delete(driver.data, key)
	return nil
}

func (driver *memoryCache) TTL(key string) (time.Duration, error) {
	record, exists := driver.read(key)
	if !exists {
		return 0, nil
	}

	if record.expiry == nil {
		return time.Duration(math.MaxInt64), nil
	} else {
		return time.Since(*record.expiry), nil
	}
}

func (driver *memoryCache) Increment(key string, value int64) (bool, error) {
	record, exists := driver.read(key)
	if !exists {
		return false, nil
	}

	caster := gocast.NewCaster(record.data)
	num, err := caster.Int64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	record.data = num + value
	driver.data[key] = *record
	return true, nil
}

func (driver *memoryCache) Decrement(key string, value int64) (bool, error) {
	record, exists := driver.read(key)
	if !exists {
		return false, nil
	}

	caster := gocast.NewCaster(record.data)
	num, err := caster.Int64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	record.data = num - value
	driver.data[key] = *record
	return true, nil
}

func (driver *memoryCache) IncrementFloat(key string, value float64) (bool, error) {
	record, exists := driver.read(key)
	if !exists {
		return false, nil
	}

	caster := gocast.NewCaster(record.data)
	num, err := caster.Float64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	record.data = num + value
	driver.data[key] = *record
	return true, nil
}

func (driver *memoryCache) DecrementFloat(key string, value float64) (bool, error) {
	record, exists := driver.read(key)
	if !exists {
		return false, nil
	}

	caster := gocast.NewCaster(record.data)
	num, err := caster.Float64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	record.data = num - value
	driver.data[key] = *record
	return true, nil
}
