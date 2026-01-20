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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBulkRehearsalsSongSections_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.New(), Rehearsals: 1}},
		SongID:   uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections/bulk-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkRehearsalsSongSections_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{{ID: uuid.New(), Rehearsals: 1}},
		SongID:   songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections/bulk-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkRehearsalsSongSections_WhenSuccessful_ShouldDeleteSections(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with sections and previous stats
	song := songData.Songs[0]
	oldSections := slices.Clone(song.Sections)
	request := requests.BulkRehearsalsSongSectionsRequest{
		Sections: []requests.BulkRehearsalsSongSectionRequest{
			{ID: song.Sections[0].ID, Rehearsals: 10},
			{ID: song.Sections[1].ID, Rehearsals: 0},
			{ID: song.Sections[2].ID, Rehearsals: 5},
		},
		SongID: songData.Songs[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections/bulk-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.
		Preload("Sections", func(db gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Preload("Sections.History", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Find(&newSong, song.ID)

	for i, newSection := range newSong.Sections {
		ind := slices.IndexFunc(request.Sections, func(s requests.BulkRehearsalsSongSectionRequest) bool {
			return newSection.ID == s.ID
		})
		if ind == -1 || newSection.Rehearsals == 0 {
			continue
		}
		assert.Greater(t, newSection.Rehearsals, oldSections[i].Rehearsals)
		assert.Greater(t, newSection.RehearsalsScore, oldSections[i].RehearsalsScore)
		assert.Greater(t, newSection.Progress, oldSections[i].Progress)

		assert.NotEmpty(t, newSection.History[0].ID)
		assert.Equal(t, oldSections[i].Rehearsals, newSection.History[0].From)
		assert.Equal(t, oldSections[i].Rehearsals+request.Sections[ind].Rehearsals, newSection.History[0].To)
		assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)
	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
