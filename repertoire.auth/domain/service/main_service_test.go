package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/auth/api/requests"
	"repertoire/auth/data/logger"
	"repertoire/auth/data/repository"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
	"strings"
	"testing"
)

func TestRefresh_WhenValidateJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewMainService(jwtService, nil, nil, logger.NewLoggerMock())

	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	jwtService.On("Validate", request.Token).Return(uuid.Nil, internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

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
	_uut := NewMainService(jwtService, nil, userRepository, nil)

	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	internalError := errors.New("something went wrong")
	userRepository.On("Get", new(model.User), userID).Return(internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

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
	_uut := NewMainService(jwtService, nil, userRepository, nil)

	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()
	userRepository.On("Get", new(model.User), userID).Return(nil).Once()

	// when
	token, errCode := _uut.Refresh(request)

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
	_uut := NewMainService(jwtService, nil, userRepository, nil)

	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &model.User{ID: userID}
	userRepository.On("Get", new(model.User), userID).Return(nil, user).Once()

	internalError := wrapper.InternalServerError(errors.New("something went wrong"))
	jwtService.On("CreateToken", *user).Return("", internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

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
	_uut := NewMainService(jwtService, nil, userRepository, nil)

	// given - mocking
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &model.User{ID: userID}
	userRepository.On("Get", new(model.User), userID).Return(nil, user).Once()

	expectedToken := "This is the new token"
	jwtService.On("CreateToken", *user).Return(expectedToken, nil).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestMainService_SignIn_WhenGetUserByEmailFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := NewMainService(nil, nil, userRepository, nil)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	internalError := errors.New("something went wrong")
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestMainService_SignIn_WhenUserIsEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := NewMainService(nil, nil, userRepository, nil)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestMainService_SignIn_WhenPasswordsAreNotTheSame_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := NewMainService(nil, bCryptService, userRepository, nil)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &model.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(errors.New("")).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestMainService_SignIn_WhenCreateTokenFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := NewMainService(jwtService, bCryptService, userRepository, nil)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &model.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(nil).Once()

	internalError := wrapper.InternalServerError(errors.New("something went wrong"))
	jwtService.On("CreateToken", *user).Return("", internalError).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestMainService_SignIn_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := NewMainService(jwtService, bCryptService, userRepository, nil)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &model.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(nil).Once()

	expectedToken := "This is the generated token"
	jwtService.On("CreateToken", *user).Return(expectedToken, nil).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
