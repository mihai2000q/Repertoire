package instrument

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/instrument"
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

func TestMoveInstrument_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := instrument.NewMoveInstrument(nil, jwtService)

	request := requests.MoveInstrumentRequest{
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

func TestMoveInstrument_WhenGetInstrumentsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewMoveInstrument(userDataRepository, jwtService)

	request := requests.MoveInstrumentRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
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

func TestMoveInstrument_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewMoveInstrument(userDataRepository, jwtService)

	request := requests.MoveInstrumentRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "instrument not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveInstrument_WhenOverInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewMoveInstrument(userDataRepository, jwtService)

	request := requests.MoveInstrumentRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: request.ID},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over instrument not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestMoveInstrument_WhenUpdateAllInstrumentsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewMoveInstrument(userDataRepository, jwtService)

	request := requests.MoveInstrumentRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: request.ID},
		{ID: request.OverID},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	internalError := errors.New("internal error")
	userDataRepository.On("UpdateAllInstruments", mock.IsType(instruments)).
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

func TestMoveInstrument_WhenSuccessful_ShouldReturnInstruments(t *testing.T) {
	tests := []struct {
		name        string
		instruments *[]model.Instrument
		index       uint
		overIndex   uint
	}{
		{
			"Use case 1",
			&[]model.Instrument{
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
			&[]model.Instrument{
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
			_uut := instrument.NewMoveInstrument(userDataRepository, jwtService)

			request := requests.MoveInstrumentRequest{
				ID:     (*tt.instruments)[tt.index].ID,
				OverID: (*tt.instruments)[tt.overIndex].ID,
			}
			token := "this is a token"

			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
				Return(nil, tt.instruments).
				Once()

			userDataRepository.On("UpdateAllInstruments", mock.IsType(tt.instruments)).
				Run(func(args mock.Arguments) {
					newInstruments := args.Get(0).(*[]model.Instrument)
					instruments := slices.Clone(*newInstruments)
					slices.SortFunc(instruments, func(a, b model.Instrument) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, instruments[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, instruments[tt.overIndex+1].ID, request.OverID)
					}
					for i, tune := range instruments {
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
