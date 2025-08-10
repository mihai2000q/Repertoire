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

func TestUpdateSongSectionsOccurrences_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionsOccurrencesRequest{
		SongID:   uuid.New(),
		Sections: []requests.UpdateSectionOccurrencesRequest{{ID: uuid.New()}},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/occurrences", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSectionsOccurrences_WhenSuccessful_ShouldUpdateSongSectionsOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[4]
	request := requests.UpdateSongSectionsOccurrencesRequest{
		SongID: song.ID,
		Sections: []requests.UpdateSectionOccurrencesRequest{
			{
				ID:          song.Sections[0].ID,
				Occurrences: uint(4),
			},
			{
				ID:          song.Sections[1].ID,
				Occurrences: uint(7),
			},
			{
				ID:          song.Sections[2].ID,
				Occurrences: uint(8),
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/occurrences", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&song, request.SongID)

	for i, section := range song.Sections {
		assert.Equal(t, request.Sections[i].ID, section.ID)
		assert.Equal(t, request.Sections[i].Occurrences, section.Occurrences)
	}
}
