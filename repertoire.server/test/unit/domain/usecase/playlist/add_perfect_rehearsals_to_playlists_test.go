package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
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

func TestAddPerfectPlaylistRehearsals_WhenGetPlaylistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
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

func TestAddPerfectPlaylistRehearsals_WhenPlaylistsLenIs0_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlists not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, songProcessor, transactionManager)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := []model.Playlist{{ID: request.IDs[0]}}
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil, &mockPlaylists).
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

func TestAddPerfectPlaylistRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := []model.Playlist{{
		ID:            request.IDs[0],
		PlaylistSongs: []model.PlaylistSong{{ID: uuid.New(), Song: model.Song{ID: uuid.New()}}},
	}}
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil, &mockPlaylists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(internalError, false).
		Times(len(mockPlaylists[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := []model.Playlist{{
		ID:            request.IDs[0],
		PlaylistSongs: []model.PlaylistSong{{ID: uuid.New(), Song: model.Song{ID: uuid.New()}}},
	}}
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil, &mockPlaylists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, true).
		Times(len(mockPlaylists[0].Songs))

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
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := []model.Playlist{
		{
			ID: request.IDs[0],
			PlaylistSongs: []model.PlaylistSong{
				{ID: uuid.New(), Song: model.Song{ID: uuid.New()}},
				{ID: uuid.New(), Song: model.Song{ID: uuid.New()}},
			},
		},
	}
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil, &mockPlaylists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, false).
		Times(len(mockPlaylists[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectPlaylistRehearsals_WhenSuccessful_ShouldUpdatePlaylists(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := playlist.NewAddPerfectRehearsalsToPlaylists(playlistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToPlaylistsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockPlaylists := []model.Playlist{
		{
			ID: request.IDs[0],
			PlaylistSongs: []model.PlaylistSong{
				{ID: uuid.New(), Song: model.Song{ID: uuid.New()}},
				{ID: uuid.New(), Song: model.Song{ID: uuid.New()}},
			},
		},
		{
			ID: request.IDs[1],
			PlaylistSongs: []model.PlaylistSong{
				{ID: uuid.New(), Song: model.Song{ID: uuid.New()}},
			},
		},
		{ID: request.IDs[2]},
	}
	playlistRepository.On("GetAllByIDsWithSongSections", new([]model.Playlist), request.IDs).
		Return(nil, &mockPlaylists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	var playlistSongs []model.Song
	for _, a := range mockPlaylists {
		for _, ps := range a.PlaylistSongs {
			songProcessor.On("AddPerfectRehearsal", &ps.Song, transactionSongSectionRepository).
				Return(nil, true).
				Once()
			playlistSongs = append(playlistSongs, ps.Song)
		}
	}

	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(playlistSongs))
			assert.ElementsMatch(t, *newSongs, playlistSongs)
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
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
