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

func TestAddPerfectSongRehearsal_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	songRepository.On("GetWithSections", new(model.Song), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, songProcessor, transactionManager)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{ID: request.ID}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
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

func TestAddPerfectSongRehearsal_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{ID: request.ID}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", &mockSong, transactionSongSectionRepository).
		Return(internalError, false).
		Once()

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

func TestAddPerfectSongRehearsal_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{ID: request.ID}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", &mockSong, transactionSongSectionRepository).
		Return(nil, true).
		Once()

	internalError := errors.New("internal error")
	transactionSongRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
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

func TestAddPerfectSongRehearsal_WhenSongIsNotUpdated_ShouldNotUpdateSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{ID: request.ID}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", &mockSong, transactionSongSectionRepository).
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
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenSuccessful_ShouldUpdateSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{ID: request.ID}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", &mockSong, transactionSongSectionRepository).
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
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
