package playlist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddPerfectRehearsalsToPlaylists_WhenPlaylistsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectRehearsalsToPlaylists_WhenSuccessful_ShouldUpdateSongsAndSectionsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{
			playlistData.Playlists[0].ID,
			playlistData.Playlists[1].ID,
		},
	}

	var playlists []model.Playlist
	db := utils.GetDatabase(t)
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB { return db.Order("song_track_no") }).
		Preload("PlaylistSongs.Song").
		Preload("PlaylistSongs.Song.Sections").
		Preload("PlaylistSongs.Song.Sections.History").
		Find(&playlists, request.IDs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newPlaylists []model.Playlist
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB { return db.Order("song_track_no") }).
		Preload("PlaylistSongs.Song").
		Preload("PlaylistSongs.Song.Sections").
		Preload("PlaylistSongs.Song.Sections.History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).
		Find(&playlists, request.IDs)

	for i, playlist := range newPlaylists {
		for j := range playlist.PlaylistSongs {
			song := playlists[i].PlaylistSongs[j].Song
			newSong := newPlaylists[i].PlaylistSongs[j].Song
			totalOccurrences := uint(0)
			for _, section := range newSong.Sections {
				totalOccurrences += section.Occurrences
			}

			if totalOccurrences > 0 {
				assertion.PerfectSongRehearsal(t, song, newSong)
			} else {
				assert.Equal(t, song.Rehearsals, newSong.Rehearsals)
				assert.Equal(t, song.Progress, newSong.Progress)
				assert.Equal(t, song.LastTimePlayed, newSong.LastTimePlayed)
			}
		}
	}
}
