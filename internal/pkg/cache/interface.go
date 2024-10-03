package cache

import "time"

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
	Delete(key string) error
	GetAllHashFields(key string) (map[string]string, error)
	SetHash(key string, values ...interface{}) error
	SetExpire(key string, expiration time.Duration) error
}
