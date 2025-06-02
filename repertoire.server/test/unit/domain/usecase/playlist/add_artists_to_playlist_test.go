package playlist

import (
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

func TestAddArtistToPlaylist_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, nil)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New()},
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

func TestAddArtistsToPlaylist_WhenGetArtistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), request.ID).
		Return(nil, playlistSongs).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
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
	artistRepository.AssertExpectations(t)
}

func TestAddArtistsToPlaylist_WhenAddSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), request.ID).
		Return(nil, playlistSongs).
		Once()

	artists := &[]model.Artist{}
	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(artists), request.ArtistIDs).
		Return(nil, artists).
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
	artistRepository.AssertExpectations(t)
}

func TestAddArtistsToPlaylist_WhenWithoutDuplicatesButWithForceAdd_ShouldReturnBadRequestError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.AddArtistsToPlaylistRequest
	}{
		{
			"with force add false",
			requests.AddArtistsToPlaylistRequest{
				ID:        uuid.New(),
				ArtistIDs: []uuid.UUID{uuid.New()},
				ForceAdd:  &[]bool{false}[0],
			},
		},
		{
			"with force add true",
			requests.AddArtistsToPlaylistRequest{
				ID:        uuid.New(),
				ArtistIDs: []uuid.UUID{uuid.New()},
				ForceAdd:  &[]bool{true}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			artistRepository := new(repository.ArtistRepositoryMock)
			_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

			// given - mocking
			playlistSongs := &[]model.PlaylistSong{}
			playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), tt.request.ID).
				Return(nil, playlistSongs).
				Once()

			artists := &[]model.Artist{}
			artistRepository.On("GetAllByIDsWithSongs", mock.IsType(artists), tt.request.ArtistIDs).
				Return(nil, artists).
				Once()

			// when
			res, errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, res)
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
			assert.Equal(t, "force adding when there are no duplicates", errCode.Error.Error())

			playlistRepository.AssertExpectations(t)
			artistRepository.AssertExpectations(t)
		})
	}
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldReturnNoSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: uuid.New()},
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	artists := []model.Artist{
		{ID: request.ArtistIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.ArtistIDs[1], Songs: []model.Song{{ID: playlistSongs[1].SongID}}},
	}
	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
		Return(nil, &artists).
		Once()

	duplicatedArtistIDs, duplicatedSongIDs := getArtistAndSongDuplicates(artists, playlistSongs)

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.False(t, res.Success)
	assert.ElementsMatch(t, res.DuplicateArtistIDs, duplicatedArtistIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.Empty(t, res.AddedSongIDs)

	playlistRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestAddArtistsToPlaylist_WhenWithoutDuplicatesNorForceAdd_ShouldReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	artists := []model.Artist{
		{ID: request.ArtistIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.ArtistIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
		Return(nil, &artists).
		Once()

	var newSongs []model.Song
	for _, artist := range artists {
		newSongs = append(newSongs, artist.Songs...)
	}

	var addedSongIDs []uuid.UUID
	oldPlaylistSongsLength := len(playlistSongs) + 1
	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(newSongs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, uint(oldPlaylistSongsLength+i), playlistSong.SongTrackNo)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, newSongs[i].ID, playlistSong.SongID)
				addedSongIDs = append(addedSongIDs, playlistSong.SongID)
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
	assert.Empty(t, res.DuplicateArtistIDs)
	assert.Empty(t, res.DuplicateSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddDuplicateSongsAndReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New(), uuid.New()},
		ForceAdd:  &[]bool{true}[0],
	}
	artists := []model.Artist{
		{ID: request.ArtistIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.ArtistIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: artists[0].Songs[0].ID},
		{SongID: uuid.New()},
		{SongID: artists[1].Songs[0].ID},
	}

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
		Return(nil, &artists).
		Once()

	duplicatedArtistIDs, duplicatedSongIDs := getArtistAndSongDuplicates(artists, playlistSongs)
	var addedSongIDs []uuid.UUID
	for _, artist := range artists {
		for _, song := range artist.Songs {
			addedSongIDs = append(addedSongIDs, song.ID)
		}
	}

	oldPlaylistSongsLength := len(playlistSongs) + 1
	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(addedSongIDs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, uint(oldPlaylistSongsLength+i), playlistSong.SongTrackNo)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, addedSongIDs[i], playlistSong.SongID)
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
	assert.ElementsMatch(t, res.DuplicateArtistIDs, duplicatedArtistIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldSkipDuplicateSongsAndReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

	request := requests.AddArtistsToPlaylistRequest{
		ID:        uuid.New(),
		ArtistIDs: []uuid.UUID{uuid.New(), uuid.New()},
		ForceAdd:  &[]bool{false}[0],
	}

	artists := []model.Artist{
		{ID: request.ArtistIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.ArtistIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: artists[0].Songs[0].ID},
		{SongID: uuid.New()},
		{SongID: artists[1].Songs[0].ID},
	}

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
		Return(nil, &artists).
		Once()

	duplicatedArtistIDs, duplicatedSongIDs := getArtistAndSongDuplicates(artists, playlistSongs)
	var addedSongIDs []uuid.UUID
	for _, artist := range artists {
		for _, song := range artist.Songs {
			if !slices.Contains(duplicatedSongIDs, song.ID) {
				addedSongIDs = append(addedSongIDs, song.ID)
			}
		}
	}

	oldPlaylistSongsLength := len(playlistSongs) + 1
	playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newPlaylistSongs, len(addedSongIDs))
			for i, playlistSong := range *newPlaylistSongs {
				assert.NotEmpty(t, playlistSong.ID)
				assert.Equal(t, uint(oldPlaylistSongsLength+i), playlistSong.SongTrackNo)
				assert.Equal(t, request.ID, playlistSong.PlaylistID)
				assert.Equal(t, addedSongIDs[i], playlistSong.SongID)
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
	assert.ElementsMatch(t, res.DuplicateArtistIDs, duplicatedArtistIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func getArtistAndSongDuplicates(artists []model.Artist, playlistSongs []model.PlaylistSong) ([]uuid.UUID, []uuid.UUID) {
	var duplicatedAlbumIDs []uuid.UUID
	var duplicatedSongIDs []uuid.UUID
	for _, artist := range artists {
		var songIDs []uuid.UUID
		for _, song := range artist.Songs {
			for _, ps := range playlistSongs {
				if song.ID == ps.SongID {
					songIDs = append(songIDs, song.ID)
				}
			}
		}
		duplicatedSongIDs = append(duplicatedSongIDs, songIDs...)
		if len(songIDs) == len(artist.Songs) {
			duplicatedAlbumIDs = append(duplicatedAlbumIDs, artist.ID)
		}
	}

	return duplicatedAlbumIDs, duplicatedSongIDs
}
