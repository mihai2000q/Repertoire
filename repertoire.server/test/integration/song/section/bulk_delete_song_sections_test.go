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
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk-delete", request)

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
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkDeleteSongSections_WhenSuccessful_ShouldDeleteSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with sections and previous stats
	song := songData.Songs[0]
	request := requests.BulkDeleteSongSectionsRequest{
		IDs:    []uuid.UUID{songData.SongSections[0].ID, songData.SongSections[1].ID},
		SongID: song.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).Find(&newSong, song.ID)

	for _, id := range request.IDs {
		assert.False(t,
			slices.ContainsFunc(newSong.Sections, func(s model.SongSection) bool {
				return s.ID == id
			}),
			"song section %v should have been deleted", id,
		)
	}

	// sections reordered
	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
	}
}

func TestBulkDeleteSongSections_WhenSuccessfulWithParts_ShouldDeleteSectionsAndParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.BulkDeleteSongSectionsRequest{
		IDs:       []uuid.UUID{songData.SongSections[0].ID, songData.SongSections[1].ID},
		SongID:    song.ID,
		WithParts: true,
	}

	// Collect unique part IDs that are expected to be deleted
	partIDSet := make(map[uuid.UUID]bool)
	expectedDeletedPartIDs := make([]uuid.UUID, 0)
	for _, sp := range songData.SongSectionParts {
		if slices.Contains(request.IDs, sp.SectionID) && !partIDSet[sp.PartID] {
			partIDSet[sp.PartID] = true
			expectedDeletedPartIDs = append(expectedDeletedPartIDs, sp.PartID)
		}
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections/bulk-delete", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_order")
		}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(&newSong, song.ID)

	for _, id := range request.IDs {
		assert.False(t,
			slices.ContainsFunc(newSong.Sections, func(s model.SongSection) bool {
				return s.ID == id
			}),
			"song section %v should have been deleted", id,
		)
	}

	// sections reordered
	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
	}

	// Verify parts are deleted
	for _, partID := range expectedDeletedPartIDs {
		assert.True(t,
			slices.IndexFunc(newSong.Parts, func(p model.SongPart) bool {
				return p.ID == partID
			}) == -1,
			"song part %v should have been deleted", partID,
		)
	}

	// parts reordered
	for i, s := range newSong.Parts {
		assert.Equal(t, uint(i), s.SongOrder)
	}

	// Stats should decrease because parts are deleted
	assert.Less(t, newSong.Confidence, song.Confidence)
	assert.Less(t, newSong.Rehearsals, song.Rehearsals)
	assert.Less(t, newSong.Progress, song.Progress)
}
