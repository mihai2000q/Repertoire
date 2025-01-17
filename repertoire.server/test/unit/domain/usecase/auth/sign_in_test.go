package auth

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/auth"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"strings"
	"testing"
)

func TestSignIn_WhenGetUserByEmailFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := auth.NewSignIn(nil, nil, userRepository)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	internalError := errors.New("something went wrong")
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestSignIn_WhenUserIsEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := auth.NewSignIn(nil, nil, userRepository)

	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestSignIn_WhenPasswordsAreNotTheSame_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := auth.NewSignIn(nil, bCryptService, userRepository)

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
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestSignIn_WhenCreateTokenFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := auth.NewSignIn(jwtService, bCryptService, userRepository)

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
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestSignIn_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := auth.NewSignIn(jwtService, bCryptService, userRepository)

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
	token, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
