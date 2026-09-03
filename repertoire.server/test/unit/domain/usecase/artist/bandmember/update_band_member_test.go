package bandmember

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist/bandmember"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateBandMember_WhenGetBandMembersFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, nil)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenBandMembersIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, nil)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, transactionManager)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestUpdateBandMember_WhenGetRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txArtistRepo := new(repository.ArtistRepositoryMock)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(txArtistRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	txArtistRepo.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txArtistRepo.AssertExpectations(t)
}

func TestUpdateBandMember_WhenReplaceRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txArtistRepo := new(repository.ArtistRepositoryMock)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(txArtistRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txArtistRepo.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txArtistRepo.
		On(
			"ReplaceRolesFromBandMember",
			mock.IsType([]model.BandMemberRole{}),
			mock.IsType(new(model.BandMember)),
		).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txArtistRepo.AssertExpectations(t)
}

func TestUpdateBandMember_WhenUpdateBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txArtistRepo := new(repository.ArtistRepositoryMock)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(txArtistRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txArtistRepo.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(nil).
		Once()

	txArtistRepo.
		On(
			"ReplaceRolesFromBandMember",
			mock.IsType([]model.BandMemberRole{}),
			mock.IsType(new(model.BandMember)),
		).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txArtistRepo.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txArtistRepo.AssertExpectations(t)
}

func TestUpdateBandMember_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := bandmember.NewUpdateBandMember(artistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txArtistRepo := new(repository.ArtistRepositoryMock)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(txArtistRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	roles := []model.BandMemberRole{
		{ID: request.RoleIDs[0]},
	}
	txArtistRepo.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(nil, &roles).
		Once()

	txArtistRepo.On("ReplaceRolesFromBandMember", roles, mockBandMember).
		Return(nil).
		Once()

	txArtistRepo.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Run(func(args mock.Arguments) {
			newBandMember := args.Get(0).(*model.BandMember)
			assertUpdatedBandMember(t, request, *newBandMember)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txArtistRepo.AssertExpectations(t)
}

func assertUpdatedBandMember(
	t *testing.T,
	request requests.UpdateBandMemberRequest,
	member model.BandMember,
) {
	assert.Equal(t, request.Name, member.Name)
	assert.Equal(t, request.Color, member.Color)
}
