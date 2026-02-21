package song

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSong_WhenSuccessful_ShouldCreateSong(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongRequest
	}{
		{
			"Minimal",
			requests.CreateSongRequest{
				Title: "New Song",
			},
		},
		{
			"Maximal",
			requests.CreateSongRequest{
				Title:          "New Song",
				Description:    "New Song Description",
				Bpm:            &[]uint{123}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-song"}[0],
				YoutubeLink:    &[]string{"https://youtu.be/9DyxtUCW84o?si=2pNX8eaV4KwKfOaF"}[0],
				ReleaseDate:    &[]internal.Date{internal.Date(time.Now())}[0],
				Difficulty:     &[]enums.Difficulty{enums.Hard}[0],
				GuitarTuningID: &[]uuid.UUID{songData.Users[0].GuitarTunings[0].ID}[0],
			},
		},
		{
			"With Artist",
			requests.CreateSongRequest{
				Title:    "New Song with Artist",
				ArtistID: &[]uuid.UUID{songData.Artists[0].ID}[0],
			},
		},
		{
			"With New Artist",
			requests.CreateSongRequest{
				Title:      "New Song with new Artist",
				ArtistName: &[]string{"New Artist Name"}[0],
			},
		},
		{
			"With Album",
			requests.CreateSongRequest{
				Title:   "New Song with Artist",
				AlbumID: &[]uuid.UUID{songData.Albums[0].ID}[0],
			},
		},
		{
			"With New Album",
			requests.CreateSongRequest{
				Title:       "New Song with new Artist",
				ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
				AlbumTitle:  &[]string{"New Album Title"}[0],
			},
		},
		{
			"With New Artist and New Album",
			requests.CreateSongRequest{
				Title:      "New Song with new Artist and album",
				ArtistName: &[]string{"New Artist Name"}[0],
				AlbumTitle: &[]string{"New Album Title"}[0],
			},
		},
		{
			"With SectionOccurrences",
			requests.CreateSongRequest{
				Title: "New Song with new Artist and album",
				Sections: []requests.CreateSectionRequest{
					{
						Name:   "Song Section 1",
						TypeID: songData.Users[0].SongSectionTypes[0].ID,
					},
					{
						Name:   "Song Section 2",
						TypeID: songData.Users[0].SongSectionTypes[1].ID,
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

			user := songData.Users[0]

			messages := utils.SubscribeToTopic(topics.SongCreatedTopic)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithUser(user).
				POST(w, "/api/songs", test.request)

			// then
			var response struct{ ID uuid.UUID }
			_ = json.Unmarshal(w.Body.Bytes(), &response)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.NotEmpty(t, response)

			db := utils.GetDatabase(t)
			var song model.Song
			db.
				Joins("Settings").
				Joins("Artist").
				Joins("Album").
				Joins("GuitarTuning").
				Preload("Album.Songs").
				Preload("Sections").
				Preload("Sections.SongSectionType").
				Find(&song, response.ID)
			assertCreatedSong(t, test.request, song, user.ID)

			assertion.AssertMessage(t, messages, func(payloadSong model.Song) {
				assert.Equal(t, song.ID, payloadSong.ID)
			})
		})
	}
}

func assertCreatedSong(
	t *testing.T,
	request requests.CreateSongRequest,
	song model.Song,
	userID uuid.UUID,
) {
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.False(t, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assert.Nil(t, song.LastTimePlayed)
	if request.ReleaseDate != nil {
		assertion.Date(t, request.ReleaseDate, song.ReleaseDate)
	}
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Nil(t, song.ImageURL)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
	assert.Equal(t, userID, song.UserID)

	// assert settings
	assert.NotEmpty(t, song.Settings.ID)

	// assert arrangement
	assert.NotEmpty(t, song.DefaultArrangementID)

	assert.Len(t, song.Arrangements, 1)
	assert.Equal(t, song.Arrangements[0].ID, *song.DefaultArrangementID)
	assert.NotEmpty(t, song.Arrangements[0].ID)
	assert.Equal(t, model.DefaultSongArrangementName, song.Arrangements[0].Name)
	assert.Equal(t, uint(0), song.Arrangements[0].Order)
	assert.Equal(t, song.ID, song.Arrangements[0].SongID)

	assert.Len(t, request.Sections, len(song.Sections))
	for i, sectionRequest := range request.Sections {
		assert.NotEmpty(t, song.Sections[i].ID)
		assert.Equal(t, sectionRequest.Name, song.Sections[i].Name)
		assert.Zero(t, song.Sections[i].Rehearsals)
		assert.Equal(t, model.DefaultSongSectionConfidence, song.Sections[i].Confidence)
		assert.Zero(t, song.Sections[i].RehearsalsScore)
		assert.Zero(t, song.Sections[i].ConfidenceScore)
		assert.Zero(t, song.Sections[i].Progress)
		assert.Equal(t, uint(i), song.Sections[i].Order)
		assert.Equal(t, sectionRequest.TypeID, song.Sections[i].SongSectionTypeID)
		assert.Equal(t, song.ID, song.Sections[i].SongID)

		// assert section occurrences on arrangement
		assert.Zero(t, song.Arrangements[0].SectionOccurrences[i].Occurrences)
		assert.Equal(t, song.Sections[i].ID, song.Arrangements[0].SectionOccurrences[i].SectionID)
		assert.Equal(t, song.Arrangements[0].ID, song.Arrangements[0].SectionOccurrences[i].ArrangementID)
	}

	if request.ArtistID != nil {
		assert.Equal(t, request.ArtistID, song.ArtistID)
		assert.NotNil(t, song.Artist)
	}

	if request.ArtistName != nil {
		assert.NotNil(t, song.Artist)
		assert.Equal(t, *song.ArtistID, song.Artist.ID)
		assert.NotEmpty(t, song.Artist.ID)
		assert.Equal(t, *request.ArtistName, song.Artist.Name)
		assert.Equal(t, song.UserID, song.Artist.UserID)
	}

	if request.AlbumTitle != nil {
		assert.NotEmpty(t, song.Album.ID)
		assert.Equal(t, *request.AlbumTitle, song.Album.Title)
		assert.Equal(t, song.ArtistID, song.Album.ArtistID)
		assert.Equal(t, song.UserID, song.Album.UserID)
		assertion.Date(t, request.ReleaseDate, song.ReleaseDate)
		assert.Equal(t, uint(1), *song.AlbumTrackNo)
	}

	if request.AlbumID != nil {
		assert.NotNil(t, song.Album)
		assert.Equal(t, uint(len(song.Album.Songs)), *song.AlbumTrackNo)
		assert.Equal(t, song.Album.ArtistID, song.ArtistID)
		if request.ReleaseDate == nil {
			assertion.Date(t, song.Album.ReleaseDate, song.ReleaseDate)
		}
	}

	if request.AlbumID == nil && request.AlbumTitle == nil {
		assert.Nil(t, song.AlbumTrackNo)
	}
}
