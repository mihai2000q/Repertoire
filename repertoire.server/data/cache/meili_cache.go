package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type MeiliCache interface {
	Get(k string) (interface{}, bool)
	Set(k string, x interface{}, d time.Duration)
}

func NewMeiliCache() MeiliCache {
	return cache.New(1*time.Minute, 10*time.Minute)
}
