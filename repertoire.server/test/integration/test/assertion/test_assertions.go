package assertion

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Token(t *testing.T, actual string) {
	env := utils.GetEnv()

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(env.JwtPublicKey))
	assert.NoError(t, err)
	token, err := jwt.Parse(actual, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, token.Valid)
}

func Time(t *testing.T, expected *time.Time, actual *time.Time) {
	if expected != nil {
		assert.WithinDuration(t, *expected, *actual, 1*time.Minute)
	} else {
		assert.Nil(t, actual)
	}
}

func AssertMessage[T any](
	t *testing.T,
	message utils.SubscribedToTopic,
	assertFunc func(T),
) {
	select {
	case msg := <-message.Messages:
		if msg.Metadata.Get("topic") != string(message.Topic) {
			return
		}
		var unmarshalledPayload T
		_ = json.Unmarshal(msg.Payload, &unmarshalledPayload)
		assertFunc(unmarshalledPayload)
	case <-time.After(5 * time.Second):
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
	assert.Equal(t, artist.IsBand, response.IsBand)
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
		ResponseSongSettings(t, song.Settings, response.Settings)

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

func ResponseSongSettings(t *testing.T, settings model.SongSettings, response model.SongSettings) {
	assert.Equal(t, settings.ID, response.ID)
	if settings.DefaultInstrument != nil {
		ResponseInstrument(t, *settings.DefaultInstrument, *response.DefaultInstrument)
	} else {
		assert.Nil(t, response.DefaultInstrument)
	}
	if settings.DefaultBandMember != nil {
		ResponseBandMember(t, *settings.DefaultBandMember, *response.DefaultBandMember, false)
	} else {
		assert.Nil(t, response.DefaultBandMember)
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

func ResponseEnhancedPlaylist(t *testing.T, playlist model.Playlist, response model.EnhancedPlaylist) {
	assert.Equal(t, playlist.ID, response.ID)
	assert.Equal(t, playlist.Title, response.Title)
	assert.Equal(t, playlist.Description, response.Description)
	assert.Equal(t, playlist.ImageURL, response.ImageURL)

	assert.Len(t, response.SongsIDs, len(playlist.Songs))
	assert.Equal(t, response.SongsCount, len(playlist.Songs))
	for i := range playlist.Songs {
		assert.Equal(t, playlist.Songs[i].ID, response.SongsIDs[i])
	}
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
	if album.ReleaseDate != nil {
		Time(t, album.ReleaseDate, albumSearch.ReleaseDate)
	} else {
		assert.Nil(t, albumSearch.ReleaseDate)
	}
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
	if song.ReleaseDate != nil {
		Time(t, song.ReleaseDate, songSearch.ReleaseDate)
	} else {
		assert.Nil(t, songSearch.ReleaseDate)
	}
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
		if song.Album.ReleaseDate != nil {
			Time(t, song.Album.ReleaseDate, songSearch.Album.ReleaseDate)
		} else {
			assert.Nil(t, songSearch.Album.ReleaseDate)
		}
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
