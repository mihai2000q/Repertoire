package section

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestCreateSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Chorus 1-New",
		TypeID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/songs/sections", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSongSection_WhenRequestHasBandMemberIDButItIsNotAssociated_ShouldReturnBadRequestError(t *testing.T) {
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

			request := requests.CreateSongSectionRequest{
				SongID:       test.song.ID,
				Name:         "Chorus 1-New",
				TypeID:       uuid.New(),
				BandMemberID: &[]uuid.UUID{uuid.New()}[0],
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/sections", request)

			// then
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestCreateSongSection_WhenSuccessful_ShouldCreateSection(t *testing.T) {
	tests := []struct {
		name         string
		song         model.Song
		bandMemberID *uuid.UUID
		instrumentID *uuid.UUID
	}{
		{
			"Without Band Member or Instrument",
			songData.Songs[0],
			nil,
			nil,
		},
		{
			"With Band Member",
			songData.Songs[0],
			&songData.Artists[0].BandMembers[0].ID,
			nil,
		},
		{
			"With Instrument",
			songData.Songs[0],
			nil,
			&songData.Users[0].Instruments[0].ID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// song with sections and previous stats
			song := songData.Songs[0]
			request := requests.CreateSongSectionRequest{
				SongID:       song.ID,
				Name:         "Chorus 1-New",
				TypeID:       songData.Users[0].SongSectionTypes[0].ID,
				BandMemberID: test.bandMemberID,
				InstrumentID: test.instrumentID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().POST(w, "/api/songs/sections", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase(t)

			var section model.SongSection
			db.Preload("Song").Find(&section, &model.SongSection{Name: request.Name})

			assert.LessOrEqual(t, section.Song.Confidence, song.Confidence)
			assert.LessOrEqual(t, section.Song.Rehearsals, song.Rehearsals)
			assert.LessOrEqual(t, section.Song.Progress, song.Progress)

			assertCreatedSongSection(t, section, request, len(song.Sections))
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
	assert.Equal(t, request.BandMemberID, songSection.BandMemberID)
	assert.Equal(t, request.InstrumentID, songSection.InstrumentID)
	assert.Zero(t, songSection.Rehearsals)
	assert.Equal(t, model.DefaultSongSectionConfidence, songSection.Confidence)
	assert.Zero(t, songSection.RehearsalsScore)
	assert.Zero(t, songSection.ConfidenceScore)
	assert.Zero(t, songSection.Progress)
	assert.Equal(t, uint(order), songSection.Order)
}
