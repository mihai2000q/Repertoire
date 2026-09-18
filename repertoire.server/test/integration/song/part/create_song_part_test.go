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

func TestCreateSongPart_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:       songData.Songs[0].ID,
		Name:         "Chorus 1-New",
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongPart_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:    songData.Songs[0].ID,
		Name:      "Chorus 1-New",
		SectionID: &[]uuid.UUID{uuid.New()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongPart_WhenSectionDoesNotBelongToSong_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:    songData.Songs[0].ID,
		Name:      "Chorus 1-New",
		SectionID: &songData.SongSections[4].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateSongPart_WhenBandMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:       songData.Songs[0].ID,
		Name:         "Chorus 1-New",
		SectionID:    &songData.SongSections[0].ID,
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongPart_WhenRequestHasBandMemberIDButItIsNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongPartRequest{
		SongID:       songData.Songs[0].ID,
		Name:         "Chorus 1-New",
		SectionID:    &songData.SongSections[0].ID,
		BandMemberID: &songData.Artists[1].BandMembers[0].ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/parts", request)

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateSongPart_WhenSuccessful_ShouldCreatePart(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongPartRequest
		song    model.Song
	}{
		{
			"Without Band Member or Instrument",
			requests.CreateSongPartRequest{
				SongID: songData.Songs[0].ID,
				Name:   "Chorus 1-New",
			},
			songData.Songs[0],
		},
		{
			"With Instrument",
			requests.CreateSongPartRequest{
				SongID:       songData.Songs[0].ID,
				Name:         "Chorus 1-New",
				InstrumentID: &songData.Users[0].Instruments[0].ID,
			},
			songData.Songs[0],
		},
		{
			"With Band Member and Section",
			requests.CreateSongPartRequest{
				SongID:       songData.Songs[0].ID,
				Name:         "Chorus 1-New",
				SectionID:    &songData.SongSections[0].ID,
				BandMemberID: &songData.Artists[0].BandMembers[0].ID,
			},
			songData.Songs[0],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			db := utils.GetDatabase(t)
			var oldArrangements []model.SongArrangement
			var partsCount int64
			var sectionPartsCount int64
			db.Preload("PartOccurrences").
				Where(&model.SongArrangement{SongID: test.song.ID}).
				Order("\"order\"").
				Find(&oldArrangements)
			db.Model(&model.SongPart{}).Where(&model.SongPart{SongID: test.song.ID}).Count(&partsCount)
			if test.request.SectionID != nil {
				db.Model(&model.SongSectionPart{}).Where(&model.SongSectionPart{SectionID: *test.request.SectionID}).Count(&sectionPartsCount)
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/parts", test.request)

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
				Find(&part, &model.SongPart{Name: test.request.Name})

			assertCreatedSongPart(t, part, test.request, partsCount)

			// updateSong
			assert.Less(t, part.Song.Confidence, test.song.Confidence)
			assert.Less(t, part.Song.Rehearsals, test.song.Rehearsals)
			assert.Less(t, part.Song.Progress, test.song.Progress)

			// updateArrangements
			for i, arrangement := range part.Song.Arrangements {
				assert.Len(t, arrangement.PartOccurrences, len(oldArrangements[i].PartOccurrences)+1)
				newOccurrence := arrangement.PartOccurrences[0]
				assert.Equal(t, part.ID, newOccurrence.PartID)
				assert.Equal(t, arrangement.ID, newOccurrence.ArrangementID)
				assert.Zero(t, newOccurrence.Occurrences)
			}

			// updateBandMember
			if test.request.SectionID == nil {
				assert.Empty(t, part.SectionParts)
				return
			}

			assert.Len(t, part.SectionParts, 1)
			newSectionPart := part.SectionParts[0]
			assert.Equal(t, part.ID, newSectionPart.PartID)
			assert.Equal(t, *test.request.SectionID, newSectionPart.SectionID)
			assert.Equal(t, test.request.BandMemberID, newSectionPart.BandMemberID)
			assert.Equal(t, uint(sectionPartsCount), newSectionPart.Order)
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
	assert.Equal(t, request.InstrumentID, songPart.InstrumentID)
	assert.Zero(t, songPart.Rehearsals)
	assert.Equal(t, model.DefaultSongPartConfidence, songPart.Confidence)
	assert.Zero(t, songPart.RehearsalsScore)
	assert.Zero(t, songPart.ConfidenceScore)
	assert.Zero(t, songPart.Progress)
	assert.Equal(t, uint(order), songPart.SongOrder)
}
