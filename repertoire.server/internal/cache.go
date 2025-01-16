package internal

import (
	"github.com/patrickmn/go-cache"
	"time"
)

func NewCache() *cache.Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}
