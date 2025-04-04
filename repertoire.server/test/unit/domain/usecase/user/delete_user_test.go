package user

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	_uut := user.NewDeleteUser(nil, jwtService, nil)

	token := "This is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteUser_WhenDeleteUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewDeleteUser(userRepository, jwtService, nil)

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestDeleteUser_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := user.NewDeleteUser(userRepository, jwtService, messagePublisherService)

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	userRepository.On("Delete", id).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.UserDeletedTopic, id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteUser_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                 string
		deleteDirectoryError *wrapper.ErrorCode
	}{
		{
			"Without Files",
			wrapper.NotFoundError(errors.New("cannot delete the directory as it's not found")),
		},
		{
			"With Files",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			jwtService := new(service.JwtServiceMock)
			userRepository := new(repository.UserRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := user.NewDeleteUser(userRepository, jwtService, messagePublisherService)

			token := "This is a token"

			id := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

			userRepository.On("Delete", id).Return(nil).Once()

			messagePublisherService.On("Publish", topics.UserDeletedTopic, id).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(token)

			// then
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			userRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
