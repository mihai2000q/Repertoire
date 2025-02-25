package tuning

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/udata/guitar/tuning"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteGuitarTuning_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := tuning.NewDeleteGuitarTuning(nil, jwtService)

	id := uuid.New()
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteGuitarTuning_WhenGetGuitarTuningsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewDeleteGuitarTuning(userDataRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteGuitarTuning_WhenGuitarTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewDeleteGuitarTuning(userDataRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: uuid.New()},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "guitar tuning not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteGuitarTuning_WhenUpdateAllGuitarTuningsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewDeleteGuitarTuning(userDataRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: id},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	internalError := errors.New("internal error")
	userDataRepository.On("UpdateAllGuitarTunings", mock.IsType(tunings)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteGuitarTuning_WhenDeleteGuitarTuningFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewDeleteGuitarTuning(userDataRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: id},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	userDataRepository.On("UpdateAllGuitarTunings", mock.IsType(tunings)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	userDataRepository.On("DeleteGuitarTuning", id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteGuitarTuning_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewDeleteGuitarTuning(userDataRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: id},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	userDataRepository.On("UpdateAllGuitarTunings", mock.IsType(tunings)).
		Run(func(args mock.Arguments) {
			newGuitarTunings := args.Get(0).(*[]model.GuitarTuning)
			guitarTunings := slices.DeleteFunc(*newGuitarTunings, func(t model.GuitarTuning) bool {
				return t.ID == id
			})
			for i, tune := range guitarTunings {
				assert.Equal(t, i, tune.Order)
			}
		}).
		Return(nil).
		Once()

	userDataRepository.On("DeleteGuitarTuning", id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}
