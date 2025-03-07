package assertion

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"repertoire/server/internal/enums"
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

	expiresIn, _ := time.ParseDuration(env.JwtExpirationTime)

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
	assert.WithinDuration(t, time.Now().Add(expiresIn).UTC(), exp.Time, 10*time.Second)
}

func AssertMessage[T any](
	t *testing.T,
	messages <-chan *message.Message,
	assertFunc func(T),
) {
	select {
	case msg := <-messages:
		var unmarshalledPayload T
		_ = json.Unmarshal(msg.Payload, &unmarshalledPayload)
		assertFunc(unmarshalledPayload)
	case <-time.After(5 * time.Millisecond):
		t.Fatal("Timed out waiting for message")
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
			ResponseArtist(t, *album.Artist, *response.Artist, false)
		} else {
			assert.Nil(t, response.Artist)
		}
	}

	if withSongs {
		for i := 0; i < len(album.Songs); i++ {
			ResponseSong(
				t,
				album.Songs[i],
				response.Songs[i],
				false,
				false,
				false,
				false,
			)
		}
	}
}

func ResponseArtist(t *testing.T, artist model.Artist, response model.Artist, withBandMembers bool) {
	assert.Equal(t, artist.ID, response.ID)
	assert.Equal(t, artist.Name, response.Name)
	assert.Equal(t, artist.ImageURL, response.ImageURL)

	if withBandMembers {
		for i := 0; i < len(artist.BandMembers); i++ {
			ResponseBandMember(t, artist.BandMembers[i], response.BandMembers[i], true)
		}
	}
}

func ResponseBandMember(t *testing.T, bandMember model.BandMember, response model.BandMember, withRoles bool) {
	assert.Equal(t, bandMember.ID, response.ID)
	assert.Equal(t, bandMember.Name, response.Name)
	assert.Equal(t, bandMember.Color, response.Color)
	assert.Equal(t, bandMember.ImageURL, response.ImageURL)
	if withRoles {
		for i := 0; i < len(bandMember.Roles); i++ {
			ResponseBandMemberRole(t, bandMember.Roles[i], response.Roles[i])
		}
	}
}

func ResponseBandMemberRole(t *testing.T, bandMemberRole model.BandMemberRole, response model.BandMemberRole) {
	assert.Equal(t, bandMemberRole.ID, response.ID)
	assert.Equal(t, bandMemberRole.Name, response.Name)
}

func ResponseSong(
	t *testing.T,
	song model.Song,
	response model.Song,
	withAlbum bool,
	withArtist bool,
	withAssociations bool,
	withSongSectionsDetails bool,
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
	assert.Equal(t, song.Rehearsals, response.Rehearsals)
	assert.Equal(t, song.Confidence, response.Confidence)
	assert.Equal(t, song.Progress, response.Progress)

	if withAlbum {
		if song.Album != nil {
			ResponseAlbum(t, *song.Album, *response.Album, false, false)
		} else {
			assert.Nil(t, response.Album)
		}
	}

	if withArtist {
		if song.Artist != nil {
			ResponseArtist(t, *song.Artist, *response.Artist, true)
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
			ResponseSongSection(t, song.Sections[i], response.Sections[i], withSongSectionsDetails)
		}

		for i := range song.Playlists {
			ResponsePlaylist(t, song.Playlists[i], response.Playlists[i], false)
		}
	}
}

func ResponseGuitarTuning(t *testing.T, guitarTuning model.GuitarTuning, response model.GuitarTuning) {
	assert.Equal(t, guitarTuning.ID, response.ID)
	assert.Equal(t, guitarTuning.Name, response.Name)
}

func ResponseInstrument(t *testing.T, instrument model.Instrument, response model.Instrument) {
	assert.Equal(t, instrument.ID, response.ID)
	assert.Equal(t, instrument.Name, response.Name)
}

