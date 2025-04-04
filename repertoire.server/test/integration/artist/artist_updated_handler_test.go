package artist

import (
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	artistData "repertoire/server/test/integration/test/data/artist"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestArtistUpdated_WhenSuccessful_ShouldPublishMessage(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, artistData.Users, artistData.SeedData)

	messages := utils.SubscribeToTopic(topics.UpdateFromSearchEngineTopic)

	artistID := artistData.Artists[0].ID

	// when
	err := utils.PublishToTopic(topics.ArtistUpdatedTopic, artistID)

	// then
	assert.NoError(t, err)

	db := utils.GetDatabase(t)
	var artist model.Artist
	db.
		Preload("Albums").
		Preload("Albums.Artist").
		Preload("Songs").
		Preload("Songs.Artist").
		Preload("Songs.Album").
		Find(&artist, artistID)

	assertion.AssertMessage(t, messages, func(documents []any) {
		assert.Len(t, documents, len(artist.Albums)+len(artist.Songs)+1)

		artistSearch := utils.UnmarshallDocument[model.ArtistSearch](documents[0])
		assertion.ArtistSearch(t, artistSearch, artist)

		for i, song := range artist.Songs {
			songSearch := utils.UnmarshallDocument[model.SongSearch](documents[1+i])
			assertion.SongSearch(t, songSearch, song)
		}

		for i, album := range artist.Albums {
			albumSearch := utils.UnmarshallDocument[model.AlbumSearch](documents[1+len(artist.Songs)+i])
			assertion.AlbumSearch(t, albumSearch, album)
		}
	})
}
