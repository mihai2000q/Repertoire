package playlist

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
)

func TestAddArtistsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldReturnNoSuccess(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		ArtistIDs: []uuid.UUID{
			playlistData.Artists[0].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddArtistsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Success)
	assert.Empty(t, response.DuplicateArtistIDs)
	assert.ElementsMatch(t, response.DuplicateSongIDs, []uuid.UUID{playlistData.Songs[0].ID}) // explicit
	assert.Empty(t, response.AddedSongIDs)
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_AndOneWholeArtistIsAlreadyOnPlaylist_ShouldReturnNoSuccess(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	artist := playlistData.Artists[0]  // whole artist
	artist2 := playlistData.Artists[1] // only partial

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[0].ID,
		ArtistIDs: []uuid.UUID{
			artist.ID,
			artist2.ID,
		},
	}

	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&artist, artist.ID)
	db.Preload("Songs").Find(&artist2, artist2.ID)

	var expectedSongIDs []uuid.UUID
	for _, song := range artist.Songs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}
	// other songs from artist2
	expectedSongIDs = append(expectedSongIDs, playlistData.Songs[2].ID, playlistData.Songs[3].ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddArtistsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Success)
	assert.Len(t, response.DuplicateArtistIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, expectedSongIDs)
	assert.Empty(t, response.AddedSongIDs)
}

func TestAddArtistsToPlaylist_WhenWithoutDuplicatesNorWithoutForceAdd_ShouldAddArtistSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		ArtistIDs: []uuid.UUID{
			playlistData.Artists[1].ID,
		},
	}

	artistSongs := getArtistSongs(t, request.ArtistIDs)
	var expectedSongIDs []uuid.UUID
	for _, song := range artistSongs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddArtistsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateArtistIDs)
	assert.Empty(t, response.DuplicateSongIDs)
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	assertArtistsAddedToPlaylist(t, request, artistSongs)
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddDuplicatesSongsTooOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		ArtistIDs: []uuid.UUID{
			playlistData.Artists[0].ID,
			playlistData.Artists[1].ID,
		},
		ForceAdd: &[]bool{true}[0],
	}

	artistSongs := getArtistSongs(t, request.ArtistIDs)
	var expectedSongIDs []uuid.UUID
	for _, song := range artistSongs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddArtistsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateArtistIDs)
	assert.Len(t, response.DuplicateSongIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, []uuid.UUID{playlistData.Songs[0].ID}) // explicit
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	assertArtistsAddedToPlaylist(t, request, artistSongs)
}

func TestAddArtistsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldSkipDuplicateSongsWhenAddingToPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddArtistsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		ArtistIDs: []uuid.UUID{
			playlistData.Artists[0].ID,
			playlistData.Artists[1].ID,
		},
		ForceAdd: &[]bool{false}[0],
	}

	artistSongs := getArtistSongs(t, request.ArtistIDs)
	duplicatedSongIDs := []uuid.UUID{playlistData.Songs[0].ID}
	var expectedSongIDs []uuid.UUID
	for _, song := range artistSongs {
		if slices.Contains(duplicatedSongIDs, song.ID) {
			continue
		}
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-artists", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddArtistsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateArtistIDs)
	assert.Len(t, response.DuplicateSongIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.PlaylistSongs), len(expectedSongIDs))

	sizeDiff := len(playlist.PlaylistSongs) - len(expectedSongIDs)
	for i := 0; i < len(playlist.PlaylistSongs)-sizeDiff; i++ {
		assert.Equal(t, expectedSongIDs[i], playlist.PlaylistSongs[i+sizeDiff].SongID)
	}

	for i, song := range playlist.PlaylistSongs {
		assert.Equal(t, uint(i+1), song.SongTrackNo)
	}
}

// the artists from the test data do not have the songs attached
// thus the database query
func getArtistSongs(t *testing.T, artistIDs []uuid.UUID) []model.Song {
	db := utils.GetDatabase(t)
	var artistSongs []model.Song
	for _, artistID := range artistIDs {
		var songs []model.Song
		db.Model(model.Song{}).Where("artist_id = ?", artistID).Find(&songs)
		artistSongs = append(artistSongs, songs...)
	}
	return artistSongs
}

func assertArtistsAddedToPlaylist(t *testing.T, request requests.AddArtistsToPlaylistRequest, artistSongs []model.Song) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.PlaylistSongs), len(artistSongs))

	sizeDiff := len(playlist.PlaylistSongs) - len(artistSongs)
	for i := 0; i < len(playlist.PlaylistSongs)-sizeDiff; i++ {
		assert.Equal(t, artistSongs[i].ID, playlist.PlaylistSongs[i+sizeDiff].SongID)
	}

	for i, song := range playlist.PlaylistSongs {
		assert.Equal(t, uint(i+1), song.SongTrackNo)
	}
}
