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
)

func TestGetPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylist(playlistRepository)

	request := requests.GetPlaylistRequest{
		ID:           uuid.New(),
		SongsOrderBy: []string{"ordering"},
	}

	internalError := errors.New("internal error")
	playlistRepository.On("GetWithAssociations", new(model.Playlist), request.ID, request.SongsOrderBy).
		Return(internalError).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultPlaylist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylist(playlistRepository)

	request := requests.GetPlaylistRequest{
		ID:           uuid.New(),
		SongsOrderBy: []string{"ordering"},
	}

	playlistRepository.On("GetWithAssociations", new(model.Playlist), request.ID, request.SongsOrderBy).
		Return(nil).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultPlaylist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylist_WhenSuccessful_ShouldReturnPlaylist(t *testing.T) {
	tests := []struct {
		name    string
		request requests.GetPlaylistRequest
	}{
		{
			"With Songs Order By",
			requests.GetPlaylistRequest{
				ID:           uuid.New(),
				SongsOrderBy: []string{"ordering"},
			},
		},
		{
			"Without Songs Order By - Will have default value",
			requests.GetPlaylistRequest{ID: uuid.New()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			_uut := playlist.NewGetPlaylist(playlistRepository)

			expectedPlaylist := &model.Playlist{
				ID:    tt.request.ID,
				Title: "Some Playlist",
				PlaylistSongs: []model.PlaylistSong{
					{
						Song:        model.Song{ID: uuid.New()},
						SongTrackNo: 1,
					},
					{
						Song:        model.Song{ID: uuid.New()},
						SongTrackNo: 2,
					},
				},
			}

			if len(tt.request.SongsOrderBy) != 0 {
				playlistRepository.
					On(
						"GetWithAssociations",
						new(model.Playlist),
						tt.request.ID,
						tt.request.SongsOrderBy,
					).
					Return(nil, expectedPlaylist).
					Once()
			} else {
				playlistRepository.
					On(
						"GetWithAssociations",
						new(model.Playlist),
						tt.request.ID,
						[]string{"song_track_no"},
					).
					Return(nil, expectedPlaylist).
					Once()
			}

			// when
			resultPlaylist, errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)
			assert.NotEmpty(t, resultPlaylist)

			assert.Equal(t, expectedPlaylist.ID, resultPlaylist.ID)
			assert.Equal(t, expectedPlaylist.Title, resultPlaylist.Title)
			assert.Equal(t, expectedPlaylist.Description, resultPlaylist.Description)
			assert.Equal(t, expectedPlaylist.ImageURL, resultPlaylist.ImageURL)

			for i := range resultPlaylist.Songs {
				assert.Equal(t, expectedPlaylist.PlaylistSongs[i].Song.ID, resultPlaylist.Songs[i].ID)
				assert.Equal(t, expectedPlaylist.PlaylistSongs[i].SongTrackNo, resultPlaylist.Songs[i].PlaylistTrackNo)
				assert.Equal(t, expectedPlaylist.PlaylistSongs[i].CreatedAt, resultPlaylist.Songs[i].PlaylistCreatedAt)
			}

			playlistRepository.AssertExpectations(t)
		})
	}
}
