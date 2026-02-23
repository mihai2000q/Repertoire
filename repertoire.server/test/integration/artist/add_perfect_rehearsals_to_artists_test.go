package artist

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddPerfectRehearsalsToArtists_WhenArtistsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddPerfectRehearsalsToArtists_WhenSuccessful_ShouldUpdateSongsAndSectionsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{
			artistData.Artists[1].ID,
			artistData.Artists[2].ID,
		},
	}
	
	getArtistsQuery := func(db *gorm.DB, artists *[]model.Artist) {
		db.Preload("Songs", func(db *gorm.DB) *gorm.DB { return db.Order("songs.title") }).
			Preload("Songs.Sections", func(db *gorm.DB) *gorm.DB { return db.Order("song_sections.order") }).
			Preload("Songs.Sections.History", func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") }).
			Preload("Songs.Sections.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_sections ON song_sections.id = section_id").
					Joins("LEFT JOIN songs ON songs.id = song_sections.song_id").
					Where("arrangement_id = default_arrangement_id")
			}).
			Find(&artists, request.IDs)
	}

	var artists []model.Artist
	db := utils.GetDatabase(t)
	getArtistsQuery(db, &artists)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/artists/perfect-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newArtists []model.Artist
	db = db.Session(&gorm.Session{NewDB: true})
	getArtistsQuery(db, &newArtists)

	for i, artist := range newArtists {
		for j := range artist.Songs {
			totalOccurrences := uint(0)
			for _, section := range newArtists[i].Songs[j].Sections {
				totalOccurrences += section.ArrangementOccurrences[0].Occurrences
			}

			if totalOccurrences > 0 {
				assertion.PerfectSongRehearsal(t, artists[i].Songs[j], newArtists[i].Songs[j])
			} else {
				assert.Equal(t, artists[i].Songs[j].Rehearsals, newArtists[i].Songs[j].Rehearsals)
				assert.Equal(t, artists[i].Songs[j].Progress, newArtists[i].Songs[j].Progress)
				assert.Equal(t, artists[i].Songs[j].LastTimePlayed, newArtists[i].Songs[j].LastTimePlayed)
			}
		}
	}
}
