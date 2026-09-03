package arrangement

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song/arrangement"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDeleteSongArrangement_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, nil)

	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, nil)

	id := uuid.New()
	songID := uuid.New()

	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, nil)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Arrangements: []model.SongArrangement{
			{ID: uuid.New(), Order: 0},
		},
	}
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangement not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, transactionManager)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Arrangements: []model.SongArrangement{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Arrangements: []model.SongArrangement{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil, song).
		Once()

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	txSongRepo.On("UpdateWithAssociations", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenDeleteArrangementFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Arrangements: []model.SongArrangement{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil, song).
		Once()

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(song)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txSongArrangementRepo.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}

func TestDeleteSongArrangement_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := arrangement.NewDeleteSongArrangement(songRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txSongRepo := new(repository.SongRepositoryMock)
	txSongArrangementRepo := new(repository.SongArrangementRepositoryMock)

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	mockSong := model.Song{
		ID: songID,
		Arrangements: []model.SongArrangement{
			{ID: uuid.New(), Order: 0},
			{ID: id, Order: 1},
			{ID: uuid.New(), Order: 2},
		},
	}
	songRepository.On("GetWithArrangements", new(model.Song), songID).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	repositoryFactory.On("NewSongArrangementRepository").Return(txSongArrangementRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txSongRepo.On("UpdateWithAssociations", mock.IsType(&mockSong)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)

			// arrangements ordered
			arrangements := slices.Clone(newSong.Arrangements)
			arrangements = slices.DeleteFunc(arrangements, func(a model.SongArrangement) bool {
				return a.ID == id
			})
			for i, s := range arrangements {
				assert.Equal(t, uint(i), s.Order)
			}
		}).
		Return(nil).
		Once()

	txSongArrangementRepo.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
	txSongArrangementRepo.AssertExpectations(t)
}
