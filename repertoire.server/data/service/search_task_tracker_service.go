package service

import (
	goCache "github.com/patrickmn/go-cache"
	"repertoire/server/data/cache"
)

type SearchTaskTrackerService interface {
	Track(taskID string, userID string)
	GetUserID(taskID string) (string, bool)
}

type searchTaskTrackerService struct {
	cache cache.MeiliCache
}

func NewSearchTaskTrackerService(cache cache.MeiliCache) SearchTaskTrackerService {
	return searchTaskTrackerService{cache: cache}
}

func (d searchTaskTrackerService) Track(taskID string, userID string) {
	d.cache.Set("task-"+taskID, userID, goCache.DefaultExpiration)
}

func (d searchTaskTrackerService) GetUserID(taskID string) (string, bool) {
	userID, found := d.cache.Get("task-" + taskID)
	if found {
		return userID.(string), found
	}
	return "", found
}
