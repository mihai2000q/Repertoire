package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type StorageCache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

func NewStorageCache() StorageCache {
	return cache.New(5*time.Minute, 10*time.Minute)
}