func ResponseSongSection(
	t *testing.T,
	songSection model.SongSection,
	response model.SongSection,
	withBandMember bool,
) {
	assert.Equal(t, songSection.ID, response.ID)
	assert.Equal(t, songSection.Name, response.Name)
	assert.Equal(t, songSection.Occurrences, response.Occurrences)
	assert.Equal(t, songSection.Rehearsals, response.Rehearsals)
	assert.Equal(t, songSection.Confidence, response.Confidence)
	assert.Equal(t, songSection.RehearsalsScore, response.RehearsalsScore)
	assert.Equal(t, songSection.ConfidenceScore, response.ConfidenceScore)
	assert.Equal(t, songSection.Progress, response.Progress)

	ResponseSongSectionType(t, songSection.SongSectionType, response.SongSectionType)
	if songSection.Instrument != nil {
		ResponseInstrument(t, *songSection.Instrument, *response.Instrument)
	} else {
		assert.Nil(t, response.Instrument)
	}
	if withBandMember {
		if songSection.BandMember != nil {
			ResponseBandMember(t, *songSection.BandMember, *response.BandMember, true)
		} else {
			assert.Nil(t, response.BandMember)
		}
	}
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
		ResponseSong(
			t,
			playlist.Songs[i],
			response.Songs[i],
			withSongsMetadata,
			withSongsMetadata,
			false,
			false,
		)
		if withSongsMetadata {
			// making sure the After Find hook works
			assert.Equal(t, playlist.PlaylistSongs[i].SongID, response.Songs[i].ID)
			assert.Equal(t, playlist.PlaylistSongs[i].SongTrackNo, response.Songs[i].PlaylistTrackNo)
			Time(t, &playlist.PlaylistSongs[i].CreatedAt, &response.Songs[i].PlaylistCreatedAt)
		}
	}
}

func ResponseUser(t *testing.T, user model.User, response model.User) {
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.ProfilePictureURL, response.ProfilePictureURL)
}

// Search

func ArtistSearch(t *testing.T, artistSearch model.ArtistSearch, artist model.Artist) {
	assert.Equal(t, "artist-"+artist.ID.String(), artistSearch.ID)
	assert.Equal(t, artist.Name, artistSearch.Name)
	assert.Equal(t, artist.ImageURL.StripURL(), artistSearch.ImageUrl)
	assert.Equal(t, artist.UpdatedAt, artistSearch.UpdatedAt.Local())
	assert.Equal(t, enums.Artist, artistSearch.Type)
}

func AlbumSearch(t *testing.T, albumSearch model.AlbumSearch, album model.Album) {
	assert.Equal(t, "album-"+album.ID.String(), albumSearch.ID)
	assert.Equal(t, album.Title, albumSearch.Title)
	assert.Equal(t, album.ImageURL.StripURL(), albumSearch.ImageUrl)
	assert.Equal(t, album.UpdatedAt, albumSearch.UpdatedAt.Local())
	assert.Equal(t, enums.Album, albumSearch.Type)

	if album.Artist != nil {
		assert.Equal(t, album.Artist.ID, albumSearch.Artist.ID)
		assert.Equal(t, album.Artist.Name, albumSearch.Artist.Name)
		assert.Equal(t, album.Artist.UpdatedAt, albumSearch.Artist.UpdatedAt.Local())
		assert.Equal(t, album.Artist.ImageURL.StripURL(), albumSearch.Artist.ImageUrl)
	} else {
		assert.Nil(t, albumSearch.Artist)
	}
}

func SongSearch(t *testing.T, songSearch model.SongSearch, song model.Song) {
	assert.Equal(t, "song-"+song.ID.String(), songSearch.ID)
	assert.Equal(t, song.Title, songSearch.Title)
	assert.Equal(t, song.ImageURL.StripURL(), songSearch.ImageUrl)
	assert.Equal(t, song.UpdatedAt, songSearch.UpdatedAt.Local())
	assert.Equal(t, enums.Song, songSearch.Type)

	if song.Artist != nil {
		assert.Equal(t, song.Artist.ID, songSearch.Artist.ID)
		assert.Equal(t, song.Artist.Name, songSearch.Artist.Name)
		assert.Equal(t, song.Artist.UpdatedAt, songSearch.Artist.UpdatedAt.Local())
		assert.Equal(t, song.Artist.ImageURL.StripURL(), songSearch.Artist.ImageUrl)
	} else {
		assert.Nil(t, songSearch.Artist)
	}

	if song.Album != nil {
		assert.Equal(t, song.Album.ID, songSearch.Album.ID)
		assert.Equal(t, song.Album.Title, songSearch.Album.Title)
		assert.Equal(t, song.Album.UpdatedAt, songSearch.Album.UpdatedAt.Local())
		assert.Equal(t, song.Album.ImageURL.StripURL(), songSearch.Album.ImageUrl)
	} else {
		assert.Nil(t, songSearch.Album)
	}
}

func PlaylistSearch(t *testing.T, playlistSearch model.PlaylistSearch, playlist model.Playlist) {
	assert.Equal(t, "playlist-"+playlist.ID.String(), playlistSearch.ID)
	assert.Equal(t, playlist.Title, playlistSearch.Title)
	assert.Equal(t, playlist.ImageURL.StripURL(), playlistSearch.ImageUrl)
	assert.Equal(t, playlist.UpdatedAt, playlistSearch.UpdatedAt.Local())
	assert.Equal(t, enums.Playlist, playlistSearch.Type)
}
