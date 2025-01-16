package playlist

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveSongFromPlaylist_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewMoveSongFromPlaylist(playlistRepository)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
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

func TestMoveSongFromPlaylist_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewMoveSongFromPlaylist(playlistRepository)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestMoveSongFromPlaylist_WhenOverSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewMoveSongFromPlaylist(playlistRepository)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongID},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over song not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestMoveSongFromPlaylist_WhenUpdateAllFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewMoveSongFromPlaylist(playlistRepository)

	request := requests.MoveSongFromPlaylistRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongID},
		{SongID: request.OverSongID},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(new([]model.PlaylistSong))).
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

func TestMoveSongFromPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name          string
		playlistSongs *[]model.PlaylistSong
		index         uint
		overIndex     uint
	}{
		{
			"Use case 1",
			&[]model.PlaylistSong{
				{SongID: uuid.New(), SongTrackNo: 1},
				{SongID: uuid.New(), SongTrackNo: 2},
				{SongID: uuid.New(), SongTrackNo: 3},
				{SongID: uuid.New(), SongTrackNo: 4},
				{SongID: uuid.New(), SongTrackNo: 5},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&[]model.PlaylistSong{
				{SongID: uuid.New(), SongTrackNo: 1},
				{SongID: uuid.New(), SongTrackNo: 2},
				{SongID: uuid.New(), SongTrackNo: 3},
				{SongID: uuid.New(), SongTrackNo: 4},
				{SongID: uuid.New(), SongTrackNo: 5},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			_uut := playlist.NewMoveSongFromPlaylist(playlistRepository)

			request := requests.MoveSongFromPlaylistRequest{
				ID:         uuid.New(),
				SongID:     (*tt.playlistSongs)[tt.index].SongID,
				OverSongID: (*tt.playlistSongs)[tt.overIndex].SongID,
			}

			// given - mocking
			playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
				Return(nil, tt.playlistSongs).
				Once()

			playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(new([]model.PlaylistSong))).
				Run(func(args mock.Arguments) {
					newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
					playlistSongs := slices.Clone(*newPlaylistSongs)
					slices.SortFunc(playlistSongs, func(a, b model.PlaylistSong) int {
						return cmp.Compare(a.SongTrackNo, b.SongTrackNo)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, playlistSongs[tt.overIndex-1].SongID, request.OverSongID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, playlistSongs[tt.overIndex+1].SongID, request.OverSongID)
					}
					assert.Equal(t, playlistSongs[tt.overIndex].SongID, request.SongID)
					for i, song := range playlistSongs {
						assert.Equal(t, uint(i)+1, song.SongTrackNo)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			playlistRepository.AssertExpectations(t)
		})
	}
}
