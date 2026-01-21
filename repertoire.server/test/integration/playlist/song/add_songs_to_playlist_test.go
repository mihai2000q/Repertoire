package song

import (
	"encoding/json"
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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddSongsToPlaylist_WhenWithDuplicatesButWithoutForceAdd_ShouldReturnNoSuccess(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddSongsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		SongIDs: []uuid.UUID{
			playlistData.Songs[0].ID,
			playlistData.Songs[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/add", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response responses.AddSongsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.False(t, response.Success)
	assert.Len(t, response.Duplicates, 1)
	assert.ElementsMatch(t, response.Duplicates, []uuid.UUID{request.SongIDs[0]})
	assert.Empty(t, response.Added)
}

func TestAddSongsToPlaylist_WhenWithoutDuplicatesNorForceAdd_ShouldAddSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddSongsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		SongIDs: []uuid.UUID{
			playlistData.Songs[1].ID,
			playlistData.Songs[2].ID,
			playlistData.Songs[3].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/add", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	var response responses.AddSongsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Empty(t, response.Duplicates)
	assert.ElementsMatch(t, response.Added, request.SongIDs)

	assertAddedSongsToPlaylist(t, request)
}

func TestAddSongsToPlaylist_WhenWithDuplicatesAndForceAddTrue_ShouldAddAllSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddSongsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		SongIDs: []uuid.UUID{
			playlistData.Songs[1].ID,
			playlistData.Songs[0].ID,
			playlistData.Songs[2].ID,
			playlistData.Songs[3].ID,
		},
		ForceAdd: &[]bool{true}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/add", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	var response responses.AddSongsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Duplicates, 1)
	assert.ElementsMatch(t, response.Duplicates, []uuid.UUID{request.SongIDs[1]})
	assert.ElementsMatch(t, response.Added, request.SongIDs)

	assertAddedSongsToPlaylist(t, request)
}

func TestAddSongsToPlaylist_WhenWithDuplicatesAndForceAddFalse_ShouldAddOnlyNonDuplicatedSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddSongsToPlaylistRequest{
		ID: playlistData.Playlists[1].ID,
		SongIDs: []uuid.UUID{
			playlistData.Songs[1].ID,
			playlistData.Songs[0].ID,
			playlistData.Songs[2].ID,
			playlistData.Songs[3].ID,
		},
		ForceAdd: &[]bool{false}[0],
	}

	duplicatedSongIDs := []uuid.UUID{request.SongIDs[1]}
	var expectedSongIDs []uuid.UUID
	for _, songID := range request.SongIDs {
		if slices.Contains(duplicatedSongIDs, songID) {
			continue
		}
		expectedSongIDs = append(expectedSongIDs, songID)
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/add", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	var response responses.AddSongsToPlaylistResponse
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.True(t, response.Success)
	assert.Len(t, response.Duplicates, 1)
	assert.ElementsMatch(t, response.Duplicates, duplicatedSongIDs)
	assert.ElementsMatch(t, response.Added, expectedSongIDs)

	// assert added songs on playlist without duplicates
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

func assertAddedSongsToPlaylist(t *testing.T, request requests.AddSongsToPlaylistRequest) {
	db := utils.GetDatabase(t)

	var playlist model.Playlist
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Song").Order("song_track_no")
	}).Find(&playlist, request.ID)

	assert.GreaterOrEqual(t, len(playlist.PlaylistSongs), len(request.SongIDs))

	sizeDiff := len(playlist.PlaylistSongs) - len(request.SongIDs)
	for i := 0; i < len(playlist.PlaylistSongs)-sizeDiff; i++ {
		assert.Equal(t, request.SongIDs[i], playlist.PlaylistSongs[i+sizeDiff].SongID)
	}

	for i, song := range playlist.PlaylistSongs {
		assert.Equal(t, uint(i+1), song.SongTrackNo)
	}
}
