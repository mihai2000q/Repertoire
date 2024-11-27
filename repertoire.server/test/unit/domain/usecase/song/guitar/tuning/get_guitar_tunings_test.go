package tuning

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song/guitar/tuning"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetGuitarTunings_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := tuning.NewGetGuitarTunings(nil, jwtService)

	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	tunings, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, tunings)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetGuitarTunings_WhenGetGuitarTuningsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := tuning.NewGetGuitarTunings(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).Return(internalError).Once()

	// when
	tunings, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, tunings)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestGetGuitarTunings_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := tuning.NewGetGuitarTunings(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedTunings := &[]model.GuitarTuning{
		{ID: uuid.New(), Name: "Drop D"},
		{ID: uuid.New(), Name: "Drop C"},
	}
	songRepository.On("GetGuitarTunings", mock.IsType(expectedTunings), userID).
		Return(nil, expectedTunings).
		Once()

	// when
	tunings, errCode := _uut.Handle(token)

	// then
	assert.Equal(t, expectedTunings, &tunings)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
