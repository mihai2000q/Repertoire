package album

import (
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	albumData "repertoire/server/test/integration/test/data/album"
	"repertoire/server/test/integration/test/utils"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAlbumsUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, albumData.Users, albumData.SeedData)

	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	ids := []uuid.UUID{
		albumData.Albums[0].ID,
		albumData.Albums[1].ID,
		albumData.Albums[2].ID,
	}

	// when
	err := utils.PublishToTopic(topics.AlbumsUpdatedTopic, ids)

	// then
	assert.NoError(t, err)

	db := utils.GetDatabase(t)
	var albums []model.Album
	db.Model(&model.Album{}).
		Joins("Artist").
		Preload("Songs").
		Preload("Songs.Artist").
		Preload("Songs.Album").
		Find(&albums, ids)

	// maps for lookup
	albumsMap := make(map[string]model.Album)
	songsMap := make(map[string]model.Song)
	for _, album := range albums {
		albumsMap[album.ID.String()] = album
		for _, song := range album.Songs {
			songsMap[song.ID.String()] = song
		}
	}

	assertion.AssertMessage(t, messages, func(documents []map[string]any) {
		assert.Len(t, documents, len(albumsMap)+len(songsMap))
		for _, doc := range documents {
			id := doc["id"].(string)
			if strings.HasPrefix(id, "album-") {
				albumSearch := utils.UnmarshalDocument[model.AlbumSearch](doc)
				album := albumsMap[strings.Replace(albumSearch.ID, "album-", "", 1)]
				assertion.AlbumSearch(t, albumSearch, album)
			} else if strings.HasPrefix(id, "song-") {
				songSearch := utils.UnmarshalDocument[model.SongSearch](doc)
				song := songsMap[strings.Replace(songSearch.ID, "song-", "", 1)]
				assertion.SongSearch(t, songSearch, song)
			} else {
				assert.Fail(t, "Document was not found in songs or albums")
			}
		}
	})
}
