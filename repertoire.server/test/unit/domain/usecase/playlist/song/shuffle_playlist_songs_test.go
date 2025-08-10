package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShufflePlaylistSongs_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewShufflePlaylistSongs(playlistRepository)

	request := requests.ShufflePlaylistSongsRequest{ID: uuid.New()}

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetPlaylistSongs",
			new([]model.PlaylistSong),
			request.ID,
		).
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

func TestShufflePlaylistSongs_WhenUpdateAllPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewShufflePlaylistSongs(playlistRepository)

	request := requests.ShufflePlaylistSongsRequest{ID: uuid.New()}

	playlistSongs := []model.PlaylistSong{{ID: uuid.New()}}

	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, &playlistSongs).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(&playlistSongs)).
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

func TestShufflePlaylistSongs_WhenSuccessful_ShouldReturnSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewShufflePlaylistSongs(playlistRepository)

	request := requests.ShufflePlaylistSongsRequest{ID: uuid.New()}

	playlistSongs := []model.PlaylistSong{
		{ID: uuid.New(), SongTrackNo: 1},
		{ID: uuid.New(), SongTrackNo: 2},
		{ID: uuid.New(), SongTrackNo: 3},
		{ID: uuid.New(), SongTrackNo: 4},
	}
	oldPlaylistSongs := slices.Clone(playlistSongs)

	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, &playlistSongs).
		Once()

	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(&playlistSongs)).
		Run(func(args mock.Arguments) {
			songs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *songs, len(oldPlaylistSongs))
			for i := range *songs {
				assert.Equal(t, uint(i+1), (*songs)[i].SongTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
