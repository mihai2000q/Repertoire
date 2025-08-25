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

	albums := []model.Album{
		albumData.Albums[0],
		albumData.Albums[1],
	}
	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{
			albums[0].ID,
			albums[1].ID,
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/albums/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newAlbums []model.Album
	db := utils.GetDatabase(t)
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB { return db.Order("songs.album_track_no") }).
		Preload("Songs.Sections").
		Preload("Songs.Sections.History").
		Find(&newAlbums, request.IDs)

	for i, album := range newAlbums {
		for j := range album.Songs {
			totalOccurrences := uint(0)
			for _, section := range newAlbums[i].Songs[j].Sections {
				totalOccurrences += section.Occurrences
			}

			if totalOccurrences > 0 {
				assertion.PerfectSongRehearsal(t, albums[i].Songs[j], newAlbums[i].Songs[j])
			} else {
				assert.Equal(t, newAlbums[i].Songs[j].Rehearsals, albums[i].Songs[j].Rehearsals)
				assert.Equal(t, newAlbums[i].Songs[j].Progress, albums[i].Songs[j].Progress)
				assert.Equal(t, newAlbums[i].Songs[j].LastTimePlayed, albums[i].Songs[j].LastTimePlayed)
			}
		}
	}
}
