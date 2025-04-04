package search

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"repertoire/server/domain/usecase/search"
	"repertoire/server/test/unit/data/service"
	"strconv"
	"testing"
)

func TestSearchMeiliWebhook_WhenTaskHasNotSucceeded_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewMeiliWebhook(searchEngineService, nil, nil)

	taskID := 23
	var task = struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}{
		UID:    int64(taskID),
		Status: "success",
	}

	searchEngineService.On("HasTaskSucceeded", task.Status).Return(false).Once()

	// when
	errCode := _uut.Handle(writeRequestBody(task))

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Contains(t, errCode.Error.Error(), "meilisearch task")
	assert.Contains(t, errCode.Error.Error(), "failed")

	searchEngineService.AssertExpectations(t)
}

func TestSearchMeiliWebhook_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	searchTaskTrackerService := new(service.SearchTaskTrackerServiceMock)
	realTimeService := new(service.RealTimeServiceMock)
	_uut := search.NewMeiliWebhook(searchEngineService, searchTaskTrackerService, realTimeService)

	taskID := 23
	var task = struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}{
		UID:    int64(taskID),
		Status: "success",
	}

	searchEngineService.On("HasTaskSucceeded", task.Status).Return(true).Once()

	userID := "some-user-id"
	searchTaskTrackerService.On("GetUserID", strconv.Itoa(taskID)).Return(userID, true).Once()

	internalError := errors.New("internal error")
	realTimeService.
		On("Publish", "search", userID, map[string]any{"action": "SEARCH_CACHE_INVALIDATION"}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(writeRequestBody(task))

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	searchEngineService.AssertExpectations(t)
	searchTaskTrackerService.AssertExpectations(t)
	realTimeService.AssertExpectations(t)
}

func TestSearchMeiliWebhook_WhenTaskIsNotTracked_ShouldNotReturnAnyError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	searchTaskTrackerService := new(service.SearchTaskTrackerServiceMock)
	_uut := search.NewMeiliWebhook(searchEngineService, searchTaskTrackerService, nil)

	taskID := 23
	var task = struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}{
		UID:    int64(taskID),
		Status: "success",
	}

	searchEngineService.On("HasTaskSucceeded", task.Status).Return(true).Once()
	searchTaskTrackerService.On("GetUserID", strconv.Itoa(taskID)).Return("", false).Once()

	// when
	errCode := _uut.Handle(writeRequestBody(task))

	// then
	assert.Nil(t, errCode)

	searchEngineService.AssertExpectations(t)
	searchTaskTrackerService.AssertExpectations(t)
}

func TestSearchMeiliWebhook_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	searchTaskTrackerService := new(service.SearchTaskTrackerServiceMock)
	realTimeService := new(service.RealTimeServiceMock)
	_uut := search.NewMeiliWebhook(searchEngineService, searchTaskTrackerService, realTimeService)

	taskID := 23
	var task = struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}{
		UID:    int64(taskID),
		Status: "success",
	}

	searchEngineService.On("HasTaskSucceeded", task.Status).Return(true).Once()

	userID := "some-user-id"
	searchTaskTrackerService.On("GetUserID", strconv.Itoa(taskID)).Return(userID, true).Once()

	realTimeService.
		On("Publish", "search", userID, map[string]any{"action": "SEARCH_CACHE_INVALIDATION"}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(writeRequestBody(task))

	// then
	assert.Nil(t, errCode)

	searchEngineService.AssertExpectations(t)
	searchTaskTrackerService.AssertExpectations(t)
	realTimeService.AssertExpectations(t)
}

func writeRequestBody(payload any) io.ReadCloser {
	marshalledPayload, _ := json.Marshal(payload)
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(marshalledPayload)
	_ = gw.Close()
	body := io.NopCloser(bytes.NewReader(buf.Bytes()))
	return body
}
