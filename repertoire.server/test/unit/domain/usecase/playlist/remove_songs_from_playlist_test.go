package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveSongsFromPlaylist_WhenGetPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, nil)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
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

func TestRemoveSongsFromPlaylist_WhenNotAllSongsFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, nil)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongIDs[0], SongTrackNo: 1},
		{SongID: uuid.New(), SongTrackNo: 2},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "could not find all songs", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestRemoveSongsFromPlaylist_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, transactionManager)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongIDs[0], SongTrackNo: 1},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
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
	transactionManager.AssertExpectations(t)
}

func TestRemoveSongsFromPlaylist_WhenRemoveSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionPlaylistRepository := new(repository.PlaylistRepositoryMock)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongIDs[0], SongTrackNo: 1},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	repositoryFactory.On("NewPlaylistRepository").Return(transactionPlaylistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionPlaylistRepository.On("RemoveSongs", mock.IsType(playlistSongs)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	transactionPlaylistRepository.AssertExpectations(t)
}

func TestRemoveSongsFromPlaylist_WhenUpdateAllPlaylistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionPlaylistRepository := new(repository.PlaylistRepositoryMock)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: request.SongIDs[0], SongTrackNo: 1},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	repositoryFactory.On("NewPlaylistRepository").Return(transactionPlaylistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	transactionPlaylistRepository.On("RemoveSongs", mock.IsType(playlistSongs)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	transactionPlaylistRepository.On("UpdateAllPlaylistSongs", mock.IsType(playlistSongs)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	transactionPlaylistRepository.AssertExpectations(t)
}

func TestRemoveSongsFromPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	transactionManager := new(transaction.ManagerMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewRemoveSongsFromPlaylist(playlistRepository, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionPlaylistRepository := new(repository.PlaylistRepositoryMock)

	request := requests.RemoveSongsFromPlaylistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	playlistSongs := &[]model.PlaylistSong{
		{SongID: uuid.New(), SongTrackNo: 1},
		{SongID: request.SongIDs[0], SongTrackNo: 2},
		{SongID: uuid.New(), SongTrackNo: 3},
		{SongID: uuid.New(), SongTrackNo: 4},
		{SongID: request.SongIDs[1], SongTrackNo: 5},
	}
	playlistRepository.On("GetPlaylistSongs", new([]model.PlaylistSong), request.ID).
		Return(nil, playlistSongs).
		Once()

	repositoryFactory.On("NewPlaylistRepository").Return(transactionPlaylistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	transactionPlaylistRepository.On("RemoveSongs", mock.IsType(playlistSongs)).
		Run(func(args mock.Arguments) {
			songs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *songs, len(request.SongIDs))
			for _, song := range *songs {
				assert.Contains(t, request.SongIDs, song.SongID)
			}
		}).
		Return(nil).
		Once()

	transactionPlaylistRepository.On("UpdateAllPlaylistSongs", mock.IsType(playlistSongs)).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.Len(t, *newSongs, len(*playlistSongs)-len(request.SongIDs))
			for i, song := range *newSongs {
				assert.NotContains(t, request.SongIDs, song.SongID)
				assert.Equal(t, uint(i)+1, song.SongTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	transactionPlaylistRepository.AssertExpectations(t)
}
