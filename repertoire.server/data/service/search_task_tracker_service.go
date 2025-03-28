package service

import (
	"repertoire/server/data/cache"
	"time"
)

type SearchTaskTrackerService interface {
	Track(taskID string, userID string)
	GetUserID(taskID string) (string, bool)
}

type searchTaskTrackerService struct {
	cache cache.Cache
}

func NewSearchTaskTrackerService(cache cache.Cache) SearchTaskTrackerService {
	return searchTaskTrackerService{cache: cache}
}

func (d searchTaskTrackerService) Track(taskID string, userID string) {
	d.cache.Set("meiliTask-"+taskID, userID, time.Minute)
}

func (d searchTaskTrackerService) GetUserID(taskID string) (string, bool) {
	userID, found := d.cache.Get("meiliTask-" + taskID)
	if found {
		return userID.(string), found
	}
	return "", found
}
