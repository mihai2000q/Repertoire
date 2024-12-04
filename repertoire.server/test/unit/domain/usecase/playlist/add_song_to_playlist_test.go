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

func TestAddSongToPlaylist_WhenCountSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
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
}

func TestAddSongToPlaylist_WhenAddSongToPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	count := &[]int64{12}[0]
	playlistRepository.On("CountSongs", mock.IsType(count), request.ID).
		Return(nil, count).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("AddSong", mock.IsType(new(model.PlaylistSong))).
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

func TestAddSongToPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddSongToPlaylist(playlistRepository)

	request := requests.AddSongToPlaylistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	count := &[]int64{12}[0]
	playlistRepository.On("CountSongs", mock.IsType(count), request.ID).
		Return(nil, count).
		Once()

	playlistRepository.On("AddSong", mock.IsType(new(model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			playlistSong := args.Get(0).(*model.PlaylistSong)
			assert.Equal(t, request.ID, playlistSong.PlaylistID)
			assert.Equal(t, request.SongID, playlistSong.SongID)
			assert.Equal(t, uint(*count), playlistSong.SongTrackNo)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
