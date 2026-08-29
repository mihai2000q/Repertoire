package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist/song"
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

func TestAddPerfectPlaylistSongRehearsals_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, nil, nil)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs:        []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenSongsLenIsNotTheSameAsRequest_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, nil, nil)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs:        []uuid.UUID{uuid.New()},
	}

	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist songs not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, songProcessor, transactionManager)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs:        []uuid.UUID{uuid.New()},
	}

	mockPlaylistSongs := []model.PlaylistSong{{ID: request.IDs[0]}}
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil, &mockPlaylistSongs).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs:        []uuid.UUID{uuid.New()},
	}

	mockPlaylistSongs := []model.PlaylistSong{{ID: request.IDs[0], Song: model.Song{ID: uuid.New()}}}
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil, &mockPlaylistSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongPartRepository).
		Return(internalError, false).
		Times(len(mockPlaylistSongs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs:        []uuid.UUID{uuid.New()},
	}

	mockPlaylistSongs := []model.PlaylistSong{{ID: request.IDs[0], Song: model.Song{ID: uuid.New()}}}
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil, &mockPlaylistSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongPartRepository).
		Return(nil, true).
		Times(len(mockPlaylistSongs))

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

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	mockPlaylistSongs := []model.PlaylistSong{
		{ID: request.IDs[0], Song: model.Song{ID: uuid.New()}},
		{ID: request.IDs[1], Song: model.Song{ID: uuid.New()}},
	}
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil, &mockPlaylistSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongPartRepository).
		Return(nil, false).
		Times(len(mockPlaylistSongs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistSongRehearsals_WhenSuccessful_ShouldUpdateSongs(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := song.NewAddPerfectPlaylistSongRehearsals(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongPartRepository := new(repository.SongPartRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectPlaylistSongRehearsalsRequest{
		PlaylistID: uuid.New(),
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
		},
	}

	mockPlaylistSongs := []model.PlaylistSong{
		{ID: request.IDs[0], Song: model.Song{ID: uuid.New()}},
		{ID: request.IDs[1], Song: model.Song{ID: uuid.New()}},
	}
	var mockSongs []model.Song
	for _, playlistSong := range mockPlaylistSongs {
		mockSongs = append(mockSongs, playlistSong.Song)
	}
	playlistRepository.
		On(
			"GetPlaylistSongsByIDsWithPartsAndDefaultOccurrences",
			new([]model.PlaylistSong),
			request.IDs,
			request.PlaylistID,
		).
		Return(nil, &mockPlaylistSongs).
		Once()

	repositoryFactory.On("NewSongPartRepository").Return(transactionSongPartRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	for _, s := range mockSongs {
		songProcessor.On("AddPerfectRehearsal", &s, transactionSongPartRepository).
			Return(nil, true).
			Once()
	}

	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(mockPlaylistSongs))
			assert.ElementsMatch(t, mockSongs, *newSongs)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongPartRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
