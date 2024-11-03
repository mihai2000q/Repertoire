package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveSongFromPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository: playlistRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("Get", mock.IsType(new(model.Playlist)), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestRemoveSongFromPlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository: playlistRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	playlistRepository.On("Get", mock.IsType(new(model.Playlist)), id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestRemoveSongFromPlaylist_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository:     playlistRepository,
		songRepository: songRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id}
	playlistRepository.On("Get", mock.IsType(playlist), id).
		Return(nil, playlist).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromPlaylist_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository:     playlistRepository,
		songRepository: songRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	playlist := &model.Playlist{ID: id}
	playlistRepository.On("Get", mock.IsType(playlist), id).
		Return(nil, playlist).
		Once()

	songRepository.On("Get", mock.IsType(new(model.Song)), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromPlaylist_WhenRemoveSongFromPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository:     playlistRepository,
		songRepository: songRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{ID: songID}
	songRepository.On("Get", mock.IsType(song), songID).
		Return(nil, song).
		Once()

	playlist := &model.Playlist{ID: id}
	playlistRepository.On("Get", mock.IsType(playlist), id).
		Return(nil, playlist).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("RemoveSong", playlist, song).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromPlaylist{
		repository:     playlistRepository,
		songRepository: songRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{ID: songID}
	songRepository.On("Get", mock.IsType(song), songID).
		Return(nil, song).
		Once()

	playlist := &model.Playlist{ID: id}
	playlistRepository.On("Get", mock.IsType(playlist), id).
		Return(nil, playlist).
		Once()

	playlistRepository.On("RemoveSong", playlist, song).Return(nil).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
