package album

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToAlbum_WhenGetAlbumFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewSaveImageToAlbum(albumRepository, nil, nil, nil)

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
	_uut := album.NewSaveImageToAlbum(albumRepository, nil, nil, nil)

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
	_uut := album.NewSaveImageToAlbum(albumRepository, storageFilePathProvider, storageService, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *mockAlbum).Return(imagePath).Once()

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
	_uut := album.NewSaveImageToAlbum(albumRepository, storageFilePathProvider, storageService, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *mockAlbum).Return(imagePath).Once()

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

func TestSaveImageToAlbum_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewSaveImageToAlbum(
		albumRepository,
		storageFilePathProvider,
		storageService,
		messagePublisherService,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *mockAlbum).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(new(model.Album))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
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
	messagePublisherService.AssertExpectations(t)
}

func TestSaveImageToAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewSaveImageToAlbum(
		albumRepository,
		storageFilePathProvider,
		storageService,
		messagePublisherService,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: nil}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	imagePath := "albums file path"
	storageFilePathProvider.On("GetAlbumImagePath", file, *mockAlbum).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(new(model.Album))).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, imagePath, string(*newAlbum.ImageURL))
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
