package role

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/userdata/band/member/role"
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

func TestDeleteBandMemberRole_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := role.NewDeleteBandMemberRole(nil, jwtService, nil)

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

func TestDeleteBandMemberRole_WhenGetBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, nil)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
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

func TestDeleteBandMemberRole_WhenRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, nil)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: uuid.New()},
	}
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member role not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userDataRepository.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, transactionManager)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
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

func TestDeleteBandMemberRole_WhenUpdateAllBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	txUserDataRepo.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
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

func TestDeleteBandMemberRole_WhenDeleteBandMemberRoleFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txUserDataRepo.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txUserDataRepo.On("DeleteBandMemberRole", id).
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

func TestDeleteBandMemberRole_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userDataRepository := new(repository.UserDataRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := role.NewDeleteBandMemberRole(userDataRepository, jwtService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txUserDataRepo := new(repository.UserDataRepositoryMock)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	userDataRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	repositoryFactory.On("NewUserDataRepository").Return(txUserDataRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txUserDataRepo.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
		Run(func(args mock.Arguments) {
			newBandMemberRoles := args.Get(0).(*[]model.BandMemberRole)
			sortedBandMemberRoles := slices.DeleteFunc(*newBandMemberRoles, func(s model.BandMemberRole) bool {
				return s.ID == id
			})
			for i, bandMemberRole := range sortedBandMemberRoles {
				assert.Equal(t, i, bandMemberRole.Order)
			}
		}).
		Return(nil).
		Once()

	txUserDataRepo.On("DeleteBandMemberRole", id).
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
