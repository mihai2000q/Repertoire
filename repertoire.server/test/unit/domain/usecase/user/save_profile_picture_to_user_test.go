package user

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveProfilePictureToUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := user.NewSaveProfilePictureToUser(
		nil,
		nil,
		jwtService,
		nil,
	)

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

func TestSaveProfilePictureToUser_WhenGetUserFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		nil,
		jwtService,
		nil,
	)

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

func TestSaveProfilePictureToUser_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		nil,
		jwtService,
		nil,
	)

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

func TestSaveProfilePictureToUser_WhenStorageDeleteFileFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		nil,
		jwtService,
		storageService,
	)

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id, ProfilePictureURL: &[]internal.FilePath{"file-path"}[0]}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockUser.ProfilePictureURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveProfilePictureToUser_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		storageFilePathProvider,
		jwtService,
		storageService,
	)

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, mock.IsType(*mockUser)).
		Return(imagePath).
		Once()

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
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		storageFilePathProvider,
		jwtService,
		storageService,
	)

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, mock.IsType(*mockUser)).
		Return(imagePath).
		Once()

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

func TestSaveProfilePictureToUser_WhenWithoutOldProfilePicture_ShouldSaveNewOneAndNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		storageFilePathProvider,
		jwtService,
		storageService,
	)

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id, ProfilePictureURL: nil}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, mock.IsType(*mockUser)).
		Run(func(args mock.Arguments) {
			newUser := args.Get(1).(model.User)
			assert.Equal(t, newUser.ID, id)
			assert.WithinDuration(t, newUser.UpdatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(imagePath).
		Once()

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

func TestSaveProfilePictureToUser_WhenWithOldProfilePicture_ShouldDeleteOldPictureSaveNewOneAndNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := user.NewSaveProfilePictureToUser(
		userRepository,
		storageFilePathProvider,
		jwtService,
		storageService,
	)

	file := new(multipart.FileHeader)
	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id, ProfilePictureURL: &[]internal.FilePath{"file-path"}[0]}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	storageService.On("DeleteFile", *mockUser.ProfilePictureURL).Return(nil).Once()

	imagePath := "users file path"
	storageFilePathProvider.On("GetUserProfilePicturePath", file, mock.IsType(*mockUser)).
		Run(func(args mock.Arguments) {
			newUser := args.Get(1).(model.User)
			assert.Equal(t, newUser.ID, id)
			assert.WithinDuration(t, newUser.UpdatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(imagePath).
		Once()
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
