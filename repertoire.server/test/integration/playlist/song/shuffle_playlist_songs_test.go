package song

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	playlistData "repertoire/server/test/integration/test/data/playlist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestShufflePlaylistSongs_WhenSuccessful_ShouldShuffleSongsOnPlaylist(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, playlistData.Users, playlistData.SeedData)

	playlist := playlistData.Playlists[0]
	request := requests.ShufflePlaylistSongsRequest{
		ID: playlist.ID,
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().POST(w, "/api/playlists/songs/shuffle", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response wrapper.WithTotalCount[model.Song]
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	var playlistSongs []model.PlaylistSong
	db := utils.GetDatabase(t)
	db.Find(&playlistSongs, model.PlaylistSong{PlaylistID: playlist.ID})

	for i := range playlistSongs {
		assert.Equal(t, uint(i+1), playlistSongs[i].SongTrackNo)
	}
}
