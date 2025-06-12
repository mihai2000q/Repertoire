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

func TestRemoveAlbumsFromArtist_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewRemoveAlbumsFromArtist(albumRepository, nil)

	request := requests.RemoveAlbumsFromArtistRequest{
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

func TestRemoveAlbumsFromArtist_WhenOneAlbumArtistDoesNotMatch_ShouldReturnConflictError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewRemoveAlbumsFromArtist(albumRepository, nil)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{{ID: request.AlbumIDs[0]}}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "album "+request.AlbumIDs[0].String()+" is not owned by this artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveAlbumsFromArtist_WhenUpdateAllAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := artist.NewRemoveAlbumsFromArtist(albumRepository, nil)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: &request.ID,
		},
	}
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

func TestRemoveAlbumsFromArtist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewRemoveAlbumsFromArtist(albumRepository, messagePublisherService)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: &request.ID,
		},
	}
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

func TestRemoveAlbumsFromArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewRemoveAlbumsFromArtist(albumRepository, messagePublisherService)

	request := requests.RemoveAlbumsFromArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	albums := []model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: &request.ID,
			Songs:    []model.Song{{ArtistID: &request.ID}, {ArtistID: &request.ID}, {ArtistID: &request.ID}},
		},
		{
			ID:       request.AlbumIDs[1],
			ArtistID: &request.ID,
		},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(&albums), request.AlbumIDs).
		Return(nil, &albums).
		Once()

	albumRepository.On("UpdateAllWithSongs", mock.IsType(&albums)).
		Run(func(args mock.Arguments) {
			newAlbums := args.Get(0).(*[]model.Album)
			for _, album := range *newAlbums {
				assert.Nil(t, album.ArtistID)
				for _, song := range album.Songs {
					assert.Nil(t, song.ArtistID)
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
