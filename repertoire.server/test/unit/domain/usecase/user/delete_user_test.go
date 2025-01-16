package user

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/wrapper"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	_uut := user.NewDeleteUser(nil, jwtService, nil, nil)

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

func TestDeleteUser_WhenDeleteDirectoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := user.NewDeleteUser(nil, jwtService, storageService, storageFilePathProvider)

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetUserDirectoryPath", id).Return(directoryPath).Once()

	internalError := errors.New("internal error")
	storageService.On("DeleteDirectory", directoryPath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteUser_WhenDeleteUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := user.NewDeleteUser(userRepository, jwtService, storageService, storageFilePathProvider)

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetUserDirectoryPath", id).Return(directoryPath).Once()

	storageService.On("DeleteDirectory", directoryPath).Return(nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestDeleteUser_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// give
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := user.NewDeleteUser(userRepository, jwtService, storageService, storageFilePathProvider)

	token := "This is a token"

	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetUserDirectoryPath", id).Return(directoryPath).Once()

	storageService.On("DeleteDirectory", directoryPath).Return(nil).Once()

	userRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}
