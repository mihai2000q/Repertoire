package member

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestMoveBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	request := requests.MoveBandMemberRequest{
		ArtistID: uuid.New(),
		ID:       uuid.New(),
		OverID:   uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveBandMember_WhenMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	song := artistData.Artists[0]
	request := requests.MoveBandMemberRequest{
		ArtistID: song.ID,
		ID:       uuid.New(),
		OverID:   uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveBandMember_WhenOverMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	artist := artistData.Artists[0]
	request := requests.MoveBandMemberRequest{
		ArtistID: artist.ID,
		ID:       artist.BandMembers[0].ID,
		OverID:   uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/artists/band-members/move", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveBandMemberType_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		artist    model.Artist
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			artistData.Artists[0],
			2,
			0,
		},
		{
			"From lower position to upper",
			artistData.Artists[0],
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

			request := requests.MoveBandMemberRequest{
				ArtistID: test.artist.ID,
				ID:       test.artist.BandMembers[test.index].ID,
				OverID:   test.artist.BandMembers[test.overIndex].ID,
			}

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/artists/band-members/move", request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var sections []model.BandMember
			db := utils.GetDatabase(t)
			db.Order("\"order\"").Find(&sections, &model.BandMember{ArtistID: test.artist.ID})

			assertMovedTunings(t, request, sections, test.index, test.overIndex)
		})
	}
}

func assertMovedTunings(
	t *testing.T,
	request requests.MoveBandMemberRequest,
	members []model.BandMember,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, members[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, members[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, members[overIndex].ID, request.ID)
	for i, member := range members {
		assert.Equal(t, uint(i), member.Order)
	}
}
