package artist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBulkDeleteArtists_WhenGetArtistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, nil, nil)

	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
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

func TestBulkDeleteArtists_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, nil, nil)

	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
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

func TestBulkDeleteArtists_WhenDeleteArtistAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.BulkDeleteArtistsRequest{
		IDs:        []uuid.UUID{uuid.New()},
		WithAlbums: true,
	}

	mockArtists := &[]model.Artist{{ID: request.IDs[0]}}
	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtists).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("DeleteAlbums", request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionArtistRepository.AssertExpectations(t)
}

func TestBulkDeleteArtists_WhenDeleteArtistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.BulkDeleteArtistsRequest{
		IDs:       []uuid.UUID{uuid.New()},
		WithSongs: true,
	}

	mockArtists := &[]model.Artist{{ID: request.IDs[0]}}
	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtists).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("DeleteSongs", request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionArtistRepository.AssertExpectations(t)
}

func TestBulkDeleteArtists_WhenDeleteArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := &[]model.Artist{{ID: request.IDs[0]}}
	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtists).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("Delete", request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionArtistRepository.AssertExpectations(t)
}

func TestBulkDeleteArtists_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewBulkDeleteArtists(artistRepository, messagePublisherService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.BulkDeleteArtistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockArtists := &[]model.Artist{{ID: request.IDs[0]}}
	artistRepository.
		On(
			"GetAllByIDs",
			new([]model.Artist),
			request.IDs,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtists).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	transactionArtistRepository.On("Delete", request.IDs).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.ArtistsDeletedTopic, *mockArtists).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	transactionArtistRepository.AssertExpectations(t)
}

func TestBulkDeleteArtists_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	tests := []struct {
		name       string
		artists    []model.Artist
		withAlbums bool
		withSongs  bool
	}{
		{
			"Without Albums or Songs",
			[]model.Artist{{ID: uuid.New()}},
			false,
			false,
		},
		{
			"With Albums",
			[]model.Artist{{ID: uuid.New()}},
			true,
			false,
		},
		{
			"With Songs",
			[]model.Artist{{ID: uuid.New()}},
			false,
			true,
		},
		{
			"With Songs and Albums",
			[]model.Artist{{ID: uuid.New()}},
			true,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			artistRepository := new(repository.ArtistRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			transactionManager := new(transaction.ManagerMock)
			_uut := artist.NewBulkDeleteArtists(artistRepository, messagePublisherService, transactionManager)

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			transactionArtistRepository := new(repository.ArtistRepositoryMock)

			var ids []uuid.UUID
			for _, art := range tt.artists {
				ids = append(ids, art.ID)
			}
			request := requests.BulkDeleteArtistsRequest{
				IDs:        ids,
				WithAlbums: tt.withAlbums,
				WithSongs:  tt.withSongs,
			}

			artistRepository.
				On(
					"GetAllByIDs",
					new([]model.Artist),
					request.IDs,
					request.WithSongs,
					request.WithAlbums,
				).
				Return(nil, &tt.artists).
				Once()

			repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			if tt.withAlbums {
				transactionArtistRepository.On("DeleteAlbums", request.IDs).
					Return(nil).
					Once()
			}
			if tt.withSongs {
				transactionArtistRepository.On("DeleteSongs", request.IDs).
					Return(nil).
					Once()
			}

			transactionArtistRepository.On("Delete", request.IDs).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.ArtistsDeletedTopic, tt.artists).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			artistRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			transactionArtistRepository.AssertExpectations(t)
		})
	}
}
