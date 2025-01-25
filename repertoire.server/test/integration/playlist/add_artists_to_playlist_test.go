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

func TestAddArtistsToPlaylist_WhenSuccessful_ShouldHaveArtistSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		ArtistIDs: []uuid.UUID{
			playlistData.Artists[0].ID,
			playlistData.Artists[1].ID,
		},
	}

	// the above artists from the test data do not have the songs attached
	// thus the database query
	db := utils.GetDatabase(t)
	var artistSongs []model.Song
	for _, artistID := range request.ArtistIDs {
		var songs []model.Song
		db.Model(model.Song{}).Where("artist_id = ?", artistID).Find(&songs)
		artistSongs = append(artistSongs, songs...)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assertArtistsAddedToPlaylist(t, request, artistSongs)
}

func assertArtistsAddedToPlaylist(t *testing.T, request requests.AddArtistsToPlaylistRequest, artistSongs []model.Song) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.Songs), len(request.ArtistIDs))

	// it is not mandatory that an artist song wasn't already part of the playlist
	var matchSongs []model.Song
	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i+1), song.PlaylistTrackNo) // check correct order
		if slices.ContainsFunc(artistSongs, func(s model.Song) bool {
			return s.ID == song.ID
		}) {
			matchSongs = append(matchSongs, song)
		}
	}
	assert.Len(t, matchSongs, len(artistSongs)) // check if all the artist songs were added
}
