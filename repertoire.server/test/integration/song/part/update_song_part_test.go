package part

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

func TestUpdateSongPart_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongPartRequest{
		ID:   uuid.New(),
		Name: "New Chorus Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongPart_WhenBandMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	part := songData.SongParts[0]

	request := requests.UpdateSongPartRequest{
		ID:           part.ID,
		Name:         "New Chorus Name",
		Rehearsals:   part.Rehearsals,
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongPart_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	part := songData.SongParts[0]

	request := requests.UpdateSongPartRequest{
		ID:           part.ID,
		Name:         "New Chorus Name",
		Rehearsals:   part.Rehearsals,
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSongPart_WhenRehearsalsAreDecreasing_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	part := songData.SongParts[0]

	request := requests.UpdateSongPartRequest{
		ID:         part.ID,
		Name:       "New Chorus Name",
		Rehearsals: part.Rehearsals - 1,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateSongPart_WhenRequestHasBandMemberIDButItIsNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongPartRequest{
		ID:           songData.SongParts[0].ID,
		Name:         "Chorus 1-New",
		BandMemberID: &songData.Artists[1].BandMembers[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateSongPart_WhenSuccessful_ShouldUpdatePart(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongPartRequest{
		ID:   songData.SongParts[2].ID,
		Name: "New Chorus Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var part model.SongPart
	db.Find(&part, &model.SongPart{ID: request.ID})

	assertUpdatedSongPart(t, part, request)
}

func TestUpdateSongPart_WhenSuccessfulWithRehearsals_ShouldUpdatePartUpdateSongAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	part := songData.SongParts[0]
	request := requests.UpdateSongPartRequest{
		ID:         part.ID,
		Name:       "New Chorus Name",
		Rehearsals: 15,
		Confidence: part.Confidence,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newPart model.SongPart
	db.Preload("Song").
		Preload("History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newPart, &model.SongPart{ID: request.ID})

	assertUpdatedSongPart(t, newPart, request)

	assert.Greater(t, newPart.Rehearsals, part.Rehearsals)
	assert.Greater(t, newPart.RehearsalsScore, part.RehearsalsScore)
	assert.Greater(t, newPart.Progress, part.Progress)

	assert.NotEmpty(t, newPart.History[0].ID)
	assert.Equal(t, part.Rehearsals, newPart.History[0].From)
	assert.Equal(t, request.Rehearsals, newPart.History[0].To)
	assert.Equal(t, model.RehearsalsProperty, newPart.History[0].Property)

	assert.Greater(t, newPart.Song.Rehearsals, song.Rehearsals)
	assert.Greater(t, newPart.Song.Progress, song.Progress)

	assert.NotNil(t, newPart.Song.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newPart.Song.LastTimePlayed, 1*time.Minute)
}

func TestUpdateSongPart_WhenSuccessfulWithConfidenceIncreasing_ShouldUpdatePartUpdateSongAddHistoryAndChangeScore(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	part := songData.SongParts[0]
	request := requests.UpdateSongPartRequest{
		ID:         part.ID,
		Name:       "New Chorus Name",
		Rehearsals: part.Rehearsals,
		Confidence: 25,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newPart model.SongPart
	db.Preload("Song").
		Preload("History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newPart, &model.SongPart{ID: request.ID})

	assertUpdatedSongPart(t, newPart, request)

	assert.Greater(t, newPart.Confidence, part.Confidence)
	assert.Greater(t, newPart.ConfidenceScore, part.ConfidenceScore)
	assert.Greater(t, newPart.Progress, part.Progress)

	assert.NotEmpty(t, newPart.History[0].ID)
	assert.Equal(t, part.Confidence, newPart.History[0].From)
	assert.Equal(t, request.Confidence, newPart.History[0].To)
	assert.Equal(t, model.ConfidenceProperty, newPart.History[0].Property)

	assert.Greater(t, newPart.Song.Confidence, song.Confidence)
	assert.Greater(t, newPart.Song.Progress, song.Progress)
}

func TestUpdateSongPart_WhenSuccessfulWithBandMember_ShouldUpdatePart(t *testing.T) {
	tests := []struct {
		name         string
		part         model.SongPart
		bandMemberID *uuid.UUID
	}{
		{
			"to Nil Band Member",
			songData.SongParts[1],
			nil,
		},
		{
			"from member to Another Band Member",
			songData.SongParts[1],
			&songData.Artists[0].BandMembers[1].ID,
		},
		{
			"from nil to Another Band Member",
			songData.SongParts[2],
			&songData.Artists[0].BandMembers[1].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			request := requests.UpdateSongPartRequest{
				ID:           test.part.ID,
				Name:         test.part.Name,
				Rehearsals:   test.part.Rehearsals,
				Confidence:   test.part.Confidence,
				BandMemberID: test.bandMemberID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/parts", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var part model.SongPart
			db.Find(&part, &model.SongPart{ID: request.ID})

			assertUpdatedSongPart(t, part, request)
		})
	}
}

func assertUpdatedSongPart(
	t *testing.T,
	songPart model.SongPart,
	request requests.UpdateSongPartRequest,
) {
	assert.Equal(t, request.Name, songPart.Name)
	assert.Equal(t, request.Confidence, songPart.Confidence)
	assert.Equal(t, request.Rehearsals, songPart.Rehearsals)
	assert.Equal(t, request.BandMemberID, songPart.BandMemberID)
	assert.Equal(t, request.InstrumentID, songPart.InstrumentID)
}
