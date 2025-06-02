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

func TestAddAlbumsToPlaylist_WhenWithoutDuplicatesButWithForceAdd_ShouldReturnBadRequest(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[1].ID,
		},
		ForceAdd: &[]bool{false}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldAddAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[0].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddAlbumsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Success)
	assert.Empty(t, response.DuplicateAlbumIDs)
	assert.ElementsMatch(t, response.DuplicateSongIDs, []uuid.UUID{playlistData.Songs[0].ID}) // explicit
	assert.Empty(t, response.AddedSongIDs)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_AndOneWholeAlbumIsAlreadyOnPlaylist_ShouldAddAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	album := playlistData.Albums[0]  // whole album
	album2 := playlistData.Albums[1] // only partial

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[0].ID,
		AlbumIDs: []uuid.UUID{
			album.ID,
			album2.ID,
		},
	}

	db := utils.GetDatabase(t)
	db.Preload("Songs").Find(&album, album.ID)
	db.Preload("Songs").Find(&album2, album2.ID)

	var expectedSongIDs []uuid.UUID
	for _, song := range album.Songs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}
	// other songs from album2
	expectedSongIDs = append(expectedSongIDs, playlistData.Songs[2].ID, playlistData.Songs[3].ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddAlbumsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Success)
	assert.Len(t, response.DuplicateAlbumIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, expectedSongIDs)
	assert.Empty(t, response.AddedSongIDs)
}

func TestAddAlbumsToPlaylist_WhenWithoutDuplicatesNorWithoutForceAdd_ShouldAddAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[1].ID,
		},
	}

	albumSongs := getAlbumSongs(t, request.AlbumIDs)
	var expectedSongIDs []uuid.UUID
	for _, song := range albumSongs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddAlbumsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateAlbumIDs)
	assert.Empty(t, response.DuplicateSongIDs)
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	assertAlbumsAddedToPlaylist(t, request, albumSongs)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddDuplicatesAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[0].ID,
			playlistData.Albums[1].ID,
		},
		ForceAdd: &[]bool{true}[0],
	}

	albumSongs := getAlbumSongs(t, request.AlbumIDs)
	var expectedSongIDs []uuid.UUID
	for _, song := range albumSongs {
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddAlbumsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateAlbumIDs)
	assert.Len(t, response.DuplicateSongIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, []uuid.UUID{playlistData.Songs[0].ID}) // explicit
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	assertAlbumsAddedToPlaylist(t, request, albumSongs)
}

func TestAddAlbumsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldSkipDuplicatesAlbumSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddAlbumsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		AlbumIDs: []uuid.UUID{
			playlistData.Albums[0].ID,
			playlistData.Albums[1].ID,
		},
		ForceAdd: &[]bool{false}[0],
	}

	albumSongs := getAlbumSongs(t, request.AlbumIDs)
	duplicatedSongIDs := []uuid.UUID{playlistData.Songs[0].ID}
	var expectedSongIDs []uuid.UUID
	for _, song := range albumSongs {
		if slices.Contains(duplicatedSongIDs, song.ID) {
			continue
		}
		expectedSongIDs = append(expectedSongIDs, song.ID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/add-albums", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddAlbumsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.DuplicateAlbumIDs)
	assert.Len(t, response.DuplicateSongIDs, 1)
	assert.ElementsMatch(t, response.DuplicateSongIDs, duplicatedSongIDs)
	assert.ElementsMatch(t, response.AddedSongIDs, expectedSongIDs)

	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.Songs), len(expectedSongIDs))

	sizeDiff := len(playlist.Songs) - len(expectedSongIDs)
	for i := 0; i < len(playlist.Songs)-sizeDiff; i++ {
		assert.Equal(t, expectedSongIDs[i], playlist.Songs[i+sizeDiff].ID)
	}

	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i+1), song.PlaylistTrackNo)
	}
}

// the albums from the test data do not have the songs attached
// thus the database query
func getAlbumSongs(t *testing.T, albumIDs []uuid.UUID) []model.Song {
	db := utils.GetDatabase(t)
	var albumSongs []model.Song
	for _, albumID := range albumIDs {
		var songs []model.Song
		db.Model(model.Song{}).Where("album_id = ?", albumID).Find(&songs)
		albumSongs = append(albumSongs, songs...)
	}
	return albumSongs
}

func assertAlbumsAddedToPlaylist(t *testing.T, request requests.AddAlbumsToPlaylistRequest, albumSongs []model.Song) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.Songs), len(albumSongs))

	sizeDiff := len(playlist.Songs) - len(albumSongs)
	for i := 0; i < len(playlist.Songs)-sizeDiff; i++ {
		assert.Equal(t, albumSongs[i].ID, playlist.Songs[i+sizeDiff].ID)
	}

	for i, song := range playlist.Songs {
		assert.Equal(t, uint(i+1), song.PlaylistTrackNo)
	}
}
