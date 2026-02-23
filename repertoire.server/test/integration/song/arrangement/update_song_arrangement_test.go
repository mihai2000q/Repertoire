package arrangement

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUpdateSongArrangement_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: "New Chorus Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongArrangement_WhenSuccessful_ShouldUpdateArrangement(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongArrangementRequest{
		ID:   songData.SongArrangements[0].ID,
		Name: "New Chorus Name",
		Occurrences: []requests.UpdateSongSectionOccurrencesRequest{
			{
				SectionID:   songData.SongSections[5].ID,
				Occurrences: 1,
			},
			{
				SectionID:   songData.SongSections[4].ID,
				Occurrences: 7,
			},
			{
				SectionID:   songData.SongSections[6].ID,
				Occurrences: 2,
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var arrangement model.SongArrangement
	db.
		Preload("SectionOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_sections ON song_sections.id = song_section_occurrences.section_id").
				Order("song_sections.order")
		}).
		Find(&arrangement, &model.SongArrangement{ID: request.ID})

	assertUpdatedSongArrangement(t, arrangement, request)
}

func assertUpdatedSongArrangement(
	t *testing.T,
	songArrangement model.SongArrangement,
	request requests.UpdateSongArrangementRequest,
) {
	assert.Equal(t, request.Name, songArrangement.Name)

	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Occurrences {
		sectionsOccurrencesMap[s.SectionID] = s.Occurrences
	}
	for i := range songArrangement.SectionOccurrences {
		occurrences := sectionsOccurrencesMap[songArrangement.SectionOccurrences[i].SectionID]
		assert.Equal(t, songArrangement.SectionOccurrences[i].Occurrences, occurrences)
	}
}
