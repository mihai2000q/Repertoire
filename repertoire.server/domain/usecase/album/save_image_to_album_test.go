package album

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToAlbum_WhenGetAlbumFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := SaveImageToAlbum{
		repository: albumRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestSaveImageToAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := SaveImageToAlbum{
		repository: albumRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	albumRepository.On("Get", new(model.Album), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestSaveImageToAlbum_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToAlbum{
		repository:              albumRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *album).Return(imagePath).Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToAlbum{
		repository:              albumRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *album).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	albumRepository.On("Update", mock.IsType(new(model.Album))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToAlbum{
		repository:              albumRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *album).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(new(model.Album))).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, imagePath, string(*newAlbum.ImageURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}