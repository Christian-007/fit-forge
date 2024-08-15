package cache

import "time"

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
	Delete(key string) error
}
