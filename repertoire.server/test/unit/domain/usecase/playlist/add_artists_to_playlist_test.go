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
	errCode := _uut.Handle(request)

	// then
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
	errCode := _uut.Handle(request)

	// then
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
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestAddArtistToPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	mutualId := uuid.New()

	tests := []struct {
		name          string
		playlistSongs []model.PlaylistSong
		artists       []model.Artist
	}{
		{
			"Use Case 1",
			[]model.PlaylistSong{
				{SongID: uuid.New()}, {SongID: mutualId},
			},
			[]model.Artist{
				{
					ID:    uuid.New(),
					Songs: []model.Song{{ID: uuid.New()}, {ID: uuid.New()}},
				},
				{
					ID:    uuid.New(),
					Songs: []model.Song{{ID: mutualId}, {ID: uuid.New()}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			artistRepository := new(repository.ArtistRepositoryMock)
			_uut := playlist.NewAddArtistsToPlaylist(playlistRepository, artistRepository)

			request := requests.AddArtistsToPlaylistRequest{
				ID:        uuid.New(),
				ArtistIDs: []uuid.UUID{uuid.New(), uuid.New()},
			}

			// given - mocking
			playlistRepository.On("GetPlaylistSongs", mock.IsType(new([]model.PlaylistSong)), request.ID).
				Return(nil, &tt.playlistSongs).
				Once()

			artistRepository.On("GetAllByIDsWithSongs", mock.IsType(new([]model.Artist)), request.ArtistIDs).
				Return(nil, &tt.artists).
				Once()

			var newSongs []model.Song
			for _, artist := range tt.artists {
				newSongs = append(newSongs, artist.Songs...)
			}
			newSongs = slices.DeleteFunc(newSongs, func(s model.Song) bool {
				return s.ID == mutualId
			})

			oldPlaylistSongsLength := len(tt.playlistSongs) + 1
			playlistRepository.On("AddSongs", mock.IsType(new([]model.PlaylistSong))).
				Run(func(args mock.Arguments) {
					newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
					assert.Len(t, *newPlaylistSongs, len(newSongs))
					for i, playlistSong := range *newPlaylistSongs {
						assert.Equal(t, uint(oldPlaylistSongsLength+i), playlistSong.SongTrackNo)
						assert.Equal(t, request.ID, playlistSong.PlaylistID)
						assert.Equal(t, newSongs[i].ID, playlistSong.SongID)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			playlistRepository.AssertExpectations(t)
			artistRepository.AssertExpectations(t)
		})
	}
}
