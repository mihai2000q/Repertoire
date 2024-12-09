package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSong_WhenSuccessful_ShouldDeleteSong(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
	}{
		{
			"Normal delete, without album or files",
			songData.Songs[4],
		},
		{
			"With Image",
			songData.Songs[3],
		},
		{
			"With Album",
			songData.Songs[2],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/songs/"+test.song.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db := utils.GetDatabase()

			var deletedSong model.Song
			db.Find(&deletedSong, test.song.ID)

			assert.Empty(t, deletedSong)

			if test.song.AlbumID != nil {
				var albumSongs []model.Song
				db.Order("album_track_no").Find(&albumSongs, model.Song{AlbumID: test.song.AlbumID})

				for i, song := range albumSongs {
					assert.Equal(t, uint(i+1), *song.AlbumTrackNo)
				}
			}
		})
	}

}
