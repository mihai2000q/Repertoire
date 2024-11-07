package user

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	_uut := DeleteUser{jwtService: jwtService}

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
	_uut := DeleteUser{
		repository: userRepository,
		jwtService: jwtService,
	}

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInsufficientStorage, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestDeleteUser_WhenSuccessfull_ShouldNotReturnAnyError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := DeleteUser{
		repository: userRepository,
		jwtService: jwtService,
	}

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	userRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}
