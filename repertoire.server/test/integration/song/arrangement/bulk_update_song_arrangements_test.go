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
				Occurrences: []requests.UpdateSongPartOccurrencesRequest{
					{
						PartID:      songData.SongParts[5].ID,
						Occurrences: 1,
					},
					{
						PartID:      songData.SongParts[4].ID,
						Occurrences: 7,
					},
					{
						PartID:      songData.SongParts[6].ID,
						Occurrences: 2,
					},
				},
			},
			{
				ID:   songData.SongArrangements[2].ID,
				Name: "New Chorus Name",
				Occurrences: []requests.UpdateSongPartOccurrencesRequest{
					{
						PartID:      songData.SongParts[5].ID,
						Occurrences: 5,
					},
					{
						PartID:      songData.SongParts[4].ID,
						Occurrences: 0,
					},
					{
						PartID:      songData.SongParts[6].ID,
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
	db.Preload("PartOccurrences").Find(&arrangements, ids)

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

	partsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, o := range request.Occurrences {
		partsOccurrencesMap[o.PartID] = o.Occurrences
	}
	for _, po := range songArrangement.PartOccurrences {
		occurrences := partsOccurrencesMap[po.PartID]
		assert.Equal(t, po.Occurrences, occurrences)
	}
}
