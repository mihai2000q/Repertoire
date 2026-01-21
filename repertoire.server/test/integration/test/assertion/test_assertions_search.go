package assertion

import (
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ArtistSearchID(t *testing.T, id uuid.UUID, searchID string) {
	assert.Equal(t, "artist-"+id.String(), searchID)
}

func AlbumSearchID(t *testing.T, id uuid.UUID, searchID string) {
	assert.Equal(t, "album-"+id.String(), searchID)
}

func SongSearchID(t *testing.T, id uuid.UUID, searchID string) {
	assert.Equal(t, "song-"+id.String(), searchID)
}

func PlaylistSearchID(t *testing.T, id uuid.UUID, searchID string) {
	assert.Equal(t, "playlist-"+id.String(), searchID)
}

func ArtistSearch(t *testing.T, artistSearch model.ArtistSearch, artist model.Artist) {
	ArtistSearchID(t, artist.ID, artistSearch.ID)
	assert.Equal(t, artist.Name, artistSearch.Name)
	assert.Equal(t, artist.ImageURL.StripURL(), artistSearch.ImageUrl)
	Time(t, &artist.UpdatedAt, &artistSearch.UpdatedAt)
	Time(t, &artist.CreatedAt, &artistSearch.CreatedAt)
	assert.Equal(t, enums.Artist, artistSearch.Type)
}

func AlbumSearch(t *testing.T, albumSearch model.AlbumSearch, album model.Album) {
	AlbumSearchID(t, album.ID, albumSearch.ID)
	assert.Equal(t, album.Title, albumSearch.Title)
	dateAndString(t, album.ReleaseDate, albumSearch.ReleaseDate)
	assert.Equal(t, album.ImageURL.StripURL(), albumSearch.ImageUrl)
	Time(t, &album.UpdatedAt, &albumSearch.UpdatedAt)
	Time(t, &album.CreatedAt, &albumSearch.CreatedAt)
	assert.Equal(t, enums.Album, albumSearch.Type)

	if album.Artist != nil {
		assert.Equal(t, album.Artist.ID, albumSearch.Artist.ID)
		assert.Equal(t, album.Artist.Name, albumSearch.Artist.Name)
		Time(t, &album.Artist.UpdatedAt, &albumSearch.Artist.UpdatedAt)
		assert.Equal(t, album.Artist.ImageURL.StripURL(), albumSearch.Artist.ImageUrl)
	} else {
		assert.Nil(t, albumSearch.Artist)
	}
}

func SongSearch(t *testing.T, songSearch model.SongSearch, song model.Song) {
	SongSearchID(t, song.ID, songSearch.ID)
	assert.Equal(t, song.Title, songSearch.Title)
	dateAndString(t, song.ReleaseDate, songSearch.ReleaseDate)
	assert.Equal(t, song.ImageURL.StripURL(), songSearch.ImageUrl)
	Time(t, &song.UpdatedAt, &songSearch.UpdatedAt)
	Time(t, &song.CreatedAt, &songSearch.CreatedAt)
	assert.Equal(t, enums.Song, songSearch.Type)

	if song.Artist != nil {
		assert.Equal(t, song.Artist.ID, songSearch.Artist.ID)
		assert.Equal(t, song.Artist.Name, songSearch.Artist.Name)
		Time(t, &song.Artist.UpdatedAt, &songSearch.Artist.UpdatedAt)
		assert.Equal(t, song.Artist.ImageURL.StripURL(), songSearch.Artist.ImageUrl)
	} else {
		assert.Nil(t, songSearch.Artist)
	}

	if song.Album != nil {
		assert.Equal(t, song.Album.ID, songSearch.Album.ID)
		assert.Equal(t, song.Album.Title, songSearch.Album.Title)
		dateAndString(t, song.Album.ReleaseDate, songSearch.Album.ReleaseDate)
		Time(t, &song.Album.UpdatedAt, &songSearch.Album.UpdatedAt)
		assert.Equal(t, song.Album.ImageURL.StripURL(), songSearch.Album.ImageUrl)
	} else {
		assert.Nil(t, songSearch.Album)
	}
}

func PlaylistSearch(t *testing.T, playlistSearch model.PlaylistSearch, playlist model.Playlist) {
	PlaylistSearchID(t, playlist.ID, playlistSearch.ID)
	assert.Equal(t, playlist.Title, playlistSearch.Title)
	assert.Equal(t, playlist.ImageURL.StripURL(), playlistSearch.ImageUrl)
	Time(t, &playlist.UpdatedAt, &playlist.UpdatedAt)
	Time(t, &playlist.CreatedAt, &playlist.CreatedAt)
	assert.Equal(t, enums.Playlist, playlistSearch.Type)
}

func dateAndString(t *testing.T, date *internal.Date, str *string) {
	if date != nil {
		assert.NotNil(t, str)
		assert.Equal(t, (*time.Time)(date).Format("2006-01-02"), *str)
	} else {
		assert.Nil(t, str)
	}
}
