package user

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveProfilePictureToUser_WhenGetUserFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := SaveProfilePictureToUser{
		jwtService: jwtService,
		repository: userRepository,
	}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := SaveProfilePictureToUser{jwtService: jwtService}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	err := wrapper.ForbiddenError(errors.New("internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, err).Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, err, errCode)

	jwtService.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := SaveProfilePictureToUser{
		jwtService: jwtService,
		repository: userRepository,
	}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	userRepository.On("Get", new(model.User), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveProfilePictureToUser{
		jwtService:              jwtService,
		repository:              userRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, *user).Return(imagePath).Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenUpdateUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveProfilePictureToUser{
		jwtService:              jwtService,
		repository:              userRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, *user).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Update", mock.IsType(new(model.User))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveProfilePictureToUser{
		jwtService:              jwtService,
		repository:              userRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, *user).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	userRepository.On("Update", mock.IsType(new(model.User))).
		Run(func(args mock.Arguments) {
			newUser := args.Get(0).(*model.User)
			assert.Equal(t, imagePath, string(*newUser.ProfilePictureURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
