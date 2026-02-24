package song

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

func TestAddPerfectPlaylistSongRehearsals_WhenPlaylistsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TODO: Fail Safe - thorough validation would demand this to be a Bad Request instead
func TestAddPerfectPlaylistSongRehearsals_WhenNotAllSongsAreFromTheSamePlaylist_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: playlistData.PlaylistsSongs[0].PlaylistID,
		IDs: []uuid.UUID{
			playlistData.PlaylistsSongs[1].ID,
			playlistData.PlaylistsSongs[5].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectPlaylistSongRehearsals_WhenSuccessful_ShouldUpdateSongsAndSectionsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: playlistData.PlaylistsSongs[0].PlaylistID,
		IDs: []uuid.UUID{
			playlistData.PlaylistsSongs[1].ID,
			playlistData.PlaylistsSongs[0].ID,
			playlistData.PlaylistsSongs[2].ID,
			playlistData.PlaylistsSongs[4].ID,
		},
	}

	getPlaylistSongsQuery := func(db *gorm.DB, playlistSongs *[]model.PlaylistSong) {
		db.
			Preload("Song").
			Preload("Song.Sections", func(db *gorm.DB) *gorm.DB {
				return db.Order("song_sections.order")
			}).
			Preload("Song.Sections.History", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at desc, property desc")
			}).
			Preload("Song.Sections.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Order("song_track_no").
			Find(&playlistSongs, request.IDs)
	}

	var playlistSongs []model.PlaylistSong
	db := utils.GetDatabase(t)
	getPlaylistSongsQuery(db, &playlistSongs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newPlaylistSongs []model.PlaylistSong
	db = db.Session(&gorm.Session{NewDB: true})
	getPlaylistSongsQuery(db, &newPlaylistSongs)

	songDuplicatesMap := make(map[uuid.UUID]int)
	for _, playlistSong := range newPlaylistSongs {
		_, ok := songDuplicatesMap[playlistSong.SongID]
		if ok {
			songDuplicatesMap[playlistSong.SongID] += 1
		} else {
			songDuplicatesMap[playlistSong.SongID] = 0
		}
	}

	for i := range newPlaylistSongs {
		assertion.PerfectSongRehearsalWithDuplicates(
			t,
			playlistSongs[i].Song,
			newPlaylistSongs[i].Song,
			songDuplicatesMap[playlistSongs[i].SongID],
		)
	}
}
