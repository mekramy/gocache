package gocache

import "time"

type verification struct {
	name  string
	ttl   time.Duration
	cache Cache
}

func (driver verification) Set(code string) error {
	exists, err := driver.cache.Set(driver.name, code)
	if err != nil {
		return err
	}

	if !exists {
		return driver.cache.Put(driver.name, code, &driver.ttl)
	}

	return nil
}

func (driver verification) Generate(count uint) (string, error) {
	code, err := randomString(count, "0123456789")
	if err != nil {
		return "", err
	}

	err = driver.Set(code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (driver verification) Clear() error {
	return driver.cache.Forget(driver.name)
}

func (driver verification) Get() (string, error) {
	caster, err := driver.cache.Cast(driver.name)
	if err != nil {
		return "", err
	}

	if caster.IsNil() {
		return "", nil
	}

	code, err := caster.String()
	if err != nil {
		return "", err
	}

	return code, nil
}

func (driver verification) Exists() (bool, error) {
	return driver.cache.Exists(driver.name)
}

func (driver verification) TTL() (time.Duration, error) {
	ttl, err := driver.cache.TTL(driver.name)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
