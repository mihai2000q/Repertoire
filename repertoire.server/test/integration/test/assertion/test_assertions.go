package assertion

import (
	"repertoire/server/model"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Time(t *testing.T, expected *time.Time, actual *time.Time) {
	if expected != nil {
		assert.WithinDuration(t, *expected, *actual, 1*time.Minute)
	} else {
		assert.Nil(t, actual)
	}
}

func Token(t *testing.T, actual string, user model.User) {
	env := utils.GetEnv()

	// get token
	token, err := jwt.Parse(actual, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.JwtSecretKey), nil
	})
	assert.NoError(t, err)

	jtiClaim := token.Claims.(jwt.MapClaims)["jti"].(string)
	jti, err := uuid.Parse(jtiClaim)
	assert.NoError(t, err)
	sub, err := token.Claims.GetSubject()
	assert.NoError(t, err)
	aud, err := token.Claims.GetAudience()
	assert.NoError(t, err)
	iss, err := token.Claims.GetIssuer()
	assert.NoError(t, err)
	iat, err := token.Claims.GetIssuedAt()
	assert.NoError(t, err)
	exp, err := token.Claims.GetExpirationTime()
	assert.NoError(t, err)

	assert.Equal(t, jwt.SigningMethodHS256, token.Method)
	assert.NotEmpty(t, jti)
	assert.Equal(t, user.ID.String(), sub)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.JwtAudience, aud[0])
	assert.Equal(t, env.JwtIssuer, iss)
	assert.WithinDuration(t, time.Now().UTC(), iat.Time, 10*time.Second)
	assert.WithinDuration(t, time.Now().Add(time.Hour).UTC(), exp.Time, 10*time.Second)
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
		if song.GuitarTuning != nil {
			ResponseGuitarTuning(t, *song.GuitarTuning, *response.GuitarTuning)
		} else {
			assert.Nil(t, response.GuitarTuning)
		}

		for i := range song.Sections {
			ResponseSongSection(t, song.Sections[i], response.Sections[i])
		}
	}
}

func ResponseGuitarTuning(t *testing.T, guitarTuning model.GuitarTuning, response model.GuitarTuning) {
	assert.Equal(t, guitarTuning.ID, response.ID)
	assert.Equal(t, guitarTuning.Name, response.Name)
}

func ResponseSongSection(t *testing.T, songSection model.SongSection, response model.SongSection) {
	assert.Equal(t, songSection.ID, response.ID)
	assert.Equal(t, songSection.Name, response.Name)
	assert.Equal(t, songSection.Rehearsals, response.Rehearsals)

	ResponseSongSectionType(t, songSection.SongSectionType, response.SongSectionType)
}

func ResponseSongSectionType(t *testing.T, songSectionType model.SongSectionType, response model.SongSectionType) {
	assert.Equal(t, songSectionType.ID, response.ID)
	assert.Equal(t, songSectionType.Name, response.Name)
}

func ResponsePlaylist(t *testing.T, playlist model.Playlist, response model.Playlist, withSongsMetadata bool) {
	assert.Equal(t, playlist.ID, response.ID)
	assert.Equal(t, playlist.Title, response.Title)
	assert.Equal(t, playlist.Description, response.Description)
	assert.Equal(t, playlist.ImageURL, response.ImageURL)

	for i := range playlist.Songs {
		ResponseSong(t, playlist.Songs[i], response.Songs[i], withSongsMetadata, withSongsMetadata, false)
		if withSongsMetadata {
			// making sure the After Find hook works
			assert.Equal(t, playlist.PlaylistSongs[i].SongID, response.Songs[i].ID)
			assert.Equal(t, playlist.PlaylistSongs[i].SongTrackNo, response.Songs[i].PlaylistTrackNo)
			assert.Equal(t, playlist.PlaylistSongs[i].CreatedAt, response.Songs[i].PlaylistCreatedAt)
		}
	}
}

func ResponseUser(t *testing.T, user model.User, response model.User) {
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.ProfilePictureURL, response.ProfilePictureURL)
}
