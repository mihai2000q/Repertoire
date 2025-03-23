package search

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/test/unit/data/logger"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestAddToSearchEngineHandler_WhenAddFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewAddToSearchEngineHandler(nil, searchEngineService)

	data := []any{"something"}

	internalError := errors.New("internal error")
	searchEngineService.On("Add", data).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(&data)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
}

func TestAddToSearchEngineHandler_WhenSuccessful_ShouldAddDataToSearchEngine(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := search.NewAddToSearchEngineHandler(logger.NewLoggerMock(), searchEngineService)

	data := []any{"something"}

	searchEngineService.On("Add", data).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(&data)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
}
