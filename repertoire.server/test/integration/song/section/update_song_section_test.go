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

func TestUpdateSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "New Chorus Name",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSection_WhenPartsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:      songData.SongSections[2].ID,
		Name:    "New Chorus Name",
		TypeID:  songData.Users[0].SongSectionTypes[0].ID,
		PartIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongSection_WhenPartDoesNotBelongToSong_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:      songData.SongSections[2].ID,
		Name:    "New Chorus Name",
		TypeID:  songData.Users[0].SongSectionTypes[0].ID,
		PartIDs: []uuid.UUID{songData.SongParts[4].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateSongSection_WhenSuccessful_ShouldUpdateSection(t *testing.T) {
	tests := []struct {
		name    string
		request requests.UpdateSongSectionRequest
	}{
		{
			"Without Parts - keep",
			requests.UpdateSongSectionRequest{
				ID:      songData.SongSections[0].ID,
				Name:    "New Chorus Name",
				TypeID:  songData.Users[0].SongSectionTypes[0].ID,
				PartIDs: []uuid.UUID{songData.SongParts[0].ID, songData.SongParts[1].ID},
			},
		},
		{
			"With Parts Changed",
			requests.UpdateSongSectionRequest{
				ID:      songData.SongSections[0].ID,
				Name:    "New Chorus Name",
				TypeID:  songData.Users[0].SongSectionTypes[0].ID,
				PartIDs: []uuid.UUID{songData.SongParts[1].ID, songData.SongParts[2].ID},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/sections", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var section model.SongSection
			db.
				Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
					return db.Order("\"order\"")
				}).
				Find(&section, &model.SongSection{ID: test.request.ID})

			assertUpdatedSongSection(t, section, test.request)
		})
	}
}

func assertUpdatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.UpdateSongSectionRequest,
) {
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)

	assert.Len(t, songSection.SectionParts, len(request.PartIDs))
	for i, sectionPart := range songSection.SectionParts {
		assert.Equal(t, songSection.ID, sectionPart.SectionID)
		assert.True(t,
			slices.Contains(request.PartIDs, sectionPart.PartID),
			"parts have not been updated",
		)
		assert.Equal(t, uint(i), sectionPart.Order)
	}
}
