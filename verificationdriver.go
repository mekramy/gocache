package gocache

import "time"

type verification struct {
	name  string
	ttl   time.Duration
	cache Cache
}

func (v verification) Set(code string) error {
	exists, err := v.cache.Set(v.name, code)
	if err != nil {
		return err
	}

	if !exists {
		return v.cache.Put(v.name, code, &v.ttl)
	}

	return nil
}

func (v verification) Generate(count uint) (string, error) {
	code, err := randomString(count, "0123456789")
	if err != nil {
		return "", err
	}

	err = v.Set(code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (v verification) Clear() error {
	return v.cache.Forget(v.name)
}

func (v verification) Get() (string, error) {
	caster, err := v.cache.Cast(v.name)
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

func (v verification) Exists() (bool, error) {
	return v.cache.Exists(v.name)
}

func (v verification) TTL() (time.Duration, error) {
	ttl, err := v.cache.TTL(v.name)
	if err != nil {
		return 0, err
	}

	return ttl, nil
}
