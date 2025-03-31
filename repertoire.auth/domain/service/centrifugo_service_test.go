package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"testing"
)

func TestCentrifugoService_Token_WhenGetUserIDFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewCentrifugoService(jwtService)

	userToken := "initial-user-token"
	
	unauthorizedError := wrapper.UnauthorizedError(errors.New("you are not authorized"))
	jwtService.On("GetUserIDFromJwt", userToken).Return(uuid.Nil, unauthorizedError).Once()

	// when
	token, expiresIn, errCode := _uut.Token(userToken)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, unauthorizedError, errCode)
	assert.Empty(t, token)
	assert.Empty(t, expiresIn)

	jwtService.AssertExpectations(t)
}

func TestCentrifugoService_Token_WhenCreateTokenFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewCentrifugoService(jwtService)

	userToken := "initial-user-token"

	userID := uuid.New()
	jwtService.On("GetUserIDFromJwt", userToken).Return(userID, nil).Once()

	unauthorizedError := wrapper.UnauthorizedError(errors.New("you are not authorized"))
	jwtService.On("CreateCentrifugoToken", userID).Return("", "", unauthorizedError).Once()

	// when
	token, expiresIn, errCode := _uut.Token(userToken)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, unauthorizedError, errCode)
	assert.Empty(t, token)
	assert.Empty(t, expiresIn)

	jwtService.AssertExpectations(t)
}

func TestCentrifugoService_Token_WhenSuccessful_ShouldReturnToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewCentrifugoService(jwtService)

	userToken := "initial-user-token"

	expectedToken := "some-token"
	expectedExpiresIn := "1h"

	userID := uuid.New()
	jwtService.On("GetUserIDFromJwt", userToken).Return(userID, nil).Once()

	jwtService.On("CreateCentrifugoToken", userID).Return(expectedToken, expectedExpiresIn, nil).Once()

	// when
	token, expiresIn, errCode := _uut.Token(userToken)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)
	assert.Equal(t, expectedExpiresIn, expiresIn)

	jwtService.AssertExpectations(t)
}
