package playlist

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestAddAlbumsToPlaylist_WhenSuccessful_ShouldHaveAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[0].ID,
			playlistData.Albums[1].ID,
		},
	}

	// the above albums from the test data do not have the songs attached
	// thus the database query
	db := utils.GetDatabase(t)
	var albumSongs []model.Song
	for _, albumID := range request.AlbumIDs {
		var songs []model.Song
		db.Model(model.Song{}).Where("album_id = ?", albumID).Find(&songs)
		albumSongs = append(albumSongs, songs...)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assertAlbumsAddedToPlaylist(t, request, albumSongs)
}

func assertAlbumsAddedToPlaylist(t *testing.T, request requests.AddAlbumsToPlaylistRequest, albumSongs []model.Song) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.Songs), len(request.AlbumIDs))

	// it is not mandatory that an album song wasn't already part of the playlist
	var matchSongs []model.Song
	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i+1), song.PlaylistTrackNo) // check correct order
		if slices.ContainsFunc(albumSongs, func(s model.Song) bool {
			return s.ID == song.ID
		}) {
			matchSongs = append(matchSongs, song)
		}
	}
	assert.Len(t, matchSongs, len(albumSongs)) // check if all the album songs were added
}
