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
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.Description, newSong.Description)
			assert.False(t, newSong.IsRecorded)
			assert.Equal(t, request.Bpm, newSong.Bpm)
			assert.Equal(t, request.SongsterrLink, newSong.SongsterrLink)
			assert.Equal(t, request.GuitarTuningID, newSong.GuitarTuningID)
			assert.Equal(t, userID, newSong.UserID)
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
	songRepository.On("Create", mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.Description, newSong.Description)
			assert.False(t, newSong.IsRecorded)
			assert.Equal(t, request.Bpm, newSong.Bpm)
			assert.Equal(t, request.SongsterrLink, newSong.SongsterrLink)
			assert.Equal(t, request.GuitarTuningID, newSong.GuitarTuningID)
			assert.Equal(t, userID, newSong.UserID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
