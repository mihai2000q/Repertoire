package artist

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestAddAlbumsToArtist_WhenGetAlbumsWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewAddAlbumsToArtist(albumRepository, nil)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongs", mock.Anything, request.AlbumIDs).
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

func TestAddAlbumsToArtist_WhenOneAlbumAlreadyHasArtist_ShouldReturnConflictError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewAddAlbumsToArtist(albumRepository, nil)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: &[]uuid.UUID{uuid.New()}[0],
		},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "album "+request.AlbumIDs[0].String()+" already has an artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToArtist_WhenUpdateAllAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewAddAlbumsToArtist(albumRepository, nil)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{{ID: request.AlbumIDs[0]}}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateAllWithSongs", mock.IsType(albums)).
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

func TestAddAlbumsToArtist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewAddAlbumsToArtist(albumRepository, messagePublisherService)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{{ID: request.AlbumIDs[0]}}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	albumRepository.On("UpdateAllWithSongs", mock.IsType(albums)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, request.AlbumIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestAddAlbumsToArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewAddAlbumsToArtist(albumRepository, messagePublisherService)

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: nil,
			Songs:    []model.Song{{}, {}, {}},
		},
		{
			ID:       request.AlbumIDs[1],
			ArtistID: nil,
		},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	albumRepository.On("UpdateAllWithSongs", mock.IsType(albums)).
		Run(func(args mock.Arguments) {
			newAlbums := args.Get(0).(*[]model.Album)
			for _, album := range *newAlbums {
				assert.Equal(t, &request.ID, album.ArtistID)
				for _, song := range album.Songs {
					assert.Equal(t, &request.ID, song.ArtistID)
				}
			}
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.AlbumsUpdatedTopic, request.AlbumIDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
