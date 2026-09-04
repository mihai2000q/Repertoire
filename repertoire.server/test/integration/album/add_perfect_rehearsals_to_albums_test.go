package album

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddPerfectRehearsalsToAlbums_WhenAlbumsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectRehearsalsToAlbums_WhenSuccessful_ShouldUpdateSongsAndPartsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{
			albumData.Albums[0].ID,
			albumData.Albums[1].ID,
		},
	}

	getAlbumsQuery := func(db *gorm.DB, albums *[]model.Album) {
		db.Preload("Songs", func(db *gorm.DB) *gorm.DB { return db.Order("songs.album_track_no") }).
			Preload("Songs.Parts", func(db *gorm.DB) *gorm.DB { return db.Order("song_parts.song_order") }).
			Preload("Songs.Parts.History", func(db *gorm.DB) *gorm.DB {
				return db.
					Where(&model.SongPartHistory{Property: model.RehearsalsProperty}).
					Order("created_at desc")
			}).
			Preload("Songs.Parts.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Find(&albums, request.IDs)
	}

	var albums []model.Album
	db := utils.GetDatabase(t)
	getAlbumsQuery(db, &albums)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newAlbums []model.Album
	db = db.Session(&gorm.Session{NewDB: true})
	getAlbumsQuery(db, &newAlbums)

	for i, album := range newAlbums {
		for j := range album.Songs {
			assertion.PerfectSongRehearsal(t, albums[i].Songs[j], newAlbums[i].Songs[j])
		}
	}
}
