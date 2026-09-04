package arrangement

import (
	"encoding/json"
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
	var response struct{ ID uuid.UUID }
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)

	db = db.Session(&gorm.Session{NewDB: true})

	var arrangement model.SongArrangement
	db.
		Preload("PartOccurrences", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("LEFT JOIN song_parts ON song_parts.id = song_part_occurrences.part_id").
				Order("song_parts.song_order")
		}).
		Order("\"order\"").
		Find(&arrangement, response.ID)

	var songParts []model.SongPart
	db.Where(&model.SongPart{SongID: song.ID}).Order("song_order").Find(&songParts)

	assertCreatedSongArrangement(t, arrangement, request, arrangementsCount, songParts)
}

func assertCreatedSongArrangement(
	t *testing.T,
	songArrangement model.SongArrangement,
	request requests.CreateSongArrangementRequest,
	order int64,
	songParts []model.SongPart,
) {
	assert.NotEmpty(t, songArrangement.ID)
	assert.Equal(t, request.SongID, songArrangement.SongID)
	assert.Equal(t, request.Name, songArrangement.Name)
	assert.Equal(t, uint(order), songArrangement.Order)
	assert.Len(t, songArrangement.PartOccurrences, len(songParts))
	for i, part := range songParts {
		assert.Equal(t, part.ID, songArrangement.PartOccurrences[i].PartID)
		assert.Equal(t, songArrangement.ID, songArrangement.PartOccurrences[i].ArrangementID)
		assert.Zero(t, songArrangement.PartOccurrences[i].Occurrences)
	}
}
