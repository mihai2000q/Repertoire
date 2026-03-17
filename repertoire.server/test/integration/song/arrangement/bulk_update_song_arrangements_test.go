package arrangement

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
)

func TestBulkUpdateSongArrangements_WhenArrangementsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkUpdateSongArrangementsRequest{
		SongID: songData.SongArrangements[1].SongID,
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   songData.SongArrangements[1].ID,
				Name: "New Chorus Name",
			},
			{
				ID:   songData.SongArrangements[0].ID, // not on the same song
				Name: "New Chorus Name",
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/bulk", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBulkUpdateSongArrangements_WhenSuccessful_ShouldUpdateArrangements(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.BulkUpdateSongArrangementsRequest{
		SongID: songData.SongArrangements[1].SongID,
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   songData.SongArrangements[1].ID,
				Name: "New Chorus Name",
				Occurrences: []requests.UpdateSongSectionOccurrencesRequest{
					{
						SectionID:   songData.SongSections[5].ID,
						Occurrences: 1,
					},
					{
						SectionID:   songData.SongSections[4].ID,
						Occurrences: 7,
					},
					{
						SectionID:   songData.SongSections[6].ID,
						Occurrences: 2,
					},
				},
			},
			{
				ID:   songData.SongArrangements[2].ID,
				Name: "New Chorus Name",
				Occurrences: []requests.UpdateSongSectionOccurrencesRequest{
					{
						SectionID:   songData.SongSections[5].ID,
						Occurrences: 5,
					},
					{
						SectionID:   songData.SongSections[4].ID,
						Occurrences: 0,
					},
					{
						SectionID:   songData.SongSections[6].ID,
						Occurrences: 1,
					},
				},
			},
		},
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/arrangements/bulk", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	requestsMap := make(map[uuid.UUID]requests.UpdateSongArrangementRequest)
	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
		requestsMap[r.ID] = r
	}

	db := utils.GetDatabase(t)

	var arrangements []model.SongArrangement
	db.Preload("SectionOccurrences").Find(&arrangements, ids)

	for _, a := range arrangements {
		assertUpdatedSongArrangement(t, a, requestsMap[a.ID])
	}
}

func assertUpdatedSongArrangement(
	t *testing.T,
	songArrangement model.SongArrangement,
	request requests.UpdateSongArrangementRequest,
) {
	assert.Equal(t, request.Name, songArrangement.Name)

	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Occurrences {
		sectionsOccurrencesMap[s.SectionID] = s.Occurrences
	}
	for i := range songArrangement.SectionOccurrences {
		occurrences := sectionsOccurrencesMap[songArrangement.SectionOccurrences[i].SectionID]
		assert.Equal(t, songArrangement.SectionOccurrences[i].Occurrences, occurrences)
	}
}
