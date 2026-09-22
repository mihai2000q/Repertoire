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
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &[]uuid.UUID{uuid.New()}[0]},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
		Parts: []requests.CreateSongSectionPartRequest{
			{
				NewPart: &requests.CreateNewSongPartRequest{
					Name:         "New Part",
					InstrumentID: &[]uuid.UUID{uuid.New()}[0],
				},
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenBandMembersAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &songData.SongParts[0].ID, BandMemberIDs: []uuid.UUID{uuid.New()}},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenPartsDoesNotBelong_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: songData.Songs[0].ID,
		Name:   "Chorus 1-New",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &songData.SongParts[4].ID},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateSongSection_WhenSuccessful_ShouldCreateSection(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongSectionRequest
	}{
		{
			"Minimal",
			requests.CreateSongSectionRequest{
				SongID: songData.Songs[0].ID,
				Name:   "Chorus 1-New",
				TypeID: songData.Users[0].SongSectionTypes[0].ID,
			},
		},
		{
			"With Parts",
			requests.CreateSongSectionRequest{
				SongID: songData.Songs[0].ID,
				Name:   "Chorus 1-New",
				TypeID: songData.Users[0].SongSectionTypes[0].ID,
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: &songData.SongParts[0].ID},
					{
						NewPart: &requests.CreateNewSongPartRequest{
							Name:         "New Part",
							InstrumentID: &songData.Users[0].Instruments[0].ID,
						},
					},
					{PartID: &songData.SongParts[1].ID, BandMemberIDs: []uuid.UUID{songData.Artists[0].BandMembers[0].ID}},
					{
						NewPart: &requests.CreateNewSongPartRequest{Name: "New Part"},
						BandMemberIDs: []uuid.UUID{
							songData.Artists[0].BandMembers[0].ID,
							songData.Artists[0].BandMembers[1].ID,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			sectionsCount := len(slices.Clone(slices.DeleteFunc(slices.Clone(songData.SongSections), func(section model.SongSection) bool {
				return section.SongID != test.request.SongID
			})))

			partsCount := len(slices.Clone(slices.DeleteFunc(slices.Clone(songData.SongParts), func(part model.SongPart) bool {
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
					return db.Preload("BandMembers").Joins("Part").Order("\"order\"")
				}).
				Find(&section, &model.SongSection{Name: test.request.Name})

			assertCreatedSongSection(t, section, test.request, sectionsCount, partsCount)
		})
	}
}

func assertCreatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.CreateSongSectionRequest,
	order int,
	songOrder int,
) {
	assert.NotEmpty(t, songSection.ID)
	assert.Equal(t, request.SongID, songSection.SongID)
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)
	assert.Equal(t, uint(order), songSection.Order)

	for i, sectionPart := range songSection.SectionParts {
		req := request.Parts[i]
		if req.PartID != nil {
			assert.Equal(t, *req.PartID, sectionPart.PartID)
		}
		if req.NewPart != nil {
			assert.Equal(t, req.NewPart.Name, sectionPart.Part.Name)
			assert.Equal(t, req.NewPart.InstrumentID, sectionPart.Part.InstrumentID)
			assert.Equal(t, uint(songOrder), sectionPart.Part.SongOrder)
			songOrder++
		}
		assert.Equal(t, uint(i), sectionPart.Order)
		// bandMembers
		assert.Len(t, sectionPart.BandMembers, len(req.BandMemberIDs))
		bandMemberIDs := make([]uuid.UUID, len(sectionPart.BandMembers))
		for j, bm := range sectionPart.BandMembers {
			bandMemberIDs[j] = bm.ID
		}
		assert.ElementsMatch(t, req.BandMemberIDs, bandMemberIDs)
	}
}
