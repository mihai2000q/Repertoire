package search

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestUpdateFromSearchEngineHandler_WhenUpdateFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewUpdateFromSearchEngineHandler(searchEngineService)

	documents := []any{"id1", "id2", "id3"}

	internalError := errors.New("internal error")
	searchEngineService.On("Update", documents).
		Return(internalError).
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

func TestUpdateFromSearchEngineHandler_WhenSuccessful_ShouldUpdateDataFromSearchEngine(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewUpdateFromSearchEngineHandler(searchEngineService)

	documents := []any{"id1", "id2", "id3"}

	searchEngineService.On("Update", documents).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(&documents)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
}
