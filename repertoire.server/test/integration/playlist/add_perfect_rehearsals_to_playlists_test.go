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
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectRehearsalsToPlaylists_WhenSuccessful_ShouldUpdateSongsAndPartsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{
			playlistData.Playlists[0].ID,
			playlistData.Playlists[1].ID,
			playlistData.Playlists[2].ID,
		},
	}

	getPlaylistsQuery := func(db *gorm.DB, playlists *[]model.Playlist) {
		db.Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB { return db.Order("song_track_no") }).
			Preload("PlaylistSongs.Song").
			Preload("PlaylistSongs.Song.Parts", func(db *gorm.DB) *gorm.DB {
				return db.Order("song_parts.song_order")
			}).
			Preload("PlaylistSongs.Song.Parts.History", func(db *gorm.DB) *gorm.DB {
				return db.
					Where(&model.SongPartHistory{Property: model.RehearsalsProperty}).
					Order("created_at desc")
			}).
			Preload("PlaylistSongs.Song.Parts.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Find(&playlists, request.IDs)
	}

	var playlists []model.Playlist
	db := utils.GetDatabase(t)
	getPlaylistsQuery(db, &playlists)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newPlaylists []model.Playlist
	db = db.Session(&gorm.Session{NewDB: true})
	getPlaylistsQuery(db, &newPlaylists)

	songDuplicatesMap := make(map[uuid.UUID]int)
	for _, playlist := range newPlaylists {
		for _, playlistSong := range playlist.PlaylistSongs {
			_, ok := songDuplicatesMap[playlistSong.SongID]
			if ok {
				songDuplicatesMap[playlistSong.SongID] += 1
			} else {
				songDuplicatesMap[playlistSong.SongID] = 0
			}
		}
	}

	for i, playlist := range newPlaylists {
		for j := range playlist.PlaylistSongs {
			assertion.PerfectSongRehearsalWithDuplicates(
				t,
				playlists[i].PlaylistSongs[j].Song,
				newPlaylists[i].PlaylistSongs[j].Song,
				songDuplicatesMap[playlists[i].PlaylistSongs[j].SongID],
			)
		}
	}
}
