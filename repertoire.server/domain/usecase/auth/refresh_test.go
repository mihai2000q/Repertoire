package auth

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
	"testing"
)

func TestRefresh_WhenValidateJwtFails_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &Refresh{
		jwtService: jwtService,
	}
	request := request.RefreshRequest{
		Token: "This is a token",
	}

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	jwtService.On("Validate", request.Token).Return(uuid.Nil, internalError).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
}

func TestRefresh_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &Refresh{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := request.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	internalError := errors.New("something went wrong")
	userRepository.On("Get", new(model.User), userID).Return(internalError).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestRefresh_WhenUserIsEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &Refresh{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := request.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()
	userRepository.On("Get", new(model.User), userID).Return(nil).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "not authorized", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestRefresh_WhenCreateTokenFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &Refresh{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := request.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &model.User{ID: userID}
	userRepository.On("Get", new(model.User), userID).Return(nil, user).Once()

	internalError := wrapper.InternalServerError(errors.New("something went wrong"))
	jwtService.On("CreateToken", *user).Return("", internalError).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestRefresh_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &Refresh{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := request.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &model.User{ID: userID}
	userRepository.On("Get", new(model.User), userID).Return(nil, user).Once()

	expectedToken := "This is the new token"
	jwtService.On("CreateToken", *user).Return(expectedToken, nil).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}
