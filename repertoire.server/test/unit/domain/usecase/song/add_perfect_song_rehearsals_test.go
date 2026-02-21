package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddPerfectSongRehearsals_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, nil, nil)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenSongsLenIsNotTheSameAsRequest_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, nil, nil)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "songs not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, songProcessor, transactionManager)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := []model.Song{{ID: request.IDs[0]}}
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil, &mockSongs).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := []model.Song{{ID: request.IDs[0]}}
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(internalError, false).
		Times(len(mockSongs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := []model.Song{{ID: request.IDs[0]}}
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, true).
		Times(len(mockSongs))

	internalError := errors.New("internal error")
	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	mockSongs := []model.Song{
		{ID: request.IDs[0]},
		{ID: request.IDs[1]},
	}
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, false).
		Times(len(mockSongs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsals_WhenSuccessful_ShouldUpdateSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsals(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	mockSongs := []model.Song{
		{ID: request.IDs[0]},
		{ID: request.IDs[1]},
	}
	songRepository.On("GetAllByIDsWithSectionsAndOccurrences", new([]model.Song), request.IDs).
		Return(nil, &mockSongs).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	for _, s := range mockSongs {
		songProcessor.On("AddPerfectRehearsal", &s, transactionSongSectionRepository).
			Return(nil, true).
			Once()
	}

	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(mockSongs))
			assert.ElementsMatch(t, mockSongs, *newSongs)
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
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
