package tuning

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
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

func TestMoveGuitarTuning_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := tuning.NewMoveGuitarTuning(nil, jwtService)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestMoveGuitarTuning_WhenGetGuitarTuningsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewMoveGuitarTuning(userDataRepository, jwtService)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveGuitarTuning_WhenGuitarTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewMoveGuitarTuning(userDataRepository, jwtService)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "tuning not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveGuitarTuning_WhenOverGuitarTuningIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewMoveGuitarTuning(userDataRepository, jwtService)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: request.ID},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over tuning not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveGuitarTuning_WhenUpdateAllGuitarTuningsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := tuning.NewMoveGuitarTuning(userDataRepository, jwtService)

	request := requests.MoveGuitarTuningRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.GuitarTuning{
		{ID: request.ID},
		{ID: request.OverID},
	}
	userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
		Return(nil, tunings).
		Once()

	internalError := errors.New("internal error")
	userDataRepository.On("UpdateAllGuitarTunings", mock.IsType(tunings)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveGuitarTuning_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	tests := []struct {
		name      string
		tunings   *[]model.GuitarTuning
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&[]model.GuitarTuning{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&[]model.GuitarTuning{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			jwtService := new(service.JwtServiceMock)
			userDataRepository := new(repository.UserDataRepositoryMock)
			_uut := tuning.NewMoveGuitarTuning(userDataRepository, jwtService)

			request := requests.MoveGuitarTuningRequest{
				ID:     (*tt.tunings)[tt.index].ID,
				OverID: (*tt.tunings)[tt.overIndex].ID,
			}
			token := "this is a token"

			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			userDataRepository.On("GetGuitarTunings", new([]model.GuitarTuning), userID).
				Return(nil, tt.tunings).
				Once()

			userDataRepository.On("UpdateAllGuitarTunings", mock.IsType(tt.tunings)).
				Run(func(args mock.Arguments) {
					newGuitarTunings := args.Get(0).(*[]model.GuitarTuning)
					tunings := slices.Clone(*newGuitarTunings)
					slices.SortFunc(tunings, func(a, b model.GuitarTuning) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, tunings[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, tunings[tt.overIndex+1].ID, request.OverID)
					}
					for i, tune := range tunings {
						assert.Equal(t, uint(i), tune.Order)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request, token)

			// then
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			userDataRepository.AssertExpectations(t)
		})
	}
}
