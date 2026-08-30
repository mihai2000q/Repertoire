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
	"github.com/stretchr/testify/require"
)

func TestAddCustomSongRehearsal_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, nil, nil)

	request := requests.AddCustomSongRehearsalRequest{
		ID:            uuid.New(),
		ArrangementID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
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

func TestAddCustomSongRehearsal_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, nil, nil)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsal_WhenSongHasArrangement_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, nil, nil)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID:    request.ID,
		Parts: []model.SongPart{{ID: uuid.New()}},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangement not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsal_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, songProcessor, transactionManager)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.ID,
		Parts: []model.SongPart{
			{ID: uuid.New(), ArrangementOccurrences: []model.SongPartOccurrences{{}}},
		},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
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

func TestAddCustomSongRehearsal_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.ID,
		Parts: []model.SongPart{
			{ID: uuid.New(), ArrangementOccurrences: []model.SongPartOccurrences{{}}},
		},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddCustomRehearsal", &mockSong, transactionSongPartRepository, []*uuid.UUID{nil}[0]).
		Return(internalError, false).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsal_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.ID,
		Parts: []model.SongPart{
			{ID: uuid.New(), ArrangementOccurrences: []model.SongPartOccurrences{{}}},
		},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddCustomRehearsal", &mockSong, transactionSongPartRepository, []*uuid.UUID{nil}[0]).
		Return(nil, true).
		Once()

	internalError := errors.New("internal error")
	transactionSongRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
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
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsal_WhenSongIsNotUpdated_ShouldNotUpdateSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.ID,
		Parts: []model.SongPart{
			{ID: uuid.New(), ArrangementOccurrences: []model.SongPartOccurrences{{}}},
		},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddCustomRehearsal", &mockSong, transactionSongPartRepository, []*uuid.UUID{nil}[0]).
		Return(nil, false).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddCustomSongRehearsal_WhenSuccessful_ShouldUpdateSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddCustomSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddCustomSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.ID,
		Parts: []model.SongPart{
			{ID: uuid.New(), ArrangementOccurrences: []model.SongPartOccurrences{{}}},
		},
	}
	songRepository.
		On(
			"GetWithPartsAndArrangementOccurrences",
			new(model.Song),
			request.ID,
			request.ArrangementID,
		).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddCustomRehearsal", &mockSong, transactionSongPartRepository, []*uuid.UUID{nil}[0]).
		Return(nil, true).
		Once()

	transactionSongRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
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
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
