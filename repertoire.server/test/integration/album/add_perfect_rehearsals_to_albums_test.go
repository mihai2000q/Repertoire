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
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectRehearsalsToAlbums_WhenSuccessful_ShouldUpdateSongsAndSectionsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{
			albumData.Albums[0].ID,
			albumData.Albums[1].ID,
		},
	}

	var albums []model.Album
	db := utils.GetDatabase(t)
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB { return db.Order("songs.album_track_no") }).
		Preload("Songs.Sections").
		Preload("Songs.Sections.History").
		Find(&albums, request.IDs)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newAlbums []model.Album
	db = db.Session(&gorm.Session{NewDB: true})
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB { return db.Order("songs.album_track_no") }).
		Preload("Songs.Sections", func(db *gorm.DB) *gorm.DB { return db.Order("song_sections.order") }).
		Preload("Songs.Sections.History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).
		Preload("Songs.Sections.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.Where("song_section_occurrences.arrangement_id = songs.default_arrangement_id")
		}).
		Find(&newAlbums, request.IDs)

	for i, album := range newAlbums {
		for j := range album.Songs {
			totalOccurrences := uint(0)
			for _, section := range newAlbums[i].Songs[j].Sections {
				totalOccurrences += section.ArrangementOccurrences[0].Occurrences
			}

			if totalOccurrences > 0 {
				assertion.PerfectSongRehearsal(t, albums[i].Songs[j], newAlbums[i].Songs[j])
			} else {
				assert.Equal(t, albums[i].Songs[j].Rehearsals, newAlbums[i].Songs[j].Rehearsals)
				assert.Equal(t, albums[i].Songs[j].Progress, newAlbums[i].Songs[j].Progress)
				assert.Equal(t, albums[i].Songs[j].LastTimePlayed, newAlbums[i].Songs[j].LastTimePlayed)
			}
		}
	}
}
