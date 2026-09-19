package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/date"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/database/transaction"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewUpdateAlbum(albumRepository, nil, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenTransactionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, nil)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}
	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	internalError := errors.New("internal error")
	transactionManager.On("Execute", mock.Anything).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
}

func TestUpdateAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, nil)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txAlbumRepo := new(repository.AlbumRepositoryMock)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	internalError := errors.New("internal error")
	txAlbumRepo.On("Update", mock.IsType(mockAlbum)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txAlbumRepo.AssertExpectations(t)
}

func TestUpdateAlbum_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, nil)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txAlbumRepo := new(repository.AlbumRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	request := requests.UpdateAlbumRequest{
		ID:       uuid.New(),
		Title:    "New Album",
		ArtistID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txAlbumRepo.On("Update", mock.IsType(mockAlbum)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	txSongRepo.On("GetAllByAlbum", new([]model.Song), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txAlbumRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}

func TestUpdateAlbum_WhenUpdateSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, nil)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txAlbumRepo := new(repository.AlbumRepositoryMock)
	txSongRepo := new(repository.SongRepositoryMock)

	request := requests.UpdateAlbumRequest{
		ID:       uuid.New(),
		Title:    "New Album",
		ArtistID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
	repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txAlbumRepo.On("Update", mock.IsType(mockAlbum)).
		Return(nil).
		Once()

	songs := []model.Song{{ID: uuid.New()}}

	txSongRepo.On("GetAllByAlbum", new([]model.Song), request.ID).
		Return(nil, &songs).
		Once()

	internalError := errors.New("internal error")
	txSongRepo.On("UpdateAll", mock.IsType(&[]model.Song{})).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txAlbumRepo.AssertExpectations(t)
	txSongRepo.AssertExpectations(t)
}

func TestUpdateAlbum_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, messagePublisherService)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txAlbumRepo := new(repository.AlbumRepositoryMock)

	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txAlbumRepo.On("Update", mock.IsType(mockAlbum)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txAlbumRepo.AssertExpectations(t)
}

func TestUpdateAlbum_WhenArtistHasNotChanged_ShouldUpdateOnlyAlbumAndNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	transactionManager := new(transaction.ManagerMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewUpdateAlbum(albumRepository, transactionManager, messagePublisherService)

	repositoryFactory := new(transaction.RepositoryFactoryMock)
	txAlbumRepo := new(repository.AlbumRepositoryMock)

	request := requests.UpdateAlbumRequest{
		ID:          uuid.New(),
		Title:       "New Album",
		ReleaseDate: &[]date.Date{date.Date(time.Now())}[0],
	}

	mockAlbum := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
	transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

	txAlbumRepo.On("Update", mock.IsType(mockAlbum)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assertUpdatedAlbum(t, *newAlbum, request)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{mockAlbum.ID}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	transactionManager.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
	repositoryFactory.AssertExpectations(t)
	txAlbumRepo.AssertExpectations(t)
}

func TestUpdateAlbum_WhenArtistHasChanged_ShouldUpdateAlbumAndSongsAndNotReturnAnyError(t *testing.T) {
	albumID := uuid.New()

	tests := []struct {
		name    string
		request requests.UpdateAlbumRequest
		album   model.Album
	}{
		{
			"Artist has changed",
			requests.UpdateAlbumRequest{
				ID:       albumID,
				Title:    "New Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			model.Album{
				ID:       albumID,
				Title:    "Some Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"Artist has been added",
			requests.UpdateAlbumRequest{
				ID:       albumID,
				Title:    "New Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			model.Album{
				ID:    albumID,
				Title: "Some Album",
			},
		},
		{
			"Artist has been removed",
			requests.UpdateAlbumRequest{
				ID:    albumID,
				Title: "New Album",
			},
			model.Album{
				ID:       albumID,
				Title:    "Some Album",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			transactionManager := new(transaction.ManagerMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewUpdateAlbum(albumRepository, transactionManager, messagePublisherService)

			repositoryFactory := new(transaction.RepositoryFactoryMock)
			txAlbumRepo := new(repository.AlbumRepositoryMock)
			txSongRepo := new(repository.SongRepositoryMock)

			albumRepository.On("Get", new(model.Album), tt.request.ID).
				Return(nil, &tt.album).
				Once()

			repositoryFactory.On("NewAlbumRepository").Return(txAlbumRepo).Once()
			repositoryFactory.On("NewSongRepository").Return(txSongRepo).Once()
			transactionManager.On("Execute", mock.Anything).Return(nil, repositoryFactory).Once()

			txAlbumRepo.On("Update", mock.IsType(&tt.album)).
				Run(func(args mock.Arguments) {
					newAlbum := args.Get(0).(*model.Album)
					assertUpdatedAlbum(t, *newAlbum, tt.request)
				}).
				Return(nil).
				Once()

			songs := []model.Song{{ID: uuid.New()}, {ID: uuid.New()}, {ID: uuid.New()}}

			txSongRepo.On("GetAllByAlbum", new([]model.Song), tt.request.ID).
				Return(nil, &songs).
				Once()

			txSongRepo.On("UpdateAll", mock.IsType(&[]model.Song{})).
				Run(func(args mock.Arguments) {
					newSongs := args.Get(0).(*[]model.Song)
					for _, song := range *newSongs {
						assert.Equal(t, tt.request.ArtistID, song.ArtistID)
					}
				}).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, []uuid.UUID{tt.request.ID}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
			transactionManager.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
			repositoryFactory.AssertExpectations(t)
			txAlbumRepo.AssertExpectations(t)
			txSongRepo.AssertExpectations(t)
		})
	}
}

func assertUpdatedAlbum(
	t *testing.T,
	album model.Album,
	request requests.UpdateAlbumRequest,
) {
	assert.Equal(t, request.Title, album.Title)
	assert.Equal(t, request.ReleaseDate, album.ReleaseDate)
	assert.Equal(t, request.ArtistID, album.ArtistID)
}
