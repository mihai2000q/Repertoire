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

func TestAddAlbumsToPlaylist_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, nil)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
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

func TestAddAlbumsToPlaylist_WhenGetAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), request.ID).
		Return(nil, playlistSongs).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Album)), request.AlbumIDs).
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
	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToPlaylist_WhenAddSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), request.ID).
		Return(nil, playlistSongs).
		Once()

	albums := &[]model.Album{}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
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
	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToPlaylist_WhenWithoutDuplicatesButWithForceAdd_ShouldReturnBadRequestError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.AddAlbumsToPlaylistRequest
	}{
		{
			"with force add false",
			requests.AddAlbumsToPlaylistRequest{
				ID:       uuid.New(),
				AlbumIDs: []uuid.UUID{uuid.New()},
				ForceAdd: &[]bool{false}[0],
			},
		},
		{
			"with force add true",
			requests.AddAlbumsToPlaylistRequest{
				ID:       uuid.New(),
				AlbumIDs: []uuid.UUID{uuid.New()},
				ForceAdd: &[]bool{true}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

			// given - mocking
			playlistSongs := &[]model.PlaylistSong{}
			playlistRepository.On("GetPlaylistSongs", mock.IsType(playlistSongs), tt.request.ID).
				Return(nil, playlistSongs).
				Once()

			albums := &[]model.Album{}
			albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), tt.request.AlbumIDs).
				Return(nil, albums).
				Once()

			// when
			res, errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, res)
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
			assert.Equal(t, "force adding when there are no duplicates", errCode.Error.Error())

			playlistRepository.AssertExpectations(t)
			albumRepository.AssertExpectations(t)
		})
	}
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldReturnNoSuccess(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
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

	albums := []model.Album{
		{ID: request.AlbumIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.AlbumIDs[1], Songs: []model.Song{{ID: playlistSongs[1].SongID}}},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Album)), request.AlbumIDs).
		Return(nil, &albums).
		Once()

	duplicatedAlbumIDs, duplicatedSongIDs := getAlbumAndSongDuplicates(albums, playlistSongs)

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, res)
	assert.False(t, res.Success)
	assert.ElementsMatch(t, res.DuplicateAlbumIDs, duplicatedAlbumIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.Empty(t, res.AddedSongIDs)

	playlistRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToPlaylist_WhenWithoutDuplicatesNorForceAdd_ShouldReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
	}
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	albums := []model.Album{
		{ID: request.AlbumIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.AlbumIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Album)), request.AlbumIDs).
		Return(nil, &albums).
		Once()

	var newSongs []model.Song
	for _, album := range albums {
		newSongs = append(newSongs, album.Songs...)
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
	assert.Empty(t, res.DuplicateAlbumIDs)
	assert.Empty(t, res.DuplicateSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddDuplicateSongsAndReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
		ForceAdd: &[]bool{true}[0],
	}
	albums := []model.Album{
		{ID: request.AlbumIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.AlbumIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: albums[0].Songs[0].ID},
		{SongID: uuid.New()},
		{SongID: albums[1].Songs[0].ID},
	}

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Album)), request.AlbumIDs).
		Return(nil, &albums).
		Once()

	duplicatedAlbumIDs, duplicatedSongIDs := getAlbumAndSongDuplicates(albums, playlistSongs)
	var addedSongIDs []uuid.UUID
	for _, album := range albums {
		for _, song := range album.Songs {
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
	assert.ElementsMatch(t, res.DuplicateAlbumIDs, duplicatedAlbumIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldSkipDuplicateSongsAndReturnSuccessResponse(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := playlist.NewAddAlbumsToPlaylist(playlistRepository, albumRepository)

	request := requests.AddAlbumsToPlaylistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
		ForceAdd: &[]bool{false}[0],
	}
	albums := []model.Album{
		{ID: request.AlbumIDs[0], Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}}},
		{ID: request.AlbumIDs[1], Songs: []model.Song{{ID: uuid.New()}}},
	}
	playlistSongs := []model.PlaylistSong{
		{SongID: uuid.New()},
		{SongID: albums[0].Songs[0].ID},
		{SongID: uuid.New()},
		{SongID: albums[1].Songs[0].ID},
	}

	// given - mocking
	playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
		Return(nil, &playlistSongs).
		Once()

	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Album)), request.AlbumIDs).
		Return(nil, &albums).
		Once()

	duplicatedAlbumIDs, duplicatedSongIDs := getAlbumAndSongDuplicates(albums, playlistSongs)
	var addedSongIDs []uuid.UUID
	for _, album := range albums {
		for _, song := range album.Songs {
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
	assert.ElementsMatch(t, res.DuplicateAlbumIDs, duplicatedAlbumIDs)
	assert.ElementsMatch(t, res.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, res.AddedSongIDs, addedSongIDs)

	playlistRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func getAlbumAndSongDuplicates(albums []model.Album, playlistSongs []model.PlaylistSong) ([]uuid.UUID, []uuid.UUID) {
	var duplicatedAlbumIDs []uuid.UUID
	var duplicatedSongIDs []uuid.UUID
	for _, album := range albums {
		var songIDs []uuid.UUID
		for _, song := range album.Songs {
			for _, ps := range playlistSongs {
				if song.ID == ps.SongID {
					songIDs = append(songIDs, song.ID)
				}
			}
		}
		duplicatedSongIDs = append(duplicatedSongIDs, songIDs...)
		if len(songIDs) == len(album.Songs) {
			duplicatedAlbumIDs = append(duplicatedAlbumIDs, album.ID)
		}
	}

	return duplicatedAlbumIDs, duplicatedSongIDs
}
