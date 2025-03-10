package song

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"
	"time"
)

func TestUpdateSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Title",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSong_WhenSuccessful_ShouldUpdateSong(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	request := requests.UpdateSongRequest{
		ID:          song.ID,
		Title:       "New Title",
		ReleaseDate: &[]time.Time{time.Now()}[0],
		AlbumID:     song.AlbumID,
		ArtistID:    song.ArtistID,
	}

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var updatedSong model.Song
	db := utils.GetDatabase(t)
	db.Find(&updatedSong, request.ID)

	assertUpdatedSong(t, request, updatedSong)

	assertion.AssertMessage(t, messages, func(payloadSong model.Song) {
		assert.Equal(t, song.ID, payloadSong.ID)
	})
}

func TestUpdateSong_WhenRequestHasAlbum_ShouldUpdateSongAndReorderOldAlbum(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	album := songData.Albums[2]
	request := requests.UpdateSongRequest{
		ID:          song.ID,
		Title:       "New Title",
		ReleaseDate: &[]time.Time{time.Now()}[0],
		AlbumID:     &album.ID,
		ArtistID:    album.ArtistID,
	}

	oldAlbumID := *song.AlbumID
	oldAlbumSongsCount := len(slices.DeleteFunc(slices.Clone(songData.Songs), func(s model.Song) bool {
		return s.AlbumID == nil || *s.AlbumID != oldAlbumID
	}))
	newAlbumSongsCount := len(slices.DeleteFunc(slices.Clone(songData.Songs), func(s model.Song) bool {
		return s.AlbumID == nil || *s.AlbumID != *request.AlbumID
	}))

	messages := utils.SubscribeToTopic(topics.SongsUpdatedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/songs", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var newSong model.Song
	db := utils.GetDatabase(t)
	db.Find(&newSong, request.ID)

	assertUpdatedSong(t, request, newSong)
	assert.Equal(t, uint(newAlbumSongsCount+1), *newSong.AlbumTrackNo)

	// check the old album to be ordered
	var newAlbum model.Album
	db.Preload("Songs", func(db *gorm.DB) *gorm.DB {
		return db.Order("album_track_no")
	}).Find(&newAlbum, oldAlbumID)

	assert.Len(t, newAlbum.Songs, oldAlbumSongsCount-1)
	for i, s := range newAlbum.Songs {
		assert.Equal(t, uint(i+1), *s.AlbumTrackNo)
	}

	assertion.AssertMessage(t, messages, func(payloadIDs uuid.UUID) {
		assert.Len(t, payloadIDs, 1)
		assert.Equal(t, song.ID, payloadIDs[0])
	})
}

func assertUpdatedSong(t *testing.T, request requests.UpdateSongRequest, song model.Song) {
	assert.Equal(t, request.ID, song.ID)
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.Equal(t, request.IsRecorded, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assertion.Time(t, request.ReleaseDate, song.ReleaseDate)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
	assert.Equal(t, request.ArtistID, song.ArtistID)
	assert.Equal(t, request.AlbumID, song.AlbumID)
}
