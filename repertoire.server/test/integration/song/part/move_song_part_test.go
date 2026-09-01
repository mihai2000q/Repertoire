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
)

func TestMoveSongPartInSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.MoveSongPartInSongRequest{
		SongID: uuid.New(),
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/move-in-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongPartInSong_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.MoveSongPartInSongRequest{
		SongID: song.ID,
		ID:     uuid.New(),
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/move-in-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongPartInSong_WhenOverPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.MoveSongPartInSongRequest{
		SongID: song.ID,
		ID:     songData.SongParts[0].ID,
		OverID: uuid.New(),
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs/parts/move-in-song", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMoveSongPartInSong_WhenSuccessful_ShouldMoveTypes(t *testing.T) {
	tests := []struct {
		name      string
		request   requests.MoveSongPartInSongRequest
		index     int
		overIndex int
	}{
		{
			"From upper position to lower",
			requests.MoveSongPartInSongRequest{
				SongID: songData.SongParts[0].SongID,
				ID:     songData.SongParts[2].ID,
				OverID: songData.SongParts[0].ID,
			},
			2,
			0,
		},
		{
			"From lower position to upper",
			requests.MoveSongPartInSongRequest{
				SongID: songData.SongParts[0].SongID,
				ID:     songData.SongParts[0].ID,
				OverID: songData.SongParts[2].ID,
			},
			0,
			2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().PUT(w, "/api/songs/parts/move-in-song", test.request)

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			var parts []model.SongPart
			db := utils.GetDatabase(t)
			db.Where(&model.SongPart{SongID: test.request.SongID}).Order("\"song_order\"").Find(&parts)

			assertMovedParts(t, test.request, parts, test.index, test.overIndex)
		})
	}
}

func assertMovedParts(
	t *testing.T,
	request requests.MoveSongPartInSongRequest,
	parts []model.SongPart,
	index int,
	overIndex int,
) {
	if index < overIndex {
		assert.Equal(t, parts[overIndex-1].ID, request.OverID)
	} else if index > overIndex {
		assert.Equal(t, parts[overIndex+1].ID, request.OverID)
	}

	assert.Equal(t, parts[overIndex].ID, request.ID)
	for i, part := range parts {
		assert.Equal(t, uint(i), part.SongOrder)
	}
}
