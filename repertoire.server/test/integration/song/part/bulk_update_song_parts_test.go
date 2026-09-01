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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBulkUpdateSongParts_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkUpdateSongPartsRequest{
		Requests: []requests.BulkUpdateSongPartRequest{{ID: uuid.New(), Rehearsals: 1}},
		SongID:   uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-update", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkUpdateSongParts_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkUpdateSongPartsRequest{
		Requests: []requests.BulkUpdateSongPartRequest{{ID: uuid.New(), Rehearsals: 1}},
		SongID:   songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-update", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkUpdateSongParts_WhenRehearsalsAreDecreased_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	part := songData.SongParts[0]
	request := requests.BulkUpdateSongPartsRequest{
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: part.ID, Rehearsals: part.Rehearsals - 1, Confidence: part.Confidence},
		},
		SongID: songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-update", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestBulkUpdateSongParts_WhenNoChanges_ShouldNotUpdateAnything(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.BulkUpdateSongPartsRequest{
		SongID: song.ID,
		Requests: []requests.BulkUpdateSongPartRequest{
			{
				ID:         songData.SongParts[0].ID,
				Rehearsals: songData.SongParts[0].Rehearsals,
				Confidence: songData.SongParts[0].Confidence,
			},
			{
				ID:         songData.SongParts[1].ID,
				Rehearsals: songData.SongParts[1].Rehearsals,
				Confidence: songData.SongParts[1].Confidence,
			},
		},
	}

	db := utils.GetDatabase(t)
	db.Preload("Parts").Find(&song, song.ID)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-update", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db = db.Session(&gorm.Session{NewDB: true})
	var updatedSong model.Song
	db.Preload("Parts").Find(&updatedSong, song.ID)

	for i, part := range song.Parts {
		updatedPart := updatedSong.Parts[i]
		assert.Equal(t, part.Rehearsals, updatedPart.Rehearsals)
		assert.Equal(t, part.Confidence, updatedPart.Confidence)
		assert.Equal(t, part.Progress, updatedPart.Progress)
		assert.Equal(t, part.RehearsalsScore, updatedPart.RehearsalsScore)
		assert.Equal(t, part.ConfidenceScore, updatedPart.ConfidenceScore)
	}
	assert.Equal(t, song.Confidence, updatedSong.Confidence)
	assert.Equal(t, song.Rehearsals, updatedSong.Rehearsals)
	assert.Equal(t, song.Progress, updatedSong.Progress)

	if song.LastTimePlayed != nil {
		assert.Equal(t, song.LastTimePlayed, updatedSong.LastTimePlayed)
	} else {
		assert.Nil(t, updatedSong.LastTimePlayed)
	}
}

func TestBulkUpdateSongParts_WhenSuccessful_ShouldUpdateParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with parts and previous stats
	song := songData.Songs[0]
	oldParts := slices.Clone(slices.DeleteFunc(slices.Clone(songData.SongParts), func(part model.SongPart) bool {
		return part.SongID != song.ID
	}))

	request := requests.BulkUpdateSongPartsRequest{
		Requests: []requests.BulkUpdateSongPartRequest{
			{ID: oldParts[0].ID, Rehearsals: oldParts[0].Rehearsals, Confidence: oldParts[0].Confidence},
			{ID: oldParts[1].ID, Rehearsals: oldParts[1].Rehearsals, Confidence: oldParts[1].Confidence + 15},
			{ID: oldParts[2].ID, Rehearsals: oldParts[2].Rehearsals + 5, Confidence: oldParts[2].Confidence},
			{ID: oldParts[3].ID, Rehearsals: oldParts[3].Rehearsals + 10, Confidence: oldParts[3].Confidence + 5},
		},
		SongID: song.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/bulk-update", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.
		Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Order("song_order")
		}).
		Preload("Parts.History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newSong, song.ID)

	for i, newPart := range newSong.Parts {
		oldPart := oldParts[i]

		ind := slices.IndexFunc(request.Requests, func(r requests.BulkUpdateSongPartRequest) bool {
			return newPart.ID == r.ID
		})
		if ind == -1 {
			// part wasn't in the request, nothing should have changed
			assert.Equal(t, oldPart.Rehearsals, newPart.Rehearsals)
			assert.Equal(t, oldPart.Confidence, newPart.Confidence)
			assert.Equal(t, oldPart.Progress, newPart.Progress)
			assert.Empty(t, newPart.History)
			continue
		}
		req := request.Requests[ind]

		rehearsalsChanged := req.Rehearsals != oldPart.Rehearsals
		confidenceChanged := req.Confidence != oldPart.Confidence

		if rehearsalsChanged {
			assert.Equal(t, req.Rehearsals, newPart.Rehearsals)
			assert.Greater(t, newPart.RehearsalsScore, oldPart.RehearsalsScore)

			historyInd := slices.IndexFunc(newPart.History, func(h model.SongPartHistory) bool {
				return h.Property == model.RehearsalsProperty
			})
			assert.NotEmpty(t, newPart.History[historyInd].ID)
			assert.Equal(t, oldPart.Rehearsals, newPart.History[historyInd].From)
			assert.Equal(t, req.Rehearsals, newPart.History[historyInd].To)
		} else {
			assert.Equal(t, oldPart.Rehearsals, newPart.Rehearsals)
		}

		if confidenceChanged {
			assert.Equal(t, req.Confidence, newPart.Confidence)
			assert.NotEqual(t, oldPart.ConfidenceScore, newPart.ConfidenceScore)

			historyInd := slices.IndexFunc(newPart.History, func(h model.SongPartHistory) bool {
				return h.Property == model.ConfidenceProperty
			})
			assert.NotEmpty(t, newPart.History[historyInd].ID)
			assert.Equal(t, oldPart.Confidence, newPart.History[historyInd].From)
			assert.Equal(t, req.Confidence, newPart.History[historyInd].To)
		} else {
			assert.Equal(t, oldPart.Confidence, newPart.Confidence)
		}
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Confidence, song.Confidence)
	assert.Greater(t, newSong.Progress, song.Progress)
	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
