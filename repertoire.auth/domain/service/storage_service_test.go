package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
	"testing"
)

func TestStorageService_Token_WhenValidateCredentialsFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewStorageService(jwtService)

	clientCredentials := model.ClientCredentials{}
	userID := uuid.New()

	unauthorizedError := wrapper.UnauthorizedError(errors.New("you are not authorized"))
	jwtService.On("ValidateCredentials", clientCredentials).Return(unauthorizedError).Once()

	// when
	token, expiresIn, errCode := _uut.Token(clientCredentials, userID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, unauthorizedError, errCode)
	assert.Empty(t, token)
	assert.Empty(t, expiresIn)

	jwtService.AssertExpectations(t)
}

func TestStorageService_Token_WhenCreateTokenFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewStorageService(jwtService)

	clientCredentials := model.ClientCredentials{}
	userID := uuid.New()

	jwtService.On("ValidateCredentials", clientCredentials).Return(nil).Once()

	unauthorizedError := wrapper.UnauthorizedError(errors.New("you are not authorized"))
	jwtService.On("CreateStorageToken", userID).Return("", "", unauthorizedError).Once()

	// when
	token, expiresIn, errCode := _uut.Token(clientCredentials, userID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, unauthorizedError, errCode)
	assert.Empty(t, token)
	assert.Empty(t, expiresIn)

	jwtService.AssertExpectations(t)
}

func TestStorageService_Token_WhenSuccessful_ShouldReturnToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := NewStorageService(jwtService)

	clientCredentials := model.ClientCredentials{}
	userID := uuid.New()

	expectedToken := "some-token"
	expectedExpiresIn := "1h"

	jwtService.On("ValidateCredentials", clientCredentials).Return(nil).Once()
	jwtService.On("CreateStorageToken", userID).Return(expectedToken, expectedExpiresIn, nil).Once()

	// when
	token, expiresIn, errCode := _uut.Token(clientCredentials, userID)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)
	assert.Equal(t, expectedExpiresIn, expiresIn)

	jwtService.AssertExpectations(t)
}
