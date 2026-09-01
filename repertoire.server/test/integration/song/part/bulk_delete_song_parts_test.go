package part

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

func TestBulkDeleteSongParts_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkDeleteSongPartsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongParts_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkDeleteSongPartsRequest{
		IDs:    []uuid.UUID{uuid.New()},
		SongID: songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongParts_WhenSuccessful_ShouldDeleteParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with parts and previous stats
	song := songData.Songs[0]
	request := requests.BulkDeleteSongPartsRequest{
		IDs:    []uuid.UUID{songData.SongParts[0].ID, songData.SongParts[1].ID},
		SongID: song.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"song_order\"")
		}).
		Preload("Sections").
		Preload("Sections.SectionParts", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(&song, song.ID)

	for i, s := range newSong.Parts {
		assert.Equal(t, uint(i), s.SongOrder)
		assert.True(t,
			slices.IndexFunc(newSong.Parts, func(t model.SongPart) bool {
				return t.ID == s.ID
			}) == -1,
			"Song Part with id:"+s.ID.String()+", has not been deleted",
		)
	}

	assert.LessOrEqual(t, newSong.Confidence, song.Confidence)
	assert.LessOrEqual(t, newSong.Rehearsals, song.Rehearsals)
	assert.LessOrEqual(t, newSong.Progress, song.Progress)

	for _, sec := range newSong.Sections {
		for i, sp := range sec.SectionParts {
			assert.Equal(t, uint(i), sp.Order)
		}
	}
}
