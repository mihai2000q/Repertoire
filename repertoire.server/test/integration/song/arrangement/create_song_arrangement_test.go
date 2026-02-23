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

func TestCreateSongArrangement_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "New Arrangement",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/arrangements", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongArrangement_WhenSuccessful_ShouldCreateArrangement(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with arrangements
	song := songData.Songs[4]
	request := requests.CreateSongArrangementRequest{
		SongID: song.ID,
		Name:   "New Arrangement",
	}

	db := utils.GetDatabase(t)
	var arrangementsCount int64
	db.Model(&model.SongArrangement{}).Where(&model.SongArrangement{SongID: song.ID}).Count(&arrangementsCount)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/arrangements", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db = db.Session(&gorm.Session{NewDB: true})

	var arrangement model.SongArrangement
	db.
		Preload("SectionOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_sections ON song_sections.id = song_section_occurrences.section_id").
				Order("song_sections.order")
		}).
		Order("\"order\"").
		Find(&arrangement, &model.SongArrangement{Name: request.Name})

	var songSections []model.SongSection
	db.Where(&model.SongSection{SongID: song.ID}).Order("\"order\"").Find(&songSections)

	assertCreatedSongArrangement(t, arrangement, request, arrangementsCount, songSections)
}

func assertCreatedSongArrangement(
	t *testing.T,
	songArrangement model.SongArrangement,
	request requests.CreateSongArrangementRequest,
	order int64,
	songSections []model.SongSection,
) {
	assert.NotEmpty(t, songArrangement.ID)
	assert.Equal(t, request.SongID, songArrangement.SongID)
	assert.Equal(t, request.Name, songArrangement.Name)
	assert.Equal(t, uint(order), songArrangement.Order)
	assert.Len(t, songArrangement.SectionOccurrences, len(songSections))
	for i, section := range songSections {
		assert.Equal(t, section.ID, songArrangement.SectionOccurrences[i].SectionID)
		assert.Equal(t, songArrangement.ID, songArrangement.SectionOccurrences[i].ArrangementID)
		assert.Zero(t, songArrangement.SectionOccurrences[i].Occurrences)
	}
}
