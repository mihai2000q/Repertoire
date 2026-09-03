package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddCustomSongRehearsals_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, nil, nil)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenSongsAreNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, nil, nil)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "songs not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, songProcessor, transactionManager)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{{ID: request.Requests[0].ID}}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{{ID: request.Requests[1].ID}}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "songs not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{{ID: request.Requests[0].ID}}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := httperror.InternalServerError(errors.New("internal error"))
	songProcessor.
		On(
			"AddCustomRehearsal",
			mock.IsType(new(model.Song)),
			txSongPartRepo,
			mock.IsType(new(uuid.UUID)),
		).
		Return(internalError, false).
		Times(len(request.Requests))

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: uuid.New(), ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{{ID: request.Requests[0].ID}}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.
		On(
			"AddCustomRehearsal",
			mock.IsType(new(model.Song)),
			txSongPartRepo,
			mock.IsType(new(uuid.UUID)),
		).
		Return(nil, true).
		Times(len(request.Requests))

	internalError := errors.New("internal error")
	txSongRepo.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	sameId := uuid.New()
	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: sameId, ArrangementID: uuid.New()},
			{ID: uuid.New(), ArrangementID: uuid.New()},
			{ID: sameId, ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{
		{ID: request.Requests[0].ID},
		{ID: request.Requests[1].ID},
	}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.
		On(
			"AddCustomRehearsal",
			mock.IsType(new(model.Song)),
			txSongPartRepo,
			mock.IsType(new(uuid.UUID)),
		).
		Return(nil, false).
		Times(len(request.Requests))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}

func TestAddCustomSongRehearsals_WhenSuccessful_ShouldUpdateSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongPartRepo := new(repository.SongPartRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	sameId := uuid.New()
	request := requests.AddCustomSongRehearsalsRequest{
		Requests: []requests.AddCustomSongRehearsalRequest{
			{ID: sameId, ArrangementID: uuid.New()},
			{ID: uuid.New(), ArrangementID: uuid.New()},
			{ID: sameId, ArrangementID: uuid.New()},
		},
	}

	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}

	mockSongs := []model.Song{
		{ID: request.Requests[0].ID},
		{ID: request.Requests[1].ID},
	}
	songRepository.On("GetAllByIDsWithPartsAndArrangementOccurrences", new([]model.Song), ids).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(txSongPartRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songsMap := make(map[uuid.UUID]model.Song)
	for _, s := range mockSongs {
		songsMap[s.ID] = s
	}

	for _, r := range request.Requests {
		songProcessor.
			On(
				"AddCustomRehearsal",
				&[]model.Song{songsMap[r.ID]}[0],
				txSongPartRepo,
				&r.ArrangementID,
			).
			Return(nil, true).
			Once()
	}

	txSongRepo.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(request.Requests))
			for _, s := range *newSongs {
				assert.Equal(t, songsMap[s.ID], s)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongPartRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}
