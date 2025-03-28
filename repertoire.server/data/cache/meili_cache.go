package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Cache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

func NewCache() Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}
