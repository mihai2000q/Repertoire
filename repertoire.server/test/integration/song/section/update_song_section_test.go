package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

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

func TestUpdateSongSection_WhenRehearsalsAreDecreasing_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	section := songData.Songs[0].Sections[0]

	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: section.Rehearsals - 1,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateSongSection_WhenRequestChangesBandMemberIDButItIsNotAssociated_ShouldReturnConflictError(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
	}{
		{
			"Song without artist",
			songData.Songs[4],
		},
		{
			"Song with artist but without that member",
			songData.Songs[0],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			section := test.song.Sections[0]
			request := requests.UpdateSongSectionRequest{
				ID:           section.ID,
				TypeID:       section.SongSectionTypeID,
				Rehearsals:   section.Rehearsals,
				Confidence:   section.Confidence,
				Name:         "Chorus 1-New",
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/sections", request)

			// then
			assert.Equal(t, http.StatusConflict, w.Code)
		})
	}
}

func TestUpdateSongSection_WhenSuccessful_ShouldUpdateSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongSectionRequest{
		ID:     songData.Songs[0].Sections[2].ID,
		Name:   "New Chorus Name",
		TypeID: songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var section model.SongSection
	db.Find(&section, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, section, request)
}

func TestUpdateSongSection_WhenSuccessfulWithRehearsals_ShouldUpdateSectionUpdateSongAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	section := song.Sections[0]
	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: 15,
		Confidence: section.Confidence,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSection model.SongSection
	db.Preload("Song").
		Preload("History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newSection, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, newSection, request)

	assert.Greater(t, newSection.Rehearsals, section.Rehearsals)
	assert.Greater(t, newSection.RehearsalsScore, section.RehearsalsScore)
	assert.Greater(t, newSection.Progress, section.Progress)

	assert.NotEmpty(t, newSection.History[0].ID)
	assert.Equal(t, section.Rehearsals, newSection.History[0].From)
	assert.Equal(t, request.Rehearsals, newSection.History[0].To)
	assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)

	assert.Greater(t, newSection.Song.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSection.Song.Progress, song.Progress)

	assert.NotNil(t, newSection.Song.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSection.Song.LastTimePlayed, 1*time.Minute)
}

func TestUpdateSongSection_WhenSuccessfulWithConfidenceIncreasing_ShouldUpdateSectionUpdateSongAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	section := song.Sections[0]
	request := requests.UpdateSongSectionRequest{
		ID:         section.ID,
		Name:       "New Chorus Name",
		Rehearsals: section.Rehearsals,
		Confidence: 25,
		TypeID:     songData.Users[0].SongSectionTypes[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSection model.SongSection
	db.Preload("Song").
		Preload("History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newSection, &model.SongSection{ID: request.ID})

	assertUpdatedSongSection(t, newSection, request)

	assert.Greater(t, newSection.Confidence, section.Confidence)
	assert.Greater(t, newSection.ConfidenceScore, section.ConfidenceScore)
	assert.Greater(t, newSection.Progress, section.Progress)

	assert.NotEmpty(t, newSection.History[0].ID)
	assert.Equal(t, section.Confidence, newSection.History[0].From)
	assert.Equal(t, request.Confidence, newSection.History[0].To)
	assert.Equal(t, model.ConfidenceProperty, newSection.History[0].Property)

	assert.Greater(t, newSection.Song.Confidence, song.Confidence)
	assert.Greater(t, newSection.Song.Progress, song.Progress)
}

func TestUpdateSongSection_WhenSuccessfulWithBandMember_ShouldUpdateSection(t *testing.T) {
	tests := []struct {
		name         string
		section      model.SongSection
		bandMemberID *uuid.UUID
	}{
		{
			"to Nil Band Member",
			songData.Songs[0].Sections[1],
			nil,
		},
		{
			"from member to Another Band Member",
			songData.Songs[0].Sections[1],
			&songData.Artists[0].BandMembers[1].ID,
		},
		{
			"from nil to Another Band Member",
			songData.Songs[0].Sections[2],
			&songData.Artists[0].BandMembers[1].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			request := requests.UpdateSongSectionRequest{
				ID:           test.section.ID,
				Name:         test.section.Name,
				Rehearsals:   test.section.Rehearsals,
				Confidence:   test.section.Confidence,
				TypeID:       test.section.SongSectionTypeID,
				BandMemberID: test.bandMemberID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/sections", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var section model.SongSection
			db.Find(&section, &model.SongSection{ID: request.ID})

			assertUpdatedSongSection(t, section, request)
		})
	}
}

func assertUpdatedSongSection(
	t *testing.T,
	songSection model.SongSection,
	request requests.UpdateSongSectionRequest,
) {
	assert.Equal(t, request.Name, songSection.Name)
	assert.Equal(t, request.Confidence, songSection.Confidence)
	assert.Equal(t, request.Rehearsals, songSection.Rehearsals)
	assert.Equal(t, request.TypeID, songSection.SongSectionTypeID)
	assert.Equal(t, request.BandMemberID, songSection.BandMemberID)
	assert.Equal(t, request.InstrumentID, songSection.InstrumentID)
}
