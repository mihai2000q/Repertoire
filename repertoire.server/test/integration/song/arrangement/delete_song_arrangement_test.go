package arrangement

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteSongArrangement_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/arrangements/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongArrangement_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/arrangements/"+uuid.New().String()+"/from/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongArrangement_WhenSuccessful_ShouldDeleteArrangement(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with arrangements
	arrangement := songData.SongArrangements[2]
	songID := arrangement.SongID

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/arrangements/"+arrangement.ID.String()+"/from/"+songID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Preload("Arrangements", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&newSong, songID)

	assert.True(t,
		slices.IndexFunc(newSong.Arrangements, func(t model.SongArrangement) bool {
			return t.ID == arrangement.ID
		}) == -1,
		"Song Arrangement has not been deleted",
	)

	for i, s := range newSong.Arrangements {
		assert.Equal(t, uint(i), s.Order)
	}
}
