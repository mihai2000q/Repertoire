package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSong_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewCreateSong(jwtService, nil, nil)

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
	_uut := song.NewCreateSong(jwtService, nil, albumRepository)

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
	_uut := song.NewCreateSong(jwtService, songRepository, albumRepository)

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
	_uut := song.NewCreateSong(jwtService, songRepository, nil)

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

func TestCreateSong_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongRequest
	}{
		{
			"Create song only with title",
			requests.CreateSongRequest{
				Title: "Some Song",
			},
		},
		{
			"Create song with more fields",
			requests.CreateSongRequest{
				Title:          "Some Song",
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-song"}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"Create song with new album and artist",
			requests.CreateSongRequest{
				Title:      "Some Song",
				AlbumTitle: &[]string{"New Album Title"}[0],
				ArtistName: &[]string{"New Artist Name"}[0],
			},
		},
		{
			"Create song with existing album",
			requests.CreateSongRequest{
				Title:   "Some Song",
				AlbumID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"Create song with new album and existing artist",
			requests.CreateSongRequest{
				Title:      "Some Song",
				AlbumTitle: &[]string{"New Album Title"}[0],
				ArtistID:   &[]uuid.UUID{uuid.New()}[0],
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			jwtService := new(service.JwtServiceMock)
			_uut := song.NewCreateSong(jwtService, songRepository, albumRepository)

			token := "this is a token"

			// given - mocks
			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			var album *model.Album
			if tt.request.AlbumID != nil {
				album = &model.Album{
					ID:       *tt.request.AlbumID,
					Songs:    []model.Song{{}, {}, {}, {}, {}},
					ArtistID: &[]uuid.UUID{uuid.New()}[0],
				}
				albumRepository.On("GetWithSongs", mock.IsType(album), *tt.request.AlbumID).
					Return(nil, album).
					Once()
			}

			var songID uuid.UUID
			songRepository.On("Create", mock.IsType(new(model.Song))).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					songID = newSong.ID
					assertCreatedSong(t, tt.request, *newSong, userID, album)
				}).
				Return(nil).
				Once()

			// when
			id, errCode := _uut.Handle(tt.request, token)

			// then
			assert.Nil(t, errCode)
			assert.Equal(t, id, songID)

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
	assert.Equal(t, request.ReleaseDate, song.ReleaseDate)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Nil(t, song.ImageURL)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
	assert.Equal(t, request.AlbumID, song.AlbumID)
	assert.Equal(t, userID, song.UserID)
	assert.Len(t, request.Sections, len(song.Sections))
	for i, section := range request.Sections {
		assert.NotEmpty(t, song.Sections[i].ID)
		assert.Equal(t, section.Name, song.Sections[i].Name)
		assert.Zero(t, song.Sections[i].Rehearsals)
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
		assert.Equal(t, uint(1), *song.AlbumTrackNo)
	}
	if request.ArtistName != nil {
		assert.NotNil(t, song.Artist)
		assert.NotEmpty(t, song.Artist.ID)
		assert.Equal(t, song.Artist.ID, *song.ArtistID)
		assert.Equal(t, *request.ArtistName, song.Artist.Name)
		assert.Equal(t, song.UserID, song.Artist.UserID)
	}
	if request.ArtistID != nil {
		assert.Equal(t, request.ArtistID, song.ArtistID)
	}
	if request.AlbumID != nil {
		assert.Equal(t, &[]uint{uint(len(album.Songs)) + 1}[0], song.AlbumTrackNo)
		assert.Equal(t, album.ArtistID, song.ArtistID)
	}
	if request.AlbumID == nil && request.AlbumTitle == nil {
		assert.Nil(t, song.AlbumTrackNo)
	}
}
