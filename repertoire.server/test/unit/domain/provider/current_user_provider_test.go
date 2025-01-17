package provider

import (
	"errors"
	"net/http"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Get
func TestCurrentUserProvider_Get_WhenJwtServiceReturnsAnErrorCode_ShouldReturnErrorCode(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := provider.NewCurrentUserProvider(jwtService, nil)

	token := "this is a token"

	internalErrorCode := wrapper.InternalServerError(errors.New("internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, internalErrorCode).Once()

	// when
	user, errCode := _uut.Get(token)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalErrorCode, errCode)

	jwtService.AssertExpectations(t)
}

func TestCurrentUserProvider_Get_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := provider.NewCurrentUserProvider(jwtService, userRepository)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), userID).Return(internalError).Once()

	// when
	user, errCode := _uut.Get(token)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestCurrentUserProvider_Get_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := provider.NewCurrentUserProvider(jwtService, userRepository)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	userRepository.On("Get", new(model.User), userID).Return(nil).Once()

	// when
	user, errCode := _uut.Get(token)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestCurrentUserProvider_Get_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := provider.NewCurrentUserProvider(jwtService, userRepository)

	token := "this is a token"

	userID := uuid.New()
	expectedUser := &model.User{
		ID:    userID,
		Name:  "Samuel",
		Email: "samuel.samuel@gmail.com",
	}

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	userRepository.On("Get", new(model.User), userID).Return(nil, expectedUser).Once()

	// when
	user, errCode := _uut.Get(token)

	// then
	assert.NotEmpty(t, user)
	assert.Equal(t, expectedUser, &user)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}
