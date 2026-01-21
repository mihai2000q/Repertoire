package artist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
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

func TestAddPerfectArtistRehearsals_WhenGetArtistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenArtistsLenIs0_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, nil, nil)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artists not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenTransactionExecuteFails_ShouldReturnError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, songProcessor, transactionManager)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := []model.Artist{{ID: request.IDs[0]}}
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil, &mockArtists).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenProcessorFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := []model.Artist{{
		ID:    request.IDs[0],
		Songs: []model.Song{{ID: uuid.New()}},
	}}
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil, &mockArtists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(internalError, false).
		Times(len(mockArtists[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	artistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := []model.Artist{{
		ID:    request.IDs[0],
		Songs: []model.Song{{ID: uuid.New()}},
	}}
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil, &mockArtists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, true).
		Times(len(mockArtists[0].Songs))

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

	artistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenSongsAreNotUpdated_ShouldNotUpdateSongs(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := []model.Artist{
		{
			ID: request.IDs[0],
			Songs: []model.Song{
				{ID: uuid.New()},
				{ID: uuid.New()},
			},
		},
	}
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil, &mockArtists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	songProcessor.On("AddPerfectRehearsal", mock.IsType(new(model.Song)), transactionSongSectionRepository).
		Return(nil, false).
		Times(len(mockArtists[0].Songs))

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}

func TestAddPerfectArtistRehearsals_WhenSuccessful_ShouldUpdateArtists(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	songProcessor := new(processor.SongProcessorMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewAddPerfectRehearsalsToArtists(artistRepository, songProcessor, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionSongSectionRepository := new(repository.SongSectionRepositoryMock)
	transactionSongRepository := new(repository.SongRepositoryMock)

	request := requests.AddPerfectRehearsalsToArtistsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockArtists := []model.Artist{
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
	artistRepository.On("GetAllByIDsWithSongSections", new([]model.Artist), request.IDs).
		Return(nil, &mockArtists).
		Once()

	repositoryFactory.On("NewSongSectionRepository").Return(transactionSongSectionRepository).Once()
	repositoryFactory.On("NewSongRepository").Return(transactionSongRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	var artistSongs []model.Song
	for _, a := range mockArtists {
		for _, s := range a.Songs {
			songProcessor.On("AddPerfectRehearsal", &s, transactionSongSectionRepository).
				Return(nil, true).
				Once()
			artistSongs = append(artistSongs, s)
		}
	}

	transactionSongRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.Song))).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			assert.Len(t, *newSongs, len(artistSongs))
			assert.ElementsMatch(t, *newSongs, artistSongs)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	songProcessor.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionSongSectionRepository.AssertExpectations(t)
	transactionSongRepository.AssertExpectations(t)
}
