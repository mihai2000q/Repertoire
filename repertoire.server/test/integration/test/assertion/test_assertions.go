package assertion

import (
	"github.com/stretchr/testify/assert"
	"repertoire/server/model"
	"testing"
	"time"
)

func Time(t *testing.T, expected *time.Time, actual *time.Time) {
	if expected != nil {
		assert.WithinDuration(t, *expected, *actual, 10*time.Second)
	} else {
		assert.Nil(t, actual)
	}
}

// models

func ResponseAlbum(t *testing.T, album model.Album, response model.Album, withArtist bool, withSongs bool) {
	assert.Equal(t, album.ID, response.ID)
	assert.Equal(t, album.Title, response.Title)
	Time(t, album.ReleaseDate, response.ReleaseDate)
	assert.Equal(t, album.ImageURL, response.ImageURL)

	if withArtist {
		if album.Artist != nil {
			ResponseArtist(t, *album.Artist, *response.Artist)
		} else {
			assert.Nil(t, response.Artist)
		}
	}

	if withSongs {
		for i := 0; i < len(album.Songs); i++ {
			ResponseSong(t, album.Songs[i], response.Songs[i], false, false, false)
		}
	}
}

func ResponseArtist(t *testing.T, artist model.Artist, response model.Artist) {
	assert.Equal(t, artist.ID, response.ID)
	assert.Equal(t, artist.Name, response.Name)
	assert.Equal(t, artist.ImageURL, response.ImageURL)
}

func ResponseSong(
	t *testing.T,
	song model.Song,
	response model.Song,
	withAlbum bool,
	withArtist bool,
	withAssociations bool,
) {
	assert.Equal(t, song.ID, response.ID)
	assert.Equal(t, song.Title, response.Title)
	assert.Equal(t, song.Description, response.Description)
	Time(t, song.ReleaseDate, response.ReleaseDate)
	assert.Equal(t, song.ImageURL, response.ImageURL)
	assert.Equal(t, song.IsRecorded, response.IsRecorded)
	assert.Equal(t, song.Bpm, response.Bpm)
	assert.Equal(t, song.Difficulty, response.Difficulty)
	assert.Equal(t, song.SongsterrLink, response.SongsterrLink)
	assert.Equal(t, song.YoutubeLink, response.YoutubeLink)
	assert.Equal(t, song.AlbumTrackNo, response.AlbumTrackNo)

	if withAlbum {
		if song.Album != nil {
			ResponseAlbum(t, *song.Album, *response.Album, false, false)
		} else {
			assert.Nil(t, response.Album)
		}
	}

	if withArtist {
		if song.Artist != nil {
			ResponseArtist(t, *song.Artist, *response.Artist)
		} else {
			assert.Nil(t, response.Artist)
		}
	}

	if withAssociations {
		assert.Equal(t, song.GuitarTuning.ID, response.GuitarTuning.ID)
		assert.Equal(t, song.GuitarTuning.Name, response.GuitarTuning.Name)

		for i := range song.Sections {
			ResponseSongSection(t, song.Sections[i], response.Sections[i])
		}
	}
}

func ResponseSongSection(t *testing.T, songSection model.SongSection, response model.SongSection) {
	assert.Equal(t, songSection.ID, response.ID)
	assert.Equal(t, songSection.Name, response.Name)
	assert.Equal(t, songSection.Rehearsals, response.Rehearsals)

	assert.Equal(t, songSection.SongSectionType.ID, response.SongSectionType.ID)
	assert.Equal(t, songSection.SongSectionType.Name, response.SongSectionType.Name)
}

func ResponsePlaylist(t *testing.T, playlist model.Playlist, response model.Playlist, withSongs bool) {
	assert.Equal(t, playlist.ID, response.ID)
	assert.Equal(t, playlist.Title, response.Title)
	assert.Equal(t, playlist.Description, response.Description)
	assert.Equal(t, playlist.ImageURL, response.ImageURL)

	if withSongs {
		for i := range playlist.Songs {
			ResponseSong(t, playlist.Songs[i], response.Songs[i], true, true, false)
		}
	}
}

func ResponseUser(t *testing.T, user model.User, response model.User) {
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.ProfilePictureURL, response.ProfilePictureURL)
}
