package tuning

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/guitar/tuning"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGuitarTuning_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := tuning.NewCreateGuitarTuning(nil, jwtService)

	request := requests.CreateGuitarTuningRequest{
		Name: "New Guitar Tuning",
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

func TestCreateGuitarTuning_WhenGetGuitarTuningsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := tuning.NewCreateGuitarTuning(songRepository, jwtService)

	request := requests.CreateGuitarTuningRequest{
		Name: "New Guitar Tuning",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetGuitarTuningsCount", new(int64), userID).
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

func TestCreateGuitarTuning_WhenCreateGuitarTuningFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := tuning.NewCreateGuitarTuning(songRepository, jwtService)

	request := requests.CreateGuitarTuningRequest{
		Name: "New Guitar Tuning",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{12}[0]
	songRepository.On("GetGuitarTuningsCount", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("CreateGuitarTuning", mock.IsType(new(model.GuitarTuning))).
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

func TestCreateGuitarTuning_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := tuning.NewCreateGuitarTuning(songRepository, jwtService)

	request := requests.CreateGuitarTuningRequest{
		Name: "New Guitar Tuning",
	}
	token := "this is a token"

	// given - mocking
	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{12}[0]
	songRepository.On("GetGuitarTuningsCount", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	songRepository.On("CreateGuitarTuning", mock.IsType(new(model.GuitarTuning))).
		Run(func(args mock.Arguments) {
			newGuitarTuning := args.Get(0).(*model.GuitarTuning)
			assertCreatedGuitarTuning(t, *newGuitarTuning, request, count, userID)
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

func assertCreatedGuitarTuning(
	t *testing.T,
	guitarTuning model.GuitarTuning,
	request requests.CreateGuitarTuningRequest,
	count *int64,
	userID uuid.UUID,
) {
	assert.NotEmpty(t, guitarTuning.ID)
	assert.Equal(t, request.Name, guitarTuning.Name)
	assert.Equal(t, userID, guitarTuning.UserID)
	assert.Equal(t, uint(*count), guitarTuning.Order)
}
