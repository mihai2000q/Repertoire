package song

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
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
		{
			"With Playlist",
			songData.Songs[6],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			db := utils.GetDatabase(t)
			db.Preload("Playlists").Find(&test.song, test.song.ID)

			messages := utils.SubscribeToTopic(topics.SongsDeletedTopic)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().DELETE(w, "/api/songs/"+test.song.ID.String())

			// then
			assert.Equal(t, http.StatusOK, w.Code)

			db = db.Session(&gorm.Session{NewDB: true})

			var deletedSong model.Song
			db.Find(&deletedSong, test.song.ID)
			assert.Empty(t, deletedSong)

			// assert ordered album
			if test.song.AlbumID != nil {
				var albumSongs []model.Song
				db.Order("album_track_no").Find(&albumSongs, model.Song{AlbumID: test.song.AlbumID})

				for i, song := range albumSongs {
					assert.NotEqual(t, test.song.ID, song.ID)
					assert.Equal(t, uint(i+1), *song.AlbumTrackNo)
				}
			}

			// get playlists from song
			var playlistsIds []uuid.UUID
			for _, playlist := range test.song.Playlists {
				playlistsIds = append(playlistsIds, playlist.ID)
			}

			var playlists []model.Playlist
			db.
				Preload("PlaylistSongs", func(db *gorm.DB) *gorm.DB {
					return db.Order("song_track_no")
				}).
				Find(&playlists, playlistsIds)

			// assert ordered playlists
			for _, playlist := range playlists {
				for i, playlistSong := range playlist.PlaylistSongs {
					assert.NotEqual(t, test.song.ID, playlistSong.SongID)
					assert.Equal(t, uint(i+1), playlistSong.SongTrackNo)
				}
			}

			assertion.AssertMessage(t, messages, func(payloadSongs []model.Song) {
				assert.Len(t, payloadSongs, 1)
				assert.Equal(t, test.song.ID, payloadSongs[0].ID)
			})
		})
	}
}
