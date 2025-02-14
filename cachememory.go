package gocache

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/mekramy/gocast"
)

type memoryRecord struct {
	data   any
	expiry *time.Time
}

type memoryCache struct {
	data  map[string]memoryRecord
	mutex sync.RWMutex
}

func (m *memoryCache) read(key string) (*memoryRecord, bool) {
	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Read key
	val, ok := m.data[key]
	if !ok {
		return nil, false
	}

	// Delete key if expired
	if val.expiry != nil && val.expiry.Before(time.Now()) {
		delete(m.data, key)
		return nil, false
	}

	return &val, true
}

func (m *memoryCache) Put(key string, value any, ttl *time.Duration) error {
	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Calculate expiry
	var expiry *time.Time = nil
	if ttl != nil {
		exp := time.Now().Add(*ttl)
		expiry = &exp
	}

	// Store data
	m.data[key] = memoryRecord{
		data:   value,
		expiry: expiry,
	}
	return nil
}

func (m *memoryCache) Set(key string, value any) (bool, error) {
	// Check existence
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store data
	record.data = value
	m.data[key] = *record
	return true, nil
}

func (m *memoryCache) Override(key string, value any, ttl *time.Duration) error {
	ok, err := m.Set(key, value)
	if err != nil {
		return err
	}

	if !ok {
		return m.Put(key, value, ttl)
	}

	return nil
}

func (m *memoryCache) Get(key string) (any, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return nil, nil
	}

	return record.data, nil
}

func (m *memoryCache) Pull(key string) (any, error) {
	// Read
	val, err := m.Get(key)
	if err != nil {
		return nil, err
	}

	// Delete
	err = m.Forget(key)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (m *memoryCache) Cast(key string) (gocast.Caster, error) {
	// Read
	val, err := m.Get(key)
	return gocast.NewCaster(val), err
}

func (m *memoryCache) Exists(key string) (bool, error) {
	_, exists := m.read(key)
	return exists, nil
}

func (m *memoryCache) Forget(key string) error {
	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Delete
	delete(m.data, key)
	return nil
}

func (m *memoryCache) TTL(key string) (time.Duration, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return 0, nil
	}

	// Calculate ttl
	if record.expiry == nil {
		return time.Duration(math.MaxInt64), nil
	} else {
		return time.Since(*record.expiry), nil
	}
}

func (m *memoryCache) Increment(key string, value int64) (bool, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	// Cast
	caster := gocast.NewCaster(record.data)
	num, err := caster.Int64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store
	record.data = num + value
	m.data[key] = *record
	return true, nil
}

func (m *memoryCache) Decrement(key string, value int64) (bool, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	// Cast
	caster := gocast.NewCaster(record.data)
	num, err := caster.Int64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store
	record.data = num - value
	m.data[key] = *record
	return true, nil
}

func (m *memoryCache) IncrementFloat(key string, value float64) (bool, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	// Cast
	caster := gocast.NewCaster(record.data)
	num, err := caster.Float64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store
	record.data = num + value
	m.data[key] = *record
	return true, nil
}

func (m *memoryCache) DecrementFloat(key string, value float64) (bool, error) {
	// Read
	record, exists := m.read(key)
	if !exists {
		return false, nil
	}

	// Cast
	caster := gocast.NewCaster(record.data)
	num, err := caster.Float64()
	if err != nil {
		return false, errors.New("value is not numeric")
	}

	// Safe race condition
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store
	record.data = num - value
	m.data[key] = *record
	return true, nil
}
