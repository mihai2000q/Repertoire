package song

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddCustomSongRehearsals_WhenSongsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: songData.SongArrangements[0].SongID, ArrangementID: songData.SongArrangements[0].ID},
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAddCustomSongRehearsals_WhenArrangementsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: songData.Songs[0].ID, ArrangementID: uuid.New()},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsals", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test case does not cover the usage of the duplicate songs with different arrangements
func TestAddCustomSongRehearsals_WhenSuccessful_ShouldUpdateSongAndPartsIfTheyHaveOccurrences(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: songData.SongArrangements[0].SongID, ArrangementID: songData.SongArrangements[0].ID},
			{ID: songData.SongArrangements[1].SongID, ArrangementID: songData.SongArrangements[1].ID},
			{ID: songData.SongArrangements[4].SongID, ArrangementID: songData.SongArrangements[4].ID},
			{ID: songData.SongArrangements[1].SongID, ArrangementID: songData.SongArrangements[1].ID},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	getSongsQuery := func(db *gorm.DB, songs *[]model.Song) {
		db.Preload("Parts", func(db *gorm.DB) *gorm.DB { return db.Order("song_parts.song_order") }).
			Preload("Parts.History", func(db *gorm.DB) *gorm.DB {
				return db.
					Where(&model.SongPartHistory{Property: model.RehearsalsProperty}).
					Order("created_at desc")
			}).
			Preload("Parts.ArrangementOccurrences", func(db *gorm.DB) *gorm.DB {
				return db.Joins("LEFT JOIN song_arrangements ON id = arrangement_id").Order("\"order\"")
			}).
			Find(&songs, ids)
	}

	var songs []model.Song
	db := utils.GetDatabase(t)
	getSongsQuery(db, &songs)

	songsMap := make(map[uuid.UUID]model.Song)
	for _, s := range songs {
		songsMap[s.ID] = s
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/custom-rehearsals", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSongs []model.Song
	db = db.Session(&gorm.Session{NewDB: true})
	getSongsQuery(db, &newSongs)

	newSongsMap := make(map[uuid.UUID]model.Song)
	for _, s := range newSongs {
		newSongsMap[s.ID] = s
	}

	songDuplicatesMap := make(map[uuid.UUID]int)
	for _, r := range request.Requests {
		_, ok := songDuplicatesMap[r.ID]
		if ok {
			songDuplicatesMap[r.ID] += 1
		} else {
			songDuplicatesMap[r.ID] = 0
		}
	}

	for _, r := range request.Requests {
		assertion.CustomSongRehearsalWithDuplicates(
			t,
			songsMap[r.ID],
			newSongsMap[r.ID],
			r.ArrangementID,
			songDuplicatesMap[r.ID],
		)
	}
}
