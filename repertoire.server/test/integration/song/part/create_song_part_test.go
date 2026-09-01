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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCreateSongPart_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID: uuid.New(),
		Name:   "Chorus 1-New",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongPart_WhenSectionsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:     songData.Songs[0].ID,
		Name:       "Chorus 1-New",
		SectionIDs: []uuid.UUID{uuid.New()},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongPart_WhenSectionsDoNotBelongToSong_ShouldReturnNotConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:     songData.Songs[0].ID,
		Name:       "Chorus 1-New",
		SectionIDs: []uuid.UUID{songData.SongSections[4].ID},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateSongPart_WhenRequestHasBandMemberIDButItIsNotAssociated_ShouldReturnConflictError(t *testing.T) {
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

			request := requests.CreateSongPartRequest{
				SongID:       test.song.ID,
				Name:         "Chorus 1-New",
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/parts", request)

			// then
			assert.Equal(t, http.StatusConflict, w.Code)
		})
	}
}

func TestCreateSongPart_WhenSuccessful_ShouldCreatePart(t *testing.T) {
	tests := []struct {
		name         string
		song         model.Song
		bandMemberID *uuid.UUID
		instrumentID *uuid.UUID
		sectionIDs   []uuid.UUID
	}{
		{
			"Without Band Member or Instrument",
			songData.Songs[0],
			nil,
			nil,
			[]uuid.UUID{},
		},
		{
			"With Band Member",
			songData.Songs[0],
			&songData.Artists[0].BandMembers[0].ID,
			nil,
			[]uuid.UUID{},
		},
		{
			"With Instrument",
			songData.Songs[0],
			nil,
			&songData.Users[0].Instruments[0].ID,
			[]uuid.UUID{},
		},
		{
			"With Sections",
			songData.Songs[0],
			nil,
			&songData.Users[0].Instruments[0].ID,
			[]uuid.UUID{songData.SongSections[0].ID, songData.SongSections[1].ID},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// song with parts and previous stats
			song := songData.Songs[0]
			request := requests.CreateSongPartRequest{
				SongID:       song.ID,
				Name:         "Chorus 1-New",
				BandMemberID: test.bandMemberID,
				InstrumentID: test.instrumentID,
				SectionIDs:   test.sectionIDs,
			}

			db := utils.GetDatabase(t)
			var oldArrangements []model.SongArrangement
			var partsCount int64
			db.Preload("PartOccurrences").
				Where(&model.SongArrangement{SongID: song.ID}).
				Order("\"order\"").
				Find(&oldArrangements)
			db.Model(&model.SongPart{}).Where(&model.SongPart{SongID: song.ID}).Count(&partsCount)

			expectedOrders := make(map[uuid.UUID]uint)
			for _, sectionID := range request.SectionIDs {
				var count int64
				db.Model(&model.SongSectionPart{}).
					Where(&model.SongSectionPart{SectionID: sectionID}).
					Count(&count)
				expectedOrders[sectionID] = uint(count)
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/parts", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db = db.Session(&gorm.Session{NewDB: true})

			var part model.SongPart
			db.Preload("Song").
				Preload("Song.Arrangements", func(db *gorm.DB) *gorm.DB { return db.Order("\"order\"") }).
				Preload("Song.Arrangements.PartOccurrences", func(db *gorm.DB) *gorm.DB {
					return db.
						Joins("LEFT JOIN song_parts ON song_parts.id = song_part_occurrences.part_id").
						Order("song_parts.song_order DESC")
				}).
				Preload("Song.Arrangements.PartOccurrences.Part").
				Preload("SectionParts", func(db *gorm.DB) *gorm.DB {
					return db.Order("\"order\"")
				}).
				Find(&part, &model.SongPart{Name: request.Name})

			assertCreatedSongPart(t, part, request, partsCount)

			// updateSong
			assert.LessOrEqual(t, part.Song.Confidence, song.Confidence)
			assert.LessOrEqual(t, part.Song.Rehearsals, song.Rehearsals)
			assert.LessOrEqual(t, part.Song.Progress, song.Progress)

			// updateArrangements
			for i, arrangement := range part.Song.Arrangements {
				assert.Len(t, arrangement.PartOccurrences, len(oldArrangements[i].PartOccurrences)+1)
				newOccurrence := arrangement.PartOccurrences[0]
				assert.Equal(t, part.ID, newOccurrence.PartID)
				assert.Equal(t, arrangement.ID, newOccurrence.ArrangementID)
				assert.Zero(t, newOccurrence.Occurrences)
			}

			// createSectionParts
			assert.Len(t, part.SectionParts, len(request.SectionIDs))

			sectionPartsMap := make(map[uuid.UUID]model.SongSectionPart)
			for _, sp := range part.SectionParts {
				assert.Equal(t, part.ID, sp.PartID)
				sectionPartsMap[sp.SectionID] = sp
			}

			// Verify that we have exactly the requested sections and their orders match expectations
			for _, sectionID := range request.SectionIDs {
				sp, ok := sectionPartsMap[sectionID]
				assert.True(t, ok, "section %v missing from SectionParts", sectionID)
				expectedOrder := expectedOrders[sectionID]
				assert.Equal(t, expectedOrder, sp.Order)
			}
		})
	}
}

func assertCreatedSongPart(
	t *testing.T,
	songPart model.SongPart,
	request requests.CreateSongPartRequest,
	order int64,
) {
	assert.NotEmpty(t, songPart.ID)
	assert.Equal(t, request.SongID, songPart.SongID)
	assert.Equal(t, request.Name, songPart.Name)
	assert.Equal(t, request.BandMemberID, songPart.BandMemberID)
	assert.Equal(t, request.InstrumentID, songPart.InstrumentID)
	assert.Zero(t, songPart.Rehearsals)
	assert.Equal(t, model.DefaultSongPartConfidence, songPart.Confidence)
	assert.Zero(t, songPart.RehearsalsScore)
	assert.Zero(t, songPart.ConfidenceScore)
	assert.Zero(t, songPart.Progress)
	assert.Equal(t, uint(order), songPart.SongOrder)
}
