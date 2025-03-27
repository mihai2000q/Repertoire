package service

import (
	"repertoire/server/data/cache"
	"time"
)

type MeiliTaskTrackerService interface {
	Track(taskID string, userID string)
	GetUserID(taskID string) (string, bool)
}

type meiliTaskTrackerService struct {
	cache cache.Cache
}

func NewMeiliTaskTrackerService(cache cache.Cache) MeiliTaskTrackerService {
	return meiliTaskTrackerService{cache: cache}
}

func (d meiliTaskTrackerService) Track(taskID string, userID string) {
	d.cache.Set("meiliTask-"+taskID, userID, time.Minute)
}

func (d meiliTaskTrackerService) GetUserID(taskID string) (string, bool) {
	userID, found := d.cache.Get("meiliTask-" + taskID)
	if found {
		return userID.(string), found
	}
	return "", found
}
