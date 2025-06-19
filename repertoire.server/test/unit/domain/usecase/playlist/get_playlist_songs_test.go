package playlist

import (
	"errors"
	"net/http"
	"os"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistSongs_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylistSongs(playlistRepository)

	request := requests.GetPlaylistSongsRequest{
		ID:      uuid.New(),
		OrderBy: []string{"title asc"},
	}

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetPlaylistSongsWithSongs",
			new([]model.PlaylistSong),
			request.ID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(internalError).
		Once()

	// when
	resultPlaylistSongs, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultPlaylistSongs)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylistSongs_WhenSuccessful_ShouldReturnPlaylist(t *testing.T) {
	_ = os.Setenv("STORAGE_FETCH_URL", "the_storage_url")

	tests := []struct {
		name    string
		request requests.GetPlaylistSongsRequest
	}{
		{
			"With Order By",
			requests.GetPlaylistSongsRequest{
				ID:      uuid.New(),
				OrderBy: []string{"ordering"},
			},
		},
		{
			"Without Order By - Will have default value",
			requests.GetPlaylistSongsRequest{ID: uuid.New()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			playlistRepository := new(repository.PlaylistRepositoryMock)
			_uut := playlist.NewGetPlaylistSongs(playlistRepository)

			expectedPlaylistSongs := []model.PlaylistSong{
				{
					Song:        model.Song{ID: uuid.New(), ImageURL: &[]internal.FilePath{"file_path"}[0]},
					SongTrackNo: 1,
					CreatedAt:   time.Now(),
				},
				{
					Song:        model.Song{ID: uuid.New()},
					SongTrackNo: 2,
				},
			}

			if len(tt.request.OrderBy) != 0 {
				playlistRepository.
					On(
						"GetPlaylistSongsWithSongs",
						new([]model.PlaylistSong),
						tt.request.ID,
						tt.request.CurrentPage,
						tt.request.PageSize,
						tt.request.OrderBy,
					).
					Return(nil, &expectedPlaylistSongs).
					Once()
			} else {
				playlistRepository.
					On(
						"GetPlaylistSongsWithSongs",
						new([]model.PlaylistSong),
						tt.request.ID,
						tt.request.CurrentPage,
						tt.request.PageSize,
						[]string{"song_track_no"},
					).
					Return(nil, &expectedPlaylistSongs).
					Once()
			}

			// when
			resultPlaylistSongs, errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)
			assert.NotEmpty(t, resultPlaylistSongs)

			for i := range resultPlaylistSongs {
				assert.Equal(t, expectedPlaylistSongs[i].Song.ID, resultPlaylistSongs[i].ID)
				if expectedPlaylistSongs[i].Song.ImageURL == nil {
					assert.Nil(t, resultPlaylistSongs[i].ImageURL)
				} else {
					assert.Equal(t, *expectedPlaylistSongs[i].Song.ImageURL.ToFullURL(), *resultPlaylistSongs[i].ImageURL)
				}
				assert.Equal(t, expectedPlaylistSongs[i].ID, resultPlaylistSongs[i].PlaylistSongID)
				assert.Equal(t, expectedPlaylistSongs[i].SongTrackNo, resultPlaylistSongs[i].PlaylistTrackNo)
				assert.Equal(t, expectedPlaylistSongs[i].CreatedAt, resultPlaylistSongs[i].PlaylistCreatedAt)
			}

			playlistRepository.AssertExpectations(t)
		})
	}
}
