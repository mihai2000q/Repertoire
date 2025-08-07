package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateSongSectionsPartialOccurrences_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID:   uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.New()}},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/partial-occurrences", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSectionsPartialOccurrences_WhenSuccessful_ShouldUpdateSongSectionsPartialOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[4]
	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID: song.ID,
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{
			{
				ID:                 song.Sections[0].ID,
				PartialOccurrences: uint(4),
			},
			{
				ID:                 song.Sections[1].ID,
				PartialOccurrences: uint(7),
			},
			{
				ID:                 song.Sections[2].ID,
				PartialOccurrences: uint(8),
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/partial-occurrences", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&song, request.SongID)

	for i, section := range song.Sections {
		assert.Equal(t, request.Sections[i].ID, section.ID)
		assert.Equal(t, request.Sections[i].PartialOccurrences, section.PartialOccurrences)
	}
}
