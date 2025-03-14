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

func TestDeleteFromSearchEngineHandler_WhenDeleteFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewDeleteFromSearchEngineHandler(searchEngineService)

	ids := []string{"id1", "id2", "id3"}

	internalError := errors.New("internal error")
	searchEngineService.On("Delete", ids).
		Return(internalError).
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
	_uut := search.NewDeleteFromSearchEngineHandler(searchEngineService)

	ids := []string{"id1", "id2", "id3"}

	searchEngineService.On("Delete", ids).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(&ids)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
}
