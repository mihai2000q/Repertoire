package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImageFromPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, nil)

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
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, nil)

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
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*mockPlaylist.ImageURL)).Return(internalError).Once()

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
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("Delete", string(*mockPlaylist.ImageURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
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
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("Delete", string(*mockPlaylist.ImageURL)).Return(nil).Once()

	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
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
