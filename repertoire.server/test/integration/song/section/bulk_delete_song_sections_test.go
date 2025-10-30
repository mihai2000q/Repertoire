package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
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

func TestBulkDeleteSongSections_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongSections_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongSections_WhenSuccessful_ShouldDeleteSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with sections and previous stats
	song := songData.Songs[0]
	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{songData.Songs[0].Sections[0].ID, songData.Songs[0].Sections[1].ID},
		SongID: songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Preload("Sections", func(db gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&song, song.ID)

	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
		assert.True(t,
			slices.IndexFunc(newSong.Sections, func(t model.SongSection) bool {
				return t.ID == s.ID
			}) == -1,
			"Song Section with id:"+s.ID.String()+", has not been deleted",
		)
	}

	assert.LessOrEqual(t, newSong.Confidence, song.Confidence)
	assert.LessOrEqual(t, newSong.Rehearsals, song.Rehearsals)
	assert.LessOrEqual(t, newSong.Progress, song.Progress)
}
