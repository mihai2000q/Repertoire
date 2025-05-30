package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSong_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewCreateSong(jwtService, nil, nil, nil)

	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"

	forbiddenError := wrapper.UnauthorizedError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreateSong_WhenGetAlbumWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewCreateSong(jwtService, nil, albumRepository, nil)

	request := requests.CreateSongRequest{
		Title:   "Some Song",
		AlbumID: &[]uuid.UUID{uuid.New()}[0],
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", mock.Anything, *request.AlbumID).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestCreateSong_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewCreateSong(jwtService, songRepository, albumRepository, nil)

	request := requests.CreateSongRequest{
		Title:   "Some Song",
		AlbumID: &[]uuid.UUID{uuid.New()}[0],
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	albumRepository.On("GetWithSongs", mock.Anything, *request.AlbumID).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSong_WhenCreateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewCreateSong(jwtService, songRepository, nil, nil)

	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	songRepository.On("Create", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSong_WhenCreateSongFails_ShouldReturnBadRequestError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewCreateSong(jwtService, songRepository, nil, messagePublisherService)

	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	songRepository.On("Create", mock.IsType(new(model.Song))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongCreatedTopic, mock.IsType(model.Song{})).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestCreateSong_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	requestAlbumID := uuid.New()

	tests := []struct {
		name    string
		request requests.CreateSongRequest
		album   *model.Album
	}{
		{
			"Create song only with title",
			requests.CreateSongRequest{
				Title: "Some Song",
			},
			nil,
		},
		{
			"Create song with more fields",
			requests.CreateSongRequest{
				Title:          "Some Song",
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-song"}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
			},
			nil,
		},
		{
			"Create song with new album and artist",
			requests.CreateSongRequest{
				Title:       "Some Song",
				ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
				AlbumTitle:  &[]string{"New Album Title"}[0],
				ArtistName:  &[]string{"New Artist Name"}[0],
			},
			nil,
		},
		{
			"Create song with existing album",
			requests.CreateSongRequest{
				Title:   "Some Song",
				AlbumID: &requestAlbumID,
			},
			&model.Album{
				ID:          requestAlbumID,
				ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
				Songs:       []model.Song{{}, {}, {}, {}, {}},
			},
		},
		{
			"Create song with existing album that has an artist",
			requests.CreateSongRequest{
				Title:   "Some Song",
				AlbumID: &requestAlbumID,
			},
			&model.Album{
				ID:       requestAlbumID,
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"Create song with existing artist",
			requests.CreateSongRequest{
				Title:    "Some Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			nil,
		},
		{
			"Create song with new album and existing artist",
			requests.CreateSongRequest{
				Title:      "Some Song",
				AlbumTitle: &[]string{"New Album Title"}[0],
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
			},
			nil,
		},
		{
			"Create song with sections",
			requests.CreateSongRequest{
				Title: "Some Song",
				Sections: []requests.CreateSectionRequest{
					{Name: "First Section", TypeID: uuid.New()},
					{Name: "Second Section", TypeID: uuid.New()},
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			jwtService := new(service.JwtServiceMock)
			songRepository := new(repository.SongRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := song.NewCreateSong(
				jwtService,
				songRepository,
				albumRepository,
				messagePublisherService,
			)

			token := "this is a token"

			// given - mocks
			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			if tt.request.AlbumID != nil {
				albumRepository.On("GetWithSongs", mock.IsType(tt.album), *tt.request.AlbumID).
					Return(nil, tt.album).
					Once()
			}

			var createdSong model.Song
			songRepository.On("Create", mock.IsType(new(model.Song))).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					assertCreatedSong(t, tt.request, *newSong, userID, tt.album)
					createdSong = *newSong
				}).
				Return(nil).
				Once()

			if tt.request.AlbumID != nil {
				mockAlbum := model.Album{
					ID:    *tt.request.AlbumID,
					Title: "Some Title",
				}
				albumRepository.On("Get", new(model.Album), *tt.request.AlbumID).
					Return(nil, &mockAlbum).
					Once()
			}

			messagePublisherService.On("Publish", topics.SongCreatedTopic, mock.IsType(createdSong)).
				Run(func(args mock.Arguments) {
					assert.Equal(t, createdSong, args.Get(1).(model.Song))
				}).
				Return(nil).
				Once()

			// when
			id, errCode := _uut.Handle(tt.request, token)

			// then
			assert.Nil(t, errCode)
			assert.Equal(t, id, createdSong.ID)

			jwtService.AssertExpectations(t)
			songRepository.AssertExpectations(t)
		})
	}
}

func assertCreatedSong(
	t *testing.T,
	request requests.CreateSongRequest,
	song model.Song,
	userID uuid.UUID,
	album *model.Album,
) {
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.False(t, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Nil(t, song.LastTimePlayed)
	if request.ReleaseDate != nil {
		assert.Equal(t, request.ReleaseDate, song.ReleaseDate)
	}
	assert.Nil(t, song.ImageURL)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
	assert.Equal(t, request.AlbumID, song.AlbumID)
	assert.Equal(t, userID, song.UserID)
	assert.Len(t, request.Sections, len(song.Sections))

	assert.NotEmpty(t, song.Settings.ID)

	for i, section := range request.Sections {
		assert.NotEmpty(t, song.Sections[i].ID)
		assert.Equal(t, section.Name, song.Sections[i].Name)
		assert.Zero(t, song.Sections[i].Rehearsals)
		assert.Equal(t, model.DefaultSongSectionConfidence, song.Sections[i].Confidence)
		assert.Zero(t, song.Sections[i].RehearsalsScore)
		assert.Zero(t, song.Sections[i].ConfidenceScore)
		assert.Zero(t, song.Sections[i].Progress)
		assert.Equal(t, uint(i), song.Sections[i].Order)
		assert.Equal(t, section.TypeID, song.Sections[i].SongSectionTypeID)
		assert.Equal(t, song.ID, song.Sections[i].SongID)
	}
	if request.AlbumTitle != nil {
		assert.NotNil(t, song.Album)
		assert.NotEmpty(t, song.Album.ID)
		assert.Equal(t, *request.AlbumTitle, song.Album.Title)
		assert.Equal(t, song.ArtistID, song.Album.ArtistID)
		assert.Equal(t, song.UserID, song.Album.UserID)
		assert.Equal(t, song.ReleaseDate, song.Album.ReleaseDate)
		assert.Equal(t, uint(1), *song.AlbumTrackNo)
	}
	if request.ArtistName != nil {
		assert.NotNil(t, song.Artist)
		assert.Equal(t, *song.ArtistID, song.Artist.ID)
		assert.NotEmpty(t, song.Artist.ID)
		assert.Equal(t, *request.ArtistName, song.Artist.Name)
		assert.Equal(t, song.UserID, song.Artist.UserID)
	}
	if request.ArtistID != nil {
		assert.Equal(t, request.ArtistID, song.ArtistID)
	}
	if request.AlbumID != nil {
		assert.Equal(t, &[]uint{uint(len(album.Songs)) + 1}[0], song.AlbumTrackNo)
		assert.Equal(t, album.ArtistID, song.ArtistID)
		if request.ReleaseDate == nil {
			assert.Equal(t, album.ReleaseDate, song.ReleaseDate)
		}
	}
	if request.AlbumID == nil && request.AlbumTitle == nil {
		assert.Nil(t, song.AlbumTrackNo)
	}
}
