package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
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

func TestAddPerfectAlbumRehearsals_WhenGetAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenAlbumsLenIs0_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "albums not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, songProcessor, transactionManager)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := []model.Album{{ID: request.IDs[0]}}
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil, &mockAlbums).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := []model.Album{{
		ID:    request.IDs[0],
		Songs: []model.Song{{ID: uuid.New()}},
	}}
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil, &mockAlbums).
		Once()

	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongRepository).
		Return(internalError, false).
		Times(len(mockAlbums[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	albumRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := []model.Album{{
		ID:    request.IDs[0],
		Songs: []model.Song{{ID: uuid.New()}},
	}}
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil, &mockAlbums).
		Once()

	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongRepository).
		Return(nil, true).
		Times(len(mockAlbums[0].Songs))

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

	albumRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := []model.Album{
		{
			ID: request.IDs[0],
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
	}
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil, &mockAlbums).
		Once()

	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongRepository).
		Return(nil, false).
		Times(len(mockAlbums[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectAlbumRehearsals_WhenSuccessful_ShouldUpdateAlbums(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewAddPerfectRehearsalsToAlbums(albumRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToAlbumsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockAlbums := []model.Album{
		{
			ID: request.IDs[0],
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
		{
			ID: request.IDs[1],
			Songs: []model.Song{
				{ID: uuid.New()},
			},
		},
		{ID: request.IDs[2]},
	}
	albumRepository.On("GetAllByIDsWithSongSections", new([]model.Album), request.IDs).
		Return(nil, &mockAlbums).
		Once()

	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	var albumSongs []model.Song
	for _, a := range mockAlbums {
		for _, s := range a.Songs {
			songProcessor.On("AddPerfectRehearsal", &s, transactionSongRepository).
				Return(nil, true).
				Once()
			albumSongs = append(albumSongs, s)
		}
	}

	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(albumSongs))
			assert.ElementsMatch(t, *newSongs, albumSongs)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
