package song

import (
	"errors"
	"net/http"
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

func TestGetInstruments_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := song.NewGetInstruments(nil, jwtService)

	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	instruments, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, instruments)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetInstruments_WhenGetInstrumentsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetInstruments(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetInstruments", new([]model.Instrument), userID).Return(internalError).Once()

	// when
	instruments, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, instruments)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestGetInstruments_WhenSuccessful_ShouldReturnInstruments(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetInstruments(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedTunings := &[]model.Instrument{
		{ID: uuid.New(), Name: "Drop D"},
		{ID: uuid.New(), Name: "Drop C"},
	}
	songRepository.On("GetInstruments", mock.IsType(expectedTunings), userID).
		Return(nil, expectedTunings).
		Once()

	// when
	instruments, errCode := _uut.Handle(token)

	// then
	assert.Equal(t, expectedTunings, &instruments)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
