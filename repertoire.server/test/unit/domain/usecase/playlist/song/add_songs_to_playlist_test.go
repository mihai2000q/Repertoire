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

func TestAddSongsToPlaylist_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("GetPlaylistSongs", mock.Anything, request.ID).
		Return(internalError).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, res)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongsToPlaylist_WhenAddSongsToPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Return(internalError).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, res)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldReturnNoSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: request.SongIDs[0]},
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, &playlistSongs).
		Once()

	duplicateSongIDs := getDuplicateSongIDs(playlistSongs, request.SongIDs)

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.False(t, res.Success)
	assert.ElementsMatch(t, res.Duplicates, duplicateSongIDs)
	assert.Empty(t, res.Added)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongsToPlaylist_WhenWithoutDuplicatesNorForceAdd_ShouldReturnSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(request.SongIDs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, request.SongIDs[i], playlistSong.SongID)
				assert.Equal(t, uint(len(*playlistSongs)+1+i), playlistSong.SongTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.Empty(t, res.Duplicates)
	assert.ElementsMatch(t, res.Added, request.SongIDs)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddDuplicatesTooAndReturnSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:       uuid.New(),
		SongIDs:  []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
		ForceAdd: &[]bool{true}[0],
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: request.SongIDs[0]},
		{SongID: uuid.New()},
	}

	duplicateSongIDs := getDuplicateSongIDs(playlistSongs, request.SongIDs)

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, &playlistSongs).
		Once()

	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(request.SongIDs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, request.SongIDs[i], playlistSong.SongID)
				assert.Equal(t, uint(len(playlistSongs)+1+i), playlistSong.SongTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.ElementsMatch(t, res.Duplicates, duplicateSongIDs)
	assert.ElementsMatch(t, res.Added, request.SongIDs)

	playlistRepository.AssertExpectations(t)
}

func TestAddSongsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldSkipDuplicatesAndReturnSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddSongsToPlaylist(playlistRepository)

	request := requests.AddSongsToPlaylistRequest{
		ID:       uuid.New(),
		SongIDs:  []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
		ForceAdd: &[]bool{false}[0],
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: request.SongIDs[1]},
		{SongID: uuid.New()},
	}

	duplicateSongIDs := getDuplicateSongIDs(playlistSongs, request.SongIDs)

	var expectedSongIDs []uuid.UUID
	for _, songID := range request.SongIDs {
		if slices.Contains(duplicateSongIDs, songID) {
			continue
		}
		expectedSongIDs = append(expectedSongIDs, songID)
	}

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, &playlistSongs).
		Once()

	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(request.SongIDs)-len(duplicateSongIDs))
			assert.Len(t, *newPlaylistSongs, len(expectedSongIDs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, expectedSongIDs[i], playlistSong.SongID)
				assert.Equal(t, uint(len(playlistSongs)+1+i), playlistSong.SongTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.True(t, res.Success)
	assert.NotNil(t, res.Duplicates)
	assert.ElementsMatch(t, res.Duplicates, duplicateSongIDs)
	assert.ElementsMatch(t, res.Added, expectedSongIDs)

	playlistRepository.AssertExpectations(t)
}

func getDuplicateSongIDs(playlistSongs []model.PlaylistSong, requestSongIDs []uuid.UUID) []uuid.UUID {
	var duplicateSongIDs []uuid.UUID
	for _, playlistSong := range playlistSongs {
		if slices.Contains(requestSongIDs, playlistSong.SongID) {
			duplicateSongIDs = append(duplicateSongIDs, playlistSong.SongID)
		}
	}
	return duplicateSongIDs
}
