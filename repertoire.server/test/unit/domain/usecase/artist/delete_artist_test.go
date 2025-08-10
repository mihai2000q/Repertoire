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

func TestDeleteArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, nil)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	internalError := errors.New("internal error")
	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
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

func TestDeleteArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, nil)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
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
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteArtist_WhenDeleteArtistAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.DeleteArtistRequest{
		ID:         uuid.New(),
		WithAlbums: true,
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtist).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("DeleteAlbums", []uuid.UUID{request.ID}).
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

func TestDeleteArtist_WhenDeleteArtistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.DeleteArtistRequest{
		ID:        uuid.New(),
		WithSongs: true,
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtist).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("DeleteSongs", []uuid.UUID{request.ID}).
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

func TestDeleteArtist_WhenDeleteArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtist).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	transactionArtistRepository.On("Delete", []uuid.UUID{request.ID}).
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

func TestDeleteArtist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := artist.NewDeleteArtist(artistRepository, messagePublisherService, transactionManager)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	transactionArtistRepository := new(repository.ArtistRepositoryMock)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.
		On(
			"GetWithSongsOrAlbums",
			new(model.Artist),
			request.ID,
			request.WithSongs,
			request.WithAlbums,
		).
		Return(nil, mockArtist).
		Once()

	repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	transactionArtistRepository.On("Delete", []uuid.UUID{request.ID}).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.ArtistsDeletedTopic, *mockArtist).
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

func TestDeleteArtist_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	tests := []struct {
		name       string
		artist     model.Artist
		withAlbums bool
		withSongs  bool
	}{
		{
			"Without Albums or Songs",
			model.Artist{ID: uuid.New()},
			false,
			false,
		},
		{
			"With Albums",
			model.Artist{ID: uuid.New()},
			true,
			false,
		},
		{
			"With Songs",
			model.Artist{ID: uuid.New()},
			false,
			true,
		},
		{
			"With Songs and Albums",
			model.Artist{ID: uuid.New()},
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
			_uut := artist.NewDeleteArtist(artistRepository, messagePublisherService, transactionManager)

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			transactionArtistRepository := new(repository.ArtistRepositoryMock)

			request := requests.DeleteArtistRequest{
				ID:         tt.artist.ID,
				WithAlbums: tt.withAlbums,
				WithSongs:  tt.withSongs,
			}

			artistRepository.
				On(
					"GetWithSongsOrAlbums",
					new(model.Artist),
					request.ID,
					request.WithSongs,
					request.WithAlbums,
				).
				Return(nil, &tt.artist).
				Once()

			repositoryFactory.On("NewArtistRepository").Return(transactionArtistRepository).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			if tt.withAlbums {
				transactionArtistRepository.On("DeleteAlbums", []uuid.UUID{request.ID}).
					Return(nil).
					Once()
			}
			if tt.withSongs {
				transactionArtistRepository.On("DeleteSongs", []uuid.UUID{request.ID}).
					Return(nil).
					Once()
			}

			transactionArtistRepository.On("Delete", []uuid.UUID{request.ID}).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.ArtistsDeletedTopic, tt.artist).
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
