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

func TestDeleteFromSearchEngineHandler_WhenGetDocumentFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewDeleteFromSearchEngineHandler(nil, searchEngineService, nil)

	ids := []string{"id1", "id2", "id3"}

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocument", ids[0]).
		Return(map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(&ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
}

func TestDeleteFromSearchEngineHandler_WhenDeleteFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewDeleteFromSearchEngineHandler(nil, searchEngineService, nil)

	ids := []string{"id1", "id2", "id3"}

	searchEngineService.On("GetDocument", ids[0]).
		Return(map[string]any{}, nil).
		Once()

	internalError := errors.New("internal error")
	searchEngineService.On("Delete", ids).
		Return(int64(0), internalError).
		Once()

	// when
	payload, _ := json.Marshal(&ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
}

func TestDeleteFromSearchEngineHandler_WhenSuccessful_ShouldDeleteDataFromSearchEngineById(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	searchTaskTrackerService := new(service.SearchTaskTrackerServiceMock)
	_uut := search.NewDeleteFromSearchEngineHandler(logger.NewLoggerMock(), searchEngineService, searchTaskTrackerService)

	ids := []string{"id1", "id2", "id3"}

	userID := "some_user_ID"
	searchEngineService.On("GetDocument", ids[0]).
		Return(map[string]any{"userId": userID}, nil).
		Once()

	var taskID int64 = 23
	searchEngineService.On("Delete", ids).
		Return(taskID, nil).
		Once()

	searchTaskTrackerService.On("Track", strconv.FormatInt(taskID, 10), userID).Once()

	// when
	payload, _ := json.Marshal(&ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
	searchTaskTrackerService.AssertExpectations(t)
}
