package instrument

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/instrument"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInstrument_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := instrument.NewCreateInstrument(nil, jwtService)

	request := requests.CreateInstrumentRequest{
		Name: "New Instrument",
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

func TestCreateInstrument_WhenGetInstrumentsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewCreateInstrument(userDataRepository, jwtService)

	request := requests.CreateInstrumentRequest{
		Name: "New Instrument",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetInstrumentsCount", new(int64), userID).
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

func TestCreateInstrument_WhenCreateInstrumentFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewCreateInstrument(userDataRepository, jwtService)

	request := requests.CreateInstrumentRequest{
		Name: "New Instrument",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{12}[0]
	userDataRepository.On("GetInstrumentsCount", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	internalError := errors.New("internal error")
	userDataRepository.On("CreateInstrument", mock.IsType(new(model.Instrument))).
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

func TestCreateInstrument_WhenSuccessful_ShouldReturnInstruments(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewCreateInstrument(userDataRepository, jwtService)

	request := requests.CreateInstrumentRequest{
		Name: "New Instrument",
	}
	token := "this is a token"

	// given - mocking
	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{12}[0]
	userDataRepository.On("GetInstrumentsCount", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	userDataRepository.On("CreateInstrument", mock.IsType(new(model.Instrument))).
		Run(func(args mock.Arguments) {
			newInstrument := args.Get(0).(*model.Instrument)
			assertCreatedInstrument(t, *newInstrument, request, count, userID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func assertCreatedInstrument(
	t *testing.T,
	instrument model.Instrument,
	request requests.CreateInstrumentRequest,
	count *int64,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, instrument.ID)
	assert.Equal(t, request.Name, instrument.Name)
	assert.Equal(t, userID, instrument.UserID)
	assert.Equal(t, uint(*count), instrument.Order)
}
