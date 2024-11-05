package playlist

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

func TestSaveImageToPlaylist_WhenGetPlaylistFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := SaveImageToPlaylist{
		repository: playlistRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := SaveImageToPlaylist{
		repository: playlistRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToPlaylist{
		repository:              playlistRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, *playlist).Return(imagePath).Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenUpdatePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToPlaylist{
		repository:              playlistRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, *playlist).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToPlaylist{
		repository:              playlistRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, *playlist).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Equal(t, imagePath, string(*newPlaylist.ImageURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
