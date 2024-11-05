package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImageFromPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := DeleteImageFromPlaylist{repository: playlistRepository}

	id := uuid.New()

	// given - mocking
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

func TestDeleteImageFromPlaylist_WhenGetPlaylistFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := DeleteImageFromPlaylist{repository: playlistRepository}

	id := uuid.New()

	// given - mocking
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromPlaylist{
		repository:     playlistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*playlist.ImageURL)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenUpdatePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromPlaylist{
		repository:     playlistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	storageService.On("Delete", string(*playlist.ImageURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(playlist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromPlaylist{
		repository:     playlistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, playlist).Once()

	storageService.On("Delete", string(*playlist.ImageURL)).Return(nil).Once()

	playlistRepository.On("Update", mock.IsType(playlist)).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Nil(t, newPlaylist.ImageURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}