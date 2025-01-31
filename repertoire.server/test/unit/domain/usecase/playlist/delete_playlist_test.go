package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeletePlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil, nil)

	id := uuid.New()

	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil, nil)

	id := uuid.New()

	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenDeleteDirectoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, storageService, storageFilePathProvider)

	id := uuid.New()

	mockPlaylist := &model.Playlist{
		ID: id,
	}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageFilePathProvider.On("HasPlaylistFiles", *mockPlaylist).Return(true).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetPlaylistDirectoryPath", *mockPlaylist).Return(directoryPath).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteDirectory", directoryPath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeletePlaylist_WhenDeletePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil, storageFilePathProvider)

	id := uuid.New()

	mockPlaylist := &model.Playlist{
		ID: id,
	}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageFilePathProvider.On("HasPlaylistFiles", *mockPlaylist).Return(false).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeletePlaylist_WhenSuccessful_ShouldDeletePlaylist(t *testing.T) {
	tests := []struct {
		name     string
		playlist model.Playlist
		hasFiles bool
	}{
		{
			"Without Files",
			model.Playlist{ID: uuid.New()},
			false,
		},
		{
			"With Files",
			model.Playlist{ID: uuid.New()},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			storageService := new(service.StorageServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := playlist.NewDeletePlaylist(playlistRepository, storageService, storageFilePathProvider)

			id := tt.playlist.ID

			playlistRepository.On("Get", new(model.Playlist), id).Return(nil, &tt.playlist).Once()

			storageFilePathProvider.On("HasPlaylistFiles", tt.playlist).Return(tt.hasFiles).Once()

			if tt.hasFiles {
				directoryPath := "some directory path"
				storageFilePathProvider.On("GetPlaylistDirectoryPath", tt.playlist).
					Return(directoryPath).
					Once()

				storageService.On("DeleteDirectory", directoryPath).
					Return(nil).
					Once()
			}

			playlistRepository.On("Delete", id).Return(nil).Once()

			// when
			errCode := _uut.Handle(id)

			// then
			assert.Nil(t, errCode)

			playlistRepository.AssertExpectations(t)
			storageService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}
