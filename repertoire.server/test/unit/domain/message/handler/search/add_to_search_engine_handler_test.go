package search

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/test/unit/data/logger"
	"repertoire/server/test/unit/data/service"
	"strconv"
	"testing"
)

func TestAddToSearchEngineHandler_WhenAddFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewAddToSearchEngineHandler(nil, searchEngineService, nil)

	documents := []map[string]any{{"property1": "value2"}}

	internalError := errors.New("internal error")
	searchEngineService.On("Add", documents).
		Return(int64(0), internalError).
		Once()

	// when
	payload, _ := json.Marshal(&documents)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
}

func TestAddToSearchEngineHandler_WhenSuccessful_ShouldAddDocumentsToSearchEngine(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	searchTaskTrackerService := new(service.SearchTaskTrackerServiceMock)
	_uut := search.NewAddToSearchEngineHandler(logger.NewLoggerMock(), searchEngineService, searchTaskTrackerService)

	userID := "some-user-id"
	documents := []map[string]any{{"property1": "value2", "userId": userID}}

	var taskID int64 = 12
	searchEngineService.On("Add", documents).
		Return(taskID, nil).
		Once()

	searchTaskTrackerService.On("Track", strconv.FormatInt(taskID, 10), userID).Once()

	// when
	payload, _ := json.Marshal(&documents)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
	searchTaskTrackerService.AssertExpectations(t)
}
