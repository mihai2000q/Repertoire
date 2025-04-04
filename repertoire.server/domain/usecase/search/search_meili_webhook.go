package search

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"strconv"
	"strings"
)

type MeiliWebhook struct {
	searchEngineService      service.SearchEngineService
	searchTaskTrackerService service.SearchTaskTrackerService
	realTimeService          service.RealTimeService
}

func NewMeiliWebhook(
	searchEngineService service.SearchEngineService,
	searchTaskTrackerService service.SearchTaskTrackerService,
	realTimeService service.RealTimeService,
) MeiliWebhook {
	return MeiliWebhook{
		searchEngineService:      searchEngineService,
		searchTaskTrackerService: searchTaskTrackerService,
		realTimeService:          realTimeService,
	}
}

func (m MeiliWebhook) Handle(requestBody io.ReadCloser) *wrapper.ErrorCode {
	gz, err := gzip.NewReader(requestBody)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	defer func(gz *gzip.Reader) {
		_ = gz.Close()
	}(gz)
	body, err := io.ReadAll(gz)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	tasks := strings.Split(string(body), "\n")

	var task struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}
	if err = json.Unmarshal([]byte(tasks[0]), &task); err != nil {
		return wrapper.InternalServerError(err)
	}

	taskID := strconv.FormatInt(task.UID, 10)

	if !m.searchEngineService.HasTaskSucceeded(task.Status) {
		return wrapper.InternalServerError(
			errors.New("meilisearch task, " + taskID + ", failed"),
		)
	}

	userID, isTaskTracked := m.searchTaskTrackerService.GetUserID(taskID)
	if !isTaskTracked {
		return nil
	}

	err = m.realTimeService.Publish("search", userID, map[string]any{"action": "SEARCH_CACHE_INVALIDATION"})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
