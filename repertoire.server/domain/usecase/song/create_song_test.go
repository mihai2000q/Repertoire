package song

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSong_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateSong{
		jwtService: jwtService,
	}
	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"

	forbiddenError := wrapper.UnauthorizedError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreateSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateSong{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.CreateSongRequest{
		Title: "Some Song",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	songRepository.On("Create", mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assertSong(t, request, *newSong, userID)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
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
			"Create song with more fields and album ID and artist ID",
			requests.CreateSongRequest{
				Title:          "Some Song",
				Bpm:            &[]uint{120}[0],
				SongsterrLink:  &[]string{"https://songsterr.com/some-song"}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				AlbumID:        &[]uuid.UUID{uuid.New()}[0],
				ArtistID:       &[]uuid.UUID{uuid.New()}[0],
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
			"Create song with sections",
			requests.CreateSongRequest{
				Title: "Some Song",
				Sections: []requests.CreateSongSectionRequest{
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
			jwtService := new(service.JwtServiceMock)
			_uut := &CreateSong{
				repository: songRepository,
				jwtService: jwtService,
			}
			token := "this is a token"

			// given - mocks
			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
			songRepository.On("Create", mock.IsType(new(model.Song))).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					assertSong(t, tt.request, *newSong, userID)
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request, token)

			// then
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			songRepository.AssertExpectations(t)
		})
	}
}

func assertSong(t *testing.T, request requests.CreateSongRequest, newSong model.Song, userID uuid.UUID) {
	assert.Equal(t, request.Title, newSong.Title)
	assert.Equal(t, request.Description, newSong.Description)
	assert.False(t, newSong.IsRecorded)
	assert.Equal(t, request.Bpm, newSong.Bpm)
	assert.Equal(t, request.SongsterrLink, newSong.SongsterrLink)
	assert.Equal(t, request.GuitarTuningID, newSong.GuitarTuningID)
	assert.Equal(t, request.AlbumID, newSong.AlbumID)
	assert.Equal(t, request.ArtistID, newSong.ArtistID)
	assert.Equal(t, userID, newSong.UserID)
	assert.Len(t, request.Sections, len(newSong.Sections))
	for i, section := range request.Sections {
		assert.NotEmpty(t, newSong.Sections[i].ID)
		assert.Equal(t, section.Name, newSong.Sections[i].Name)
		assert.Equal(t, section.TypeID, newSong.Sections[i].SongSectionTypeID)
	}
	if request.AlbumTitle != nil {
		assert.NotNil(t, newSong.Album)
		assert.NotEmpty(t, newSong.Album.ID)
		assert.Equal(t, *request.AlbumTitle, newSong.Album.Title)
	}
	if request.ArtistName != nil {
		assert.NotNil(t, newSong.Artist)
		assert.NotEmpty(t, newSong.Artist.ID)
		assert.Equal(t, *request.ArtistName, newSong.Artist.Name)
	}
}
