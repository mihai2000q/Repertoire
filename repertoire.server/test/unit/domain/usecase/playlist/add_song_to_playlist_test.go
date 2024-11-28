package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSongToPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(nil, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("Get", mock.IsType(new(model.Playlist)), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenPlaylistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(nil, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	playlistRepository.On("Get", mock.IsType(new(model.Playlist)), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(songRepository, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	mockPlaylist := &model.Playlist{ID: request.ID}
	playlistRepository.On("Get", mock.IsType(mockPlaylist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(songRepository, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	mockPlaylist := &model.Playlist{ID: request.ID}
	playlistRepository.On("Get", mock.IsType(mockPlaylist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	songRepository.On("Get", mock.IsType(new(model.Song)), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenCountSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(songRepository, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	mockPlaylist := &model.Playlist{ID: request.ID}
	playlistRepository.On("Get", mock.IsType(mockPlaylist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(mockSong), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("CountSongs", mock.Anything, request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenAddSongToPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(songRepository, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	mockPlaylist := &model.Playlist{ID: request.ID}
	playlistRepository.On("Get", mock.IsType(mockPlaylist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(mockSong), request.SongID).
		Return(nil, mockSong).
		Once()

	count := &[]int64{12}[0]
	playlistRepository.On("CountSongs", mock.IsType(count), request.ID).
		Return(nil, count).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("AddSong", mockPlaylist, mockSong).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(songRepository, playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	mockPlaylist := &model.Playlist{ID: request.ID}
	playlistRepository.On("Get", mock.IsType(mockPlaylist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(mockSong), request.SongID).
		Return(nil, mockSong).
		Once()

	count := &[]int64{12}[0]
	playlistRepository.On("CountSongs", mock.IsType(count), request.ID).
		Return(nil, count).
		Once()

	playlistRepository.On("AddSong", mockPlaylist, mockSong).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}