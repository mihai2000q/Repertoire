package user

import (
	"encoding/json"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"repertoire/server/domain/message/handler/user"
	"repertoire/server/internal/message/topics"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestUserDeletedHandler_WhenGetDocumentsFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := user.NewUserDeletedHandler(searchEngineService, nil)

	userID := uuid.New()

	internalError := errors.New("internal error")
	searchEngineService.On("GetDocuments", "userId = "+userID.String()).
		Return([]map[string]any{}, internalError).
		Once()

	// when
	payload, _ := json.Marshal(userID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
}

func TestUserDeletedHandler_WhenDocumentsLengthIs0_ShouldNotReturnAnyError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := user.NewUserDeletedHandler(searchEngineService, nil)

	userID := uuid.New()

	searchEngineService.On("GetDocuments", "userId = "+userID.String()).
		Return([]map[string]any{}, nil).
		Once()

	// when
	payload, _ := json.Marshal(userID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
}

func TestUserDeletedHandler_WhenPublishFails_ShouldReturnError(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := user.NewUserDeletedHandler(searchEngineService, messagePublisherService)

	userID := uuid.New()

	documents := []map[string]any{
		{"id": uuid.New().String()},
	}
	searchEngineService.On("GetDocuments", "userId = "+userID.String()).
		Return(documents, nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Return(internalError).
		Once()

	// when
	payload, _ := json.Marshal(userID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.Error(t, err)
	assert.Equal(t, err, internalError)

	searchEngineService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestUserDeletedHandler_WhenSuccessful_ShouldPublishMessageToDeleteFromSearchEngine(t *testing.T) {
	// given
	searchEngineService := new(service.SearchEngineServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := user.NewUserDeletedHandler(searchEngineService, messagePublisherService)

	userID := uuid.New()

	documents := []map[string]any{
		{"id": uuid.New().String()},
		{"id": uuid.New().String()},
	}
	searchEngineService.On("GetDocuments", "userId = "+userID.String()).
		Return(documents, nil).
		Once()

	messagePublisherService.On("Publish", topics.DeleteFromSearchEngineTopic, mock.IsType([]string{})).
		Run(func(args mock.Arguments) {
			ids := args.Get(1).([]string)
			assert.Len(t, ids, len(documents))
			for i := range ids {
				assert.Equal(t, documents[i]["id"], ids[i])
			}
		}).
		Return(nil).
		Once()

	// when
	payload, _ := json.Marshal(userID)
	msg := message.NewMessage("1", payload)
	err := _uut.Handle(msg)

	// then
	assert.NoError(t, err)

	searchEngineService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
