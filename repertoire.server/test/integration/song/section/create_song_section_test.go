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

func TestCreateSongSection_WhenTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenPartsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID:  songData.Songs[0].ID,
		Name:    "Chorus 1-New",
		TypeID:  songData.Users[0].SongSectionTypes[0].ID,
		PartIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenSuccessful_ShouldCreateSection(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongSectionRequest
	}{
		{
			"Without Parts",
			requests.CreateSongSectionRequest{
				SongID: songData.Songs[0].ID,
				Name:   "Chorus 1-New",
				TypeID: songData.Users[0].SongSectionTypes[0].ID,
			},
		},
		{
			"With Parts",
			requests.CreateSongSectionRequest{
				SongID:  songData.Songs[0].ID,
				Name:    "Chorus 1-New",
				TypeID:  songData.Users[0].SongSectionTypes[0].ID,
				PartIDs: []uuid.UUID{songData.SongParts[0].ID},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			sectionsCount := len(slices.Clone(slices.DeleteFunc(slices.Clone(songData.SongParts), func(part model.SongPart) bool {
				return part.SongID != test.request.SongID
			})))

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/sections", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)
			var section model.SongSection
			db.
				Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
					return db.Order("\"order\"")
				}).
				Find(&section, &model.SongSection{Name: test.request.Name})

			assertCreatedSongSection(t, section, test.request, sectionsCount)
		})
	}
}

func assertCreatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.CreateSongSectionRequest,
	order int,
) {
	assert.NotEmpty(t, songSection.ID)
	assert.Equal(t, request.SongID, songSection.SongID)
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)
	assert.Equal(t, uint(order), songSection.Order)

	for i, sectionPart := range songSection.SectionParts {
		assert.Equal(t, request.PartIDs[i], sectionPart.PartID)
		assert.Equal(t, songSection.ID, sectionPart.SectionID)
		assert.Equal(t, uint(i), sectionPart.Order)
	}
}
