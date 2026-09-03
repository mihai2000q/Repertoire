package instrument

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/userdata/instrument"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteInstrument_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := instrument.NewDeleteInstrument(nil, jwtService, nil)

	id := uuid.New()
	token := "this is a token"

	forbiddenError := httperror.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteInstrument_WhenGetInstrumentsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, nil)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteInstrument_WhenInstrumentIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, nil)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: uuid.New()},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "instrument not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteInstrument_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, transactionManager)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: id},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestDeleteInstrument_WhenUpdateAllInstrumentsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: id},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	txUserDataRepo.On("UpdateAllInstruments", mock.IsType(instruments)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txUserDataRepo.AssertExpectations(t)
}

func TestDeleteInstrument_WhenDeleteInstrumentFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: id},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txUserDataRepo.On("UpdateAllInstruments", mock.IsType(instruments)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txUserDataRepo.On("DeleteInstrument", id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txUserDataRepo.AssertExpectations(t)
}

func TestDeleteInstrument_WhenSuccessful_ShouldReturnInstruments(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := instrument.NewDeleteInstrument(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	instruments := &[]model.Instrument{
		{ID: id},
	}
	userDataRepository.On("GetInstruments", new([]model.Instrument), userID).
		Return(nil, instruments).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txUserDataRepo.On("UpdateAllInstruments", mock.IsType(instruments)).
		Run(func(args mock.Arguments) {
			newInstruments := args.Get(0).(*[]model.Instrument)
			guitarTunings := slices.DeleteFunc(*newInstruments, func(t model.Instrument) bool {
				return t.ID == id
			})
			for i, tune := range guitarTunings {
				assert.Equal(t, i, tune.Order)
			}
		}).
		Return(nil).
		Once()

	txUserDataRepo.On("DeleteInstrument", id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txUserDataRepo.AssertExpectations(t)
}
